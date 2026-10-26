package tasks

import (
	"path/filepath"
	"time"
)

const (
	solutionsDir      = "./solutions/complete"
	artifactDir       = "artifacts"
	defaultIterations = 1000
	refreshInterval   = 1
	timeFormat        = time.RFC3339
)

var (
	// language info
	languageInfoFile        = "all-language-data.json"
	languageInfoTsvTempfile = "all-language-data.temp.tsv"
	languageInfoTsv         = "all-language-data.tsv"
	// solutions list
	dockerfileListFile = filepath.Join(artifactDir, "dockerfiles.jsonl")
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
)

type DockerfileMap = map[string][]string
type Dockerfile struct {
	Language         string  `json:"language"`
	Tag              string  `json:"tag"`
	FileLastModified string  `json:"file_last_modified"`
	ImageLastCreated *string `json:"image_last_created"`
	WillRebuild      *bool   `json:"will_rebuild"`
}

type Command struct {
	Tag     string `json:"tag"`
	Command string `json:"command"`
}

type StarResult struct {
	Language      string  `json:"language"`
	Dirname       *string `json:"dirname"`
	Stars         int     `json:"stars"`
	RepoCreatedAt string  `json:"repo_created_at"`
	LastFetchedAt string  `json:"last_fetched_at"`
}

type BuildResult struct {
	Tag             string   `json:"tag"`
	LastBuiltAt     string   `json:"last_built_at"`
	TotalS          *float64 `json:"total_s"`
	PeakRAMBytes    *int64   `json:"peak_ram_bytes"`
	CpuS            *float64 `json:"cpu_s"`
	NetRxBytes      *int64   `json:"net_rx_bytes"` // receive
	NetTxBytes      *int64   `json:"net_tx_bytes"` // transmit
	BlkReadBytes    *int64   `json:"blk_read_bytes"`
	BlkWriteBytes   *int64   `json:"blk_write_bytes"`
	CpuRamPollCount *int64   `json:"cpu_ram_poll_count"`
}

type WeighResult struct {
	Tag            string  `json:"tag"`
	SizeMB         float64 `json:"size_mb"`
	SizeBytes      int64   `json:"size_bytes"`
	ImageUpdatedAt string  `json:"image_updated_at"`
	LastWeighAt    string  `json:"last_weigh_at"`
}

type RunResult struct {
	Tag             string   `json:"image"`
	Iterations      *int     `json:"iterations"`
	SampleSize      *int     `json:"sample_size"`
	Percent         *float64 `json:"percent"`
	Seconds         *float64 `json:"seconds"`
	IPS             *int     `json:"ips"`
	ImageUpdatedAt  string   `json:"image_updated_at"`
	LastRunAt       string   `json:"last_run_at"`
	RuntimeS        *float64 `json:"runtime_s"`
	PeakRAMBytes    *int64   `json:"peak_ram_bytes"`
	CpuS            *float64 `json:"cpu_s"` // cpu time consumed
	CpuRamPollCount *int64   `json:"cpu_ram_poll_count"`
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

type RawReadmeRow struct {
	Tag           string  `json:"tag"`
	Language      string  `json:"language"`
	Year          *int    `json:"year"`
	WikiURL       *string `json:"wiki_url"`
	WikiDisplay   *string `json:"wiki_display"`
	GitHubURL     *string `json:"github_url"`
	GitHubDisplay *string `json:"github_display"`
	Stars         *int    `json:"stars"`
	SizeMB        float64 `json:"size_mb"`
	RamKB         float64 `json:"ram_kb"`
	IPS           int     `json:"ips"`
}
