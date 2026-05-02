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
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func getDockerDPID() (int, error) {
	// Try reading /var/run/docker.pid first
	data, err := os.ReadFile("/var/run/docker.pid")
	if err == nil {
		pidStr := strings.TrimSpace(string(data))
		if pid, err := strconv.Atoi(pidStr); err == nil && pid > 0 {
			return pid, nil
		}
	}
	// Fallback: use pidof
	out, err := exec.Command("pidof", "dockerd").Output()
	if err == nil {
		pidStr := strings.TrimSpace(string(out))
		if pid, err := strconv.Atoi(pidStr); err == nil && pid > 0 {
			return pid, nil
		}
	}
	return 0, fmt.Errorf("could not find dockerd PID")
}

func ShouldRebuild(fileLastModified int64, imageLastCreated int64) bool {
	if imageLastCreated == 0 {
		return true
	}
	return fileLastModified > imageLastCreated
}

func loadExistingBuildResults() map[string]BuildResult {
	results := make(map[string]BuildResult)
	file, err := os.Open(buildResultsFile)
	if err != nil {
		return results
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var result BuildResult
		if err := json.Unmarshal(scanner.Bytes(), &result); err == nil {
			results[result.Tag] = result
		}
	}
	return results
}

func BuildAll() error {
	return buildFn(true)
}

func Build() error {
	return buildFn(false)
}

func buildFn(buildAll bool) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	infile, err := os.Open(dockerfileList)
	if err != nil {
		return fmt.Errorf("open %s: %w", dockerfileList, err)
	}
	defer infile.Close()
	scanner := bufio.NewScanner(infile)

	writers, err := OpenBufferedFiles(buildTempResultsFile)
	if err != nil {
		return err
	}
	resultsFile := writers[0]
	defer CloseBufferedFiles(writers...)

	existingResults := loadExistingBuildResults()

	// Get dockerd PID for stats collection
	dockerdPID, err := getDockerDPID()
	if err != nil {
		return fmt.Errorf("get dockerd PID: %w", err)
	}

	fmt.Print("build ")
	ticker := time.NewTicker(refreshInterval * time.Second)
	defer ticker.Stop()

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			fmt.Println("Interrupt received, stopping build loop...")
			return nil
		case <-ticker.C:
			for _, w := range writers {
				if err := w.Flush(); err != nil {
					return fmt.Errorf("flush: %w", err)
				}
			}
		default:
		}

		var dockerfile Dockerfile
		if err := json.Unmarshal(scanner.Bytes(), &dockerfile); err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}

		// Check will_rebuild flag (which was computed by map-files)
		shouldRebuild := true
		if buildAll != true && dockerfile.WillRebuild != nil {
			shouldRebuild = *dockerfile.WillRebuild
		}
		// Or compute on the fly if needed
		if !shouldRebuild {
			shouldRebuild = ShouldRebuild(dockerfile.FileLastModified, dockerfile.ImageLastCreated)
		}

		tag := dockerfile.Tag
		dir := solutionsDir
		switch dockerfile.Directory {
		case "scaffolds":
			dir = scaffoldsDir
		case "hello-worlds":
			dir = helloWorldsDir
		}
		dockerfilePath := filepath.Join(dir, dockerfile.Language, dockerfile.Filename)
		solutionPath := filepath.Join(dir, dockerfile.Language)

		var buildResult BuildResult
		if !shouldRebuild {
			// Image is up-to-date, skip rebuild
			buildResult = existingResults[tag]
			fmt.Print(".")
		} else {
			buildResult = BuildResult{Tag: tag, LastBuiltAt: existingResults[tag].LastBuiltAt}
			fmt.Print("*")

			// Snapshot host stats before build
			netRxBefore, netTxBefore, err := snapshotNetDev()
			if err != nil {
				return fmt.Errorf("snapshot net before build for %s: %w", tag, err)
			}
			blkReadBefore, blkWriteBefore, err := snapshotDiskStats()
			if err != nil {
				return fmt.Errorf("snapshot disk before build for %s: %w", tag, err)
			}

			// Run docker build with stats collection
			buildStart := time.Now()
			buildCmd := exec.CommandContext(ctx, "docker", "build", "-t", tag, "-f", dockerfilePath, solutionPath)
			buildCmd.Stdout = os.Stdout
			buildCmd.Stderr = os.Stderr

			// Collect stats during build if dockerd PID available
			var statsCh chan buildStats
			var pollCancel context.CancelFunc
			if dockerdPID > 0 {
				var pollCtx context.Context
				pollCtx, pollCancel = context.WithCancel(ctx)
				statsCh = collectBuildStats(pollCtx, dockerdPID)
			}

			if err := buildCmd.Run(); err != nil {
				if pollCancel != nil {
					pollCancel()
				}
				if statsCh != nil {
					// Drain the channel
					<-statsCh
				}
				return fmt.Errorf("docker build failed for %s: %w", tag, err)
			}

			// Stop stats collection
			if pollCancel != nil {
				pollCancel()
			}

			// Snapshot host stats after build
			netRxAfter, netTxAfter, err := snapshotNetDev()
			if err != nil {
				return fmt.Errorf("snapshot net after build for %s: %w", tag, err)
			}
			blkReadAfter, blkWriteAfter, err := snapshotDiskStats()
			if err != nil {
				return fmt.Errorf("snapshot disk after build for %s: %w", tag, err)
			}

			// Collect build stats
			buildElapsed := time.Since(buildStart)
			totalS := math.Round(buildElapsed.Seconds()*1e6) / 1e6
			buildResult.TotalS = &totalS

			if statsCh != nil {
				bs := <-statsCh
				if bs.collectError != nil {
					return fmt.Errorf("collect CPU/RAM stats for %s: %w", tag, bs.collectError)
				}
				buildResult.PeakRAMBytes = &bs.peakRAM
				if bs.firstCPU > 0 || bs.cpuRamPollCount > 0 {
					cpuDelta := bs.lastCPU - bs.firstCPU
					if cpuDelta < 0 {
						cpuDelta = 0
					}
					cpuS := float64(cpuDelta) / 100.0
					buildResult.CpuS = &cpuS
				}
				buildResult.CpuRamPollCount = &bs.cpuRamPollCount
			}

			// Calculate deltas for disk/network
			netRx := netRxAfter - netRxBefore
			buildResult.NetRxBytes = &netRx
			netTx := netTxAfter - netTxBefore
			buildResult.NetTxBytes = &netTx
			blkRead := blkReadAfter - blkReadBefore
			buildResult.BlkReadBytes = &blkRead
			blkWrite := blkWriteAfter - blkWriteBefore
			buildResult.BlkWriteBytes = &blkWrite

			buildResult.LastBuiltAt = buildStart.Unix()
		}

		if err := resultsFile.Encode(buildResult); err != nil {
			return fmt.Errorf("encode result: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner: %w", err)
	}
	if err := resultsFile.Close(); err != nil {
		return fmt.Errorf("close temp results: %w", err)
	}
	if err := os.Rename(buildTempResultsFile, buildResultsFile); err != nil {
		return fmt.Errorf("finalize results: %w", err)
	}

	fmt.Printf("\nwrote: %s\n", buildResultsFile)
	return nil
}
