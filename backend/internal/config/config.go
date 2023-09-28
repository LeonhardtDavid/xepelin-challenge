package config

import (
	"errors"
	"os"
	"strconv"
)

const (
	AppPort     = "APP_PORT"
	DefaultPort = 8080
)

type Config struct {
	Port int
}

func LoadConfig() (*Config, error) {
	port, err := strconv.Atoi(os.Getenv(AppPort))
	if err != nil {
		port = DefaultPort
	}

	if port <= 0 {
		return nil, errors.New("port must be > 0")
	}

	return &Config{
		Port: port,
	}, nil
}
