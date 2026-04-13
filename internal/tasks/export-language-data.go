package tasks

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/patrickneugebauer/birthday-paradox/internal/database"
)

// FLOW
// startSession
// fetchLanguageData -> createExport (create TSV)
// endSession

func ExportLanguageData(db *database.DB) error {
	filename := "language-data.csv"
	data, finalErr := fetchLanguageData(db)
	if finalErr != nil {
		return fmt.Errorf("fetch failed: %w", finalErr)
	}
	finalErr = createExport(filename, data)
	if finalErr != nil {
		return fmt.Errorf("export failed: %w", finalErr)
	}
	// log and return
	fmt.Printf("Successfully exported %d records to %s\n", len(data), filename)
	return nil
}

func fetchLanguageData(db *database.DB) ([]database.LanguageData, error) {
	var data []database.LanguageData
	// We use Find to get all records
	err := db.DB.Find(&data).Error
	return data, err
}

func createExport(filename string, data []database.LanguageData) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	// write data
	writer := csv.NewWriter(file)
	defer writer.Flush()
	header := []string{
		"DisplayName",
		"DirectoryAlias",
		"Year",
		"HasREPL",
		"IsCompiled",
		"IsStaticTyped",
		"IsOO",
		"IsFunctional",
		"Style",
	}
	if err := writer.Write(header); err != nil {
		return err
	}
	for _, ld := range data {
		row := []string{
			ld.DisplayName,
			*ld.DirectoryAlias,
			strconv.Itoa(*ld.Year),
			fmt.Sprintf("%t", *ld.HasREPL),
			fmt.Sprintf("%t", *ld.IsCompiled),
			fmt.Sprintf("%t", *ld.IsStaticTyped),
			fmt.Sprintf("%t", *ld.IsOO),
			fmt.Sprintf("%t", *ld.IsFunctional),
			string(*ld.Style),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}
