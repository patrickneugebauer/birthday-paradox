package tasks

import (
	"encoding/json"
	"os"
)

type LIStruct struct {
	K string `json:"language"`
	V LanguageInfo
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
		err := decoder.Decode(&v) // this is 25% faster than decoding and creating
		if err != nil {
			return err
		}
		results[v.K] = v.V
	}
	return nil
}
