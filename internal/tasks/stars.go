package tasks

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

func loadExistingRawPayloads() [][]byte {
	var payloads [][]byte
	file, err := os.Open(starsRawPayloadsFile)
	if err != nil {
		return payloads
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Must copy bytes, scanner.Bytes() returns reused buffer
		payload := make([]byte, len(scanner.Bytes()))
		copy(payload, scanner.Bytes())
		payloads = append(payloads, payload)
	}
	return payloads
}

func loadExistingStarResults() map[string]StarResult {
	results := make(map[string]StarResult)
	file, err := os.Open(starsResultsFile)
	if err != nil {
		return results
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var result StarResult
		if err := json.Unmarshal(scanner.Bytes(), &result); err == nil {
			results[result.Language] = result
		}
	}
	return results
}

func toGitHubAPIURL(url string) string {
	if strings.HasPrefix(url, "https://github.com/") {
		return "https://api.github.com/repos/" + strings.TrimPrefix(url, "https://github.com/")
	}
	return url
}

type GithubData struct {
	Language   string  `json:"language"`
	Dirname    string  `json:"dirname"`
	GithubLink *string `json:"github_link"`
}

func Stars() error {
	// externals
	langInfo := languageInfoFile
	starsTempFile := starsTempResultsFile
	starsTempRawFile := starsTempRawPayloadsFile
	// ========================================
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	infile, err := os.Open(langInfo)
	if err != nil {
		return fmt.Errorf("open %s: %w", langInfo, err)
	}
	defer infile.Close()
	decoder := json.NewDecoder(infile)
	_, err = decoder.Token()
	if err != nil {
		return fmt.Errorf("error reading opening json array: %w", err)
	}

	writers, err := OpenBufferedFiles(starsTempFile, starsTempRawFile)
	if err != nil {
		return err
	}
	resultsFile, rawPayloadsFile := writers[0], writers[1]
	defer CloseBufferedFiles(writers...)

	existingResults := loadExistingStarResults()
	existingRawPayloads := loadExistingRawPayloads()
	payloadIdx := 0

	fmt.Print("stars\n")
	ticker := time.NewTicker(refreshInterval * time.Second)
	defer ticker.Stop()

	for decoder.More() {
		select {
		case <-ctx.Done():
			fmt.Println("Interrupt received, stopping stars loop...")
			return nil
		case <-ticker.C:
			for _, w := range writers {
				if err := w.Flush(); err != nil {
					return fmt.Errorf("flush: %w", err)
				}
			}
		default:
		}

		var ghData GithubData
		err := decoder.Decode(&ghData)
		if err != nil {
			return fmt.Errorf("error reading item: %w", err)
		}
		language := ghData.Language
		urlField := ghData.GithubLink

		if urlField == nil {
			continue
		}
		url := *urlField

		apiURL := toGitHubAPIURL(url)
		var existing *StarResult
		var lastFetchedAt *time.Time
		val, ok := existingResults[language]
		if ok {
			existing = &val
			lastFetched, err := time.Parse(timeFormat, val.LastFetchedAt)
			if err != nil {
				return err
			}
			lastFetchedAt = &lastFetched
		} else {
			existing = nil
			fmt.Printf("%s\n", language)
		}

		var result StarResult
		if existing == nil || (*lastFetchedAt).AddDate(0, 1, 0).Before(time.Now()) {
			fmt.Print(" *\n")
			repo, rawBody, err := getStarsWithRaw(apiURL)
			if err != nil {
				return fmt.Errorf("get stars for %s: %w", url, err)
			}
			if err := rawPayloadsFile.WriteString(string(rawBody) + "\n"); err != nil {
				return fmt.Errorf("write raw payload: %w", err)
			}

			createdTime, _ := time.Parse(time.RFC3339, repo.CreatedAt)
			result = StarResult{
				Language:      language,
				Dirname:       &ghData.Dirname,
				Stars:         repo.StargazersCount,
				RepoCreatedAt: createdTime.Format(timeFormat),
				LastFetchedAt: time.Now().Format(timeFormat),
			}
		} else {
			fmt.Print(".")
			// Use cached data, write old raw payload and update last_fetched_at
			if payloadIdx < len(existingRawPayloads) {
				// existingRawPayloads entries are raw bytes from scanner, no trailing newline
				if err := rawPayloadsFile.WriteString(string(existingRawPayloads[payloadIdx]) + "\n"); err != nil {
					return fmt.Errorf("write raw payload: %w", err)
				}
			}
			result = *existing
			result.LastFetchedAt = time.Now().Format(timeFormat)
		}
		payloadIdx++

		if err := resultsFile.Encode(result); err != nil {
			return fmt.Errorf("encode result: %w", err)
		}
	}

	_, err = decoder.Token()
	if err != nil {
		return fmt.Errorf("Error reading close bracket: %w", err)
	}
	fmt.Println()
	if err := CloseBufferedFiles(writers...); err != nil {
		return fmt.Errorf("close files: %w", err)
	}
	if err := os.Rename(starsTempRawPayloadsFile, starsRawPayloadsFile); err != nil {
		return fmt.Errorf("finalize raw payloads: %w", err)
	}
	if err := os.Rename(starsTempResultsFile, starsResultsFile); err != nil {
		return fmt.Errorf("finalize results: %w", err)
	}
	fmt.Printf("wrote: %s\n", starsResultsFile)
	return nil
}

func getStarsWithRaw(url string) (GithubRepo, []byte, error) {
	githubRepo := GithubRepo{}
	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return githubRepo, nil, fmt.Errorf("create request: %w", err)
	}
	token := os.Getenv("ghtoken")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Accept", "application/vnd.github+json")
	resp, err := client.Do(req)
	if err != nil {
		return githubRepo, nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return githubRepo, nil, fmt.Errorf("read body: %w", err)
	}
	if resp.StatusCode != 200 {
		return githubRepo, body, fmt.Errorf("http %d: %s", resp.StatusCode, string(body))
	}
	if err := json.Unmarshal(body, &githubRepo); err != nil {
		return githubRepo, body, fmt.Errorf("unmarshal: %w", err)
	}
	return githubRepo, body, nil
}

func getStars(url string) (GithubRepo, error) {
	repo, _, err := getStarsWithRaw(url)
	return repo, err
}
