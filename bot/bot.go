package bot

import (
	"log"
	"strings"

	"github.com/Zind-dev/HowardTheChad_bot/config"
	"github.com/Zind-dev/HowardTheChad_bot/users"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot represents the Telegram bot
type Bot struct {
	api         *tgbotapi.BotAPI
	config      *config.Config
	userManager *users.Manager
}

// New creates a new bot instance
func New(cfg *config.Config) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return nil, err
	}

	log.Printf("Authorized on account %s", api.Self.UserName)

	return &Bot{
		api:         api,
		config:      cfg,
		userManager: users.NewManager(),
	}, nil
}

// Start starts the bot and handles incoming messages
func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Store user information
		b.userManager.UpdateUser(update.Message.From)

		// Handle the message
		b.handleMessage(update.Message)
	}

	return nil
}

// handleMessage processes incoming messages
func (b *Bot) handleMessage(message *tgbotapi.Message) {
	// Check if the message is in a group
	if !message.Chat.IsGroup() && !message.Chat.IsSuperGroup() {
		// Respond to private messages
		b.respondToPrivateMessage(message)
		return
	}

	// In group chats, check if bot is mentioned
	if b.isBotMentioned(message) {
		b.respondToMention(message)
	}
}

// isBotMentioned checks if the bot is mentioned in the message
func (b *Bot) isBotMentioned(message *tgbotapi.Message) bool {
	// Check for @ mentions
	botMention := "@" + b.config.BotUsername
	if strings.Contains(message.Text, botMention) {
		return true
	}

	// Check entities for mentions
	for _, entity := range message.Entities {
		if entity.Type == "mention" {
			mention := message.Text[entity.Offset : entity.Offset+entity.Length]
			if mention == botMention {
				return true
			}
		}
	}

	// Check if the message is a reply to the bot
	if message.ReplyToMessage != nil && message.ReplyToMessage.From.UserName == b.config.BotUsername {
		return true
	}

	return false
}

// respondToPrivateMessage handles private messages
func (b *Bot) respondToPrivateMessage(message *tgbotapi.Message) {
	userInfo := b.userManager.GetUser(message.From.ID)

	response := "Hello! I'm HowardTheChad bot. "
	if userInfo != nil {
		response += "I can see you, " + userInfo.FirstName + "! "
	}
	response += "Add me to a group and mention me with @ to chat!"

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	msg.ReplyToMessageID = message.MessageID

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

// respondToMention handles mentions in group chats
func (b *Bot) respondToMention(message *tgbotapi.Message) {
	userInfo := b.userManager.GetUser(message.From.ID)

	// Build context-aware response
	response := b.generateResponse(message, userInfo)

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	msg.ReplyToMessageID = message.MessageID

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

// generateResponse generates a context-aware response
func (b *Bot) generateResponse(message *tgbotapi.Message, userInfo *users.User) string {
	// For now, return a simple response
	// This is where AI model integration would happen in the future

	userName := "there"
	if userInfo != nil && userInfo.FirstName != "" {
		userName = userInfo.FirstName
	}

	responses := []string{
		"Hey " + userName + "! What's up?",
		"Hello " + userName + "! I'm here to help.",
		"Hi " + userName + "! What can I do for you?",
		userName + ", I'm listening!",
		"Yo " + userName + "! How can I contribute?",
	}

	// Simple selection based on message length
	// In the future, this would be replaced with AI model
	index := len(message.Text) % len(responses)
	return responses[index]
}

// GetUserInfo retrieves information about a user
func (b *Bot) GetUserInfo(userID int64) *users.User {
	return b.userManager.GetUser(userID)
}

// GetAllUsers returns all stored users
func (b *Bot) GetAllUsers() map[int64]*users.User {
	return b.userManager.GetAllUsers()
}
