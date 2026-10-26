package tasks

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// getImageUpdatedTimes fetches LastTagTime for all tags in parallel.
func getImageUpdatedTimes(tags []string) map[string]*time.Time {
	times := make(map[string]*time.Time, len(tags))
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, tag := range tags {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			if ts := GetImageUpdatedTime(t); ts != nil {
				mu.Lock()
				times[t] = ts
				mu.Unlock()
				fmt.Print("*")
			}
		}(tag)
	}
	wg.Wait()
	// fmt.Printf("%s \n", t.Local().Format("2006-01-02 15:04:05"))
	// { t.Unio(), t.Local().Format() }
	return times
}

// GetImageUpdatedTime retrieves the last tag time of a docker image.
func GetImageUpdatedTime(tag string) *time.Time {
	cmd := exec.Command("docker", "image", "inspect", "--format={{.Metadata.LastTagTime}}", tag)
	output, err := cmd.Output()
	if err != nil {
		return nil
	}
	createdStr := strings.TrimSpace(string(output))
	if createdStr == "" {
		return nil
	}
	// Format: "2026-04-30 01:48:34.039676004 -0400 EDT"
	t, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", createdStr)
	if err != nil {
		return nil
	}
	return &t
}
