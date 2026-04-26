package tasks

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"sort"
	"strings"
)

func Readme() error {
	// Load all data
	dockerfiles, err := loadJson[Dockerfile](dockerfileList)
	if err != nil {
		return err
	}

	languageInfo, err := loadLanguageInfo(languageInfoFile)
	if err != nil {
		return err
	}

	stars, err := loadJson[StarResult](starsResultsFile)
	if err != nil {
		return err
	}

	sizes, err := loadJson[WeighResult](weighResultsFile)
	if err != nil {
		return err
	}

	runs, err := loadJson[RunResult](runResultsFile)
	if err != nil {
		return err
	}

	// Build a map of tag → dockerfile info for quick lookup
	dockerfileMap := make(map[string]*Dockerfile)
	for i := range dockerfiles {
		dockerfileMap[dockerfiles[i].Tag] = &dockerfiles[i]
	}

	// Build a map of tag → size for quick lookup
	sizeMap := make(map[string]float64)
	for _, size := range sizes {
		sizeMap[size.Tag] = size.SizeMB
	}

	// Build a map of language → stars for quick lookup
	starMap := make(map[string]int)
	for _, star := range stars {
		starMap[star.Language] = star.Stars
	}

	// Build readme rows from runs (the main dataset)
	var rows []ReadmeRow
	for _, run := range runs {
		tag := run.Tag

		// Look up dockerfile/runtime
		dockerfileInfo, ok := dockerfileMap[tag]
		if !ok {
			continue // skip if no dockerfile for this tag
		}

		language := dockerfileInfo.Language
		runtime := ""
		if dockerfileInfo.Runtime != nil {
			runtime = *dockerfileInfo.Runtime
		}

		dataStructure := ""
		if dockerfileInfo.DataStructure != nil {
			dataStructure = *dockerfileInfo.DataStructure
		}

		executionMethod := ""
		if dockerfileInfo.ExecutionMethod != nil {
			executionMethod = *dockerfileInfo.ExecutionMethod
		}

		// Look up year, wiki, github from language info
		info := languageInfo[language]
		year := info.Year
		wikiURL := info.Wiki
		githubURL := info.GitHub

		// Look up stars
		starsCount := starMap[language]

		// Look up size
		sizeMB := sizeMap[tag]

		ips := 0
		if run.IPS != nil {
			ips = *run.IPS
		}

		row := ReadmeRow{
			Tag:             tag,
			Language:        language,
			Runtime:         runtime,
			DataStructure:   dataStructure,
			ExecutionMethod: executionMethod,
			Year:            year,
			WikiURL:         wikiURL,
			WikiDisplay:     extractWikiDisplay(wikiURL),
			GitHubURL:       githubURL,
			GitHubDisplay:   extractGitHubDisplay(githubURL),
			Stars:           starsCount,
			SizeMB:          sizeMB,
			IPS:             ips,
		}

		rows = append(rows, row)
	}

	// Sort rows by IPS descending (0 IPS goes to end)
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].IPS == 0 && rows[j].IPS == 0 {
			return false
		}
		if rows[i].IPS == 0 {
			return false
		}
		if rows[j].IPS == 0 {
			return true
		}
		return rows[i].IPS > rows[j].IPS
	})

	// Open output files
	jsonlFile, err := os.Create(readmeResultsFile)
	if err != nil {
		return err
	}
	defer jsonlFile.Close()

	markdownFile, err := os.Create(readmeFile)
	if err != nil {
		return err
	}
	defer markdownFile.Close()

	// Write markdown header
	fmt.Fprintln(markdownFile, "| Language | Runtime | Data Structure | Execution Method | Year | GitHub | Stars | Size (MB) | IPS |")
	fmt.Fprintln(markdownFile, "|---|---|---|---|---|---|---|---|---|")

	// Write rows to both files
	for _, row := range rows {
		// JSONL
		jsonBytes, _ := json.Marshal(row)
		fmt.Fprintln(jsonlFile, string(jsonBytes))

		// Markdown
		line := formatMarkdownRow(row)
		fmt.Fprintln(markdownFile, line)
	}

	return nil
}

func extractGitHubDisplay(url string) string {
	if url == "" {
		return ""
	}
	parts := strings.Split(url, "github.com/")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

func extractWikiDisplay(wikiURL string) string {
	if wikiURL == "" {
		return ""
	}
	parts := strings.Split(wikiURL, "/wiki/")
	if len(parts) > 1 {
		display := parts[1]
		decoded, err := url.QueryUnescape(display)
		if err != nil {
			return display
		}
		// Remove _(programming_language) suffix, handling fragments like #ReScript
		if idx := strings.Index(decoded, "#"); idx != -1 {
			beforeHash := decoded[:idx]
			beforeHash = strings.TrimSuffix(beforeHash, "_(programming_language)")
			decoded = beforeHash + decoded[idx:]
		} else {
			decoded = strings.TrimSuffix(decoded, "_(programming_language)")
		}
		return decoded
	}
	return ""
}

func formatMarkdownRow(row ReadmeRow) string {
	runtime := row.Runtime
	if runtime == "" {
		runtime = "-"
	}

	dataStructure := "-"
	if row.DataStructure != "" {
		dataStructure = row.DataStructure
	}

	executionMethod := "-"
	if row.ExecutionMethod != "" {
		executionMethod = row.ExecutionMethod
	}

	year := "-"
	if row.Year > 0 {
		year = fmt.Sprintf("%d", row.Year)
	}

	github := "-"
	if row.GitHubURL != "" {
		github = fmt.Sprintf("[%s](%s)", row.GitHubDisplay, row.GitHubURL)
	}

	stars := "-"
	if row.Stars > 0 {
		stars = formatWithCommas(row.Stars)
	}

	size := "-"
	if row.SizeMB > 0 {
		size = formatWithCommas(int(row.SizeMB))
	}

	ips := "-"
	if row.IPS > 0 {
		ips = formatWithCommas(row.IPS)
	}

	return fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s | %s | %s |",
		row.Language, runtime, dataStructure, executionMethod, year, github, stars, size, ips)
}

func formatWithCommas(n int) string {
	str := fmt.Sprintf("%d", n)
	if len(str) <= 3 {
		return str
	}

	// Add commas from right to left
	var result strings.Builder
	for i, c := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(c)
	}
	return result.String()
}

type LanguageInfo struct {
	Year   int
	Wiki   string
	GitHub string
}

func loadLanguageInfo(path string) (map[string]LanguageInfo, error) {
	results := make(map[string]LanguageInfo)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	if _, err := reader.Read(); err != nil {
		if err == io.EOF {
			return results, nil
		}
		return nil, err
	}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		language := record[0]
		year := 0
		if record[1] != "" {
			fmt.Sscanf(record[1], "%d", &year)
		}
		results[language] = LanguageInfo{
			Year:   year,
			Wiki:   record[2],
			GitHub: record[3],
		}
	}
	return results, nil
}

func loadJson[T any](path string) ([]T, error) {
	results := make([]T, 0, 100)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		bytes := scanner.Bytes()
		var result T
		if err := json.Unmarshal(bytes, &result); err != nil {
			return nil, fmt.Errorf("Error decoding line: %v", err)
		}
		results = append(results, result)
	}
	return results, nil
}