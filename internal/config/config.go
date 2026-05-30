package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Env         string
	UpstreamURL string
	DBPath      string
}

func Load() (*Config, error) {

	_ = godotenv.Load()

	return &Config{
		Env:         checkEnv("ENV"),
		Port:        checkEnv("PROXY_PORT"),
		DBPath:      checkEnv("PROXY_DB"),
		UpstreamURL: checkEnv("KOITO_URL"),
	}, nil
}

func checkEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic("enviroment variable '" + key + "' is not set")
	}
	return val
}
