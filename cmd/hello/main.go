package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand/v2"

	"github.com/patrickneugebauer/birthday-paradox/internal/database"
)

func main() {
	// beginning
	var err error
	db, err := database.StartSession("hello")
	if err != nil {
		log.Fatalf("ERR: session not started %v\n", err)
	}
	fmt.Printf("session started\n")
	// middle
	hasError := rand.IntN(2) == 1
	if hasError {
		err = errors.New("ERR: random error")
	}
	// end
	err = database.EndSession(db, err)
	if err != nil {
		log.Fatalf("ERR: session not finished %v\n", err)
	}
	fmt.Printf("session finished\n")
}
