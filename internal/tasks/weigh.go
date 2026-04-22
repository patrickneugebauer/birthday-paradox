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
)

func Weigh() error {
	// Create a context that listens for Ctrl+C
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop() // Restore default behavior when Build returns
	infileName := buildResultsFile
	cmdFileName := weighCommandsFile
	outfileName := weighResultsFile
	transformer := func(ctx context.Context, scanner *bufio.Scanner, cmdWriter *bufio.Writer, encoder *json.Encoder) error {
		// read data
		inBytes := scanner.Bytes()
		var buildResult BuildResult
		if err := json.Unmarshal(inBytes, &buildResult); err != nil {
			return fmt.Errorf("failed to unmarshall %w", err)
		}
		// write command
		tagName := buildResult.Tag
		commandText := fmt.Sprintf("docker image inspect %s --format {{.Size}}\n", tagName)
		if _, err := cmdWriter.WriteString(commandText); err != nil {
			return fmt.Errorf("failed to write command %w", err)
		}
		// run command
		cmd := exec.Command(
			"docker", "image", "inspect",
			tagName,
			"--format", "{{.Size}}",
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("Failed to run %s: %v\n", tagName, err)
		}
		// write results
		roundedSize, err := parseSize(out)
		if err != nil {
			return fmt.Errorf("Failed to parse size %v\n", err)
		}
		imgInfo := WeighResult{Tag: tagName, SizeMB: roundedSize}
		if err := encoder.Encode(imgInfo); err != nil {
			return fmt.Errorf("failed to marshal json %w", err)
		}
		return nil
	}
	if err := transform(ctx, infileName, transformer, cmdFileName, outfileName); err != nil {
		return fmt.Errorf("failed to tansform %w", err)
	}
	// log and return
	fmt.Printf("wrote to file: %s\n", outfileName)
	return nil
}

func parseSize(bytes []byte) (float64, error) {
	sizeStr := strings.TrimSpace(string(bytes))
	sizeBytes, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		fmt.Errorf("Failed to parse size %s: %v\n", err)
	}
	sizeMB := float64(sizeBytes) / (1024 * 1024)
	var roundedSize float64
	if sizeMB < 100 {
		roundedSize = roundToPrecision(sizeMB, 1)
	} else {
		roundedSize = math.Round(sizeMB)
	}
	return roundedSize, nil
}

func roundToPrecision(f float64, precision int) float64 {
	shift := math.Pow(10, float64(precision))
	return math.Round(f*shift) / shift
}
