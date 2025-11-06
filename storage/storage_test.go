package storage

import (
	"os"
	"testing"
	"time"
)

func TestSQLiteStorage(t *testing.T) {
	// Create temporary database file
	dbPath := "test_bot.db"
	defer os.Remove(dbPath)

	// Initialize storage
	storage, err := NewSQLiteStorage(dbPath)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer storage.Close()

	if err := storage.Initialize(); err != nil {
		t.Fatalf("Failed to initialize storage: %v", err)
	}

	t.Run("Chat Operations", func(t *testing.T) {
		chat := &Chat{
			ID:           123,
			Title:        "Test Chat",
			Type:         "group",
			MessageCount: 5,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		// Save chat
		if err := storage.SaveChat(chat); err != nil {
			t.Errorf("Failed to save chat: %v", err)
		}

		// Get chat
		retrieved, err := storage.GetChat(123)
		if err != nil {
			t.Errorf("Failed to get chat: %v", err)
		}
		if retrieved == nil {
			t.Error("Chat not found")
		} else {
			if retrieved.Title != "Test Chat" {
				t.Errorf("Expected title 'Test Chat', got '%s'", retrieved.Title)
			}
			if retrieved.MessageCount != 5 {
				t.Errorf("Expected message count 5, got %d", retrieved.MessageCount)
			}
		}

		// Update message count
		if err := storage.UpdateChatMessageCount(123, 10); err != nil {
			t.Errorf("Failed to update message count: %v", err)
		}

		updated, _ := storage.GetChat(123)
		if updated.MessageCount != 10 {
			t.Errorf("Expected message count 10, got %d", updated.MessageCount)
		}

		// Get all chats
		chats, err := storage.GetAllChats()
		if err != nil {
			t.Errorf("Failed to get all chats: %v", err)
		}
		if len(chats) != 1 {
			t.Errorf("Expected 1 chat, got %d", len(chats))
		}
	})

	t.Run("User Operations", func(t *testing.T) {
		user := &User{
			ID:           456,
			UserName:     "testuser",
			FirstName:    "Test",
			LastName:     "User",
			MessageCount: 3,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		// Save user
		if err := storage.SaveUser(user); err != nil {
			t.Errorf("Failed to save user: %v", err)
		}

		// Get user
		retrieved, err := storage.GetUser(456)
		if err != nil {
			t.Errorf("Failed to get user: %v", err)
		}
		if retrieved == nil {
			t.Error("User not found")
		} else {
			if retrieved.UserName != "testuser" {
				t.Errorf("Expected username 'testuser', got '%s'", retrieved.UserName)
			}
			if retrieved.FirstName != "Test" {
				t.Errorf("Expected first name 'Test', got '%s'", retrieved.FirstName)
			}
		}

		// Update message count
		if err := storage.UpdateUserMessageCount(456, 7); err != nil {
			t.Errorf("Failed to update user message count: %v", err)
		}

		updated, _ := storage.GetUser(456)
		if updated.MessageCount != 7 {
			t.Errorf("Expected message count 7, got %d", updated.MessageCount)
		}

		// Get all users
		users, err := storage.GetAllUsers()
		if err != nil {
			t.Errorf("Failed to get all users: %v", err)
		}
		if len(users) != 1 {
			t.Errorf("Expected 1 user, got %d", len(users))
		}
	})

	t.Run("Chat Settings Operations", func(t *testing.T) {
		settings := &ChatSettings{
			ChatID:                  123,
			ResponseFrequency:       15,
			AlwaysRespondToMentions: true,
			CreatedAt:               time.Now(),
			UpdatedAt:               time.Now(),
		}

		// Save settings
		if err := storage.SaveChatSettings(123, settings); err != nil {
			t.Errorf("Failed to save chat settings: %v", err)
		}

		// Get settings
		retrieved, err := storage.GetChatSettings(123)
		if err != nil {
			t.Errorf("Failed to get chat settings: %v", err)
		}
		if retrieved == nil {
			t.Error("Settings not found")
		} else {
			if retrieved.ResponseFrequency != 15 {
				t.Errorf("Expected frequency 15, got %d", retrieved.ResponseFrequency)
			}
			if !retrieved.AlwaysRespondToMentions {
				t.Error("Expected always respond to mentions to be true")
			}
		}

		// Update settings
		settings.ResponseFrequency = 20
		if err := storage.SaveChatSettings(123, settings); err != nil {
			t.Errorf("Failed to update chat settings: %v", err)
		}

		updated, _ := storage.GetChatSettings(123)
		if updated.ResponseFrequency != 20 {
			t.Errorf("Expected frequency 20, got %d", updated.ResponseFrequency)
		}

		// Delete settings
		if err := storage.DeleteChatSettings(123); err != nil {
			t.Errorf("Failed to delete chat settings: %v", err)
		}

		deleted, _ := storage.GetChatSettings(123)
		if deleted != nil {
			t.Error("Settings should be deleted")
		}
	})

	t.Run("Message Operations", func(t *testing.T) {
		// Save messages
		messages := []*Message{
			{
				ChatID:    123,
				UserID:    456,
				Text:      "Hello",
				IsBot:     false,
				Timestamp: time.Now().Add(-5 * time.Minute),
			},
			{
				ChatID:    123,
				UserID:    456,
				Text:      "How are you?",
				IsBot:     false,
				Timestamp: time.Now().Add(-4 * time.Minute),
			},
			{
				ChatID:    123,
				UserID:    789,
				Text:      "I'm fine, thanks!",
				IsBot:     true,
				Timestamp: time.Now().Add(-3 * time.Minute),
			},
			{
				ChatID:    123,
				UserID:    456,
				Text:      "Great!",
				IsBot:     false,
				Timestamp: time.Now().Add(-2 * time.Minute),
			},
		}

		for _, msg := range messages {
			if err := storage.SaveMessage(msg); err != nil {
				t.Errorf("Failed to save message: %v", err)
			}
			if msg.ID == 0 {
				t.Error("Message ID should be set after save")
			}
		}

		// Get recent messages
		recent, err := storage.GetRecentMessages(123, 3)
		if err != nil {
			t.Errorf("Failed to get recent messages: %v", err)
		}
		if len(recent) != 3 {
			t.Errorf("Expected 3 messages, got %d", len(recent))
		}
		// Check chronological order
		if recent[0].Text != "How are you?" {
			t.Errorf("Expected first message to be 'How are you?', got '%s'", recent[0].Text)
		}
		if recent[2].Text != "Great!" {
			t.Errorf("Expected last message to be 'Great!', got '%s'", recent[2].Text)
		}

		// Get user messages in chat
		userMsgs, err := storage.GetUserMessagesInChat(123, 456, 10)
		if err != nil {
			t.Errorf("Failed to get user messages: %v", err)
		}
		if len(userMsgs) != 3 {
			t.Errorf("Expected 3 messages from user 456, got %d", len(userMsgs))
		}

		// Get messages by time range
		start := time.Now().Add(-4 * time.Minute)
		end := time.Now()
		rangeMessages, err := storage.GetMessagesByTimeRange(123, start, end)
		if err != nil {
			t.Errorf("Failed to get messages by time range: %v", err)
		}
		if len(rangeMessages) < 2 {
			t.Errorf("Expected at least 2 messages in time range, got %d", len(rangeMessages))
		}
	})

	t.Run("User Profile Operations", func(t *testing.T) {
		profile := &UserProfile{
			ChatID:           123,
			UserID:           456,
			Interests:        "gaming, music",
			Topics:           "video games, concerts",
			Personality:      "friendly, enthusiastic",
			LastInteraction:  time.Now(),
			InteractionCount: 5,
			Notes:            "Active participant",
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		// Save profile
		if err := storage.SaveUserProfile(profile); err != nil {
			t.Errorf("Failed to save user profile: %v", err)
		}

		// Get profile
		retrieved, err := storage.GetUserProfile(123, 456)
		if err != nil {
			t.Errorf("Failed to get user profile: %v", err)
		}
		if retrieved == nil {
			t.Error("Profile not found")
		} else {
			if retrieved.Interests != "gaming, music" {
				t.Errorf("Expected interests 'gaming, music', got '%s'", retrieved.Interests)
			}
			if retrieved.InteractionCount != 5 {
				t.Errorf("Expected interaction count 5, got %d", retrieved.InteractionCount)
			}
		}

		// Update profile
		updates := map[string]interface{}{
			"interaction_count": 10,
			"notes":             "Very active",
		}
		if err := storage.UpdateUserProfile(123, 456, updates); err != nil {
			t.Errorf("Failed to update user profile: %v", err)
		}

		updated, _ := storage.GetUserProfile(123, 456)
		if updated.InteractionCount != 10 {
			t.Errorf("Expected interaction count 10, got %d", updated.InteractionCount)
		}
		if updated.Notes != "Very active" {
			t.Errorf("Expected notes 'Very active', got '%s'", updated.Notes)
		}
	})

	t.Run("Get Chat Users", func(t *testing.T) {
		// Add another user
		user2 := &User{
			ID:        789,
			UserName:  "botuser",
			FirstName: "Bot",
			LastName:  "User",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		storage.SaveUser(user2)

		// Get chat users (should return users who have messages in the chat)
		users, err := storage.GetChatUsers(123)
		if err != nil {
			t.Errorf("Failed to get chat users: %v", err)
		}
		if len(users) != 2 {
			t.Errorf("Expected 2 users in chat, got %d", len(users))
		}
	})
}

func TestSQLiteStorageNonExistent(t *testing.T) {
	dbPath := "test_nonexistent.db"
	defer os.Remove(dbPath)

	storage, err := NewSQLiteStorage(dbPath)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer storage.Close()

	if err := storage.Initialize(); err != nil {
		t.Fatalf("Failed to initialize storage: %v", err)
	}

	// Test getting non-existent records
	t.Run("Non-existent Chat", func(t *testing.T) {
		chat, err := storage.GetChat(999)
		if err != nil {
			t.Errorf("Should not error on non-existent chat: %v", err)
		}
		if chat != nil {
			t.Error("Expected nil for non-existent chat")
		}
	})

	t.Run("Non-existent User", func(t *testing.T) {
		user, err := storage.GetUser(999)
		if err != nil {
			t.Errorf("Should not error on non-existent user: %v", err)
		}
		if user != nil {
			t.Error("Expected nil for non-existent user")
		}
	})

	t.Run("Non-existent Settings", func(t *testing.T) {
		settings, err := storage.GetChatSettings(999)
		if err != nil {
			t.Errorf("Should not error on non-existent settings: %v", err)
		}
		if settings != nil {
			t.Error("Expected nil for non-existent settings")
		}
	})

	t.Run("Non-existent Profile", func(t *testing.T) {
		profile, err := storage.GetUserProfile(999, 999)
		if err != nil {
			t.Errorf("Should not error on non-existent profile: %v", err)
		}
		if profile != nil {
			t.Error("Expected nil for non-existent profile")
		}
	})
}
