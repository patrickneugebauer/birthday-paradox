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

func toGitHubAPIURL(url string) string {
	if strings.HasPrefix(url, "https://github.com/") {
		return "https://api.github.com/repos/" + strings.TrimPrefix(url, "https://github.com/")
	}
	return url
}

func Stars() error {
	// Phase A — Open input and set up scanner (no output files yet)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	infile, err := os.Open(languageInfoFile)
	if err != nil {
		return fmt.Errorf("open %s: %w", languageInfoFile, err)
	}
	defer infile.Close()
	scanner := bufio.NewScanner(infile)

	// Phase B — Open all output files
	writers, err := OpenBufferedFiles(starsCommandsFile, starsDataFile, starsRawPayloadsFile, starsTempResultsFile)
	if err != nil {
		return err
	}
	cmdFile, dataFile, rawPayloadsFile, resultsFile := writers[0], writers[1], writers[2], writers[3]
	defer CloseBufferedFiles(writers...)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// Phase C — Main loop
	for scanner.Scan() {
		fmt.Print(".")
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
		language := parts[0]
		url := parts[3]

		apiURL := toGitHubAPIURL(url)
		headers := `"Accept: application/vnd.github+json" -H "Authorization: Bearer $ghtoken"`
		if err := cmdFile.WriteString(fmt.Sprintf("curl -i -H %s %s\n", headers, apiURL)); err != nil {
			return fmt.Errorf("write command: %w", err)
		}

		repo, rawBody, err := getStarsWithRaw(apiURL)
		if err != nil {
			return fmt.Errorf("get stars for %s: %w", url, err)
		}
		if err := rawPayloadsFile.WriteString(string(rawBody) + "\n"); err != nil {
			return fmt.Errorf("write raw payload: %w", err)
		}
		if err := dataFile.Encode(repo); err != nil {
			return fmt.Errorf("encode repo data: %w", err)
		}
		if err := resultsFile.Encode(StarResult{Language: language, Stars: repo.StargazersCount}); err != nil {
			return fmt.Errorf("encode result: %w", err)
		}
	}
	fmt.Println()

	// Phase D — Check scanner, close results, rename temp to final
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner: %w", err)
	}
	if err := resultsFile.Close(); err != nil {
		return fmt.Errorf("close temp results: %w", err)
	}
	if err := os.Rename(starsTempResultsFile, starsResultsFile); err != nil {
		return fmt.Errorf("finalize results: %w", err)
	}
	fmt.Printf("wrote to file: %s\n", starsResultsFile)
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
