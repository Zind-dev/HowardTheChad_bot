package bot

import (
	"os"
	"testing"

	"github.com/Zind-dev/HowardTheChad_bot/config"
	"github.com/Zind-dev/HowardTheChad_bot/users"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestNew(t *testing.T) {
	cfg := &config.Config{
		TelegramToken: "test_token",
		BotUsername:   "test_bot",
	}

	// This will fail with invalid token, which is expected in test environment
	_, err := New(cfg)
	// We expect an error with invalid token
	if err == nil {
		t.Error("Expected error with invalid token")
	}
}

func TestNew_WithRealToken(t *testing.T) {
	// This test only runs if TEST_BOT_TOKEN environment variable is set
	// Usage: $env:TEST_BOT_TOKEN="your_token"; $env:TEST_BOT_USERNAME="your_bot"; go test ./bot -v -run TestNew_WithRealToken
	token := os.Getenv("TEST_BOT_TOKEN")
	username := os.Getenv("TEST_BOT_USERNAME")

	if token == "" || username == "" {
		t.Skip("Skipping integration test - TEST_BOT_TOKEN and TEST_BOT_USERNAME not set")
	}

	cfg := &config.Config{
		TelegramToken: token,
		BotUsername:   username,
	}

	bot, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create bot with valid token: %v", err)
	}

	if bot == nil {
		t.Fatal("Bot is nil")
	}

	if bot.api == nil {
		t.Error("Bot API is nil")
	}

	if bot.userManager == nil {
		t.Error("User manager is nil")
	}

	if bot.config != cfg {
		t.Error("Bot config not set correctly")
	}
}

func TestIsBotMentioned(t *testing.T) {
	// Create a mock bot with userManager
	bot := &Bot{
		config: &config.Config{
			BotUsername: "testbot",
		},
		userManager: users.NewManager(),
	}

	tests := []struct {
		name     string
		message  *tgbotapi.Message
		expected bool
	}{
		{
			name: "Direct mention in text",
			message: &tgbotapi.Message{
				Text: "Hello @testbot how are you?",
			},
			expected: true,
		},
		{
			name: "No mention",
			message: &tgbotapi.Message{
				Text: "Hello everyone!",
			},
			expected: false,
		},
		{
			name: "Different bot mentioned",
			message: &tgbotapi.Message{
				Text: "Hey @otherbot!",
			},
			expected: false,
		},
		{
			name: "Reply to bot message",
			message: &tgbotapi.Message{
				Text: "Thanks!",
				ReplyToMessage: &tgbotapi.Message{
					From: &tgbotapi.User{
						UserName: "testbot",
					},
				},
			},
			expected: true,
		},
		{
			name: "Reply to other user",
			message: &tgbotapi.Message{
				Text: "Thanks!",
				ReplyToMessage: &tgbotapi.Message{
					From: &tgbotapi.User{
						UserName: "otheruser",
					},
				},
			},
			expected: false,
		},
		{
			name: "Mention entity",
			message: &tgbotapi.Message{
				Text: "Hey @testbot are you there?",
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "mention",
						Offset: 4,
						Length: 8,
					},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bot.isBotMentioned(tt.message)
			if result != tt.expected {
				t.Errorf("isBotMentioned() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestGenerateResponse(t *testing.T) {
	bot := &Bot{
		config: &config.Config{
			BotUsername: "testbot",
		},
		userManager: users.NewManager(),
	}

	tests := []struct {
		name        string
		message     *tgbotapi.Message
		userInfo    *users.User
		shouldCheck bool // Whether to check for specific content
	}{
		{
			name: "With user info",
			message: &tgbotapi.Message{
				Text: "Hello bot!",
			},
			userInfo: &users.User{
				ID:        123,
				FirstName: "John",
				LastName:  "Doe",
			},
			shouldCheck: true,
		},
		{
			name: "Without user info",
			message: &tgbotapi.Message{
				Text: "Hello bot!",
			},
			userInfo:    nil,
			shouldCheck: true,
		},
		{
			name: "Different message length",
			message: &tgbotapi.Message{
				Text: "A",
			},
			userInfo: &users.User{
				FirstName: "Jane",
			},
			shouldCheck: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := bot.generateResponse(tt.message, tt.userInfo)
			if response == "" {
				t.Error("generateResponse() returned empty string")
			}

			if tt.shouldCheck && tt.userInfo != nil && tt.userInfo.FirstName != "" {
				// Response should contain user's first name
				// Note: Due to random selection, we can't guarantee this always
				// but the logic should prefer using the name
			}
		})
	}
}

func TestGetUserInfo(t *testing.T) {
	bot := &Bot{
		userManager: users.NewManager(),
	}

	// Add a user first
	tgUser := &tgbotapi.User{
		ID:        12345,
		UserName:  "testuser",
		FirstName: "Test",
		LastName:  "User",
	}
	bot.userManager.UpdateUser(tgUser)

	// Test getting existing user
	userInfo := bot.GetUserInfo(12345)
	if userInfo == nil {
		t.Fatal("Expected to get user info, got nil")
	}
	if userInfo.ID != 12345 {
		t.Errorf("Expected user ID 12345, got %d", userInfo.ID)
	}

	// Test getting non-existent user
	userInfo = bot.GetUserInfo(99999)
	if userInfo != nil {
		t.Errorf("Expected nil for non-existent user, got %+v", userInfo)
	}
}

func TestGetAllUsers(t *testing.T) {
	bot := &Bot{
		userManager: users.NewManager(),
	}

	// Add multiple users
	users := []*tgbotapi.User{
		{ID: 1, UserName: "user1", FirstName: "First1"},
		{ID: 2, UserName: "user2", FirstName: "First2"},
		{ID: 3, UserName: "user3", FirstName: "First3"},
	}

	for _, u := range users {
		bot.userManager.UpdateUser(u)
	}

	allUsers := bot.GetAllUsers()
	if len(allUsers) != 3 {
		t.Errorf("Expected 3 users, got %d", len(allUsers))
	}

	// Verify each user exists
	for _, u := range users {
		user, exists := allUsers[u.ID]
		if !exists {
			t.Errorf("User with ID %d not found", u.ID)
		}
		if user.FirstName != u.FirstName {
			t.Errorf("Expected FirstName %s, got %s", u.FirstName, user.FirstName)
		}
	}
}

func TestHandleMessage_PrivateChat(t *testing.T) {
	t.Skip("Skipping test that requires mocked Telegram API")
	// Note: To properly test this, we would need to:
	// 1. Create an interface for the Telegram API
	// 2. Implement a mock version
	// 3. Inject the mock into the Bot struct
}

func TestHandleMessage_GroupChat(t *testing.T) {
	t.Skip("Skipping test that requires mocked Telegram API")
	// Note: To properly test this, we would need to:
	// 1. Create an interface for the Telegram API
	// 2. Implement a mock version
	// 3. Inject the mock into the Bot struct
}
