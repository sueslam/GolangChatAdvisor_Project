package config

import "os"

type Config struct {
	AWSRegion     string
	AdvisorsTable string
	SessionsTable string
}

func Load() Config {
	return Config{
		AWSRegion:     getEnv("AWS_REGION", "us-east-1"),
		AdvisorsTable: getEnv("ADVISORS_TABLE", "advisors"),
		SessionsTable: getEnv("SESSIONS_TABLE", "sessions"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
