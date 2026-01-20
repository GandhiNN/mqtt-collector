package database

import (
	"database/sql"
	"fmt"
)

func RunMigrations(db *sql.DB) error {
	// Check if we are using SQLite or Postgres
	var query string

	// Detect database type by trying to get version
	var version string
	err := db.QueryRow("SELECT sqlite_version()").Scan(&version)
	isSQLite := err == nil

	if isSQLite {
		query = `
		CREATE TABLE IF NOT EXISTS topics (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			broker_id TEXT NOT NULL,
			payload_type TEXT NOT NULL,
			sample_payload BLOB NOT NULL,
			last_seen TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(broker_id, topic));
			
		CREATE INDEX IF NOT EXISTS idx_topics_broker_id ON topics(broker_id);
		CREATE INDEX IF NOT EXISTS idx_topics_last_seen ON topics(last_seen DESC);
		CREATE INDEX IF NOT EXISTS idx_topics_topic ON topics(topic);
		`
	} else {
		query = `
		CREATE TABLE IF NOT EXISTS topics (
			id SERIAL PRIMARY KEY,
			broker_id TEXT NOT NULL,
			topic TEXT NOT NULL,
			payload_type TEXT NOT NULL,
			sample_payload BYTEA NOT NULL,
			last_seen TIMESTAMP NOT NULL DEFAULT NOW(),
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			UNIQUE(broker_id, topic)
		);

		CREATE INDEX IF NOT EXISTS idx_topics_broker_id ON topics(broker_id);
		CREATE INDEX IF NOT EXISTS idx_topics_last_seen ON topics(last_seen DESC);
		CREATE INDEX IF NOT EXISTS idx_topics_topic ON topics(topic);
		`
	}

	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
