package utils

import "os"

// getEnv returns the value of the environment variable identified by key.
// If the variable is not set or is empty, it returns the provided fallback value.
func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
