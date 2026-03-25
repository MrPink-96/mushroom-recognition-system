package config

import (
	"os"
)

type Config struct {
	Port     string
	MLURL    string
	InfoURL  string
	ImageURL string
}

func Load() *Config {
	return &Config{
		Port:     getEnv("PORT", ":8080"),
		MLURL:    getEnv("ML_SERVICE_URL", "http://localhost:8000"),
		InfoURL:  getEnv("INFO_SERVICE_URL", "http://localhost:8081"),
		ImageURL: getEnv("IMAGE_SERVICE_URL", "http://localhost:8082"),
	}
}

func getEnv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}
