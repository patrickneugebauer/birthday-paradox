package tasks

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Run() error {
	// 1. Archive the previous results
	if _, err := os.Stat(resultsFile); err == nil {
		os.Rename(resultsFile, previousResultsFile)
	}

	content, err := os.ReadFile(runScript)
	if err != nil {
		return fmt.Errorf("run script missing: %w", err)
	}

	// Open JSONL file for writing results
	output, err := os.Create(tempResultsFile)
	if err != nil {
		return fmt.Errorf("failed to create temp results file: %w", err)
	}
	defer output.Close()

	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		imageTag := fields[3]

		fmt.Printf("Benchmarking %s...\n", imageTag)

		cmd := exec.Command(fields[0], fields[1:]...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Skipping %s: %v\n", imageTag, err)
			continue
		}

		// 2. Parse and write
		res := parseOutput(string(out))
		res.Image = imageTag

		// 3. Write result to JSONL file immediately
		data, _ := json.Marshal(res)
		output.Write(data)
		output.WriteString("\n")
	}

	// 4. Rename temp file to final results file
	if err := os.Rename(tempResultsFile, resultsFile); err != nil {
		return fmt.Errorf("failed to finalize results: %w", err)
	}

	fmt.Println("Benchmarks complete. Results saved to", resultsFile)
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
