package tasks

import (
	"fmt"
	"os"
	"os/exec"
)

func Build() error {
	scriptName := buildScript
	// 1. Check if the script exists
	if _, err := os.Stat(scriptName); os.IsNotExist(err) {
		return fmt.Errorf("build script not found: please run 'pre-build' first to generate %s", scriptName)
	}

	fmt.Printf("Starting build process using %s...\n", scriptName)

	// 2. Execute the shell script
	// We use /bin/bash to ensure the script runs correctly regardless of file permissions
	cmd := exec.Command("/bin/bash", scriptName)

	// Stream output to the terminal so the user can see build progress
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	fmt.Println("Build complete.")
	return nil
}
