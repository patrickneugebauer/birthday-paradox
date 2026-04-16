package tasks

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PostWeigh handles the file I/O orchestration
func PostWeigh() error {
	inFile := sizeFile
	outFile := sizeFileMB
	fmt.Printf("Reading %s -> Writing %s\n", inFile, outFile)
	input, err := os.Open(inFile)
	if err != nil {
		return fmt.Errorf("open input: %w", err)
	}
	defer input.Close()

	output, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer output.Close()

	return processLines(input, output)
}

func convertToMB(rawSize string) (float64, error) {
	clean := strings.ToUpper(strings.TrimSpace(rawSize))

	if strings.HasSuffix(clean, "GB") {
		numStr := strings.TrimSuffix(clean, "GB")
		val, err := strconv.ParseFloat(strings.TrimSpace(numStr), 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number format in GB: %w", err)
		}
		return val * 1024, nil
	}

	if strings.HasSuffix(clean, "MB") {
		numStr := strings.TrimSuffix(clean, "MB")
		val, err := strconv.ParseFloat(strings.TrimSpace(numStr), 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number format in MB: %w", err)
		}
		return val, nil
	}

	return 0, fmt.Errorf("unsupported size format (expected MB or GB): %s", rawSize)
}

func processLines(r *os.File, w *os.File) error {
	scanner := bufio.NewScanner(r)
	writer := bufio.NewWriter(w)

	for scanner.Scan() {
		var img ImageInfo
		if err := json.Unmarshal(scanner.Bytes(), &img); err != nil {
			continue
		}

		sizeMB, err := convertToMB(img.Size)
		if err != nil {
			return fmt.Errorf("error processing image %s: %w", img.Repository, err)
		}
		img.SizeMB = sizeMB

		encoded, _ := json.Marshal(img)
		writer.Write(encoded)
		writer.WriteString("\n")

		if err := writer.Flush(); err != nil {
			return fmt.Errorf("failed to flush line to file: %w", err)
		}
	}

	return scanner.Err()
}
