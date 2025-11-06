package main

import (
	"log"

	"github.com/Zind-dev/HowardTheChad_bot/bot"
	"github.com/Zind-dev/HowardTheChad_bot/config"
	"github.com/Zind-dev/HowardTheChad_bot/storage"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize storage
	store, err := storage.NewSQLiteStorage("bot_data.db")
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err)
	}
	defer store.Close()

	if err := store.Initialize(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database initialized successfully")

	// Create bot instance
	b, err := bot.New(cfg, store)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Start the bot
	log.Println("Bot is starting...")
	if err := b.Start(); err != nil {
		log.Fatalf("Bot error: %v", err)
	}
}
