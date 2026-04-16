package tasks

import (
	"fmt"
	"os"
	"os/exec"
)

func Weigh() error {
	scriptName := weighScript
	if _, err := os.Stat(scriptName); os.IsNotExist(err) {
		return fmt.Errorf("script not found: please run 'pre-' command first to generate %s", scriptName)
	}
	fmt.Printf("Starting process using %s...\n", scriptName)
	cmd := exec.Command("/bin/bash", scriptName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}
	fmt.Println("complete.")
	return nil
}
