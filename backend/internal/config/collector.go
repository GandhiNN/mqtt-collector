// Loads collector configuration including broker connections,
// database URLs, and collection intervals
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type BrokerConfig struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type CollectorConfig struct {
	Brokers            []BrokerConfig
	DBServiceURL       string
	CollectionDuration time.Duration
}

// Reads broker config from JSON file and environment variables with fallback defaults
func LoadCollectorConfig() (*CollectorConfig, error) {
	configPath := getEnv("BROKERS_CONFIG", "brokers.json")
	dbServiceURL := getEnv("DB_SERVICE_URL", "http://localhost:8080")
	durationStr := getEnv("COLLECTION_DURATION", "1m")

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return nil, fmt.Errorf("invalid duration: %w", err)
	}

	brokers, err := loadBrokersConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("load brokers config: %w", err)
	}

	return &CollectorConfig{
		Brokers:            brokers,
		DBServiceURL:       dbServiceURL,
		CollectionDuration: duration,
	}, nil
}

// Parses JSON configuration file containing MQTT broker connection details
func loadBrokersConfig(configPath string) ([]BrokerConfig, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var brokers []BrokerConfig
	if err := json.Unmarshal(file, &brokers); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return brokers, nil
}
