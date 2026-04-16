package tasks

import (
	"fmt"
	"os"
)

func PreWeigh() error {
	// build command
	filter := `--filter "reference=*bday/*"`
	format := `--format '{"repository":"{{.Repository}}", "size":"{{.Size}}"}'`
	sort := "| sort"
	output := fmt.Sprintf("> %s", sizeFile)
	command := fmt.Sprintf(
		"docker image ls \\\n  %s \\\n  %s \\\n  %s \\\n  %s",
		filter,
		format,
		sort,
		output,
	)
	// write script
	err := os.WriteFile(weighScript, []byte(command), 0755)
	if err != nil {
		return err
	}
	// print and return
	fmt.Printf("Generated weigh.sh.\n")
	return nil
}
