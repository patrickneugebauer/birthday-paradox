package tasks

import (
	"encoding/json"
	"os"
)

type LanguageInfo struct {
	Year   *int    `json:"language_year"`
	Wiki   *string `json:"wiki_link"`
	GitHub *string `json:"github_link"`
}

type LIStruct struct {
	Language string  `json:"dirname"`
	Year     *int    `json:"language_year"`
	Wiki     *string `json:"wiki_link"`
	GitHub   *string `json:"github_link"`
}

// 25% faster
func GetLanguageInfoMap(results map[string]LanguageInfo) error {
	file, err := os.Open(languageInfoFile)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)

	_, err = decoder.Token()
	if err != nil {
		return err
	}

	for decoder.More() {
		var v LIStruct
		err := decoder.Decode(&v)
		if err != nil {
			return err
		}
		results[v.Language] = LanguageInfo{
			Year:   v.Year,
			Wiki:   v.Wiki,
			GitHub: v.GitHub,
		}
	}
	return nil
}
