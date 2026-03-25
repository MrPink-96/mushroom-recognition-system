package config

import "os"

type Config struct {
	Port        string
	StoragePath string
}

func Load() *Config {
	port := getEnv("PORT", "8082")
	storagePath := getEnv("STORAGE_PATH", "./images")

	return &Config{
		Port:        port,
		StoragePath: storagePath,
	}
}

func getEnv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}
