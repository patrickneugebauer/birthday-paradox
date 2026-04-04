package db

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	// The underscore import registers the driver with database/sql
	_ "github.com/mattn/go-sqlite3"
)

const dbname = "./app.db"

var db *sql.DB
var session_id int

func get_timestamp() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

func StartSession(session_type string) *sql.DB {
	// database setup
	var err error
	db, err = sql.Open("sqlite3", dbname)
	if err != nil {
		log.Fatal(err)
	}

	// create session
	timestamp := get_timestamp()
	fmt.Println(timestamp)
	query := "INSERT INTO sessions(session_type, started_at) VALUES (?, ?) RETURNING id"
	err = db.QueryRow(query, session_type, timestamp).Scan(&session_id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("session created: " + strconv.Itoa(session_id))

	// return
	return db
}

func EndSession(err error) {
	var exitCode int
	var errorMessage string
	if err == nil {
		exitCode = 0
	} else {
		exitCode = 1
		errorMessage = err.Error()
	}
	timestamp := get_timestamp()
	query := `UPDATE sessions
              SET finished_at = ?, exit_code = ?, error = ?
              WHERE id = ?`

	_, execErr := db.Exec(query, timestamp, exitCode, errorMessage, session_id)
	if execErr != nil {
		fmt.Printf("Failed to log session end: %v\n", execErr)
	}
	db.Close()
}
