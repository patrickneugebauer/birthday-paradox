package tasks

import (
	"path/filepath"
)

const (
	solutionsDir      = "./solutions"
	referenceDir      = "reference"
	artifactDir       = "artifacts"
	defaultIterations = 1000
	capacity          = 70
)

var (
	// solutions list
	dockerfileMap  = filepath.Join(artifactDir, "dockerfiles.yaml")
	dockerfileList = filepath.Join(artifactDir, "dockerfiles.jsonl")
	// language info
	githubLinksFile = filepath.Join(referenceDir, "github-links.csv")
	wikiLinksFile   = filepath.Join(referenceDir, "wiki-links.csv")
	// stars
	starScript      = filepath.Join(artifactDir, "stars.sh")
	githubStarsFile = filepath.Join(artifactDir, "github-stars.jsonl")
	// build
	buildScript    = filepath.Join(artifactDir, "build.sh")
	buildArtifacts = filepath.Join(artifactDir, "builds.jsonl")
	// size
	weighScript = filepath.Join(artifactDir, "weigh.sh")
	sizeFile    = filepath.Join(artifactDir, "image-sizes.jsonl")
	// run
	runScript           = filepath.Join(artifactDir, "run.sh")
	tempResultsFile     = filepath.Join(artifactDir, "temp-results.jsonl")
	resultsFile         = filepath.Join(artifactDir, "results.jsonl")
	previousResultsFile = filepath.Join(artifactDir, "previous-results.jsonl")
)

type BuildConfig struct {
	DirName    string
	Dockerfile string
	Tag        string
	Context    string
}

type BuildArtifact struct {
	Tag string `json:"tag"`
}

type ImageInfo struct {
	Repository string  `json:"repository"`
	SizeMB     float64 `json:"size_mb"`
}

type RunResult struct {
	Image      string  `json:"image"`
	Iterations int     `json:"iterations"`
	SampleSize int     `json:"sample_size"`
	Percent    float64 `json:"percent"`
	Seconds    float64 `json:"seconds"`
	IPS        int     `json:"ips"`
}
