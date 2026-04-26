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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	infile, err := os.Open(buildResultsFile)
	if err != nil {
		return fmt.Errorf("open %s: %w", buildResultsFile, err)
	}
	defer infile.Close()
	scanner := bufio.NewScanner(infile)

	writers, err := OpenBufferedFiles(weighCommandsFile, weighResultsFile)
	if err != nil {
		return err
	}
	cmdFile, resultsFile := writers[0], writers[1]
	defer CloseBufferedFiles(writers...)

	for scanner.Scan() {
		fmt.Print(".")
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

		cmd := fmt.Sprintf("docker image inspect --format {{.Size}} %s\n", buildResult.Tag)
		if err := cmdFile.WriteString(cmd); err != nil {
			return fmt.Errorf("write command: %w", err)
		}

		output, err := exec.CommandContext(ctx, "docker", "image", "inspect", "--format", "{{.Size}}", buildResult.Tag).CombinedOutput()
		if err != nil {
			return fmt.Errorf("docker inspect failed: %w", err)
		}

		sizeMB, sizeBytes, err := parseSize(output)
		if err != nil {
			return fmt.Errorf("parse size: %w", err)
		}
		if err := resultsFile.Encode(WeighResult{Tag: buildResult.Tag, SizeMB: sizeMB, SizeBytes: sizeBytes}); err != nil {
			return fmt.Errorf("encode result: %w", err)
		}
	}
	fmt.Println()

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner: %w", err)
	}
	fmt.Printf("wrote to file: %s\n", weighResultsFile)
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
