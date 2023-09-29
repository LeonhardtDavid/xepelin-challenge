package config

import (
	"errors"
	"os"
	"strconv"
)

const (
	AppPort     = "APP_PORT"
	DefaultPort = 8080
	DatabaseUrl = "DATABASE_URL"
)

type Config struct {
	Port        int
	DatabaseUrl string
}

func LoadConfig() (*Config, error) {
	port, err := strconv.Atoi(os.Getenv(AppPort))
	if err != nil {
		port = DefaultPort
	}
	if port <= 0 {
		return nil, errors.New("port must be > 0")
	}

	databaseUrl := os.Getenv(DatabaseUrl)
	if databaseUrl == "" {
		return nil, errors.New("database url is required")
	}

	return &Config{
		Port:        port,
		DatabaseUrl: databaseUrl,
	}, nil
}
