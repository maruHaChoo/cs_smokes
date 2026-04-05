package config

import (
    "fmt"
    "os"
)

type Config struct {
    BotToken string
}

func Load() (Config, error) {
    cfg := Config{BotToken: os.Getenv("BOT_TOKEN")}
    if cfg.BotToken == "" {
        return Config{}, fmt.Errorf("BOT_TOKEN is required")
    }
    return cfg, nil
}
