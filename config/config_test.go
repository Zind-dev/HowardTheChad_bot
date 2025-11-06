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
