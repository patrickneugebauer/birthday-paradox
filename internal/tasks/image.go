package tasks

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// GetImageUpdatedTimes fetches LastTagTime for all tags in parallel.
func GetImageUpdatedTimes(tags []string) map[string]int64 {
	times := make(map[string]int64, len(tags))
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, tag := range tags {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			if ts := GetImageUpdatedTime(t); ts != 0 {
				mu.Lock()
				times[t] = ts
				mu.Unlock()
				fmt.Print("*")
			}
		}(tag)
	}
	wg.Wait()
	return times
}

// GetImageUpdatedTime retrieves the last tag time of a docker image.
func GetImageUpdatedTime(tag string) int64 {
	cmd := exec.Command("docker", "image", "inspect", "--format={{.Metadata.LastTagTime}}", tag)
	output, err := cmd.Output()
	if err != nil {
		return 0
	}
	createdStr := strings.TrimSpace(string(output))
	if createdStr == "" {
		return 0
	}
	// Format: "2026-04-30 01:48:34.039676004 -0400 EDT"
	t, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", createdStr)
	if err != nil {
		return 0
	}
	return t.Unix()
}
