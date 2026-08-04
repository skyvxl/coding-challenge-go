package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	ServerAddr  string
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load("./.env"); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %v", err)
	}
	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		ServerAddr:  os.Getenv("SERVER_ADDR"),
	}, nil
}
