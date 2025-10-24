package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(pathToDatabase string) (*Database, error) {
	os.MkdirAll(filepath.Dir(pathToDatabase), os.ModePerm)

	db, err := sql.Open("sqlite3", pathToDatabase)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	enablePragmas(db)
	createTables(db)

	return &Database{
		db,
	}, nil
}

func enablePragmas(db *sql.DB) {
	pragmas := []string{
		"PRAGMA foreign_keys = ON",
		"PRAGMA journal_mode = WAL",
		"PRAGMA synchronous = NORMAL",
		"PRAGMA temp_store = MEMORY",
	}
	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			log.Printf("pragma %s has error %s", pragma, err.Error())
		}
	}
}

func createTables(db *sql.DB) error {
	tables := []string{
		"CREATE TABLE IF NOT EXISTS souls (id INTEGER PRIMARY KEY, address VARCHAR(39) NOT NULL UNIQUE, seed VARCHAR(32), painted_pixels INTEGER NOT NULL DEFAULT 0)",
		"CREATE TABLE IF NOT EXISTS pixels (soul_id INTEGER NOT NULL REFERENCES souls(id), color INTEGER NOT NULL, x INT NOT NULL, y INT NOT NULL, PRIMARY KEY (x, y))",
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	for _, table := range tables {
		if _, err := tx.Exec(table); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (d Database) Close() {
	d.db.Close()
}
