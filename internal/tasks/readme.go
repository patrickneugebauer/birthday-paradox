package tasks

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/url"
	"os"
	"sort"
	"strings"
)

func Readme() error {
	fmt.Print("readme\n")
	// Stage 1: Build raw results
	rows, err := buildRawReadmeResults()
	if err != nil {
		return err
	}

	// Stage 2: Build formatted results
	formattedRows, err := buildFormattedReadmeResults(rows)
	if err != nil {
		return err
	}

	// Stage 3: Generate markdown files
	if err := generateReadmeMarkdownFiles(formattedRows); err != nil {
		return err
	}

	return nil
}

func buildRawReadmeResults() ([]ReadmeRow, error) {
	// Load all data
	dockerfiles, err := loadJson[Dockerfile](dockerfileList)
	if err != nil {
		return nil, err
	}

	languageInfo, err := loadLanguageInfo(languageInfoFile)
	if err != nil {
		return nil, err
	}

	stars, err := loadJson[StarResult](starsResultsFile)
	if err != nil {
		return nil, err
	}

	sizes, err := loadJson[WeighResult](weighResultsFile)
	if err != nil {
		return nil, err
	}

	runs, err := loadJson[RunResult](runResultsFile)
	if err != nil {
		return nil, err
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
			Runtime:         dockerfileInfo.Runtime,
			DataStructure:   dockerfileInfo.DataStructure,
			ExecutionMethod: dockerfileInfo.ExecutionMethod,
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

	var buf bytes.Buffer
	for _, row := range rows {
		roundedRow := row
		roundedRow.SizeMB = math.Round(row.SizeMB*1000) / 1000
		jsonBytes, _ := json.Marshal(roundedRow)
		buf.Write(jsonBytes)
		buf.WriteByte('\n')
	}
	if err := os.WriteFile(readmeResultsFile, buf.Bytes(), 0644); err != nil {
		return nil, err
	}

	return rows, nil
}

func buildFormattedReadmeResults(rows []ReadmeRow) ([]ReadmeRow, error) {
	// Apply all formatting to rows
	formattedRows := make([]ReadmeRow, len(rows))
	for i, row := range rows {
		formattedRows[i] = row
		// Round SizeMB to 3 decimals
		formattedRows[i].SizeMB = math.Round(row.SizeMB*1000) / 1000

		// Format language with wiki link
		language := row.Language
		if row.WikiURL != "" {
			language = fmt.Sprintf("[%s](%s)", row.Language, row.WikiURL)
		}
		formattedRows[i].FormattedLanguage = language

		// Format runtime
		runtime := "-"
		if row.Runtime != nil {
			runtime = *row.Runtime
		}
		formattedRows[i].FormattedRuntime = runtime

		// Format data structure
		dataStructure := "-"
		if row.DataStructure != nil {
			dataStructure = *row.DataStructure
		}
		formattedRows[i].FormattedDataType = dataStructure

		// Format execution method
		executionMethod := "-"
		if row.ExecutionMethod != nil {
			executionMethod = *row.ExecutionMethod
		}
		formattedRows[i].FormattedExecMode = executionMethod

		// Format year
		year := "-"
		if row.Year > 0 {
			year = fmt.Sprintf("%d", row.Year)
		}
		formattedRows[i].FormattedYear = year

		// Format stars with github link
		stars := "-"
		if row.Stars > 0 {
			starStr := formatWithCommas(row.Stars)
			if row.GitHubURL != "" {
				stars = fmt.Sprintf("[%s](%s)", starStr, row.GitHubURL)
			} else {
				stars = starStr
			}
		}
		formattedRows[i].FormattedStars = stars

		// Format size MB - show at least 3 digits
		size := "-"
		if row.SizeMB > 0 {
			if row.SizeMB < 10 {
				size = fmt.Sprintf("%.2f", row.SizeMB)
			} else if row.SizeMB < 100 {
				size = fmt.Sprintf("%.1f", row.SizeMB)
			} else {
				size = formatWithCommas(int(row.SizeMB))
			}
		}
		formattedRows[i].FormattedSizeMB = size

		// Format IPS
		ips := "-"
		if row.IPS > 0 {
			ips = formatWithCommas(row.IPS)
		}
		formattedRows[i].FormattedIPS = ips
	}

	var buf bytes.Buffer
	for _, row := range formattedRows {
		jsonBytes, _ := json.Marshal(row)
		buf.Write(jsonBytes)
		buf.WriteByte('\n')
	}
	if err := os.WriteFile(formattedReadmeResultsFile, buf.Bytes(), 0644); err != nil {
		return nil, err
	}

	return formattedRows, nil
}

func generateReadmeMarkdownFiles(rows []ReadmeRow) error {
	// Sort by IPS descending (0 IPS goes to end)
	sortByIPS(rows)

	// Create sorted copies for each view
	rowsByLanguage := append([]ReadmeRow{}, rows...)
	sortByLanguage(rowsByLanguage)

	rowsByYear := append([]ReadmeRow{}, rows...)
	sortByYear(rowsByYear)

	rowsByStars := append([]ReadmeRow{}, rows...)
	sortByStars(rowsByStars)

	rowsBySize := append([]ReadmeRow{}, rows...)
	sortBySize(rowsBySize)

	// Write all markdown files
	if err := writeReadmeFile(readmeFile, "Performance (IPS, highest first)", rows); err != nil {
		return err
	}
	if err := writeReadmeFile(readmeFileByLanguage, "Language (A-Z)", rowsByLanguage); err != nil {
		return err
	}
	if err := writeReadmeFile(readmeFileByYear, "Year (newest first)", rowsByYear); err != nil {
		return err
	}
	if err := writeReadmeFile(readmeFileByStars, "GitHub Stars (highest first)", rowsByStars); err != nil {
		return err
	}
	if err := writeReadmeFile(readmeFileBySize, "Size in MB (smallest first)", rowsBySize); err != nil {
		return err
	}

	jsRows := filterByLanguageAndRuntime(rows, "javascript", "node")
	if err := writeLanguageReadmeFile(readmeFileJavaScript, "JavaScript", jsRows); err != nil {
		return err
	}

	if err := generateBuildResultsReadme(); err != nil {
		return err
	}

	if err := generateRunResultsReadme(); err != nil {
		return err
	}

	fmt.Printf("wrote: %s\n", readmeFile)
	return nil
}

func sortByIPS(rows []ReadmeRow) {
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
}

func sortByLanguage(rows []ReadmeRow) {
	sort.Slice(rows, func(i, j int) bool {
		return strings.ToLower(rows[i].Language) < strings.ToLower(rows[j].Language)
	})
}

func sortByYear(rows []ReadmeRow) {
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Year == rows[j].Year {
			return strings.ToLower(rows[i].Language) < strings.ToLower(rows[j].Language)
		}
		return rows[i].Year > rows[j].Year
	})
}

