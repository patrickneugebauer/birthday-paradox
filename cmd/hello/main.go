package main

import (
	"errors"
	"fmt"
	"math/rand/v2"

	"github.com/patrickneugebauer/birthday-paradox/internal/db"
)

func main() {
	var err error
	db.StartSession("hello")
	fmt.Println("hello")
	// defer db.EndSession(nil)
	defer func() {
		db.EndSession(err)
		fmt.Println("goodbye")
	}()

	hasError := rand.IntN(2) == 1
	if hasError {
		err = errors.New("ERR: random error")
	}
}
