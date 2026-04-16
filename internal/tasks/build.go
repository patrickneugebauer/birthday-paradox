package tasks

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Build() error {
	scriptName := buildScript
	// 1. Check if the script exists
	if _, err := os.Stat(scriptName); os.IsNotExist(err) {
		return fmt.Errorf("build script not found: please run 'pre-build' first to generate %s", scriptName)
	}

	content, err := os.ReadFile(scriptName)
	if err != nil {
		return fmt.Errorf("failed to read build script: %w", err)
	}

	// Open JSONL file for writing build artifacts
	output, err := os.Create(buildArtifacts)
	if err != nil {
		return fmt.Errorf("failed to create build artifacts file: %w", err)
	}
	defer output.Close()

	fmt.Printf("Starting build process using %s...\n", scriptName)

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Extract tag from the build command (everything after -t)
		parts := strings.Split(line, " -t ")
		var tag string
		if len(parts) >= 2 {
			tag = strings.Fields(parts[1])[0]
		}

		fmt.Printf("Building %s...\n", tag)

		cmd := exec.Command("/bin/bash", "-c", line)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("Build failed for %s: %v\n", tag, err)
			continue
		}

		// Write build artifact to JSONL
		artifact := BuildArtifact{Tag: tag}
		data, _ := json.Marshal(artifact)
		output.Write(data)
		output.WriteString("\n")
	}

	fmt.Println("Build complete. Build artifacts saved to", buildArtifacts)
	return nil
}
