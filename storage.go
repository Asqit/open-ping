package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init_db() {
	var err error
	db, err = sql.Open("sqlite3", "openping.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS pings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		target TEXT NOT NULL,
		status INTEGER NOT NULL,
		success BOOLEAN NOT NULL,
		latency INTEGER NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

func save_log(target string, status int, success bool, latency int) {
	fmt.Println(target, status, success, latency)
	if db == nil {
		log.Fatal("Database not initialized. Call InitDB() first.")
	}

	_, err := db.Exec(`
		INSERT INTO pings (target, status, success, latency)
		VALUES (?, ?, ?, ?)`,
		target, status, success, latency,
	)
	if err != nil {
		log.Printf("Failed to save log: %v", err)
	}
}

func close_db() {
	if db != nil {
		db.Close()
	}
}
