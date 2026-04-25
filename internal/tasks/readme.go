package tasks

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func Readme() error {
	// Phase A — Open input and set up scanner (no output files yet)
	// ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	// defer stop()

	runtimes, err := loadJson[Dockerfile](dockerfileList)
	if err != nil {
		return err
	}
	jsonData, err := json.MarshalIndent(runtimes, "", "    ")
	fmt.Println(string(jsonData))

	years, err := loadCsvToMap(yearsFile)
	if err != nil {
		return err
	}
	jsonData, err = json.MarshalIndent(years, "", "    ")
	fmt.Println(string(jsonData))

	wiki, err := loadCsvToMap(wikiLinksFile)
	if err != nil {
		return err
	}
	jsonData, err = json.MarshalIndent(wiki, "", "    ")
	fmt.Println(string(jsonData))

	github, err := loadCsvToMap(githubLinksFile)
	if err != nil {
		return err
	}
	jsonData, err = json.MarshalIndent(github, "", "    ")
	fmt.Println(string(jsonData))

	stars, err := loadJson[StarResult](starsResultsFile)
	if err != nil {
		return err
	}
	jsonData, err = json.MarshalIndent(stars, "", "    ")
	fmt.Println(string(jsonData))

	sizes, err := loadJson[WeighResult](weighResultsFile)
	if err != nil {
		return err
	}
	jsonData, err = json.MarshalIndent(sizes, "", "    ")
	fmt.Println(string(jsonData))

	runs, err := loadJson[RunResult](runResultsFile)
	if err != nil {
		return err
	}
	jsonData, err = json.MarshalIndent(runs, "", "    ")
	fmt.Println(string(jsonData))

	return nil
}

// need everything to be based on tag
// load languages, runtimes based on tag?

func loadJson[T any](path string) ([]T, error) {
	results := make([]T, 0, 100)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
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

func loadCsvToMap(path string) (map[string]string, error) {
	results := make(map[string]string)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	if _, err := reader.Read(); err != nil {
		if err == io.EOF {
			return results, nil // File was empty
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
		key, val := record[0], record[1]
		if err != nil {
			return nil, err
		}
		results[key] = val
	}
	return results, nil
}
