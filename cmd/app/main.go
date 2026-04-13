package main

import (
	"fmt"
	"log"
	"os"

	"github.com/patrickneugebauer/birthday-paradox/internal/database"
	"github.com/patrickneugebauer/birthday-paradox/internal/tasks"
)

func main() {
	// check args
	if len(os.Args) < 2 {
		log.Fatal("Usage: myapp <command> (e.g., list-dirs, export, sync)")
	}
	command := os.Args[1]
	// connect DB
	db, err := database.StartSession(command)
	if err != nil {
		log.Fatalf("Session start failed: %v", err)
	}
	var finalErr error
	defer func() {
		database.EndSession(db, finalErr)
	}()
	// run cpmmand
	switch command {
	case "hello", "hello-world":
		tasks.Hello()
	case "languages":
		finalErr = tasks.GetLanguages(db)
	case "runtimes":
		finalErr = tasks.GetRuntimes(db)
	case "export-language-data", "export-lang-data":
		finalErr = tasks.ExportLanguageData(db)
	case "import-language-data", "import-lang-data":
		finalErr = tasks.ImportLanguageData(db)
	default:
		finalErr = fmt.Errorf("unknown command: %s", command)
	}
	if finalErr != nil {
		log.Printf("Command failed: %v", finalErr)
		os.Exit(1)
	}
}
