package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadBrokersConfig(t *testing.T) {
	// Create temporary config file
	tmpfile, err := os.CreateTemp("", "brokers-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	content := `[
		{"id": "test-broker-1", "url": "tcp://localhost:1883"},
		{"id": "test-broker-2", "url": "tcp://localhost:1884"}
	]`

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test loading
	brokers, err := loadBrokersConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("loadBrokersConfig() error = %v", err)
	}

	if len(brokers) != 2 {
		t.Errorf("expected 2 brokers, got %d", len(brokers))
	}

	if brokers[0].ID != "test-broker-1" {
		t.Errorf("expected ID 'test-broker-1', got '%s'", brokers[0].ID)
	}

	if brokers[0].URL != "tcp://localhost:1883" {
		t.Errorf("expected URL 'tcp://localhost:1883', got '%s'", brokers[0].URL)
	}
}

func TestLoadBrokersConfigInvalidJSON(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "brokers-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	content := `[{"id", "broker-1", "url": invalid json}]`
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	_, err = loadBrokersConfig(tmpfile.Name())
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestLoadBrokersConfigNonExistentFile(t *testing.T) {
	_, err := loadBrokersConfig("/nonexistent/file.json")
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "env var set",
			key:          "TEST_VAR_1",
			defaultValue: "default",
			envValue:     "custom",
			expected:     "custom",
		},
		{
			name:         "env var not set",
			key:          "TEST_VAR_2",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnv(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getEnv() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	// Create temporary config file
	tmpfile, err := os.CreateTemp("", "brokers-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	content := `[{"id": "test-broker", "url": "tcp://localhost:1883"}]`
	tmpfile.Write([]byte(content))
	tmpfile.Close()

	// Set environment variables
	os.Setenv("BROKERS_CONFIG", tmpfile.Name())
	os.Setenv("DB_SERVICE_URL", "http://test:8080")
	os.Setenv("COLLECTION_DURATION", "2m")
	defer func() {
		os.Unsetenv("BROKERS_CONFIG")
		os.Unsetenv("DB_SERVICE_URL")
		os.Unsetenv("COLLECTION_DURATION")
	}()

	cfg, err := LoadCollectorConfig()
	if err != nil {
		t.Fatalf("Load () error = %v", err)
	}

	if cfg.DBServiceURL != "http://test:8080" {
		t.Errorf("expected DBServiceURL 'http://test:8080', got '%s'", cfg.DBServiceURL)
	}

	if cfg.CollectionDuration != 2*time.Minute {
		t.Errorf("expected duration 2m, got %v", cfg.CollectionDuration)
	}

	if len(cfg.Brokers) != 1 {
		t.Errorf("expected 1 broker, got %d", len(cfg.Brokers))
	}
}
