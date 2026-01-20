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
