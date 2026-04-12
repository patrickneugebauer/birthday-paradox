package database

import (
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	SessionType string
	ExitCode    *int64
	Error       *string
}

type Language struct {
	gorm.Model
	Name     string    `gorm:"not null;unique"`
	Runtimes []Runtime `gorm:"foreignKey:LanguageID"`
}

type LanguageData struct {
	gorm.Model
	Name           string `gorm:"not null;unique"`
	Compiled       bool
	StaticType     bool
	ObjectOriented bool
	Functional     bool
	Style          string
	LanguageID     *uint     `gorm:"index"`
	Language       *Language `gorm:"foreignKey:LanguageID"`
}

type Runtime struct {
	gorm.Model
	Filename    string       `gorm:"not null;uniqueIndex:idx_lang_file"`
	LanguageID  uint         `gorm:"not null;uniqueIndex:idx_lang_file"`
	Language    *Language    `gorm:"foreignKey:LanguageID"`
	DockerImage *DockerImage `gorm:"foreignKey:RuntimeID"`
}

type RuntimeData struct {
	gorm.Model
	GithubURL string       `gorm:"not null;unique"`
	RuntimeID uint         `gorm:"not null"`
	Runtime   *Runtime     `gorm:"foreignKey:RuntimeID"`
	Stars     []GithubStar `gorm:"foreignKey:RuntimeDataID"`
}

type GithubStar struct {
	gorm.Model
	StarCount     int          `gorm:"not null"`
	RuntimeDataID uint         `gorm:"not null"`
	RuntimeData   *RuntimeData `gorm:"foreignKey:RuntimeDataID"`
}

type DockerImage struct {
	gorm.Model
	Tag           string         `gorm:"not null;unique"`
	RuntimeID     uint           `gorm:"not null;unique"`
	Runtime       *Runtime       `gorm:"foreignKey:RuntimeID"`
	ImageSize     *ImageSize     `gorm:"foreignKey:ImageID"`
	ContainerRuns []ContainerRun `gorm:"foreignKey:ImageID"`
}

type ImageSize struct {
	gorm.Model
	Size    string       `gorm:"not null"`
	ImageID uint         `gorm:"not null;unique"`
	Image   *DockerImage `gorm:"foreignKey:ImageID"`
}

type ContainerRun struct {
	gorm.Model
	Tag                 string       `gorm:"not null"`
	Iterations          int          `gorm:"not null"`
	SampleSize          int          `gorm:"not null"`
	Percent             float64      `gorm:"not null"`
	Seconds             float64      `gorm:"not null"`
	IterationsPerSecond int          `gorm:"not null"`
	ImageID             uint         `gorm:"not null"`
	Image               *DockerImage `gorm:"foreignKey:ImageID"`
}

// ==================================================

type LanguageReport struct {
	Name                string
	Size                string
	IterationsPerSecond int
	Compiled            bool
	StaticType          bool
	ObjectOriented      bool
	Functional          bool
	Style               string
}

func ReadMeQuery(db *gorm.DB) ([]LanguageReport, error) {
	var results []LanguageReport

	tx := db.Table("languages").
		Select(`
			languages.name,
			CASE
				WHEN dockerfiles.filename LIKE '%.%'
				THEN SUBSTR(dockerfiles.filename, INSTR(dockerfiles.filename, '.') + 1)
				ELSE NULL
			END AS extension,
			image_sizes.size,
			container_runs.iterations_per_second,
			languages.compiled,
			languages.static_type,
			languages.object_oriented,
			languages.functional,
			languages.style"
		`).Joins("JOIN directories ON directories.language_id = languages.id").
		Joins("JOIN dockerfiles ON dockerfiles.directory_id = directories.id").
		Joins("JOIN docker_images ON docker_images.dockerfile_id = dockerfiles.id").
		Joins("JOIN image_sizes ON image_sizes.image_id = docker_images.id").
		Joins("JOIN container_runs ON container_runs.image_id = docker_images.id").
		Scan(&results)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return results, nil
}
