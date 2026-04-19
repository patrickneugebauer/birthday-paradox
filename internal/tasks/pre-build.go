package tasks

import (
	"bufio"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

func PreBuild() error {
	infileName := dockerfileList
	outfileName := buildScript
	var transformer Transformer = func(scanner *bufio.Scanner, writer *bufio.Writer) error {
		// read
		bytes := scanner.Bytes()
		var df Dockerfile
		err := json.Unmarshal(bytes, &df)
		if err != nil {
			return fmt.Errorf("failed to unmarshall %w", err)
		}
		// transform
		path := filepath.Join(solutionsDir, df.Language)
		filename := df.Filename
		tag := generateTag(df.Language, df.Filename)
		line := fmt.Sprintf("docker build -f %s/%s %s -t %s", path, filename, path, tag)
		// write
		_, err = writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("failed to write %w", err)
		}
		return nil
	}
	err := transform(infileName, transformer, outfileName)
	if err != nil {
		return fmt.Errorf("failed to tansform %w", err)
	}
	// log and return
	fmt.Printf("wrote to file: %s\n", outfileName)
	return nil
}

func generateTag(dir, filename string) string {
	tag := "bday/" + dir
	if strings.Contains(filename, ".") {
		extension := strings.SplitN(filename, ".", 2)[1]
		tag += "-" + extension
	}
	return tag
}
