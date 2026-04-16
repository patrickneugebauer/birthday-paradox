package tasks

import (
	"path/filepath"
)

const (
	solutionsDir      = "./solutions"
	artifactDir       = "artifacts"
	defaultIterations = 1000
)

var (
	dockerfilesList     = filepath.Join(artifactDir, "dockerfiles.yaml")
	buildScript         = filepath.Join(artifactDir, "build.sh")
	buildArtifacts      = filepath.Join(artifactDir, "builds.jsonl")
	runScript           = filepath.Join(artifactDir, "run.sh")
	weighScript         = filepath.Join(artifactDir, "weigh.sh")
	sizeFile            = filepath.Join(artifactDir, "image-sizes.jsonl")
	sizeFileMB          = filepath.Join(artifactDir, "image-sizes-mb.jsonl")
	resultsFile         = filepath.Join(artifactDir, "results.jsonl")
	previousResultsFile = filepath.Join(artifactDir, "previous-results.jsonl")
	tempResultsFile     = filepath.Join(artifactDir, "temp-results.jsonl")
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
