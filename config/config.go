package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	TelegramToken     string
	BotUsername       string
	ResponseFrequency int  // How often to respond to regular messages (e.g., every 10th message)
	RespondToMentions bool // Whether to always respond to mentions
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

	// Load response frequency (default: 10)
	frequency := 10
	if freqStr := os.Getenv("BOT_RESPONSE_FREQUENCY"); freqStr != "" {
		if f, err := strconv.Atoi(freqStr); err == nil {
			frequency = f
		}
	}

	// Load respond to mentions setting (default: true)
	respondToMentions := true
	if mentionsStr := os.Getenv("BOT_RESPOND_TO_MENTIONS"); mentionsStr != "" {
		if mentionsStr == "false" || mentionsStr == "0" {
			respondToMentions = false
		}
	}

	return &Config{
		TelegramToken:     token,
		BotUsername:       username,
		ResponseFrequency: frequency,
		RespondToMentions: respondToMentions,
	}, nil
}
