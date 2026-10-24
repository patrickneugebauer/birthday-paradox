package tasks

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func TsvToJsonl() error {
	// externals
	inputFile := "all-language-data.tsv"
	outputTempfile := "all-language-data.temp.jsonl"
	outputFile := "all-language-data.jsonl"
	// ====================
	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = '\t'

	f, err := os.Create(outputTempfile)
	if err != nil {
		return fmt.Errorf("create %s: %w", outputTempfile, err)
	}
	writer := bufio.NewWriter(f)
	defer writer.Flush()
	encoder := json.NewEncoder(writer)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		p := SafeParse{}
		encoder.Encode(LanguageData{
			Language:                record[0],
			Runtime:                 &record[1],
			Dirname:                 &record[2],
			IsCompiled:              p.Bool(record[3]),
			TranspileTo:             &record[4],
			OfficialDockerImage:     &record[5],
			IsStandalone:            p.Bool(record[6]),
			OfficialDockerImageLink: &record[7],
			LanguageYear:            p.Int(record[8]),
			RuntimeYear:             p.Int(record[9]),
			GithubLink:              &record[10],
			GithubIsMirror:          p.Bool(record[11]),
			GitlabLink:              &record[12],
			SourcehutLink:           &record[13],
			SvnLink:                 &record[14],
			SourceforgeLink:         &record[15],
			WikiLink:                &record[16],
			WebsiteLink:             &record[17],
			GithubStars:             p.Int(record[18]),
			Notes:                   &record[19],
		})
	}

	writer.Flush()
	// cleanup
	if err := os.Rename(outputTempfile, outputFile); err != nil {
		return fmt.Errorf("finalize results: %w", err)
	}
	fmt.Printf("wrote: %s\n", outputFile)
	return nil
}

type SafeParse struct{}

func (SafeParse) Int(s string) *int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &i
}
func (SafeParse) Bool(s string) *bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return nil
	}
	return &b
}
