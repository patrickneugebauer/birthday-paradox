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

func ShouldRefetchStars(now int64, lastFetched int64, createdAt int64) bool {
	if lastFetched == 0 || createdAt == 0 {
		return true
	}
	repoAge := now - createdAt
	if repoAge <= 0 {
		return true
	}
	timeSinceLastFetch := now - lastFetched
	ratio := float64(timeSinceLastFetch) / float64(repoAge)
	return ratio > 0.01
}

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

func Stars() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	infile, err := os.Open(languageInfoFile)
	if err != nil {
		return fmt.Errorf("open %s: %w", languageInfoFile, err)
	}
	defer infile.Close()
	scanner := bufio.NewScanner(infile)

	writers, err := OpenBufferedFiles(starsTempResultsFile, starsTempRawPayloadsFile)
	if err != nil {
		return err
	}
	resultsFile, rawPayloadsFile := writers[0], writers[1]
	defer CloseBufferedFiles(writers...)

	existingResults := loadExistingStarResults()
	existingRawPayloads := loadExistingRawPayloads()
	now := time.Now().Unix()
	payloadIdx := 0

	fmt.Print("stars ")
	ticker := time.NewTicker(refreshInterval * time.Second)
	defer ticker.Stop()

	for scanner.Scan() {
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

		line := scanner.Text()
		if line == "language,year,wiki,github,website" {
			continue
		}
		parts := strings.SplitN(line, ",", 5)
		if len(parts) < 4 {
			continue
		}
		language := strings.TrimSpace(parts[0])
		url := strings.TrimSpace(parts[3])

		if url == "" {
			return fmt.Errorf("empty github url for language: %s", language)
		}

		apiURL := toGitHubAPIURL(url)
		existing := existingResults[language]

		var result StarResult
		if ShouldRefetchStars(now, existing.LastFetchedAt, existing.RepoCreatedAt) {
			fmt.Print("*")
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
				Stars:         repo.StargazersCount,
				RepoCreatedAt: createdTime.Unix(),
				LastFetchedAt: now,
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
			result = existing
			result.LastFetchedAt = now
		}
		payloadIdx++

		if err := resultsFile.Encode(result); err != nil {
			return fmt.Errorf("encode result: %w", err)
		}
	}
	fmt.Println()

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner: %w", err)
	}
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
