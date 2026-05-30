package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DBPath      string
	UpstreamURL string
}

func Load() (*Config, error) {

	_ = godotenv.Load()

	return &Config{
		Port:        checkEnv("PROXY_PORT", "4112"),
		DBPath:      checkEnv("PROXY_DB", "/app/data/koito.db"),
		UpstreamURL: checkEnv("KOITO_URL", "http://localhost:4110"),
	}, nil
}

func checkEnv(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		slog.Warn("enviroment variable '" + key + "' is not set, using default value")
		return defaultValue
	}
	return val
}
