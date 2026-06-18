package config

import (
	"log/slog"
	"os"
)

type Config struct {
	Port              string
	DBPath            string
	UpstreamURL       string
	RateLimitEnabled  string // e.g., "true" or "false"
	RequestsPerSecond string
	Burst             string
}

func Load() (*Config, error) {
	return &Config{
		Port:              checkEnv("PROXY_PORT", "4112"),
		DBPath:            checkEnv("PROXY_DB", "./koito_proxy.db"),
		UpstreamURL:       checkEnv("KOITO_URL", "http://localhost:4110"),
		RateLimitEnabled:  checkEnv("RATE_LIMIT_ENABLED", "false"),
		RequestsPerSecond: checkEnv("REQUESTS_PER_SECOND", "10"),
		Burst:             checkEnv("RATE_BURST", "20"),
	}, nil
}

func checkEnv(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		slog.Info("environment variable '" + key + "' is not set, using default value")
		return defaultValue
	}
	return val
}
