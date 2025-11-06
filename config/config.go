package config

import (
	"fmt"
	"os"
)

// Config holds the application configuration
type Config struct {
	TelegramToken string
	BotUsername   string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	username := os.Getenv("BOT_USERNAME")
	if username == "" {
		return nil, fmt.Errorf("BOT_USERNAME environment variable is required")
	}

	return &Config{
		TelegramToken: token,
		BotUsername:   username,
	}, nil
}
