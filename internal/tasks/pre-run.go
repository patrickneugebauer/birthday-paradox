package tasks

import (
	"bufio"
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
)

func PreRun() error {
	infileName := buildArtifacts
	outfileName := runScript
	prevStats := loadPreviousStats(resultsFile)
	var transformer Transformer = func(scanner *bufio.Scanner, writer *bufio.Writer) error {
		// read
		bytes := scanner.Bytes()
		var tag Tag
		err := json.Unmarshal(bytes, &tag)
		if err != nil {
			return fmt.Errorf("failed to unmarshall %w", err)
		}
		// transform
		tagName := tag.Name
		iters := defaultIterations
		if ips, ok := prevStats[tagName]; ok && ips > 0 {
			iters = ips
		}
		line := fmt.Sprintf("docker run --rm %s %d", tagName, iters)
		// write
		_, err = writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("failed to write %w", err)
		}
		return nil
	}
	err := transform(infileName, transformer, outfileName)
	if err != nil {
		return fmt.Errorf("failed to transform %w", err)
	}
	// log and return
	fmt.Printf("wrote to file: %s\n", outfileName)
	return nil
}

type Tag struct {
	Name string `json:"tag"`
}

func ReadBuildList() ([]Tag, error) {
	fname := buildArtifacts
	list, err := readJsonl[Tag](fname)
	if err != nil {
		return nil, err
	}
	slices.SortFunc(list, func(a, b Tag) int {
		return cmp.Compare(a.Name, b.Name)
	})
	return list, nil
}

func loadPreviousStats(path string) map[string]int {
	stats := make(map[string]int)
	data, err := os.ReadFile(path)
	if err != nil {
		return stats
	}
	// Parse JSONL format (one result per line)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var r RunResult
		if err := json.Unmarshal([]byte(line), &r); err == nil {
			stats[r.Image] = r.IPS
		}
	}
	return stats
}
