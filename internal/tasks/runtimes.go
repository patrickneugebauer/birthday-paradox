package tasks

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/patrickneugebauer/birthday-paradox/internal/database"
	"gorm.io/gorm"
)

// FLOW
// startSession
// fetchLanguages -> getDockerfiles -> buildRuntimes -> insertRuntimes (return inserted, skipped, err)
// endSession

func GetRuntimes(db *database.DB) error {
	// get data
	langs, finalErr := fetchLanguages(db)
	if finalErr != nil {
		return finalErr
	}
	paths := make([]string, len(langs))
	for i, l := range langs {
		paths[i] = l.Directory
	}
	dockerfileMap, finalErr := getDockerfiles(paths)
	if finalErr != nil {
		return finalErr
	}
	runtimes, finalErr := buildRuntimes(dockerfileMap, langs)
	if finalErr != nil {
		return finalErr
	}
	inserted, skipped, finalErr := insertRuntimes(db, runtimes)
	if finalErr != nil {
		return finalErr
	}
	// log and return
	fmt.Printf("Skipped runtimes: %d\n", len(skipped))
	insertedNames := make([]string, 0, len(inserted))
	for _, v := range inserted {
		insertedNames = append(insertedNames, v.DockerfileName)
	}
	fmt.Printf("Inserted runtimes: %d %v\n", len(insertedNames), insertedNames)
	return nil
}

func fetchLanguages(db *database.DB) ([]database.Language, error) {
	var languages []database.Language
	err := db.DB.Find(&languages).Error
	return languages, err
}

func getDockerfiles(paths []string) (map[string][]string, error) {
	// use a map to hold dockerfile names by directory key
	dockerfiles := make(map[string][]string, len(paths))
	for _, dir := range paths {
		// add mepty slice for each directory
		dockerfiles[dir] = make([]string, 0)
		entries, err := os.ReadDir(filepath.Join(solutionsDir, dir))
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			fname := entry.Name()
			if !entry.IsDir() && strings.HasPrefix(fname, "Dockerfile") {
				// fill in the slice with dockerfiles
				dockerfiles[dir] = append(dockerfiles[dir], fname)
			}
		}
	}
	return dockerfiles, nil
}

func buildRuntimes(dockerfileMap map[string][]string, languages []database.Language) ([]database.Runtime, error) {
	// id by name lookup
	langLookup := make(map[string]uint)
	for _, l := range languages {
		langLookup[l.Directory] = l.ID
	}
	// get total count for slice capacity
	total := 0
	for _, files := range dockerfileMap {
		total += len(files)
	}
	runtimes := make([]database.Runtime, 0, total)
	// build runtimes
	for dirName, files := range dockerfileMap {
		langID, exists := langLookup[dirName]
		if !exists {
			message := fmt.Sprintf(
				"CRITICAL MAPPING ERROR: directory '%s' exists on disk but has no corresponding entry in the languages table",
				dirName,
			)
			return nil, errors.New(message)
		}
		for _, fname := range files {
			// Extract solution name (e.g., "Dockerfile.web" -> "web", "Dockerfile" -> "default")
			runtimeName := "default"
			if strings.Contains(fname, ".") {
				runtimeName = strings.SplitN(fname, ".", 2)[1]
			}
			runtimes = append(runtimes, database.Runtime{
				LanguageID:     langID,
				RuntimeName:    &runtimeName,
				DockerfileName: fname,
			})
		}
	}
	return runtimes, nil
}

func insertRuntimes(db *database.DB, runtimes []database.Runtime) ([]database.Runtime, []database.Runtime, error) {
	var inserted []database.Runtime
	var skipped []database.Runtime
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		for _, rt := range runtimes {
			// Create a copy to avoid pointer issues in the loop
			current := rt
			// Check for existing record based on the unique pair
			result := tx.Where(database.Runtime{
				LanguageID:  current.LanguageID,
				RuntimeName: current.RuntimeName,
			}).FirstOrCreate(&current)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected > 0 {
				inserted = append(inserted, current)
			} else {
				skipped = append(skipped, current)
			}
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return inserted, skipped, nil
}
