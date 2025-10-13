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
	// TODO: enable PRAGMAs. refactor var tables to INIT exec
	tables := []string{
		"CREATE TABLE IF NOT EXISTS souls (id INTEGER PRIMARY KEY, address VARCHAR(39) NOT NULL UNIQUE, seed VARCHAR(32), painted_pixels INTEGER NOT NULL DEFAULT 0)",
		"CREATE TABLE IF NOT EXISTS pixels (soul_id INTEGER NOT NULL REFERENCES souls(id), color TEXT NOT NULL, x INT NOT NULL, y INT NOT NULL, PRIMARY KEY (x, y))",
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
