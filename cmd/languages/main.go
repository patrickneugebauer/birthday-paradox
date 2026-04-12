package main

import (
	"fmt"
	"log"
	"os"

	"github.com/patrickneugebauer/birthday-paradox/internal/database"
	"gorm.io/gorm"
)

const sessionType = "languages"
const solutionsDir = "./solutions"

func main() {
	if err := run(); err != nil {
		log.Fatalf("Critical failure: %v", err)
	}
}

func run() error {
	db, err := database.StartSession(sessionType)
	if err != nil {
		return fmt.Errorf("session start: %w", err)
	}
	defer func() {
		if endErr := database.EndSession(db, nil); endErr != nil {
			log.Printf("ERR: session cleanup failed: %v", endErr)
		}
	}()
	// get data
	dirs, err := getDirs(solutionsDir)
	if err != nil {
		return fmt.Errorf("get dirs: %w", err)
	}
	// insert
	inserted, skipped, err := insertLanguages(db, dirs)
	if err != nil {
		return fmt.Errorf("insert languages: %w", err)
	}
	fmt.Printf("Skipped languages: %d\n", len(skipped))
	insertedNames := make([]string, 0, len(inserted))
	for _, v := range inserted {
		insertedNames = append(insertedNames, v.Name)
	}
	fmt.Printf("Inserted languages: %v\n", insertedNames)
	return nil
}

func getDirs(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	dirs := make([]string, 0, len(entries))
	for _, v := range entries {
		if v.IsDir() {
			dirname := v.Name()
			dirs = append(dirs, dirname)
		}
	}
	return dirs, nil
}

func insertLanguages(db *database.DB, names []string) ([]database.Language, []database.Language, error) {
	var inserted []database.Language
	var skipped []database.Language
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		for _, name := range names {
			lang := database.Language{Name: name}
			result := tx.Where(database.Language{Name: name}).FirstOrCreate(&lang)
			if result.Error != nil {
				return result.Error // Returning error rolls back the transaction
			}
			if result.RowsAffected > 0 {
				inserted = append(inserted, lang)
			} else {
				skipped = append(skipped, lang)
			}
		}
		return nil // Returning nil commits the transaction
	})
	if err != nil {
		return nil, nil, err
	}
	return inserted, skipped, nil
}
