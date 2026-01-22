// Utility functions for reading environment variables with default fallbacks
package config

import "os"

// Returns environment variable value or default if not set
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
