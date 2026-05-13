package main

import (
	"fmt"
	"log"
	"os"

	"github.com/patrickneugebauer/birthday-paradox/internal/tasks"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: myapp <command> (e.g., list-dirs, export, sync)")
	}
	command := os.Args[1]
	var err error
	switch command {
	case "map", "files", "map-files":
		err = tasks.MakeFileTree()
	case "star", "stars":
		err = tasks.Stars()
	case "github-languages":
		err = tasks.GithubLanguages()
	case "wiki-languages":
		err = tasks.WikiLanguages()
	case "bitbucket-languages":
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
