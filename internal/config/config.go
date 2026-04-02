package config

import "os"

// This config module loads environment-based settings with defaults
// Allows the application to run across different environments without code changes
type Config struct {
	AWSRegion     string
	AdvisorsTable string
	SessionsTable string
}

// Get config from environment variables
// Default if not found
func Load() Config {
	return Config{
		AWSRegion:     getEnv("AWS_REGION", "us-east-1"),
		AdvisorsTable: getEnv("ADVISORS_TABLE", "advisors"),
		SessionsTable: getEnv("SESSIONS_TABLE", "sessions"),
	}
}

// If env variable not set, fallback to defaults
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
