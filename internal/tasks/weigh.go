package tasks

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

func ShouldWeighImage(imageUpdatedAt int64, lastWeighAt int64) bool {
	if lastWeighAt == 0 {
		return true
	}
	return imageUpdatedAt > lastWeighAt
}

func loadExistingWeighResults() map[string]WeighResult {
	results := make(map[string]WeighResult)
	file, err := os.Open(weighResultsFile)
	if err != nil {
		return results
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var result WeighResult
		if err := json.Unmarshal(scanner.Bytes(), &result); err == nil {
			results[result.Tag] = result
		}
	}
	return results
}

func Weigh() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	infile, err := os.Open(buildResultsFile)
	if err != nil {
		return fmt.Errorf("open %s: %w", buildResultsFile, err)
	}
	defer infile.Close()
	scanner := bufio.NewScanner(infile)

	writers, err := OpenBufferedFiles(weighTempResultsFile)
	if err != nil {
		return err
	}
	resultsFile := writers[0]
	defer CloseBufferedFiles(writers...)

	existingResults := loadExistingWeighResults()
	dockerfiles := LoadDockerfileMap()
	now := time.Now().Unix()

	fmt.Print("weigh ")
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			fmt.Println("Interrupt received, stopping weigh loop...")
			return nil
		default:
		}

		var buildResult BuildResult
		if err := json.Unmarshal(scanner.Bytes(), &buildResult); err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}

		df, hasDockerfile := dockerfiles[buildResult.Tag]
		existing := existingResults[buildResult.Tag]

		// Check if we should skip based on image update time
		shouldWeigh := true
		if existing.LastWeighAt > 0 && hasDockerfile && df.ImageLastCreated <= existing.LastWeighAt {
			shouldWeigh = false
		}

		var result WeighResult
		if shouldWeigh {
			fmt.Print("*")
			output, err := exec.CommandContext(ctx, "docker", "image", "inspect", "--format", "{{.Size}}", buildResult.Tag).CombinedOutput()
			if err != nil {
				return fmt.Errorf("docker inspect failed: %w", err)
			}

			sizeMB, sizeBytes, err := parseSize(output)
			if err != nil {
				return fmt.Errorf("parse size: %w", err)
			}
			result = WeighResult{Tag: buildResult.Tag, SizeMB: sizeMB, SizeBytes: sizeBytes, LastWeighAt: now}
			if hasDockerfile {
				result.ImageUpdatedAt = df.ImageLastCreated
			}
		} else {
			fmt.Print(".")
			// Use cached result, update last_weigh_at
			result = existing
			result.LastWeighAt = now
		}

		// Write result for all (weighed or skipped)
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
	if err := os.Rename(weighTempResultsFile, weighResultsFile); err != nil {
		return fmt.Errorf("finalize results: %w", err)
	}
	fmt.Printf("wrote: %s\n", weighResultsFile)
	return nil
}

func parseSize(bytes []byte) (float64, int64, error) {
	sizeStr := strings.TrimSpace(string(bytes))
	sizeBytes, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parse size: %w", err)
	}
	sizeMB := float64(sizeBytes) / (1024 * 1024)
	roundedSize := roundToPrecision(sizeMB, 3)
	return roundedSize, sizeBytes, nil
}

func roundToPrecision(f float64, precision int) float64 {
	shift := math.Pow(10, float64(precision))
	return math.Round(f*shift) / shift
}
