package tasks

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

func loadPreviousStats() (map[string]int, error) {
	infileName := runResultsFile
	stats := make(map[string]int)
	_, err := os.Stat(infileName)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("no previous results: %s\n", infileName)
		return stats, nil
	}
	infile, err := os.Open(infileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	scanner := bufio.NewScanner(infile)
	// scan
	for scanner.Scan() {
		bytes := scanner.Bytes()
		var result RunResult
		if err := json.Unmarshal(bytes, &result); err == nil && result.IPS != nil && *result.IPS > 0 {
			stats[result.Tag] = *result.IPS
		}
	}
	// cleanup
	err = infile.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close infile %w", err)
	}
	return stats, nil
}

func Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	prevStats, err := loadPreviousStats()
	if err != nil {
		return err
	}

	infile, err := os.Open(buildResultsFile)
	if err != nil {
		return fmt.Errorf("open %s: %w", buildResultsFile, err)
	}
	defer infile.Close()
	scanner := bufio.NewScanner(infile)

	writers, err := OpenBufferedFiles(runCommandsFile, runTempResultsFile)
	if err != nil {
		return err
	}
	cmdFile, resultsFile := writers[0], writers[1]
	defer CloseBufferedFiles(writers...)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for scanner.Scan() {
		fmt.Print(".")
		select {
		case <-ctx.Done():
			fmt.Println("Interrupt received, stopping run loop...")
			return nil
		case <-ticker.C:
			for _, w := range writers {
				if err := w.Flush(); err != nil {
					return fmt.Errorf("flush: %w", err)
				}
			}
		default:
		}

		var buildResult BuildResult
		if err := json.Unmarshal(scanner.Bytes(), &buildResult); err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}

		iterations := defaultIterations
		if prev, ok := prevStats[buildResult.Tag]; ok && prev > 0 {
			iterations = prev
		}

		cmd := fmt.Sprintf("docker run --rm %s %d\n", buildResult.Tag, iterations)
		if err := cmdFile.WriteString(cmd); err != nil {
			return fmt.Errorf("write command: %w", err)
		}

		output, err := exec.CommandContext(ctx, "docker", "run", "--rm", buildResult.Tag, strconv.Itoa(iterations)).CombinedOutput()
		if err != nil {
			// fallback on error
			if err := resultsFile.Encode(RunResult{Tag: buildResult.Tag}); err != nil {
				return fmt.Errorf("encode fallback: %w", err)
			}
			continue
		}

		result := parseOutput(string(output))
		result.Tag = buildResult.Tag
		if err := resultsFile.Encode(result); err != nil {
			return fmt.Errorf("encode result: %w", err)
		}
	}
	fmt.Println()

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner: %w", err)
	}
	if err := resultsFile.Close(); err != nil {
		return fmt.Errorf("close temp results: %w", err)
	}
	if err := os.Rename(runTempResultsFile, runResultsFile); err != nil {
		return fmt.Errorf("finalize results: %w", err)
	}
	fmt.Printf("wrote to file: %s\n", runResultsFile)
	return nil
}

func parseOutput(raw string) RunResult {
	res := RunResult{}
	lines := strings.Split(raw, "\n")
	for _, l := range lines {
		parts := strings.SplitN(strings.TrimSpace(l), ": ", 2)
		if len(parts) < 2 {
			continue
		}
		key, val := parts[0], parts[1]
		switch key {
		case "iterations":
			if v, err := strconv.Atoi(val); err == nil {
				res.Iterations = &v
			}
		case "sample-size":
			if v, err := strconv.Atoi(val); err == nil {
				res.SampleSize = &v
			}
		case "percent":
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				res.Percent = &v
			}
		case "seconds":
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				res.Seconds = &v
				if res.Iterations != nil && v > 0 {
					ips := int(float64(*res.Iterations) / v)
					res.IPS = &ips
				}
			}
		}
	}
	return res
}
