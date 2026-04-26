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

	// Phase B — Open all three output files
	writers, err := OpenBufferedFiles(starsCommandsFile, starsDataFile, starsTempResultsFile)
	if err != nil {
		return err
	}
	cmdFile, dataFile, resultsFile := writers[0], writers[1], writers[2]
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
		if line == "language,year,wiki,github" {
			continue
		}
		parts := strings.SplitN(line, ",", 4)
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

		repo, err := getStars(apiURL)
		if err != nil {
			return fmt.Errorf("get stars for %s: %w", url, err)
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

func getStars(url string) (GithubRepo, error) {
	githubRepo := GithubRepo{}
	// curl -i -H "Accept: application/vnd.github+json" \
	// 	-H "Authorization: Bearer $ghtoken" \
	// 	https://api.github.com/repos/golang/go
	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return githubRepo, fmt.Errorf("create request: %w", err)
	}
	token := os.Getenv("ghtoken")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Accept", "application/vnd.github+json")
	resp, err := client.Do(req)
	if err != nil {
		return githubRepo, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return githubRepo, fmt.Errorf("http %d: %s", resp.StatusCode, string(body))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return githubRepo, fmt.Errorf("read body: %w", err)
	}
	if err := json.Unmarshal(body, &githubRepo); err != nil {
		return githubRepo, fmt.Errorf("unmarshal: %w", err)
	}
	return githubRepo, nil
}
