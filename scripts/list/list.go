package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	basePath := "./solutions"
	entries, err := os.ReadDir(basePath)
	if err != nil {
		log.Fatal(err)
	}
	dirs := make([]string, 0, len(entries))
	for _, v := range entries {
		if v.IsDir() {
			dirs = append(dirs, v.Name())
		}
	}
	fmt.Println(strings.Join(dirs, ", "))

	numDirs := len(dirs)
	type result struct {
		Dir        string `json:"dir"`
		Dockerfile string `json:"dockerfile"`
		Name       string `json:"name"`
	}
	results := make([]result, 0, numDirs*2)
	for _, dir := range dirs {
		entries, err := os.ReadDir(filepath.Join(basePath, dir))
		if err != nil {
			log.Fatal(err)
		}
		numFiles := len(entries)
		dockerfiles := make([]string, 0, numFiles)
		for _, entry := range entries {
			fname := entry.Name()
			if !entry.IsDir() && strings.HasPrefix(fname, "Dockerfile") {
				dockerfiles = append(dockerfiles, fname)
				_, after, found := strings.Cut(fname, ".")
				solutionName := dir
				if found {
					solutionName += "-" + after
				}
				result := result{
					Dir:        dir,
					Dockerfile: fname,
					Name:       solutionName,
				}
				results = append(results, result)
			}
		}
		fmt.Println(strings.Join(dockerfiles, ","))
	}
	fmt.Println(results)
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile("solutions.json", jsonData, 0644)

	f, err := os.Create("solutions.tsv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	// the type signature says it takes a io.Writer, but we pass a os.File
	w := csv.NewWriter(f)
	w.Comma = '\t'
	defer w.Flush()
	w.Write([]string{"dir", "dockerfile", "name"})
	rows := make([][]string, 0, len(results))
	for _, v := range results {
		row := []string{v.Dir, v.Dockerfile, v.Name}
		rows = append(rows, row)
	}
	fmt.Println(rows)
	w.WriteAll(rows)
}
