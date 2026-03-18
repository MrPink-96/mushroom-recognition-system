package config

import (
	"os"
)

type Config struct {
	Port string
	DSN  string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	//	dsn := os.Getenv("DATABASE_URL")

	return Config{
		Port: port,
		DSN:  dsn,
	}
}
