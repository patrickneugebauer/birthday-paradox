package tasks

import (
	"bufio"
	"context"
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

type Transform func(context.Context, *bufio.Scanner, *bufio.Writer, *json.Encoder) error

func transform(ctx context.Context, infileName string, transformer Transform, cmdFileName string, resultsFileName string) error {
	// get map of dockerfiles
	infile, err := os.Open(infileName)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer infile.Close()
	scanner := bufio.NewScanner(infile)
	// create cmdFile
	cmdFile, err := os.Create(cmdFileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer cmdFile.Close()
	cmdWriter := bufio.NewWriter(cmdFile)
	defer cmdWriter.Flush()
	// create cmdFile
	resultsFile, err := os.Create(resultsFileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer resultsFile.Close()
	resultsWriter := bufio.NewWriter(resultsFile)
	defer resultsWriter.Flush()
	resultsEncoder := json.NewEncoder(resultsWriter)
	// scan, transform, write
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			fmt.Println("Interrupt received, stopping build loop...")
			return nil // Return nil so flushes and closes below still run
		default:
		}

		if err := transformer(ctx, scanner, cmdWriter, resultsEncoder); err != nil {
			return err
		}
	}
	// cleanup
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}
	if err := cmdWriter.Flush(); err != nil {
		return fmt.Errorf("failed to flush to cmdFile: %w", err)
	}
	if err := resultsWriter.Flush(); err != nil {
		return fmt.Errorf("failed to flush to cmdFile: %w", err)
	}
	if err := cmdFile.Close(); err != nil {
		return fmt.Errorf("failed to close outfile: %w", err)
	}
	if err := resultsFile.Close(); err != nil {
		return fmt.Errorf("failed to close outfile: %w", err)
	}
	return nil
}
