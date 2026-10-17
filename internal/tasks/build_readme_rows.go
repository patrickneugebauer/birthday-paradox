package tasks

import "fmt"

func BuildReadmeRows(rows *[]RawReadmeRow) error {
	// Load all data
	dockerfiles, err := loadJson[Dockerfile](dockerfileListFile)
	if err != nil {
		return err
	}

	languageInfo := make(map[string]LanguageInfo)
	err = GetLanguageInfoMap(languageInfo)
	if err != nil {
		return err
	}

	stars, err := loadJson[StarResult](starsResultsFile)
	if err != nil {
		return err
	}

	sizes, err := loadJson[WeighResult](weighResultsFile)
	if err != nil {
		return err
	}

	runs, err := loadJson[RunResult](runResultsFile)
	if err != nil {
		return err
	}

	// Build a map of tag → dockerfile info for quick lookup
	dockerfileMap := make(map[string]*Dockerfile)
	for i := range dockerfiles {
		dockerfileMap[dockerfiles[i].Tag] = &dockerfiles[i]
	}

	// Build a map of tag → size for quick lookup
	sizeMap := make(map[string]float64)
	for _, size := range sizes {
		sizeMap[size.Tag] = size.SizeMB
	}

	// Build a map of language → stars for quick lookup
	starMap := make(map[string]int)
	for _, star := range stars {
		starMap[star.Language] = star.Stars
	}

	// Build readme rows from runs (the main dataset)
	// var rows []ReadmeRow
	for _, run := range runs {
		tag := run.Tag
		fmt.Printf("%v\n", run.Tag)

		// Look up dockerfile/runtime
		dockerfileInfo, ok := dockerfileMap[tag]
		if !ok {
			fmt.Println("no dockerfile, skip")
			continue // skip if no dockerfile for this tag
		}

		language := dockerfileInfo.Language

		// Look up year, wiki, github from language info
		info := languageInfo[language]
		year := info.Year
		wikiURL := info.Wiki
		githubURL := info.GitHub

		// Look up stars
		starsCount := starMap[language]

		// Look up size
		sizeMB := sizeMap[tag]

		ips := 0
		if run.IPS != nil {
			ips = *run.IPS
		}

		row := RawReadmeRow{
			Tag:             tag,
			Language:        language,
			Runtime:         dockerfileInfo.Runtime,
			DataStructure:   dockerfileInfo.DataStructure,
			ExecutionMethod: dockerfileInfo.ExecutionMethod,
			Year:            year,
			WikiURL:         wikiURL,
			WikiDisplay:     extractWikiDisplay(wikiURL),
			GitHubURL:       githubURL,
			GitHubDisplay:   extractGitHubDisplay(githubURL),
			Stars:           starsCount,
			SizeMB:          sizeMB,
			IPS:             ips,
		}

		*rows = append(*rows, row)
	}
	return nil
}
