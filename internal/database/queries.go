package database

import (
	"gorm.io/gorm"
)

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
