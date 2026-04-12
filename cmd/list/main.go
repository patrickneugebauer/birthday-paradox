package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const dbname = "./dev.db"
const basepath = "./solutions"
const outfile = "solutions.tsv"
const comma = '\t'

func main() {
	fmt.Println("hello")
	// database setup
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	stmt, err := db.Prepare("INSERT OR IGNORE INTO languages(name) VALUES(?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// get dirs
	entries, err := os.ReadDir(basepath)
	if err != nil {
		log.Fatal(err)
	}
	dirs := make([]string, 0, len(entries))
	for _, v := range entries {
		if v.IsDir() {
			dirname := v.Name()
			fmt.Println(dirname)
			dirs = append(dirs, dirname)
			_, err = stmt.Exec(dirname)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// get dockerfiles
	numDirs := len(dirs)
	type result struct {
		Dir        string `json:"dir"`
		Dockerfile string `json:"dockerfile"`
		Name       string `json:"name"`
	}
	results := make([]result, 0, numDirs*2)
	for _, dir := range dirs {
		entries, err := os.ReadDir(filepath.Join(basepath, dir))
		if err != nil {
			log.Fatal(err)
		}
		numFiles := len(entries)
		dockerfiles := make([]string, 0, numFiles)
		for _, entry := range entries {
			fname := entry.Name()
			if !entry.IsDir() && strings.HasPrefix(fname, "Dockerfile") {
				dockerfiles = append(dockerfiles, fname)
				_, after, found := strings.Cut(fname, ".")
				solutionName := dir
				if found {
					solutionName += "-" + after
				}
				result := result{
					Dir:        dir,
					Dockerfile: fname,
					Name:       solutionName,
				}
				results = append(results, result)
			}
		}
	}

	f, err := os.Create(outfile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	// the type signature says it takes a io.Writer, but we pass a os.File
	w := csv.NewWriter(f)
	w.Comma = comma
	defer w.Flush()
	w.Write([]string{"dir", "dockerfile", "name"})
	rows := make([][]string, 0, len(results))
	for _, v := range results {
		row := []string{v.Dir, v.Dockerfile, v.Name}
		rows = append(rows, row)
	}
	w.WriteAll(rows)
}
