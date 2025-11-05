package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLite(dbPath string) (*SQLiteStorage, error) {
	fmt.Print("initiating database............")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS pings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		target TEXT NOT NULL,
		status INTEGER NOT NULL,
		success BOOLEAN NOT NULL,
		latency INTEGER NOT NULL,
		target_website TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	fmt.Println("[OK]")
	return &SQLiteStorage{db: db}, nil
}

func (s *SQLiteStorage) SavePing(target, url string, status, latency int, success bool) error {
	_, err := s.db.Exec(`
		INSERT INTO pings (target, target_website, status, success, latency)
		VALUES (?, ?, ?, ?, ?)`,
		target, url, status, success, latency,
	)
	if err != nil {
		log.Printf("Failed to save log: %v", err)
	}
	return err
}

func (s *SQLiteStorage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *SQLiteStorage) DB() *sql.DB {
	return s.db
}
