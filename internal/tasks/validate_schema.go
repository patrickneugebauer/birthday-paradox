package tasks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type LanguageData struct {
	Language                string  `json:"language"`
	Runtime                 *string `json:"runtime"`
	Dirname                 *string `json:"dirname"`
	IsCompiled              *bool   `json:"is_compiled"`
	TranspileTo             *string `json:"transpile_to"`
	OfficialDockerImage     *string `json:"official_docker_image"`
	IsStandalone            *string `json:"is_standalone"`
	OfficialDockerImageLink *string `json:"official_docker_image_link"`
	LanguageYear            *int    `json:"language_year"`
	RuntimeYear             *int    `json:"runtime_year"`
	GithubLink              *string `json:"github_link"`
	GithubIsMirror          *bool   `json:"github_is_mirror"`
	GitlabLink              *string `json:"gitlab_link"`
	SourcehutLink           *string `json:"sourcehut_link"`
	SvnLink                 *string `json:"svn_link"`
	SourceforgeLink         *string `json:"sourceforge_link"`
	WikiLink                *string `json:"wiki_link"`
	WebsiteLink             *string `json:"website_link"`
	GithubStars             *int    `json:"github_stars"`
	Notes                   *string `json:"notes"`
}

func ValidateSchema() error {
	content, err := os.ReadFile("all-language-data.json")
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	keys := []string{
		"language", "runtime", "dirname", "is_compiled", "transpile_to",
		"official_docker_image", "is_standalone", "official_docker_image_link",
		"language_year", "runtime_year", "github_link", "github_is_mirror",
		"gitlab_link", "sourcehut_link", "svn_link", "sourceforge_link",
		"wiki_link", "website_link", "github_stars", "notes",
	}

	expectedKeys := make(map[string]bool)
	for _, k := range keys {
		expectedKeys[k] = true
	}

	var validationErrors []string

	// 1. Parse into raw maps first so we can iterate through items individually
	var rawData []json.RawMessage
	if err := json.Unmarshal(content, &rawData); err != nil {
		return fmt.Errorf("Critical Syntax Error (JSON is malformed): %v", err)
	}

	for i, item := range rawData {
		// Calculate a rough line number for this specific object
		// This finds where this object starts in the main file
		lineNum := bytes.Count(content[:bytes.Index(content, item)], []byte("\n")) + 1

		// A. Check for Type Mismatches
		var ld LanguageData
		if err := json.Unmarshal(item, &ld); err != nil {
			validationErrors = append(validationErrors, fmt.Sprintf("Type Error at line %d (index %d): %v", lineNum, i, err))
		}

		// B. Check for Key Count/Names
		var obj map[string]interface{}
		json.Unmarshal(item, &obj)

		langName := "unknown"
		if val, ok := obj["language"].(string); ok {
			langName = val
		}

		for key := range obj {
			if !expectedKeys[key] {
				validationErrors = append(validationErrors, fmt.Sprintf("Key Error in '%s' (line %d): unknown key %q", langName, lineNum, key))
			}
		}

		if len(obj) < len(keys) {
			for key := range expectedKeys {
				if _, exists := obj[key]; !exists {
					validationErrors = append(validationErrors, fmt.Sprintf("Key Error in '%s' (line %d): missing key %q", langName, lineNum, key))
				}
			}
		}
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("Validation failed with %d errors:\n%s", len(validationErrors), strings.Join(validationErrors, "\n"))
	}

	return nil
}
