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

func ShouldRebuild(fileLastModified int64, imageLastCreated int64) bool {
	if imageLastCreated == 0 {
		return true
	}
	return fileLastModified > imageLastCreated
}

func Build() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	infile, err := os.Open(dockerfileList)
	if err != nil {
		return fmt.Errorf("open %s: %w", dockerfileList, err)
	}
	defer infile.Close()
	scanner := bufio.NewScanner(infile)

	writers, err := OpenBufferedFiles(buildTempResultsFile)
	if err != nil {
		return err
	}
	resultsFile := writers[0]
	defer CloseBufferedFiles(writers...)

	fmt.Print("build ")
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

		// Check will_rebuild flag (which was computed by map-files)
		shouldRebuild := true
		if dockerfile.WillRebuild != nil {
			shouldRebuild = *dockerfile.WillRebuild
		}
		// Or compute on the fly if needed
		if !shouldRebuild {
			shouldRebuild = ShouldRebuild(dockerfile.FileLastModified, dockerfile.ImageLastCreated)
		}

		tag := dockerfile.Tag
		dir := solutionsDir
		switch dockerfile.Directory {
		case "scaffolds":
			dir = scaffoldsDir
		case "hello-worlds":
			dir = helloWorldsDir
		}
		dockerfilePath := filepath.Join(dir, dockerfile.Language, dockerfile.Filename)
		solutionPath := filepath.Join(dir, dockerfile.Language)

		if !shouldRebuild {
			// Image is up-to-date, skip rebuild
			fmt.Print(".")
		} else {
			fmt.Print("*")
			// Run docker build
			buildCmd := exec.CommandContext(ctx, "docker", "build", "-t", tag, "-f", dockerfilePath, solutionPath)
			buildCmd.Stdout = os.Stdout
			buildCmd.Stderr = os.Stderr
			if err := buildCmd.Run(); err != nil {
				return fmt.Errorf("docker build failed for %s: %w", tag, err)
			}
		}

		// Write result for all (built or skipped)
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

	fmt.Printf("\nwrote: %s\n", buildResultsFile)
	return nil
}
