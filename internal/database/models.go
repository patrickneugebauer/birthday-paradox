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
	Directory string    `gorm:"not null;unique"`
	Runtimes  []Runtime `gorm:"foreignKey:LanguageID"`
}

type Style string

const (
	Functional Style = "functional"
	Imperative Style = "imperative"
	Lisp       Style = "lisp"
	Scientific Style = "scientific"
)

type LanguageData struct {
	gorm.Model
	DisplayName    string  `gorm:"not null;unique"`
	DirectoryAlias *string `gorm:"uniqueIndex"`
	Year           *int
	HasREPL        *bool
	IsCompiled     *bool
	IsStaticTyped  *bool
	IsOO           *bool
	IsFunctional   *bool
	Style          *Style    `gorm:"check:style IN ('functional', 'imperative', 'lisp', 'scientific')"`
	LanguageID     *uint     `gorm:"index"`
	Language       *Language `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Runtime struct {
	gorm.Model
	RuntimeName    *string      `gorm:"uniqueIndex:idx_lang_sol"`
	DockerfileName string       `gorm:"not null"`
	LanguageID     uint         `gorm:"not null;uniqueIndex:idx_lang_sol"`
	Language       *Language    `gorm:"foreignKey:LanguageID"`
	DockerImage    *DockerImage `gorm:"foreignKey:RuntimeID"`
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
