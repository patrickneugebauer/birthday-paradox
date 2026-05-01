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

func ShouldRunBenchmark(imageUpdatedAt int64, lastRunAt int64) bool {
	if lastRunAt == 0 {
		return true
	}
	return imageUpdatedAt > lastRunAt
}

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

func loadExistingRunResults() map[string]RunResult {
	results := make(map[string]RunResult)
	file, err := os.Open(runResultsFile)
	if err != nil {
		return results
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var result RunResult
		if err := json.Unmarshal(scanner.Bytes(), &result); err == nil {
			results[result.Tag] = result
		}
	}
	return results
}

func RunAll() error {
	return RunFn(true, nil)
}

func Run() error {
	return RunFn(false, nil)
}

func RunSome(lang string) error {
	return RunFn(false, &lang)
}

func RunFn(all bool, langs *string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	prevStats, err := loadPreviousStats()
	if err != nil {
		return err
	}

	existingResults := loadExistingRunResults()
	dockerfiles := LoadDockerfileMap()
	now := time.Now().Unix()

	infile, err := os.Open(buildResultsFile)
	if err != nil {
		return fmt.Errorf("open %s: %w", buildResultsFile, err)
	}
	defer infile.Close()
	scanner := bufio.NewScanner(infile)

	writers, err := OpenBufferedFiles(runTempResultsFile)
	if err != nil {
		return err
	}
	resultsFile := writers[0]
	defer CloseBufferedFiles(writers...)

	fmt.Print("run ")
	ticker := time.NewTicker(refreshInterval * time.Second)
	defer ticker.Stop()

	for scanner.Scan() {
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

		df, hasDockerfile := dockerfiles[buildResult.Tag]
		existing := existingResults[buildResult.Tag]

		var contains bool
		if langs != nil {
			contains = strings.Contains(buildResult.Tag, *langs)
		}
		shouldRun := all || contains || buildResult.LastBuiltAt > existing.LastRunAt

		var result RunResult
		if shouldRun {
			fmt.Print("*")
			iterations := defaultIterations
			if prev, ok := prevStats[buildResult.Tag]; ok && prev > 0 {
				iterations = prev
			}

			output, err := exec.CommandContext(ctx, "docker", "run", "--rm", buildResult.Tag, strconv.Itoa(iterations)).CombinedOutput()
			if err != nil {
				// fallback on error, use existing result with updated timestamps
				result = existing
				result.Tag = buildResult.Tag
				result.LastRunAt = now
				if hasDockerfile {
					result.ImageUpdatedAt = df.ImageLastCreated
				}
			} else {
				result = parseOutput(string(output))
				if df.Directory == "solutions" && result.Iterations != nil && *result.Seconds > 0 {
					ips := int(float64(*result.Iterations) / *result.Seconds)
					result.IPS = &ips
				}
				result.Tag = buildResult.Tag
				result.LastRunAt = now
				if hasDockerfile {
					result.ImageUpdatedAt = df.ImageLastCreated
				}
			}
		} else {
			fmt.Print(".")
			// Use cached result, update last_run_at
			result = existing
			result.Tag = buildResult.Tag
			result.LastRunAt = now
			if hasDockerfile {
				result.ImageUpdatedAt = df.ImageLastCreated
			}
		}

		// Write result for all (run or skipped)
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
	fmt.Printf("wrote: %s\n", runResultsFile)
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
			}
		}
	}
	return res
}
