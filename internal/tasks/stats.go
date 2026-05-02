package tasks

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type collectedStats struct {
	collected       bool
	peakRAM         int64
	firstCPU        int64
	lastCPU         int64
	cpuRamPollCount int64
	collectError    error
}

type buildStats struct {
	collected       bool
	peakRAM         int64
	firstCPU        int64
	lastCPU         int64
	netRxBytes      int64
	netTxBytes      int64
	blkReadBytes    int64
	blkWriteBytes   int64
	cpuRamPollCount int64
	collectError    error
}

type dockerStatsResp struct {
	MemoryStats struct {
		Usage int64 `json:"usage"`
	} `json:"memory_stats"`
	CPUStats struct {
		CPUUsage struct {
			TotalUsage int64 `json:"total_usage"`
		} `json:"cpu_usage"`
	} `json:"cpu_stats"`
}

var dockerClient = &http.Client{
	Transport: &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "unix", "/var/run/docker.sock")
		},
	},
}

// collectContainerStats polls the Docker API for container stats
// until ctx is cancelled, then sends the aggregated result on the returned
// channel. Tight polling maximises coverage of short-lived containers.
func collectContainerStats(ctx context.Context, name string) chan collectedStats {
	ch := make(chan collectedStats, 1)
	go func() {
		var cs collectedStats
		cs.collected = true
		for {
			select {
			case <-ctx.Done():
				ch <- cs
				return
			default:
			}
			cs.cpuRamPollCount++
			req, err := http.NewRequestWithContext(ctx, "GET",
				fmt.Sprintf("http://localhost/containers/%s/stats?stream=false&one-shot=true", name), nil)
			if err != nil {
				continue
			}
			resp, err := dockerClient.Do(req)
			if err != nil {
				continue
			}
			if resp.StatusCode == 404 {
				resp.Body.Close()
				continue
			}
			if resp.StatusCode != 200 {
				resp.Body.Close()
				continue
			}
			var stats dockerStatsResp
			if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
				resp.Body.Close()
				continue
			}
			resp.Body.Close()

			if stats.MemoryStats.Usage >= 0 {
				if stats.MemoryStats.Usage > cs.peakRAM {
					cs.peakRAM = stats.MemoryStats.Usage
				}
				currentCPU := stats.CPUStats.CPUUsage.TotalUsage
				if cs.firstCPU == 0 || currentCPU < cs.firstCPU {
					cs.firstCPU = currentCPU
				}
				if currentCPU > cs.lastCPU {
					cs.lastCPU = currentCPU
				}
			}
		}
	}()
	return ch
}

// collectBuildStats polls dockerd process stats and host-level disk/network stats
// until ctx is cancelled, aggregating peak RAM and CPU usage.
func collectBuildStats(ctx context.Context, dockerdPID int) chan buildStats {
	ch := make(chan buildStats, 1)
	go func() {
		var bs buildStats
		bs.collected = true
		for {
			select {
			case <-ctx.Done():
				ch <- bs
				return
			default:
			}
			bs.cpuRamPollCount++
			// Poll dockerd process stats for RAM and CPU
			if ram, cpu, err := readProcStats(dockerdPID); err == nil {
				if ram > bs.peakRAM {
					bs.peakRAM = ram
				}
				if bs.firstCPU == 0 {
					bs.firstCPU = cpu
				}
				bs.lastCPU = cpu
			} else if bs.collectError == nil {
				bs.collectError = err
			}
		}
	}()
	return ch
}

// readProcStats reads RAM (VmRSS) and CPU ticks from /proc/{pid}/status and /proc/{pid}/stat
func readProcStats(pid int) (int64, int64, error) {
	// Read VmRSS from /proc/{pid}/status
	ram, err := readVmRSS(pid)
	if err != nil {
		return 0, 0, err
	}

	// Read CPU ticks from /proc/{pid}/stat
	cpu, err := readProcCPU(pid)
	if err != nil {
		return 0, 0, err
	}

	return ram, cpu, nil
}

