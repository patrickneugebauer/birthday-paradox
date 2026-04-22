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
		if err := json.Unmarshal(bytes, &result); err == nil {
			stats[result.Tag] = result.IPS
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
	// Create a context that listens for Ctrl+C
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop() // Restore default behavior when Build returns
	infileName := buildResultsFile
	cmdFileName := runCommandsFile
	tempfileName := runTempResultsFile
	savefileName := runResultsFile
	prevStats, err := loadPreviousStats()
	if err != nil {
		return err
	}
	fmt.Printf("writing to: %s...\n", tempfileName)
	transformer := func(ctx context.Context, scanner *bufio.Scanner, cmdWriter *bufio.Writer, encoder *json.Encoder) error {
		// read input
		inBytes := scanner.Bytes()
		var buildResult BuildResult
		if err := json.Unmarshal(inBytes, &buildResult); err != nil {
			return fmt.Errorf("failed to unmarshall %w", err)
		}
		// write command
		tagName := buildResult.Tag
		iters := defaultIterations
		if ips, ok := prevStats[tagName]; ok && ips > 0 {
			iters = ips
		}
		commandText := fmt.Sprintf("docker run --rm %s %d\n", tagName, iters)
		if _, err := cmdWriter.WriteString(commandText); err != nil {
			return fmt.Errorf("failed to write command %w", err)
		}
		// run command
		fmt.Printf("Benchmarking %s %d...\n", tagName, iters)
		cmd := exec.CommandContext(ctx,
			"docker", "run",
			"--rm",
			tagName,
			strconv.Itoa(iters),
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Skipping %s: %v\n", tagName, err)
			emptyResult := RunResult{Tag: tagName, IPS: defaultIterations}
			if err := encoder.Encode(emptyResult); err != nil {
				return fmt.Errorf("Failed to serialize failure %s: %v\n", tagName, err)
			}
			return nil
		}
		// print results
		result := parseOutput(string(out))
		result.Tag = tagName
		if err := encoder.Encode(result); err != nil {
			return fmt.Errorf("Failed to encode %s: %v\n", tagName, err)
		}
		return nil
	}
	if err = transform(ctx, infileName, transformer, cmdFileName, tempfileName); err != nil {
		return fmt.Errorf("failed to transform %w", err)
	}
	if err := os.Rename(tempfileName, savefileName); err != nil {
		return fmt.Errorf("failed to finalize results: %w", err)
	}
	// log and return
	fmt.Printf("renamed: %s to %s\n", tempfileName, savefileName)
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
			res.Iterations, _ = strconv.Atoi(val)
		case "sample-size":
			res.SampleSize, _ = strconv.Atoi(val)
		case "percent":
			res.Percent, _ = strconv.ParseFloat(val, 64)
		case "seconds":
			res.Seconds, _ = strconv.ParseFloat(val, 64)
			if res.Seconds > 0 {
				res.IPS = int(float64(res.Iterations) / res.Seconds)
			}
		}
	}
	return res
}
