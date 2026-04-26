package main

import (
	"fmt"
	"log"
	"os"

	"github.com/patrickneugebauer/birthday-paradox/internal/tasks"
)

func main() {
	// check args
	if len(os.Args) < 2 {
		log.Fatal("Usage: myapp <command> (e.g., list-dirs, export, sync)")
	}
	command := os.Args[1]
	var err error
	// run cpmmand
	switch command {
	case "map-files":
		err = tasks.MakeFileTree()
	case "star", "stars":
		err = tasks.Stars()
	case "build":
		err = tasks.Build()
	case "weigh":
		err = tasks.Weigh()
	case "run":
		err = tasks.Run()
	case "readme":
		err = tasks.Readme()
	default:
		err = fmt.Errorf("unknown command: %s", command)
	}
	if err != nil {
		log.Printf("Command failed: %v", err)
		os.Exit(1)
	}
}
