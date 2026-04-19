package tasks

import (
	"bufio"
	"cmp"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"go.yaml.in/yaml/v4"
)

type DockerfileMap = map[string][]string
type Dockerfile struct {
	Language string `json:"language"`
	Filename string `json:"dockerfile"`
}

func MakeFileTree() error {
	// get map of dockerfiles
	dfMap, err := getFileStructure()
	if err != nil {
		return err
	}
	// write map
	yData, err := yaml.Marshal(dfMap)
	if err != nil {
		return fmt.Errorf("failed to marshal yaml: %w", err)
	}
	err = os.WriteFile(dockerfileMap, yData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write dockerfile map: %w", err)
	}
	fmt.Printf("wrote to file: %s\n", dockerfileMap)
	// write list
	languages, err := dockerfileMapToList(dfMap)
	if err != nil {
		return fmt.Errorf("failed to write create dflist: %w", err)
	}
	rows, err := listToJsonlstrings(languages)
	if err != nil {
		return fmt.Errorf("failed to marshal jsonl: %w", err)
	}
	text := fmt.Sprintf("%s\n", strings.Join(rows, "\n"))
	err = os.WriteFile(dockerfileList, []byte(text), 0644)
	if err != nil {
		return fmt.Errorf("failed to write dockerfile list: %w", err)
	}
	fmt.Printf("wrote to file: %s\n", dockerfileList)
	// return
	return nil
}

func getFileStructure() (DockerfileMap, error) {
	root := solutionsDir
	dfMap := make(map[string][]string)
	// read root
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	// read subdirs
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dirPath := filepath.Join(root, entry.Name())
		files, _ := os.ReadDir(dirPath)
		for _, f := range files {
			if !f.IsDir() && strings.HasPrefix(f.Name(), "Dockerfile") {
				dfMap[entry.Name()] = append(dfMap[entry.Name()], f.Name())
			}
		}
	}
	return dfMap, nil
}

func dockerfileMapToList(dfMap DockerfileMap) ([]Dockerfile, error) {
	languages := make([]Dockerfile, 0, len(dfMap))
	for lang, dfList := range dfMap {
		for _, v := range dfList {
			dockerfile := Dockerfile{Language: lang, Filename: v}
			languages = append(languages, dockerfile)
		}
	}
	slices.SortFunc(languages, func(a, b Dockerfile) int {
		return cmp.Compare(a.Language, b.Language)
	})
	return languages, nil
}

func listToJsonlstrings[T any](list []T) ([]string, error) {
	results := make([]string, 0, len(list))
	for _, v := range list {
		row, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		results = append(results, string(row))
	}
	return results, nil
}

func ReadDockerfileMap() (DockerfileMap, error) {
	fname := dockerfileMap
	_, err := os.Stat(fname)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("%s not found: please run 'pre-build'", fname)
	}
	content, err := os.ReadFile(fname)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", fname, err)
	}
	var dockerFileMap DockerfileMap
	yaml.Unmarshal(content, &dockerFileMap)
	return dockerFileMap, nil
}

func ReadDockerfileList() ([]Dockerfile, error) {
	fname := dockerfileList
	list, err := readJsonl[Dockerfile](fname)
	if err != nil {
		return nil, err
	}
	slices.SortFunc(list, func(a, b Dockerfile) int {
		return cmp.Compare(a.Language, b.Language)
	})
	return list, nil
}

func readJsonl[T any](fname string) ([]T, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", fname, err)
	}
	initialCapacity := capacity
	var results = make([]T, 0, initialCapacity)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		bytes := scanner.Bytes()
		var result T
		err := json.Unmarshal(bytes, &result)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshall %w", err)
		}
		results = append(results, result)
	}
	if cap(results) > initialCapacity {
		warning := fmt.Sprintf("SLICE CAPACITY EXCEEDED. %d/%d\n", cap(results), initialCapacity)
		slog.Warn(warning)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

type Transformer func(*bufio.Scanner, *bufio.Writer) error

func transform(infileName string, transformer Transformer, outfileName string) error {
	// get map of dockerfiles
	infile, err := os.Open(infileName)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	scanner := bufio.NewScanner(infile)
	// create outfile
	outfile, err := os.Create(outfileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	writer := bufio.NewWriter(outfile)
	// scan, transform, write
	for scanner.Scan() {
		transformer(scanner, writer)
		writer.Flush()
	}
	// cleanup
	err = infile.Close()
	if err != nil {
		return fmt.Errorf("failed to close infile %w", err)
	}
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush to outfile %w", err)
	}
	err = outfile.Close()
	if err != nil {
		return fmt.Errorf("failed to close outfile %w", err)
	}
	// log and return
	return nil
}
