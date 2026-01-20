package config

import (
	"fmt"
)

type ServerConfig struct {
	ServerAddr  string
	DatabaseURL string
}

func LoadServerConfig() (*ServerConfig, error) {
	serverAddr := getEnv("SERVER_ADDR", ":8080")
	databaseURL := getEnv("DATABASE_URL", "")

	if databaseURL == "" {
		// Default to SQLite if no DATABASE_URL provided
		databaseURL = getEnv("SQLITE_PATH", "mqtt_catalog.db")
		databaseURL = fmt.Sprintf("file:%s?cache=shared&mode=rwc", databaseURL)
	}

	return &ServerConfig{
		ServerAddr:  serverAddr,
		DatabaseURL: databaseURL,
	}, nil
}
