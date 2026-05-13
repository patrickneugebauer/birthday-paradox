package tasks

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
)

type bbSearchResponse struct {
	Pagelen int               `json:"pagelen"`
	Size    int               `json:"size"`
	Page    int               `json:"page"`
	Next    string            `json:"next"`
	Values  []json.RawMessage `json:"values"`
}

func BitbucketLanguages() error {
	token := os.Getenv("bbtoken")
	if token == "" {
		return fmt.Errorf("bbtoken environment variable not set")
	}

	auth := base64.StdEncoding.EncodeToString([]byte(token))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	outDir := filepath.Join("artifacts", "bitbucket")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	// Load languages from wiki artifacts
	wikiDir := filepath.Join("artifacts", "wiki")
	langSet := make(map[string]bool)

	entries, err := os.ReadDir(wikiDir)
	if err != nil {
		return fmt.Errorf("read wiki directory: %w", err)
	}

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(wikiDir, entry.Name()))
		if err != nil {
			return fmt.Errorf("read %s: %w", entry.Name(), err)
		}

		var resp wikiCategoryResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			continue
		}

		for _, page := range resp.Query.CategoryMembers {
			langName := wikiTitleToLangName(page.Title)
			if langName != "" {
				langSet[langName] = true
			}
		}
	}

	// Convert to slice and limit to 1000
	var languages []string
	for lang := range langSet {
		languages = append(languages, lang)
		if len(languages) >= 1000 {
			break
		}
	}

	fmt.Print("bitbucket-languages ")

	processed := 0
	skipped := 0

	for _, langName := range languages {
		select {
		case <-ctx.Done():
			fmt.Printf("\n%d processed, %d skipped\n", processed, skipped)
			return nil
		default:
		}

		resp, err := fetchBitbucketRepos(ctx, langName, auth)
		if err != nil {
			return fmt.Errorf("fetch bitbucket repos for %s: %w", langName, err)
		}

		outFile := filepath.Join(outDir, langName+".json")
		respBody, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal response: %w", err)
		}

		if err := os.WriteFile(outFile, respBody, 0644); err != nil {
			return fmt.Errorf("write %s: %w", outFile, err)
		}

		if resp.Size == 0 {
			fmt.Print(".")
			skipped++
			continue
		}

		fmt.Print("*")
		processed++
	}

	fmt.Printf("\n%d processed, %d skipped\n", processed, skipped)

	return nil
}

func fetchBitbucketRepos(ctx context.Context, language string, auth string) (*bbSearchResponse, error) {
	params := url.Values{
		"q":       {fmt.Sprintf("language=\"%s\"", language)},
		"pagelen": {"1"},
	}

	apiURL := "https://api.bitbucket.org/2.0/repositories?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Basic "+auth)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http %d: %s", resp.StatusCode, string(body))
	}

	var bbResp bbSearchResponse
	if err := json.Unmarshal(body, &bbResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &bbResp, nil
}