func sortByStars(rows []ReadmeRow) {
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Stars == rows[j].Stars {
			return strings.ToLower(rows[i].Language) < strings.ToLower(rows[j].Language)
		}
		return rows[i].Stars > rows[j].Stars
	})
}

func sortBySize(rows []ReadmeRow) {
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].SizeMB == rows[j].SizeMB {
			return strings.ToLower(rows[i].Language) < strings.ToLower(rows[j].Language)
		}
		return rows[i].SizeMB < rows[j].SizeMB
	})
}

func writeReadmeFile(filename, currentSort string, rows []ReadmeRow) error {
	tmpFilename := filename + ".tmp"
	file, err := os.Create(tmpFilename)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	defer func() {
		writer.Flush()
		file.Close()
	}()

	// Build navigation links based on whether this is root README or tables/
	isRoot := filename == readmeFile
	perfLink := "readme.md"
	langLink := "results-by-language.md"
	yearLink := "results-by-year.md"
	starsLink := "results-by-stars.md"
	sizeLink := "results-by-size.md"
	jsLink := ""
	buildLink := ""
	runLink := ""

	if !isRoot {
		perfLink = "../readme.md"
		buildLink = "docker-build-results.md"
		runLink = "docker-run-results.md"
	} else {
		langLink = "tables/results-by-language.md"
		yearLink = "tables/results-by-year.md"
		starsLink = "tables/results-by-stars.md"
		sizeLink = "tables/results-by-size.md"
		jsLink = "tables/javascript-readme.md"
		buildLink = "tables/docker-build-results.md"
		runLink = "tables/docker-run-results.md"
	}

	// Write navigation bar
	nav := fmt.Sprintf("**View by:** | [Performance](%s) | [Language](%s) | [Year](%s) | [Stars](%s) | [Size](%s) |", perfLink, langLink, yearLink, starsLink, sizeLink)
	if jsLink != "" {
		nav += fmt.Sprintf(" [JavaScript](%s) |", jsLink)
	}
	nav += fmt.Sprintf(" [Build Stats](%s) | [Run Stats](%s) |", buildLink, runLink)
	fmt.Fprintln(writer, nav)
	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "| | Language | Runtime | Data Type | Exec Mode | Year | Stars | Size (MB) | IPS |")
	fmt.Fprintln(writer, "|---|---|---|---|---|---|---|---|---|")

	// Write rows with index
	for i, row := range rows {
		line := formatMarkdownRow(row)
		fmt.Fprintf(writer, "| %d %s\n", i+1, line)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("flush: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close temp: %w", err)
	}
	if err := os.Rename(tmpFilename, filename); err != nil {
		return fmt.Errorf("finalize %s: %w", filename, err)
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

func generateBuildResultsReadme() error {
	dockerfiles, err := loadJson[Dockerfile](dockerfileList)
	if err != nil {
		return err
	}
	builds, err := loadJson[BuildResult](buildResultsFile)
	if err != nil {
		return err
	}
	sizes, err := loadJson[WeighResult](weighResultsFile)
	if err != nil {
		return err
	}

	dockerfileMap := make(map[string]*Dockerfile)
	for i := range dockerfiles {
		dockerfileMap[dockerfiles[i].Tag] = &dockerfiles[i]
	}

	sizeMap := make(map[string]float64)
	for _, size := range sizes {
		sizeMap[size.Tag] = size.SizeMB
	}

	sort.Slice(builds, func(i, j int) bool {
		sizeI := sizeMap[builds[i].Tag]
		sizeJ := sizeMap[builds[j].Tag]
		if sizeI == sizeJ {
			return strings.ToLower(dockerfileMap[builds[i].Tag].Language) < strings.ToLower(dockerfileMap[builds[j].Tag].Language)
		}
		return sizeI < sizeJ
	})

	tmpFilename := dockerBuildReadme + ".tmp"
	file, err := os.Create(tmpFilename)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	defer func() {
		writer.Flush()
		file.Close()
	}()

	// Write navigation bar
	nav := "**View by:** | [Performance](../readme.md) | [Language](results-by-language.md) | [Year](results-by-year.md) | [Stars](results-by-stars.md) | [Size](results-by-size.md) | [Build Stats](docker-build-results.md) | [Run Stats](docker-run-results.md) |"
	fmt.Fprintln(writer, nav)
	fmt.Fprintln(writer)

	fmt.Fprintln(writer, "| | Language | Runtime | Exec Mode | Size (MB) | Build Time (s) | Net Activity (MB) | Disk Activity (MB) | CPU (s) |")
	fmt.Fprintln(writer, "|---|---|---|---|---|---|---|---|---|")

	for i, build := range builds {
		df := dockerfileMap[build.Tag]
		if df == nil {
			continue
		}

		language := df.Language
		runtime := "-"
		if df.Runtime != nil {
			runtime = *df.Runtime
		}
		execMode := "-"
		if df.ExecutionMethod != nil {
			execMode = *df.ExecutionMethod
		}

		size := "-"
		if sizeMB, ok := sizeMap[build.Tag]; ok && sizeMB > 0 {
			if sizeMB < 10 {
				size = fmt.Sprintf("%.2f", sizeMB)
			} else if sizeMB < 100 {
				size = fmt.Sprintf("%.1f", sizeMB)
			} else {
				size = fmt.Sprintf("%.0f", sizeMB)
			}
		}

		buildTime := "-"
		if build.TotalS != nil {
			buildTime = fmt.Sprintf("%.2f", *build.TotalS)
		}

		netActivity := "-"
		if build.NetRxBytes != nil || build.NetTxBytes != nil {
			total := int64(0)
			if build.NetRxBytes != nil {
				total += *build.NetRxBytes
			}
			if build.NetTxBytes != nil {
				total += *build.NetTxBytes
			}
			netActivity = fmt.Sprintf("%.2f", float64(total)/(1024*1024))
		}

		diskActivity := "-"
		if build.BlkReadBytes != nil || build.BlkWriteBytes != nil {
			total := int64(0)
			if build.BlkReadBytes != nil {
				total += *build.BlkReadBytes
			}
			if build.BlkWriteBytes != nil {
				total += *build.BlkWriteBytes
			}
			diskActivity = fmt.Sprintf("%.0f", float64(total)/(1024*1024))
		}

		cpuTime := "-"
		if build.CpuS != nil {
			cpuTime = fmt.Sprintf("%.2f", *build.CpuS)
		}

		fmt.Fprintf(writer, "| %d | %s | %s | %s | %s | %s | %s | %s | %s |\n",
			i+1, language, runtime, execMode, size, buildTime, netActivity, diskActivity, cpuTime)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("flush: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close temp: %w", err)
	}
	if err := os.Rename(tmpFilename, dockerBuildReadme); err != nil {
		return fmt.Errorf("finalize %s: %w", dockerBuildReadme, err)
	}
	fmt.Printf("wrote: %s\n", dockerBuildReadme)
	return nil
}

func generateRunResultsReadme() error {
	dockerfiles, err := loadJson[Dockerfile](dockerfileList)
	if err != nil {
		return err
	}
	runs, err := loadJson[RunResult](runResultsFile)
	if err != nil {
		return err
	}

	dockerfileMap := make(map[string]*Dockerfile)
	for i := range dockerfiles {
		dockerfileMap[dockerfiles[i].Tag] = &dockerfiles[i]
	}

	sort.Slice(runs, func(i, j int) bool {
		ramI := int64(0)
		if runs[i].PeakRAMBytes != nil {
			ramI = *runs[i].PeakRAMBytes
		}
		ramJ := int64(0)
		if runs[j].PeakRAMBytes != nil {
			ramJ = *runs[j].PeakRAMBytes
		}
		if ramI == ramJ {
			return strings.ToLower(dockerfileMap[runs[i].Tag].Language) < strings.ToLower(dockerfileMap[runs[j].Tag].Language)
		}
		return ramI < ramJ
	})

	tmpFilename := dockerRunReadme + ".tmp"
	file, err := os.Create(tmpFilename)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	defer func() {
		writer.Flush()
		file.Close()
	}()

	// Write navigation bar
	nav := "**View by:** | [Performance](../readme.md) | [Language](results-by-language.md) | [Year](results-by-year.md) | [Stars](results-by-stars.md) | [Size](results-by-size.md) | [Build Stats](docker-build-results.md) | [Run Stats](docker-run-results.md) |"
	fmt.Fprintln(writer, nav)
	fmt.Fprintln(writer)

	fmt.Fprintln(writer, "| | Language | Runtime | Data Structure | Exec Mode | Seconds | Runtime (s) | Infra (s) | Peak RAM (MB) | CPU (s) |")
	fmt.Fprintln(writer, "|---|---|---|---|---|---|---|---|---|---|")

	for i, run := range runs {
		df := dockerfileMap[run.Tag]
		if df == nil {
			continue
		}

		language := df.Language
		runtime := "-"
		if df.Runtime != nil {
			runtime = *df.Runtime
		}
		dataStructure := "-"
		if df.DataStructure != nil {
			dataStructure = *df.DataStructure
		}
		execMode := "-"
		if df.ExecutionMethod != nil {
			execMode = *df.ExecutionMethod
		}

		seconds := "-"
		if run.Seconds != nil {
			seconds = fmt.Sprintf("%.1f", *run.Seconds)
		}

		runtimeS := "-"
		if run.RuntimeS != nil {
			runtimeS = fmt.Sprintf("%.1f", *run.RuntimeS)
		}

		infraTime := "-"
		if run.RuntimeS != nil && run.Seconds != nil {
			infra := *run.RuntimeS - *run.Seconds
			infraTime = fmt.Sprintf("%.1f", infra)
		}

		peakRAM := "-"
		if run.PeakRAMBytes != nil {
			mb := float64(*run.PeakRAMBytes) / (1024 * 1024)
			peakRAM = fmt.Sprintf("%.2f", mb)
		}

		cpuS := "-"
		if run.CpuS != nil {
			cpuS = fmt.Sprintf("%.2f", *run.CpuS)
		}

		fmt.Fprintf(writer, "| %d | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",
			i+1, language, runtime, dataStructure, execMode, seconds, runtimeS, infraTime, peakRAM, cpuS)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("flush: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close temp: %w", err)
	}
	if err := os.Rename(tmpFilename, dockerRunReadme); err != nil {
		return fmt.Errorf("finalize %s: %w", dockerRunReadme, err)
	}
	fmt.Printf("wrote: %s\n", dockerRunReadme)
	return nil
}

func formatMarkdownRow(row ReadmeRow) string {
	return fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s | %s |",
		row.FormattedLanguage, row.FormattedRuntime, row.FormattedDataType, row.FormattedExecMode,
		row.FormattedYear, row.FormattedStars, row.FormattedSizeMB, row.FormattedIPS)
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

func filterByLanguageAndRuntime(rows []ReadmeRow, language, runtime string) []ReadmeRow {
	var filtered []ReadmeRow
	for _, row := range rows {
		if row.Language == language && row.Runtime != nil && *row.Runtime == runtime {
			filtered = append(filtered, row)
		}
	}
	return filtered
}

func writeLanguageReadmeFile(filename, language string, rows []ReadmeRow) error {
	tmpFilename := filename + ".tmp"
	file, err := os.Create(tmpFilename)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	defer func() {
		writer.Flush()
		file.Close()
	}()

	fmt.Fprintf(writer, "**[← All languages](../../readme.md)**\n\n")
	fmt.Fprintf(writer, "# %s\n\n", language)
	fmt.Fprintln(writer, "| | Data Type | IPS |")
	fmt.Fprintln(writer, "|---|---|---|")

	for i, row := range rows {
		fmt.Fprintf(writer, "| %d | %s | %s |\n",
			i+1, row.FormattedDataType, row.FormattedIPS)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("flush: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close temp: %w", err)
	}
	return os.Rename(tmpFilename, filename)
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
