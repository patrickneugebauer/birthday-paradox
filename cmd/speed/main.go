package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

const iterations = 999

// 1. Individual Inserts (No Transaction)
func insertIndividual(db *sql.DB) {
	for i := 0; i < iterations; i++ {
		_, err := db.Exec("INSERT INTO test (name) VALUES (?)", fmt.Sprintf("%d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
}

// 2. Grouped Transaction
func insertTransaction(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < iterations; i++ {
		_, err := tx.Exec("INSERT INTO test (name) VALUES (?)", fmt.Sprintf("%d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()
}

// 3. Single Bulk SQL String
func insertBulkString(db *sql.DB) {
	valueStrings := make([]string, 0, iterations)
	valueArgs := make([]interface{}, 0, iterations)
	for i := 0; i < iterations; i++ {
		valueStrings = append(valueStrings, "(?)")
		valueArgs = append(valueArgs, fmt.Sprintf("%d", i))
	}
	stmt := fmt.Sprintf("INSERT INTO test (name) VALUES %s", strings.Join(valueStrings, ","))
	_, err := db.Exec(stmt, valueArgs...)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	db, _ := sql.Open("sqlite3", "test.db")
	defer db.Close()

	db.Exec("DROP TABLE IF EXISTS test")
	db.Exec("CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT NOT NULL)")

	// Test 1: Individual
	start := time.Now()
	insertIndividual(db)
	baseSpeed := time.Since(start)
	fmt.Printf("Individual Inserts: %v\n", baseSpeed)
	db.Exec("DELETE FROM test")

	// Test 2: Transaction
	start = time.Now()
	insertTransaction(db)
	speed := time.Since(start)
	fmt.Printf("Transaction Loop:   %v\n", speed)
	speedup := baseSpeed.Seconds() / speed.Seconds()
	fmt.Printf("speedup:   %.2f\n", speedup)
	db.Exec("DELETE FROM test")

	// Test 3: Bulk String
	start = time.Now()
	insertBulkString(db)
	speed = time.Since(start)
	fmt.Printf("Single Bulk String: %v\n", speed)
	speedup = baseSpeed.Seconds() / speed.Seconds()
	fmt.Printf("speedup:   %.2f\n", speedup)
	db.Exec("DELETE FROM test")
}
