package tasks

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func MakeFileTree() error {
	outfileName := dockerfileList
	outfile, err := os.Create(outfileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	writer := bufio.NewWriter(outfile)

	// Scan solutions, scaffolds, and hello-worlds directories
	for _, dirName := range []string{solutionsDir, scaffoldsDir, helloWorldsDir} {
		dirType := filepath.Base(dirName)
		entries, err := os.ReadDir(dirName)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		if os.IsNotExist(err) {
			continue // Skip if directory doesn't exist
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			dirPath := filepath.Join(dirName, entry.Name())
			files, _ := os.ReadDir(dirPath)
			for _, f := range files {
				isDockerfile := !f.IsDir() && strings.HasPrefix(f.Name(), "Dockerfile")
				if !isDockerfile {
					continue
				}
				encoder := json.NewEncoder(writer)
				runtime, dataStructure, executionMethod, tag := getRuntimeAndTag(entry.Name(), f.Name())
				dockerfile := Dockerfile{
					Language:        entry.Name(),
					Filename:        f.Name(),
					Runtime:         runtime,
					DataStructure:   dataStructure,
					ExecutionMethod: executionMethod,
					Tag:             tag,
					Directory:       dirType,
				}
				err := encoder.Encode(dockerfile)
				if err != nil {
					return fmt.Errorf("failed to encode: %w", err)
				}
				writer.Flush()
			}
		}
	}

	// log and return
	fmt.Printf("wrote to file: %s\n", dockerfileList)
	return nil
}

func toLower(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsUpper(r) {
			return unicode.ToLower(r)
		}
		return r
	}, s)
}

func getRuntimeAndTag(dir string, filename string) (*string, *string, *string, string) {
	var runtime, dataStructure, executionMethod *string
	tag := "bday/" + dir

	parts := strings.Split(filename, ".")
	if len(parts) == 1 {
		return runtime, dataStructure, executionMethod, tag
	}

	if len(parts) >= 2 {
		runtime = &parts[1]
		tag += ":" + toLower(*runtime)
	}
	if len(parts) >= 3 {
		dataStructure = &parts[2]
		if *dataStructure != "-" {
			tag += "." + toLower(*dataStructure)
		}
	}
	if len(parts) >= 4 {
		executionMethod = &parts[3]
		tag += "." + toLower(*executionMethod)
	}

	return runtime, dataStructure, executionMethod, tag
}
