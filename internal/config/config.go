package config

import (
	"fmt"
	"os"
)

type Config struct {
	BotToken    string
	DatabaseURL string
}

func Load() (Config, error) {
	cfg := Config{
		BotToken:    os.Getenv("BOT_TOKEN"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}

	if cfg.BotToken == "" {
		return Config{}, fmt.Errorf("BOT_TOKEN is required")
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}
