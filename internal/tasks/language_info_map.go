package tasks

import (
	"encoding/json"
	"os"
)

type LanguageInfoStruct struct {
	Language string `json:"language"`
	Year     int    `json:"language_year"`
	Wiki     string `json:"wiki_link"`
	Github   string `json:"github_link"`
}

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
		var v LanguageInfoStruct
		err := decoder.Decode(&v)
		if err != nil {
			return err
		}
		results[v.Language] = LanguageInfo{
			Year:   v.Year,
			Wiki:   v.Wiki,
			GitHub: v.Github,
		}
	}
	return nil
}

type LIStruct struct {
	K string `json:"language"`
	V LanguageInfo
}

// 25% faster
func GetLanguageInfoMap2(results map[string]LanguageInfo) error {
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
		results[v.K] = v.V
	}
	return nil
}

// stick with json for faster compile time
func WriteMessagePackMap() error {
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
		// optimize this memory usage?
		var v LanguageInfoStruct
		err := decoder.Decode(&v)
		if err != nil {
			return err
		}
		// write to messagepack file
	}
	return nil
}
