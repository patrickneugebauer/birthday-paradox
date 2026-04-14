package tasks

import (
	"encoding/json"
	"fmt"
	"os"
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
	var results []RunResult
	if err := json.Unmarshal(data, &results); err == nil {
		for _, r := range results {
			stats[r.Image] = r.IPS
		}
	}
	return stats
}
