package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
)

type wikiCategoryResponse struct {
	Continue *struct {
		CmContinue string `json:"cmcontinue"`
		Continue   string `json:"continue"`
	} `json:"continue"`
	Query struct {
		CategoryMembers []wikiPage `json:"categorymembers"`
	} `json:"query"`
}

type wikiPage struct {
	PageID int    `json:"pageid"`
	NS     int    `json:"ns"`
	Title  string `json:"title"`
}

var wikiCategories = []string{
	"Category:Programming_languages",
	"Category:Functional_languages",
	"Category:Object-oriented_programming_languages",
	"Category:Scripting_languages",
	"Category:Assembly_languages",
	"Category:Concurrent_programming_languages",
	"Category:Educational_programming_languages",
	"Category:Logic_programming_languages",
	"Category:Declarative_programming_languages",
	"Category:Procedural_programming_languages",
	"Category:Esoteric_programming_languages",
	"Category:Programming_language_families",
	"Category:Domain-specific_programming_languages",
}

func WikiLanguages() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	outDir := filepath.Join("artifacts", "wiki")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	fmt.Print("wiki-languages ")

	total := 0
	for i, category := range wikiCategories {
		select {
		case <-ctx.Done():
			fmt.Printf("\n%d languages\n", total)
			return nil
		default:
		}

		cmContinue := ""
		pageNum := 1

		for {
			resp, err := fetchWikiCategoryPage(ctx, category, cmContinue)
			if err != nil {
				return fmt.Errorf("fetch %s page %d: %w", category, pageNum, err)
			}

			outFile := filepath.Join(outDir, fmt.Sprintf("%s-page%d.json", strings.ReplaceAll(category, ":", ""), pageNum))
			respBody, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				return fmt.Errorf("marshal response: %w", err)
			}

			if err := os.WriteFile(outFile, respBody, 0644); err != nil {
				return fmt.Errorf("write %s: %w", outFile, err)
			}

			total += len(resp.Query.CategoryMembers)
			fmt.Print("*")
			pageNum++

			if resp.Continue == nil || resp.Continue.CmContinue == "" {
				break
			}
			cmContinue = resp.Continue.CmContinue
			time.Sleep(2 * time.Second)
		}

		if i < len(wikiCategories)-1 {
			time.Sleep(2 * time.Second)
		}
	}

	fmt.Printf("\n%d languages\n", total)

	return nil
}

func fetchWikiCategoryPage(ctx context.Context, category string, cmContinue string) (*wikiCategoryResponse, error) {
	params := url.Values{
		"action":   {"query"},
		"list":     {"categorymembers"},
		"cmtitle":  {category},
		"cmtype":   {"page"},
		"cmlimit":  {"500"},
		"cmprop":   {"ids|title|type"},
		"format":   {"json"},
	}
	if cmContinue != "" {
		params.Set("cmcontinue", cmContinue)
	}

	url := "https://en.wikipedia.org/w/api.php?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

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

	var wikiResp wikiCategoryResponse
	if err := json.Unmarshal(body, &wikiResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &wikiResp, nil
}

func wikiTitleToLangName(title string) string {
	name := strings.TrimSuffix(title, " (programming language)")
	name = strings.TrimSuffix(name, " programming language")
	name = strings.TrimSuffix(name, " (language)")
	name = strings.ToLower(name)
	name = strings.TrimSpace(name)
	return name
}
