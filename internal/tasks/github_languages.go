package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
)

type ghSearchResponse struct {
	TotalCount        int               `json:"total_count"`
	IncompleteResults bool              `json:"incomplete_results"`
	Items             []json.RawMessage `json:"items"`
}

func GithubLanguages() error {
	token := os.Getenv("ghtoken")
	if token == "" {
		return fmt.Errorf("ghtoken environment variable not set")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	outDir := filepath.Join("artifacts", "github-languages")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	fmt.Print("github-languages ")

	total := 0
	for page := 1; page <= 10; page++ {
		select {
		case <-ctx.Done():
			fmt.Printf("\n%d repos\n", total)
			return nil
		default:
		}

		resp, _, err := fetchSearchPage(ctx, token, page)
		if err != nil {
			return fmt.Errorf("fetch page %d: %w", page, err)
		}

		outFile := filepath.Join(outDir, fmt.Sprintf("page%d.json", page))
		respBody, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal response: %w", err)
		}

		if err := os.WriteFile(outFile, respBody, 0644); err != nil {
			return fmt.Errorf("write %s: %w", outFile, err)
		}

		total += len(resp.Items)
		fmt.Print("*")

		if len(resp.Items) < 100 || page*100 >= resp.TotalCount {
			break
		}
	}

	fmt.Printf("\n%d repos\n", total)

	if err := generateReport(outDir); err != nil {
		return fmt.Errorf("generate report: %w", err)
	}

	return nil
}

func fetchSearchPage(ctx context.Context, token string, page int) (*ghSearchResponse, RateLimit, error) {
	url := fmt.Sprintf(
		"https://api.github.com/search/repositories?q=topic:programming-language&sort=stars&order=desc&per_page=100&page=%d",
		page,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, RateLimit{}, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, RateLimit{}, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	rateLimit := parseRateLimitHeaders(resp.Header)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, rateLimit, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, rateLimit, fmt.Errorf("http %d: %s", resp.StatusCode, string(body))
	}

	var searchResp ghSearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, rateLimit, fmt.Errorf("unmarshal response: %w", err)
	}

	return &searchResp, rateLimit, nil
}

type repoForReport struct {
	Name          string
	StarCount     int
	GitHubLink    string
}

func generateReport(outDir string) error {
	var repos []repoForReport

	for page := 1; page <= 10; page++ {
		filePath := filepath.Join(outDir, fmt.Sprintf("page%d.json", page))
		data, err := os.ReadFile(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				break
			}
			return fmt.Errorf("read %s: %w", filePath, err)
		}

		var resp ghSearchResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			return fmt.Errorf("unmarshal %s: %w", filePath, err)
		}

		for _, rawItem := range resp.Items {
			var item struct {
				Name            string `json:"name"`
				StargazersCount int    `json:"stargazers_count"`
				HtmlURL         string `json:"html_url"`
			}
			if err := json.Unmarshal(rawItem, &item); err != nil {
				continue
			}
			repos = append(repos, repoForReport{
				Name:       item.Name,
				StarCount:  item.StargazersCount,
				GitHubLink: item.HtmlURL,
			})
		}
	}

	sort.Slice(repos, func(i, j int) bool {
		return repos[i].StarCount > repos[j].StarCount
	})

	var report strings.Builder
	report.WriteString("name\tstar count\tgithub link\n")
	for _, r := range repos {
		fmt.Fprintf(&report, "%s\t%d\t%s\n", r.Name, r.StarCount, r.GitHubLink)
	}

	reportPath := filepath.Join(outDir, "report.tsv")
	if err := os.WriteFile(reportPath, []byte(report.String()), 0644); err != nil {
		return fmt.Errorf("write report: %w", err)
	}

	return nil
}
