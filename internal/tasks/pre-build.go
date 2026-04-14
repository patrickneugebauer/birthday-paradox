package tasks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.yaml.in/yaml/v4"
)

const solutionsDir = "./solutions"
const dockerfilesList = "dockerfiles.yaml"
const buildScript = "build.sh"
const runScript = "run.sh"

// BuildConfig holds everything we need to know about a specific image
type BuildConfig struct {
	DirName    string
	Dockerfile string
	Tag        string
	Context    string
}

func PreBuild() error {
	// get map of dockerfiles
	configs, dfMap, err := discoverBuilds(solutionsDir)
	if err != nil {
		return err
	}
	// prepare commands
	var buildLines []string
	for _, c := range configs {
		line := fmt.Sprintf("docker build -f %s/%s %s -t %s", c.Context, c.Dockerfile, c.Context, c.Tag)
		buildLines = append(buildLines, line)
	}
	// persist
	if err := saveScript(buildScript, buildLines, true); err != nil {
		return fmt.Errorf("failed to save build script: %w", err)
	}
	yData, err := yaml.Marshal(dfMap)
	if err != nil {
		return fmt.Errorf("failed to marshal yaml: %w", err)
	}
	if err := os.WriteFile(dockerfilesList, yData, 0644); err != nil {
		return fmt.Errorf("failed to write yaml manifest: %w", err)
	}
	// log
	fmt.Printf("Pre-build complete. Generated %d build targets.\n", len(configs))
	return nil
}

func discoverBuilds(root string) ([]BuildConfig, map[string][]string, error) {
	configs := []BuildConfig{}
	dfMap := make(map[string][]string)

	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dirPath := filepath.Join(root, entry.Name())
		files, _ := os.ReadDir(dirPath)

		for _, f := range files {
			if !f.IsDir() && strings.HasPrefix(f.Name(), "Dockerfile") {
				tag := generateTag(entry.Name(), f.Name())
				configs = append(configs, BuildConfig{
					DirName:    entry.Name(),
					Dockerfile: f.Name(),
					Tag:        tag,
					Context:    dirPath,
				})
				dfMap[entry.Name()] = append(dfMap[entry.Name()], f.Name())
			}
		}
	}
	return configs, dfMap, nil
}

func generateTag(dir, filename string) string {
	tag := "bday/" + dir
	if strings.Contains(filename, ".") {
		tag += "-" + strings.SplitN(filename, ".", 2)[1]
	}
	return tag
}

// saveScript wraps content in a shebang and sets executable permissions
func saveScript(name string, lines []string, failFast bool) error {
	content := "#!/bin/bash\n"
	if failFast {
		content += "set -e\n" // Stop immediately if any command fails
	}
	content += strings.Join(lines, "\n") + "\n"

	// 0755 provides read/execute permissions for everyone, and write for the owner
	return os.WriteFile(name, []byte(content), 0755)
}
