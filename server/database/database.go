package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(database string) (*Database, error) {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// create tables
	tables := []string{
		"CREATE TABLE IF NOT EXISTS fishes (seed VARCHAR(30), address VARCHAR(39), spawned_at DATETIME DEFAULT current_timestamp)",
	}
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	for _, table := range tables {
		if _, err := tx.Exec(table); err != nil {
			return nil, err
		}
	}
	tx.Commit()

	return &Database{
		db,
	}, nil
}

func (d Database) Close() {
	d.db.Close()
}
