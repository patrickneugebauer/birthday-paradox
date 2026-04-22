package tasks

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func MakeFileTree() error {
	inDirName := solutionsDir
	outfileName := dockerfileList
	outfile, err := os.Create(outfileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	writer := bufio.NewWriter(outfile)
	// get map of dockerfiles
	entries, err := os.ReadDir(inDirName)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dirPath := filepath.Join(inDirName, entry.Name())
		files, _ := os.ReadDir(dirPath)
		for _, f := range files {
			isDockerfile := !f.IsDir() && strings.HasPrefix(f.Name(), "Dockerfile")
			if !isDockerfile {
				continue
			}
			encoder := json.NewEncoder(writer)
			dockerfile := Dockerfile{Language: entry.Name(), Filename: f.Name()}
			err := encoder.Encode(dockerfile)
			if err != nil {
				return fmt.Errorf("failed to encode: %w", err)
			}
			writer.Flush()
		}
	}
	// log and return
	fmt.Printf("wrote to file: %s\n", dockerfileList)
	return nil
}

