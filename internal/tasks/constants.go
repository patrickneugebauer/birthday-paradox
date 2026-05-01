package tasks

import (
	"path/filepath"
)

const (
	solutionsDir      = "./solutions"
	scaffoldsDir      = "./scaffolds"
	helloWorldsDir    = "./hello-worlds"
	artifactDir       = "artifacts"
	defaultIterations = 1000
	refreshInterval   = 1
)

var (
	// language info
	languageInfoFile = "language-info.csv"
	// solutions list
	dockerfileList = filepath.Join(artifactDir, "dockerfiles.jsonl")
	// stars
	starsRawPayloadsFile     = filepath.Join(artifactDir, "star-raw-payloads.jsonl")
	starsTempRawPayloadsFile = filepath.Join(artifactDir, "star-temp-raw-payloads.jsonl")
	starsTempResultsFile     = filepath.Join(artifactDir, "star-temp-results.jsonl")
	starsResultsFile         = filepath.Join(artifactDir, "star-results.jsonl")
	// build
	buildTempResultsFile = filepath.Join(artifactDir, "build-temp-results.jsonl")
	buildResultsFile     = filepath.Join(artifactDir, "build-results.jsonl")
	// size
	weighTempResultsFile = filepath.Join(artifactDir, "weigh-temp-results.jsonl")
	weighResultsFile     = filepath.Join(artifactDir, "weigh-results.jsonl")
	// run
	runTempResultsFile = filepath.Join(artifactDir, "run-temp-results.jsonl")
	runResultsFile     = filepath.Join(artifactDir, "run-results.jsonl")
	// readme results
	readmeResultsFile          = filepath.Join(artifactDir, "readme-results.jsonl")
	formattedReadmeResultsFile = filepath.Join(artifactDir, "formatted-readme-results.jsonl")
	// readmes
	readmeFile           = "readme.md"
	readmeFileByLanguage = filepath.Join("tables", "results-by-language.md")
	readmeFileByYear     = filepath.Join("tables", "results-by-year.md")
	readmeFileByStars    = filepath.Join("tables", "results-by-stars.md")
	readmeFileBySize     = filepath.Join("tables", "results-by-size.md")
)

type DockerfileMap = map[string][]string
type Dockerfile struct {
	Language         string  `json:"language"`
	Filename         string  `json:"dockerfile"`
	Runtime          *string `json:"runtime"`
	DataStructure    *string `json:"data_structure"`
	ExecutionMethod  *string `json:"execution_method"`
	Tag              string  `json:"tag"`
	Directory        string  `json:"directory"`
	FileLastModified int64   `json:"file_last_modified"`
	ImageLastCreated int64   `json:"image_last_created"`
	WillRebuild      *bool   `json:"will_rebuild"`
}

type Command struct {
	Tag     string `json:"tag"`
	Command string `json:"command"`
}

type StarResult struct {
	Language      string `json:"language"`
	Stars         int    `json:"stars"`
	RepoCreatedAt int64  `json:"repo_created_at"`
	LastFetchedAt int64  `json:"last_fetched_at"`
}

type BuildResult struct {
	Tag         string `json:"tag"`
	LastBuiltAt int64  `json:"last_built_at"`
}

type WeighResult struct {
	Tag            string  `json:"tag"`
	SizeMB         float64 `json:"size_mb"`
	SizeBytes      int64   `json:"size_bytes"`
	ImageUpdatedAt int64   `json:"image_updated_at"`
	LastWeighAt    int64   `json:"last_weigh_at"`
}

type RunResult struct {
	Tag            string   `json:"image"`
	Iterations     *int     `json:"iterations"`
	SampleSize     *int     `json:"sample_size"`
	Percent        *float64 `json:"percent"`
	Seconds        *float64 `json:"seconds"`
	IPS            *int     `json:"ips"`
	ImageUpdatedAt int64    `json:"image_updated_at"`
	LastRunAt      int64    `json:"last_run_at"`
}

type GithubRepo struct {
	Name            string `json:"name"`
	FullName        string `json:"full_name"`
	HtmlUrl         string `json:"html_url"`
	Homepage        string `json:"homepage"`
	StargazersCount int    `json:"stargazers_count"`
	Language        string `json:"language"`
	MirrorUrl       string `json:"mirror_url"`
	Archived        bool   `json:"archived"`
	Disabled        bool   `json:"disabled"`
	CreatedAt       string `json:"created_at"`
}

type ReadmeRow struct {
	Tag               string  `json:"tag"`
	Language          string  `json:"language"`
	Runtime           *string `json:"runtime"`
	DataStructure     *string `json:"data_structure"`
	ExecutionMethod   *string `json:"execution_method"`
	Year              int     `json:"year"`
	WikiURL           string  `json:"wiki_url"`
	WikiDisplay       string  `json:"wiki_display"`
	GitHubURL         string  `json:"github_url"`
	GitHubDisplay     string  `json:"github_display"`
	Stars             int     `json:"stars"`
	SizeMB            float64 `json:"size_mb"`
	IPS               int     `json:"ips"`
	FormattedLanguage string  `json:"formatted_language"`
	FormattedRuntime  string  `json:"formatted_runtime"`
	FormattedDataType string  `json:"formatted_data_type"`
	FormattedExecMode string  `json:"formatted_exec_mode"`
	FormattedYear     string  `json:"formatted_year"`
	FormattedStars    string  `json:"formatted_stars"`
	FormattedSizeMB   string  `json:"formatted_size_mb"`
	FormattedIPS      string  `json:"formatted_ips"`
}
