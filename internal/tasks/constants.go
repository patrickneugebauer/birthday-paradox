package tasks

import (
	"path/filepath"
)

const (
	solutionsDir      = "./solutions"
	referenceDir      = "reference"
	artifactDir       = "artifacts"
	defaultIterations = 1000
)

var (
	// solutions list
	dockerfileList = filepath.Join(artifactDir, "dockerfiles.jsonl")
	// language info
	yearsFile       = filepath.Join(referenceDir, "years.csv")
	githubLinksFile = filepath.Join(referenceDir, "github-links.csv")
	wikiLinksFile   = filepath.Join(referenceDir, "wiki-links.csv")
	// stars
	starsCommandsFile    = filepath.Join(artifactDir, "star-commands.sh")
	starsDataFile        = filepath.Join(artifactDir, "star-data.json")
	starsTempResultsFile = filepath.Join(artifactDir, "star-temp-results.jsonl")
	starsResultsFile     = filepath.Join(artifactDir, "star-results.jsonl")
	// build
	buildCommandsFile    = filepath.Join(artifactDir, "build-commands.sh")
	buildTempResultsFile = filepath.Join(artifactDir, "build-temp-results.jsonl")
	buildResultsFile     = filepath.Join(artifactDir, "build-results.jsonl")
	// size
	weighCommandsFile = filepath.Join(artifactDir, "weigh-commands.sh")
	weighResultsFile  = filepath.Join(artifactDir, "weigh-results.jsonl")
	// run
	runCommandsFile    = filepath.Join(artifactDir, "run-commands.sh")
	runTempResultsFile = filepath.Join(artifactDir, "run-temp-results.jsonl")
	runResultsFile     = filepath.Join(artifactDir, "run-results.jsonl")
	// readme
	readmeResultsFile = filepath.Join(artifactDir, "readme-results.jsonl")
	readmeFile        = "README.md"
)

type DockerfileMap = map[string][]string
type Dockerfile struct {
	Language string  `json:"language"`
	Filename string  `json:"dockerfile"`
	Runtime  *string `json:"runtime"`
	Tag      string  `json:"tag"`
}

type Command struct {
	Tag     string `json:"tag"`
	Command string `json:"command"`
}

type StarResult struct {
	Language string `json:"language"`
	Stars    int    `json:"stars"`
}

type BuildResult struct {
	Tag string `json:"tag"`
}

type WeighResult struct {
	Tag    string  `json:"tag"`
	SizeMB float64 `json:"size_mb"`
}

type RunResult struct {
	Tag        string  `json:"image"`
	Iterations int     `json:"iterations"`
	SampleSize int     `json:"sample_size"`
	Percent    float64 `json:"percent"`
	Seconds    float64 `json:"seconds"`
	IPS        int     `json:"ips"`
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
}

type ReadmeRow struct {
	Tag           string  `json:"tag"`
	Language      string  `json:"language"`
	Runtime       string  `json:"runtime"`
	Year          int     `json:"year"`
	WikiURL       string  `json:"wiki_url"`
	WikiDisplay   string  `json:"wiki_display"`
	GitHubURL     string  `json:"github_url"`
	GitHubDisplay string  `json:"github_display"`
	Stars         int     `json:"stars"`
	SizeMB        float64 `json:"size_mb"`
	IPS           int     `json:"ips"`
}