// readVmRSS parses /proc/{pid}/status for VmRSS (resident set size in KB)
func readVmRSS(pid int) (int64, error) {
	file, err := os.Open(fmt.Sprintf("/proc/%d/status", pid))
	if err != nil {
		return 0, fmt.Errorf("open /proc/%d/status: %w", pid, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "VmRSS:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				val, err := strconv.ParseInt(parts[1], 10, 64)
				if err == nil {
					return val * 1024, nil // Convert KB to bytes
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("scan /proc/%d/status: %w", pid, err)
	}
	return 0, fmt.Errorf("VmRSS not found in /proc/%d/status", pid)
}

// readProcCPU reads CPU ticks from /proc/{pid}/stat (fields 14+15: utime+stime)
func readProcCPU(pid int) (int64, error) {
	file, err := os.Open(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return 0, fmt.Errorf("open /proc/%d/stat: %w", pid, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		// Find the end of comm (which is in parentheses)
		closeParen := strings.LastIndex(line, ")")
		if closeParen == -1 {
			return 0, fmt.Errorf("malformed /proc/%d/stat: no closing paren", pid)
		}
		// Split remaining fields
		fields := strings.Fields(line[closeParen+1:])
		// Fields are 0-indexed from after the comm; we want fields 13 and 14 (utime and stime)
		if len(fields) > 14 {
			utime, _ := strconv.ParseInt(fields[13], 10, 64)
			stime, _ := strconv.ParseInt(fields[14], 10, 64)
			return utime + stime, nil
		}
		return 0, fmt.Errorf("not enough fields in /proc/%d/stat", pid)
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("scan /proc/%d/stat: %w", pid, err)
	}
	return 0, fmt.Errorf("no data in /proc/%d/stat", pid)
}

// snapshotNetDev reads /proc/net/dev and returns rx_bytes and tx_bytes summed across all interfaces (except lo)
func snapshotNetDev() (int64, int64, error) {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return 0, 0, fmt.Errorf("open /proc/net/dev: %w", err)
	}
	defer file.Close()

	var rxBytes, txBytes int64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip header lines
		if strings.Contains(line, "|") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 10 {
			continue
		}
		// Interface name is the first part (without the colon)
		iface := strings.TrimSuffix(parts[0], ":")
		if iface == "lo" {
			continue
		}
		// rx_bytes is parts[1], tx_bytes is parts[9]
		rx, _ := strconv.ParseInt(parts[1], 10, 64)
		tx, _ := strconv.ParseInt(parts[9], 10, 64)
		rxBytes += rx
		txBytes += tx
	}
	if err := scanner.Err(); err != nil {
		return 0, 0, fmt.Errorf("scan /proc/net/dev: %w", err)
	}
	return rxBytes, txBytes, nil
}

// snapshotDiskStats reads /proc/diskstats and returns total sectors read and written
// across all non-loop block devices, converted to bytes (sector = 512 bytes)
func snapshotDiskStats() (int64, int64, error) {
	file, err := os.Open("/proc/diskstats")
	if err != nil {
		return 0, 0, fmt.Errorf("open /proc/diskstats: %w", err)
	}
	defer file.Close()

	var blkRead, blkWrite int64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 10 {
			continue
		}
		// Skip loop devices
		device := parts[2]
		if strings.HasPrefix(device, "loop") {
			continue
		}
		// sectors_read is parts[5], sectors_written is parts[9]
		sectorsRead, _ := strconv.ParseInt(parts[5], 10, 64)
		sectorsWrite, _ := strconv.ParseInt(parts[9], 10, 64)
		blkRead += sectorsRead * 512
		blkWrite += sectorsWrite * 512
	}
	if err := scanner.Err(); err != nil {
		return 0, 0, fmt.Errorf("scan /proc/diskstats: %w", err)
	}
	return blkRead, blkWrite, nil
}
