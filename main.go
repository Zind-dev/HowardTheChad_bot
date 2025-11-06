package main

import (
	"log"

	"github.com/Zind-dev/HowardTheChad_bot/bot"
	"github.com/Zind-dev/HowardTheChad_bot/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create bot instance
	b, err := bot.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Start the bot
	log.Println("Bot is starting...")
	if err := b.Start(); err != nil {
		log.Fatalf("Bot error: %v", err)
	}
}
