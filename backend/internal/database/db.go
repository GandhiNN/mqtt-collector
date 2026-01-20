package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func New(databaseURL string) (*sql.DB, error) {
	var driverName string

	if strings.HasPrefix(databaseURL, "postgres://") ||
		strings.HasPrefix(databaseURL, "postgresql://") {
		driverName = "postgres"
	} else {
		driverName = "sqlite3"
	}

	db, err := sql.Open(driverName, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return db, nil
}
