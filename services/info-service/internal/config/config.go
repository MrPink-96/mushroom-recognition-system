package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port string
	DSN  string
}

func Load() Config {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dbUser := os.Getenv("DB_USER")
		dbPass := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")

		dsn = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			dbUser, dbPass, dbHost, dbPort, dbName,
		)
	}

	return Config{
		Port: port,
		DSN:  dsn,
	}
}
