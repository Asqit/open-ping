package storage

import "database/sql"

type Storage interface {
	SavePing(target, url string, status, latency int, success bool) error
	Close() error
	DB() *sql.DB
}
