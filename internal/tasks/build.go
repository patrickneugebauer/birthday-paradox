package tasks

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
)

func getTagName(dir string, filename string) string {
	tag := "bday/" + dir
	if strings.Contains(filename, ".") {
		extension := strings.SplitN(filename, ".", 2)[1]
		tag += "-" + extension
	}
	return tag
}

func Build() error {
	infileName := dockerfileList
	cmdfileName := buildCommandsFile
	tempfileName := buildTempResultsFile
	resultsfileName := buildResultsFile
	// Create a context that listens for Ctrl+C
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop() // Restore default behavior when Build returns
	transformer := func(ctx context.Context, scanner *bufio.Scanner, cmdWriter *bufio.Writer, encoder *json.Encoder) error {
		// read
		inBytes := scanner.Bytes()
		var dockerfile Dockerfile
		err := json.Unmarshal(inBytes, &dockerfile)
		if err != nil {
			return fmt.Errorf("failed to unmarshall %w", err)
		}
		// write command
		path := filepath.Join(solutionsDir, dockerfile.Language)
		filename := dockerfile.Filename
		tagName := getTagName(dockerfile.Language, dockerfile.Filename)
		command := fmt.Sprintf("docker build -f %s/%s %s -t %s\n", path, filename, path, tagName)
		if _, err := cmdWriter.WriteString(command); err != nil {
			return fmt.Errorf("failed to write command %w", err)
		}
		// create result
		cmd := exec.CommandContext(ctx,
			"docker", "build",
			"-f", filepath.Join(path, filename),
			path,
			"-t", tagName,
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed building: %v\n", err)
		}
		// write result
		tag := BuildResult{Tag: tagName}
		err = encoder.Encode(tag)
		if err != nil {
			return fmt.Errorf("failed to marshal json %w", err)
		}
		return nil
	}
	if err := transform(ctx, infileName, transformer, cmdfileName, tempfileName); err != nil {
		return err
	}
	if err := os.Rename(tempfileName, resultsfileName); err != nil {
		return fmt.Errorf("failed to finalize results: %w", err)
	}
	// log and return
	fmt.Println("complete")
	return nil
}
