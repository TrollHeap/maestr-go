package store

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

// InitDB initialise la connexion SQLite
func InitDB(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	// Optimisations SQLite
	db.Exec("PRAGMA journal_mode=WAL")
	db.Exec("PRAGMA synchronous=NORMAL")
	db.Exec("PRAGMA cache_size=-64000")
	db.Exec("PRAGMA foreign_keys=ON")

	// Charge et exécute le schema
	schema, err := os.ReadFile("internal/store/schema.sql")
	if err != nil {
		return fmt.Errorf("read schema: %w", err)
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("exec schema: %w", err)
	}

	return nil
}

// CloseDB ferme proprement la connexion
func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// GetDB expose la connexion DB (pour requêtes custom)
func GetDB() *sql.DB {
	return db
}
