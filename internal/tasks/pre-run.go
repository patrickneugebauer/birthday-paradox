package tasks

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func PreRun() error {
	configs, _, err := discoverBuilds(solutionsDir)
	if err != nil {
		return err
	}

	prevStats := loadPreviousStats(resultsFile)

	var runLines []string
	for _, c := range configs {
		iters := defaultIterations
		if ips, ok := prevStats[c.Tag]; ok && ips > 0 {
			iters = ips
		}
		runLines = append(runLines, fmt.Sprintf("docker run --rm %s %d", c.Tag, iters))
	}

	if err := saveScript(runScript, runLines, false); err != nil {
		return fmt.Errorf("failed to save run script: %w", err)
	}

	fmt.Printf("Pre-run complete. Generated run.sh with %d targets.\n", len(configs))
	return nil
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
