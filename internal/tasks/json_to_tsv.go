package tasks

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

func JsonToTsv() error {
	// externals
	inputFile := languageInfoFile
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

	_, err = decoder.Token()
	if err != nil {
		return err
	}

	writer.Write(headers)

	for decoder.More() {
		var x LanguageData
		err := decoder.Decode(&x)
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

type SafeConv struct{}

func (SafeConv) Str(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
func (SafeConv) Int(p *int) string {
	if p == nil {
		return ""
	}
	return strconv.Itoa(*p)
}
func (SafeConv) Bool(p *bool) string {
	if p == nil {
		return ""
	}
	return strconv.FormatBool(*p)
}
