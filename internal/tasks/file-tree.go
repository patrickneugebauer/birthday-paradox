package tasks

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func LoadDockerfileMap() map[string]Dockerfile {
	dockerfiles := make(map[string]Dockerfile)
	file, err := os.Open(dockerfileListFile)
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
	// DEPENDENCIES
	solutionsDirName := solutionsDir
	getImageUpdatedAt := getImageUpdatedTimes
	type df = Dockerfile
	getFileUpdatedAt := getMaxFileModTime
	dfListFile := dockerfileListFile
	// ============

	fmt.Print("map\n")

	// Pass 1: collect all dockerfile entries from disk.
	type entry struct {
		language string
		filename string
		dirPath  string
		tag      string
	}
	var entries []entry
	dirs, err := os.ReadDir(solutionsDirName)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		dirPath := filepath.Join(solutionsDirName, dir.Name())
		files, _ := os.ReadDir(dirPath)
		for _, f := range files {
			if f.IsDir() || f.Name() != "Dockerfile" {
				continue
			}
			tag := "bday/" + dir.Name()
			entries = append(entries, entry{
				language: dir.Name(),
				filename: f.Name(),
				dirPath:  dirPath,
				tag:      tag,
			})
		}
	}

	// Pass 2: fetch all image times in parallel.
	tags := make([]string, len(entries))
	for i, e := range entries {
		tags[i] = e.tag
	}
	imageTimes := getImageUpdatedAt(tags)
	fmt.Println()

	// Pass 3: build output rows.
	var dockerfiles []df
	for _, e := range entries {
		fileLastMod := getFileUpdatedAt(e.dirPath)
		imageLastUpdated := imageTimes[e.tag]

		var imageWillRebuild bool
		if fileLastMod == nil || imageLastUpdated == nil || (*fileLastMod).After(*imageLastUpdated) {
			imageWillRebuild = true
		}

		var imageLastUpdatedAt *string
		if imageLastUpdated != nil {
			lastUpdatedString := (*imageLastUpdated).Local().Format(timeFormat)
			imageLastUpdatedAt = &lastUpdatedString
		}
		dockerfiles = append(dockerfiles, df{
			Language:         e.language,
			Tag:              e.tag,
			FileLastModified: fileLastMod.Local().Format(timeFormat),
			ImageLastCreated: *&imageLastUpdatedAt,
			WillRebuild:      &imageWillRebuild,
		})
	}

	// Pass 4: write all at once.
	var buf bytes.Buffer
	for _, df := range dockerfiles {
		jsonBytes, _ := json.Marshal(df)
		buf.Write(jsonBytes)
		buf.WriteByte('\n')
	}
	if err := os.WriteFile(dfListFile, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	fmt.Printf("wrote: %s\n", dfListFile)
	return nil
}

func getMaxFileModTime(dirPath string) *time.Time {
	var maxTime *time.Time
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			modTime := info.ModTime()
			if maxTime == nil || modTime.After(*maxTime) {
				maxTime = &modTime
			}
		}
		return nil
	})
	return maxTime
}
