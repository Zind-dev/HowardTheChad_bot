package config

import (
	"os"
	"testing"
)

func TestLoad_Success(t *testing.T) {
	// Set up environment variables
	os.Setenv("TELEGRAM_BOT_TOKEN", "test_token_123")
	os.Setenv("BOT_USERNAME", "test_bot")
	defer os.Unsetenv("TELEGRAM_BOT_TOKEN")
	defer os.Unsetenv("BOT_USERNAME")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.TelegramToken != "test_token_123" {
		t.Errorf("Expected token 'test_token_123', got '%s'", cfg.TelegramToken)
	}

	if cfg.BotUsername != "test_bot" {
		t.Errorf("Expected username 'test_bot', got '%s'", cfg.BotUsername)
	}

	// Check default values
	if cfg.ResponseFrequency != 10 {
		t.Errorf("Expected default ResponseFrequency 10, got %d", cfg.ResponseFrequency)
	}

	if !cfg.RespondToMentions {
		t.Error("Expected default RespondToMentions to be true")
	}
}

func TestLoad_MissingToken(t *testing.T) {
	// Ensure environment variables are not set
	os.Unsetenv("TELEGRAM_BOT_TOKEN")
	os.Unsetenv("BOT_USERNAME")

	cfg, err := Load()
	if err == nil {
		t.Fatal("Expected error for missing token, got nil")
	}
	if cfg != nil {
		t.Errorf("Expected nil config, got %+v", cfg)
	}
}

func TestLoad_MissingUsername(t *testing.T) {
	// Set only token
	os.Setenv("TELEGRAM_BOT_TOKEN", "test_token_123")
	os.Unsetenv("BOT_USERNAME")
	defer os.Unsetenv("TELEGRAM_BOT_TOKEN")

	cfg, err := Load()
	if err == nil {
		t.Fatal("Expected error for missing username, got nil")
	}
	if cfg != nil {
		t.Errorf("Expected nil config, got %+v", cfg)
	}
}

func TestLoad_CustomSettings(t *testing.T) {
	// Set up environment variables with custom values
	os.Setenv("TELEGRAM_BOT_TOKEN", "test_token_123")
	os.Setenv("BOT_USERNAME", "test_bot")
	os.Setenv("BOT_RESPONSE_FREQUENCY", "5")
	os.Setenv("BOT_RESPOND_TO_MENTIONS", "false")
	defer func() {
		os.Unsetenv("TELEGRAM_BOT_TOKEN")
		os.Unsetenv("BOT_USERNAME")
		os.Unsetenv("BOT_RESPONSE_FREQUENCY")
		os.Unsetenv("BOT_RESPOND_TO_MENTIONS")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.ResponseFrequency != 5 {
		t.Errorf("Expected ResponseFrequency 5, got %d", cfg.ResponseFrequency)
	}

	if cfg.RespondToMentions {
		t.Error("Expected RespondToMentions to be false")
	}
}

func TestLoad_InvalidFrequency(t *testing.T) {
	// Set up environment variables with invalid frequency
	os.Setenv("TELEGRAM_BOT_TOKEN", "test_token_123")
	os.Setenv("BOT_USERNAME", "test_bot")
	os.Setenv("BOT_RESPONSE_FREQUENCY", "invalid")
	defer func() {
		os.Unsetenv("TELEGRAM_BOT_TOKEN")
		os.Unsetenv("BOT_USERNAME")
		os.Unsetenv("BOT_RESPONSE_FREQUENCY")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should fall back to default
	if cfg.ResponseFrequency != 10 {
		t.Errorf("Expected default ResponseFrequency 10 for invalid value, got %d", cfg.ResponseFrequency)
	}
}
