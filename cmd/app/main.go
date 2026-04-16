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
	case "pre-build":
		err = tasks.PreBuild()
	case "build":
		err = tasks.Build()
	case "pre-run":
		err = tasks.PreRun()
	case "run":
		err = tasks.Run()
	case "pre-weigh":
		err = tasks.PreWeigh()
	case "weigh":
		err = tasks.Weigh()
	case "post-weigh":
		err = tasks.PostWeigh()
	default:
		err = fmt.Errorf("unknown command: %s", command)
	}
	if err != nil {
		log.Printf("Command failed: %v", err)
		os.Exit(1)
	}
}
