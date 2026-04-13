package tasks

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/patrickneugebauer/birthday-paradox/internal/database"
	"gorm.io/gorm"
)

func ImportLanguageData(db *database.DB) error {
	fileName := "language-data.csv"
	data, err := readLanguageData(fileName)
	if err != nil {
		return fmt.Errorf("read error: %w", err)
	}
	added, updated, err := persistLanguageData(db, data)
	if err != nil {
		return fmt.Errorf("persist error: %w", err)
	}
	fmt.Printf("Import complete: %d records added, %d records updated\n", len(added), len(updated))
	return nil
}

func readLanguageData(filename string) ([]database.LanguageData, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	r := csv.NewReader(file)
	if _, err := r.Read(); err != nil {
		return nil, err
	}
	// loop over rows
	var results []database.LanguageData
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		languageData := database.LanguageData{
			DisplayName:    record[0],
			DirectoryAlias: &record[1],
			Year:           ptrInt(record[2]),
			HasREPL:        ptrBool(record[3]),
			IsCompiled:     ptrBool(record[4]),
			IsStaticTyped:  ptrBool(record[5]),
			IsOO:           ptrBool(record[6]),
			IsFunctional:   ptrBool(record[7]),
			Style:          ptrStyle(record[8]),
		}
		results = append(results, languageData)
	}
	return results, nil
}

func ptrBool(s string) *bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return nil // Or return a pointer to false if preferred
	}
	return &b
}

func ptrInt(s string) *int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &i
}

func ptrStyle(s string) *database.Style {
	st := database.Style(s)
	return &st
}

func persistLanguageData(db *database.DB, data []database.LanguageData) ([]database.LanguageData, []database.LanguageData, error) {
	var added []database.LanguageData
	var updated []database.LanguageData

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		for _, ld := range data {
			// 1. Look for existing record first
			var existing database.LanguageData
			result := tx.Where("directory_alias = ?", ld.DirectoryAlias).First(&existing)
			if result.Error == nil {
				tx.Model(&existing).Updates(ld)
				updated = append(updated, existing)
			} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				tx.Create(&ld)
				added = append(added, ld)
			} else {
				return result.Error // Actual DB error
			}
		}
		return nil
	})

	return added, updated, err
}
