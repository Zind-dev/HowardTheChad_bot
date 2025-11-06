package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Zind-dev/HowardTheChad_bot/chats"
	"github.com/Zind-dev/HowardTheChad_bot/config"
	"github.com/Zind-dev/HowardTheChad_bot/settings"
	"github.com/Zind-dev/HowardTheChad_bot/storage"
	"github.com/Zind-dev/HowardTheChad_bot/users"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot represents the Telegram bot
type Bot struct {
	api             *tgbotapi.BotAPI
	config          *config.Config
	userManager     *users.Manager
	chatManager     *chats.Manager
	settingsManager *settings.Manager
	storage         storage.Storage
}

// New creates a new bot instance
func New(cfg *config.Config, store storage.Storage) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return nil, err
	}

	log.Printf("Authorized on account %s", api.Self.UserName)
	log.Printf("Configured bot username: %s", cfg.BotUsername)
	log.Printf("NOTE: These usernames must match for mentions to work!")

	// Create settings manager with defaults from config
	defaultSettings := settings.NewCustomSettings(cfg.ResponseFrequency, cfg.RespondToMentions)
	settingsMgr := settings.NewManager(defaultSettings)

	return &Bot{
		api:             api,
		config:          cfg,
		userManager:     users.NewManager(),
		chatManager:     chats.NewManager(),
		settingsManager: settingsMgr,
		storage:         store,
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
	log.Printf("Message received - Chat: %d, User: %s, Text: %s",
		message.Chat.ID, message.From.UserName, message.Text)

	// Save user to storage
	user := &storage.User{
		ID:           message.From.ID,
		UserName:     message.From.UserName,
		FirstName:    message.From.FirstName,
		LastName:     message.From.LastName,
		MessageCount: 0, // Will be updated separately
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := b.storage.SaveUser(user); err != nil {
		log.Printf("Warning: Failed to save user: %v", err)
	}

	// Save chat to storage
	chat := &storage.Chat{
		ID:           message.Chat.ID,
		Title:        message.Chat.Title,
		Type:         message.Chat.Type,
		MessageCount: 0, // Will be updated separately
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := b.storage.SaveChat(chat); err != nil {
		log.Printf("Warning: Failed to save chat: %v", err)
	}

	// Save message to storage for AI context
	msg := &storage.Message{
		ChatID:    message.Chat.ID,
		UserID:    message.From.ID,
		Text:      message.Text,
		IsBot:     message.From.IsBot,
		Timestamp: time.Now(),
	}
	if err := b.storage.SaveMessage(msg); err != nil {
		log.Printf("Warning: Failed to save message: %v", err)
	}

	// Check if the message is in a group
	if !message.Chat.IsGroup() && !message.Chat.IsSuperGroup() {
		// Respond to private messages
		b.respondToPrivateMessage(message)
		return
	}

	// Handle commands first
	if message.IsCommand() {
		b.handleCommand(message)
		return
	}

	// Track message count for this chat
	messageCount := b.chatManager.IncrementMessageCount(
		message.Chat.ID,
		message.Chat.Title,
		message.Chat.Type,
	)

	// Get settings for this specific chat
	chatSettings := b.settingsManager.GetSettings(message.Chat.ID)

	// Check if bot is mentioned
	isMentioned := b.isBotMentioned(message)

	// Determine if bot should respond
	shouldRespond := false
	if isMentioned && chatSettings.ShouldRespondToMention() {
		shouldRespond = true
	} else if !isMentioned && chatSettings.ShouldRespondToRegularMessage(messageCount) {
		shouldRespond = true
	}

	if shouldRespond {
		if isMentioned {
			b.respondToMention(message)
		} else {
			b.respondToRegularMessage(message)
		}
	}
}

// isBotMentioned checks if the bot is mentioned in the message
func (b *Bot) isBotMentioned(message *tgbotapi.Message) bool {
	botUsername := b.config.BotUsername

	log.Printf("Checking mention - Bot username: %s, Message: %s", botUsername, message.Text)

	// Check for @ mentions (with or without @)
	botMention := "@" + botUsername
	if strings.Contains(message.Text, botMention) {
		log.Printf("Detected mention via text: %s", botMention)
		return true
	}

	// Check entities for mentions
	for _, entity := range message.Entities {
		if entity.Type == "mention" {
			mention := message.Text[entity.Offset : entity.Offset+entity.Length]
			log.Printf("Found mention entity: %s", mention)
			if mention == botMention {
				return true
			}
		}
		// Also check text_mention type (when user doesn't have username)
		if entity.Type == "text_mention" {
			if entity.User != nil && entity.User.UserName == botUsername {
				log.Printf("Detected text_mention for bot")
				return true
			}
		}
	}

	// Check if the message is a reply to the bot
	if message.ReplyToMessage != nil {
		if message.ReplyToMessage.From != nil {
			if message.ReplyToMessage.From.UserName == botUsername {
				log.Printf("Detected reply to bot message")
				return true
			}
		}
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
	} else {
		b.saveResponseMessage(message.Chat.ID, response)
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
	} else {
		b.saveResponseMessage(message.Chat.ID, response)
	}
}

// respondToRegularMessage handles regular messages in group chats (periodic responses)
func (b *Bot) respondToRegularMessage(message *tgbotapi.Message) {
	userInfo := b.userManager.GetUser(message.From.ID)

	// Build context-aware response
	response := b.generateResponse(message, userInfo)

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	msg.ReplyToMessageID = message.MessageID

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	} else {
		b.saveResponseMessage(message.Chat.ID, response)
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

// saveResponseMessage saves a bot response message to storage
func (b *Bot) saveResponseMessage(chatID int64, text string) {
	msg := &storage.Message{
		ChatID:    chatID,
		UserID:    b.api.Self.ID,
		Text:      text,
		IsBot:     true,
		Timestamp: time.Now(),
	}
	if err := b.storage.SaveMessage(msg); err != nil {
		log.Printf("Warning: Failed to save bot response: %v", err)
	}
}

// GetUserInfo retrieves information about a user
func (b *Bot) GetUserInfo(userID int64) *users.User {
	return b.userManager.GetUser(userID)
}

// GetAllUsers returns all stored users
func (b *Bot) GetAllUsers() map[int64]*users.User {
	return b.userManager.GetAllUsers()
}

// GetChatInfo retrieves information about a chat
func (b *Bot) GetChatInfo(chatID int64) *chats.Chat {
	return b.chatManager.GetChat(chatID)
}

// GetAllChats returns all tracked chats
func (b *Bot) GetAllChats() map[int64]*chats.Chat {
	return b.chatManager.GetAllChats()
}

// UpdateSettings updates the bot's behavior settings for a specific chat
func (b *Bot) UpdateSettings(chatID int64, newSettings *settings.Settings) {
	b.settingsManager.SetSettings(chatID, newSettings)
}

// GetSettings returns the current bot settings for a specific chat
func (b *Bot) GetSettings(chatID int64) *settings.Settings {
	return b.settingsManager.GetSettings(chatID)
}

// isUserAdmin checks if a user is an administrator in a chat
func (b *Bot) isUserAdmin(chatID int64, userID int64) bool {
	chatConfig := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chatID,
			UserID: userID,
		},
	}

	member, err := b.api.GetChatMember(chatConfig)
	if err != nil {
		log.Printf("Error getting chat member: %v", err)
		return false
	}

	// Check if user is creator or administrator
	return member.Status == "creator" || member.Status == "administrator"
}

// handleCommand processes bot commands
func (b *Bot) handleCommand(message *tgbotapi.Message) {
	command := message.Command()
	log.Printf("Received command: /%s from user %d in chat %d", command, message.From.ID, message.Chat.ID)

	// Handle commands in private chats
	if !message.Chat.IsGroup() && !message.Chat.IsSuperGroup() {
		if command == "help" || command == "start" {
			b.handleHelpCommand(message)
		} else {
			b.respondToPrivateMessage(message)
		}
		return
	}

	// Handle commands in group chats
	switch command {
	case "settings":
		b.handleSettingsCommand(message)
	case "setfrequency":
		b.handleSetFrequencyCommand(message)
	case "togglementions":
		b.handleToggleMentionsCommand(message)
	case "resetsettings":
		b.handleResetSettingsCommand(message)
	case "help", "start":
		b.handleHelpCommand(message)
	default:
		log.Printf("Unknown command: /%s", command)
	}
} // handleSettingsCommand shows current settings for the chat
func (b *Bot) handleSettingsCommand(message *tgbotapi.Message) {
	chatSettings := b.settingsManager.GetSettings(message.Chat.ID)

	mentionsStatus := "enabled"
	if !chatSettings.AlwaysRespondToMentions {
		mentionsStatus = "disabled"
	}

	response := "📊 Current Settings:\n\n"
	response += "• Response Frequency: every " + formatFrequency(chatSettings.ResponseFrequency) + "\n"
	response += "• Respond to Mentions: " + mentionsStatus + "\n\n"
	response += "Use /help to see available commands."

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	msg.ReplyToMessageID = message.MessageID

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

// handleSetFrequencyCommand changes the response frequency
func (b *Bot) handleSetFrequencyCommand(message *tgbotapi.Message) {
	// Check if user is admin
	if !b.isUserAdmin(message.Chat.ID, message.From.ID) {
		b.sendMessage(message.Chat.ID, "❌ Only administrators can change settings.", message.MessageID)
		return
	}

	args := message.CommandArguments()
	if args == "" {
		b.sendMessage(message.Chat.ID, "Usage: /setfrequency <number>\nExample: /setfrequency 10", message.MessageID)
		return
	}

	frequency, err := strconv.Atoi(strings.TrimSpace(args))
	if err != nil || frequency < 0 {
		b.sendMessage(message.Chat.ID, "❌ Please provide a valid number (0 or greater).", message.MessageID)
		return
	}

	b.settingsManager.SetFrequency(message.Chat.ID, frequency)

	response := "✅ Response frequency updated to: every " + formatFrequency(frequency)
	b.sendMessage(message.Chat.ID, response, message.MessageID)
}

// handleToggleMentionsCommand toggles mention response setting
func (b *Bot) handleToggleMentionsCommand(message *tgbotapi.Message) {
	// Check if user is admin
	if !b.isUserAdmin(message.Chat.ID, message.From.ID) {
		b.sendMessage(message.Chat.ID, "❌ Only administrators can change settings.", message.MessageID)
		return
	}

	newValue := b.settingsManager.ToggleMentionResponse(message.Chat.ID)

	status := "enabled"
	if !newValue {
		status = "disabled"
	}

	response := "✅ Respond to mentions: " + status
	b.sendMessage(message.Chat.ID, response, message.MessageID)
}

// handleResetSettingsCommand resets settings to defaults
func (b *Bot) handleResetSettingsCommand(message *tgbotapi.Message) {
	// Check if user is admin
	if !b.isUserAdmin(message.Chat.ID, message.From.ID) {
		b.sendMessage(message.Chat.ID, "❌ Only administrators can change settings.", message.MessageID)
		return
	}

	b.settingsManager.ResetSettings(message.Chat.ID)
	b.sendMessage(message.Chat.ID, "✅ Settings reset to defaults.", message.MessageID)
}

// handleHelpCommand shows help information
func (b *Bot) handleHelpCommand(message *tgbotapi.Message) {
	response := "🤖 HowardTheChad Bot Commands\n\n"
	response += "📊 Information:\n"
	response += "/settings - Show current settings\n"
	response += "/help - Show this help message\n\n"
	response += "⚙️ Admin Commands:\n"
	response += "/setfrequency <number> - Set response frequency\n"
	response += "  Example: /setfrequency 10 (respond every 10th message)\n"
	response += "  Use 0 to only respond to mentions\n"
	response += "/togglementions - Toggle automatic response to mentions\n"
	response += "/resetsettings - Reset settings to defaults\n"

	b.sendMessage(message.Chat.ID, response, message.MessageID)
}

// sendMessage is a helper to send messages
func (b *Bot) sendMessage(chatID int64, text string, replyToMessageID int) {
	msg := tgbotapi.NewMessage(chatID, text)
	if replyToMessageID != 0 {
		msg.ReplyToMessageID = replyToMessageID
	}

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

// formatFrequency formats the frequency number for display
func formatFrequency(frequency int) string {
	if frequency == 0 {
		return "never (mentions only)"
	} else if frequency == 1 {
		return "message"
	} else {
		return fmt.Sprintf("%d messages", frequency)
	}
}
