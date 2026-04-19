package tasks

import (
	"bufio"
	"encoding/json"
	"fmt"
)

func PreWeigh() error {
	infileName := buildArtifacts
	outfileName := weighScript
	var transformer Transformer = func(scanner *bufio.Scanner, writer *bufio.Writer) error {
		// read
		bytes := scanner.Bytes()
		var tag Tag
		err := json.Unmarshal(bytes, &tag)
		if err != nil {
			return fmt.Errorf("failed to unmarshall %w", err)
		}
		// transform
		tagName := tag.Name
		line := fmt.Sprintf("docker image inspect %s --format {{.Size}}", tagName)
		// write
		_, err = writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("failed to write %w", err)
		}
		return nil
	}
	err := transform(infileName, transformer, outfileName)
	if err != nil {
		return fmt.Errorf("failed to transform %w", err)
	}
	// log and return
	fmt.Printf("wrote to file: %s\n", outfileName)
	return nil
}
