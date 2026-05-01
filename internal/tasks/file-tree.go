package tasks

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func LoadDockerfileMap() map[string]Dockerfile {
	dockerfiles := make(map[string]Dockerfile)
	file, err := os.Open(dockerfileList)
	if err != nil {
		return dockerfiles
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var df Dockerfile
		if err := json.Unmarshal(scanner.Bytes(), &df); err == nil {
			dockerfiles[df.Tag] = df
		}
	}
	return dockerfiles
}

func MakeFileTree() error {
	fmt.Print("map ")

	// Pass 1: collect all dockerfile entries from disk.
	type entry struct {
		language        string
		filename        string
		dirPath         string
		dirType         string
		tag             string
		runtime         *string
		dataStructure   *string
		executionMethod *string
	}
	var entries []entry
	for _, dirName := range []string{solutionsDir, scaffoldsDir, helloWorldsDir} {
		dirType := filepath.Base(dirName)
		dirs, err := os.ReadDir(dirName)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		for _, dir := range dirs {
			if !dir.IsDir() {
				continue
			}
			dirPath := filepath.Join(dirName, dir.Name())
			files, _ := os.ReadDir(dirPath)
			for _, f := range files {
				if f.IsDir() || !strings.HasPrefix(f.Name(), "Dockerfile") {
					continue
				}
				runtime, dataStructure, executionMethod, tag := getRuntimeAndTag(dir.Name(), f.Name())
				entries = append(entries, entry{
					language:        dir.Name(),
					filename:        f.Name(),
					dirPath:         dirPath,
					dirType:         dirType,
					tag:             tag,
					runtime:         runtime,
					dataStructure:   dataStructure,
					executionMethod: executionMethod,
				})
			}
		}
	}

	// Pass 2: fetch all image times in parallel.
	tags := make([]string, len(entries))
	for i, e := range entries {
		tags[i] = e.tag
	}
	imageTimes := GetImageUpdatedTimes(tags)
	fmt.Println()

	// Pass 3: build output rows.
	var dockerfiles []Dockerfile
	for _, e := range entries {
		fileLastMod := getMaxFileModTime(e.dirPath)
		imageLastUpdated := imageTimes[e.tag]
		willRebuild := ShouldRebuild(fileLastMod, imageLastUpdated)
		dockerfiles = append(dockerfiles, Dockerfile{
			Language:         e.language,
			Filename:         e.filename,
			Runtime:          e.runtime,
			DataStructure:    e.dataStructure,
			ExecutionMethod:  e.executionMethod,
			Tag:              e.tag,
			Directory:        e.dirType,
			FileLastModified: fileLastMod,
			ImageLastCreated: imageLastUpdated,
			WillRebuild:      &willRebuild,
		})
	}

	// Pass 4: write all at once.
	var buf bytes.Buffer
	for _, df := range dockerfiles {
		jsonBytes, _ := json.Marshal(df)
		buf.Write(jsonBytes)
		buf.WriteByte('\n')
	}
	if err := os.WriteFile(dockerfileList, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	fmt.Printf("wrote: %s\n", dockerfileList)
	return nil
}

func toLower(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsUpper(r) {
			return unicode.ToLower(r)
		}
		return r
	}, s)
}

func getRuntimeAndTag(dir string, filename string) (*string, *string, *string, string) {
	var runtime, dataStructure, executionMethod *string
	tag := "bday/" + dir

	parts := strings.Split(filename, ".")
	if len(parts) == 1 {
		return runtime, dataStructure, executionMethod, tag
	}

	if len(parts) >= 2 {
		runtime = &parts[1]
		tag += ":" + toLower(*runtime)
	}
	if len(parts) >= 3 {
		dataStructure = &parts[2]
		if *dataStructure != "-" {
			tag += "." + toLower(*dataStructure)
		}
	}
	if len(parts) >= 4 {
		executionMethod = &parts[3]
		tag += "." + toLower(*executionMethod)
	}

	return runtime, dataStructure, executionMethod, tag
}

func getMaxFileModTime(dirPath string) int64 {
	var maxTime int64
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			modTime := info.ModTime().Unix()
			if modTime > maxTime {
				maxTime = modTime
			}
		}
		return nil
	})
	return maxTime
}
