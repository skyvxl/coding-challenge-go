package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	ServerAddr  string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("load .env: %w", err)
	}
	cfg := &Config{
		DatabaseURL: strings.TrimSpace(os.Getenv("DATABASE_URL")),
		ServerAddr:  strings.TrimSpace(os.Getenv("SERVER_ADDR")),
	}
	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	if cfg.ServerAddr == "" {
		return nil, errors.New("SERVER_ADDR is required")
	}
	return cfg, nil
}
