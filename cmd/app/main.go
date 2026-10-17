package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/patrickneugebauer/birthday-paradox/internal/tasks"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: myapp <command> (e.g., list-dirs, export, sync)")
	}
	command := os.Args[1]
	var err error
	switch command {
	case "rows":
		rows := make([]tasks.RawReadmeRow, 0, 128)
		err = tasks.BuildReadmeRows(&rows)
		fileData, err := json.MarshalIndent(rows, "", "    ")
		if err != nil {
			fmt.Printf("error %w", err)
		}
		err = os.WriteFile("_rows.json", fileData, 0644)
	case "info":
		start := time.Now()
		info := make(map[string]tasks.LanguageInfo)
		err = tasks.GetLanguageInfoMap(info)
		diff := time.Since(start)
		fmt.Printf("%.9f\n", diff.Seconds())
		fileData, err := json.MarshalIndent(info, "", "    ")
		if err != nil {
			fmt.Printf("error %w", err)
		}
		fname := "_info.json"
		err = os.WriteFile(fname, fileData, 0644)
		fmt.Printf("%s\n", fname)
	case "info2":
		start := time.Now()
		info := make(map[string]tasks.LanguageInfo)
		err = tasks.GetLanguageInfoMap2(info)
		diff := time.Since(start)
		fmt.Printf("%.9f\n", diff.Seconds())
		fileData, err := json.MarshalIndent(info, "", "    ")
		if err != nil {
			fmt.Printf("error %w", err)
		}
		fname := "_info.json"
		err = os.WriteFile(fname, fileData, 0644)
		fmt.Printf("%s\n", fname)
	case "schema":
		err = tasks.ValidateSchema()
	case "map", "files", "map-files":
		err = tasks.MakeFileTree()
	case "star", "stars":
		err = tasks.Stars()
	case "github-languages", "github-langs", "gh-langs", "gh-languages":
		err = tasks.GithubLanguages()
	case "wiki-languages", "wiki-langs":
		err = tasks.WikiLanguages()
	case "bitbucket-languages", "bb-languages", "bitbucket-langs", "bb-langs":
		err = tasks.BitbucketLanguages()
	case "build":
		err = tasks.Build()
	case "build-all":
		err = tasks.BuildAll()
	case "weigh":
		err = tasks.Weigh()
	case "run":
		err = tasks.Run()
	case "run-all":
		err = tasks.RunAll()
	case "run-some":
		if len(os.Args) < 3 {
			log.Fatal("Usage: myapp run-some <search>")
		}
		err = tasks.RunSome(os.Args[2])
	case "readme":
		err = tasks.Readme()
	case "all":
		if err = tasks.MakeFileTree(); err != nil {
			break
		}
		fmt.Println()
		if err = tasks.Stars(); err != nil {
			break
		}
		fmt.Println()
		if err = tasks.Build(); err != nil {
			break
		}
		fmt.Println()
		if err = tasks.Weigh(); err != nil {
			break
		}
		fmt.Println()
		if err = tasks.Run(); err != nil {
			break
		}
		fmt.Println()
		err = tasks.Readme()
	default:
		err = fmt.Errorf("unknown command: %s", command)
	}
	if err != nil {
		log.Printf("Command failed: %v", err)
		os.Exit(1)
	}
}
