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
	"time"
)

func Build() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	infile, err := os.Open(dockerfileList)
	if err != nil {
		return fmt.Errorf("open %s: %w", dockerfileList, err)
	}
	defer infile.Close()
	scanner := bufio.NewScanner(infile)

	writers, err := OpenBufferedFiles(buildCommandsFile, buildTempResultsFile)
	if err != nil {
		return err
	}
	cmdFile, resultsFile := writers[0], writers[1]
	defer CloseBufferedFiles(writers...)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			fmt.Println("Interrupt received, stopping build loop...")
			return nil
		case <-ticker.C:
			for _, w := range writers {
				if err := w.Flush(); err != nil {
					return fmt.Errorf("flush: %w", err)
				}
			}
		default:
		}

		var dockerfile Dockerfile
		if err := json.Unmarshal(scanner.Bytes(), &dockerfile); err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}

		tag := dockerfile.Tag
		dockerfilePath := filepath.Join(solutionsDir, dockerfile.Language, dockerfile.Filename)
		solutionPath := filepath.Join(solutionsDir, dockerfile.Language)
		cmd := fmt.Sprintf("docker build -t %s -f %s %s\n", tag, dockerfilePath, solutionPath)
		if err := cmdFile.WriteString(cmd); err != nil {
			return fmt.Errorf("write command: %w", err)
		}

		// Run docker build
		buildCmd := exec.CommandContext(ctx, "docker", "build", "-t", tag, "-f", dockerfilePath, solutionPath)
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr
		if err := buildCmd.Run(); err != nil {
			return fmt.Errorf("docker build failed for %s: %w", tag, err)
		}

		if err := resultsFile.Encode(BuildResult{Tag: tag}); err != nil {
			return fmt.Errorf("encode result: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner: %w", err)
	}
	if err := resultsFile.Close(); err != nil {
		return fmt.Errorf("close temp results: %w", err)
	}
	if err := os.Rename(buildTempResultsFile, buildResultsFile); err != nil {
		return fmt.Errorf("finalize results: %w", err)
	}
	fmt.Printf("wrote to file: %s\n", buildResultsFile)
	return nil
}
