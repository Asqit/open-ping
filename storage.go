package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func Save_log() {
	db, err := sql.Open("sqlite", "openping.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// _, err := db.Exec("SELECT * FROM users")
}
