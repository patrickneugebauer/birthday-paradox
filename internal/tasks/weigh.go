package tasks

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Weigh() error {
	if _, err := os.Stat(weighScript); os.IsNotExist(err) {
		return fmt.Errorf("script not found: please run 'pre-weigh' command first to generate %s", weighScript)
	}

	content, err := os.ReadFile(weighScript)
	if err != nil {
		return fmt.Errorf("failed to read weigh script: %w", err)
	}

	// Open JSONL file for writing
	output, err := os.Create(sizeFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
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

		tag := fields[3] // tag is after "docker image inspect"
		fmt.Printf("Checking size of %s...\n", tag)

		cmd := exec.Command(fields[0], fields[1:]...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Skipping %s: %v\n", tag, err)
			continue
		}

		// Parse output as integer bytes
		sizeStr := strings.TrimSpace(string(out))
		sizeBytes, err := strconv.ParseInt(sizeStr, 10, 64)
		if err != nil {
			fmt.Printf("Failed to parse size for %s: %v\n", tag, err)
			continue
		}

		// Convert bytes to MB
		sizeMB := float64(sizeBytes) / (1024 * 1024)
		var roundedSize float64
		if sizeMB < 100 {
			roundedSize = roundToPrecision(sizeMB, 1)
		} else {
			roundedSize = math.Round(sizeMB)
		}

		// Write to JSONL file immediately
		imgInfo := ImageInfo{
			Repository: tag,
			SizeMB:     roundedSize,
		}
		data, _ := json.Marshal(imgInfo)
		output.Write(data)
		output.WriteString("\n")
	}

	fmt.Printf("Weigh complete. Image sizes saved to %s\n", sizeFile)
	return nil
}

func roundToPrecision(f float64, precision int) float64 {
	shift := math.Pow(10, float64(precision))
	return math.Round(f*shift) / shift
}
