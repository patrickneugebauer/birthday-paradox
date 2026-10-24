package tasks

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

func JsonlToTsv() error {
	// externals
	inputFile := "all-language-data.jsonl"
	outputTempfile := languageInfoTsvTempfile
	outputFile := languageInfoTsv
	// ====================
	headers := []string{
		"language",
		"runtime",
		"dirname",
		"is compiled",
		"transpile to",
		"official docker image",
		"is standalone",
		"offical docker image link",
		"language year",
		"runtime year",
		"github link",
		"github is mirror",
		"gitlab link",
		"sourcehut link",
		"svn link",
		"sourceforge link",
		"wiki link",
		"website link",
		"github stars",
		"notes",
	}
	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)

	f, err := os.Create(outputTempfile)
	if err != nil {
		return fmt.Errorf("create %s: %w", outputTempfile, err)
	}
	writer := csv.NewWriter(f)
	writer.Comma = '\t' // Sets delimiter to Tab
	defer writer.Flush()

	writer.Write(headers)

	for {
		var x LanguageData
		err := decoder.Decode(&x)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		f := SafeConv{}
		writer.Write([]string{
			x.Language,
			f.Str(x.Runtime),
			f.Str(x.Dirname),
			f.Bool(x.IsCompiled),
			f.Str(x.TranspileTo),
			f.Str(x.OfficialDockerImage),
			f.Bool(x.IsStandalone),
			f.Str(x.OfficialDockerImageLink),
			f.Int(x.LanguageYear),
			f.Int(x.RuntimeYear),
			f.Str(x.GithubLink),
			f.Bool(x.GithubIsMirror),
			f.Str(x.GitlabLink),
			f.Str(x.SourcehutLink),
			f.Str(x.SvnLink),
			f.Str(x.SourceforgeLink),
			f.Str(x.WikiLink),
			f.Str(x.WebsiteLink),
			f.Int(x.GithubStars),
			f.Str(x.Notes),
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
