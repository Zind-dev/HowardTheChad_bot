package chats

import (
	"testing"
)

func TestNewManager(t *testing.T) {
	manager := NewManager()
	if manager == nil {
		t.Fatal("NewManager returned nil")
	}
	if manager.chats == nil {
		t.Fatal("Manager chats map is nil")
	}
}

func TestIncrementMessageCount(t *testing.T) {
	manager := NewManager()

	// First message in chat
	count := manager.IncrementMessageCount(123, "Test Group", "group")
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}

	// Second message
	count = manager.IncrementMessageCount(123, "Test Group", "group")
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}

	// Third message
	count = manager.IncrementMessageCount(123, "Test Group", "group")
	if count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}
}

func TestIncrementMessageCount_MultipleChats(t *testing.T) {
	manager := NewManager()

	// Chat 1
	count1 := manager.IncrementMessageCount(100, "Chat 1", "group")
	if count1 != 1 {
		t.Errorf("Chat 1: Expected count 1, got %d", count1)
	}

	// Chat 2
	count2 := manager.IncrementMessageCount(200, "Chat 2", "supergroup")
	if count2 != 1 {
		t.Errorf("Chat 2: Expected count 1, got %d", count2)
	}

	// Chat 1 again
	count1 = manager.IncrementMessageCount(100, "Chat 1", "group")
	if count1 != 2 {
		t.Errorf("Chat 1: Expected count 2, got %d", count1)
	}

	// Verify counts are independent
	if manager.GetMessageCount(100) != 2 {
		t.Errorf("Chat 1: Expected final count 2, got %d", manager.GetMessageCount(100))
	}
	if manager.GetMessageCount(200) != 1 {
		t.Errorf("Chat 2: Expected final count 1, got %d", manager.GetMessageCount(200))
	}
}

func TestGetChat(t *testing.T) {
	manager := NewManager()

	// Create a chat
	manager.IncrementMessageCount(123, "Test Group", "group")

	chat := manager.GetChat(123)
	if chat == nil {
		t.Fatal("GetChat returned nil")
	}

	if chat.ID != 123 {
		t.Errorf("Expected ID 123, got %d", chat.ID)
	}
	if chat.Title != "Test Group" {
		t.Errorf("Expected title 'Test Group', got '%s'", chat.Title)
	}
	if chat.Type != "group" {
		t.Errorf("Expected type 'group', got '%s'", chat.Type)
	}
	if chat.MessageCount != 1 {
		t.Errorf("Expected message count 1, got %d", chat.MessageCount)
	}
}

func TestGetChat_NotFound(t *testing.T) {
	manager := NewManager()

	chat := manager.GetChat(99999)
	if chat != nil {
		t.Errorf("Expected nil for non-existent chat, got %+v", chat)
	}
}

func TestGetMessageCount(t *testing.T) {
	manager := NewManager()

	// Non-existent chat
	count := manager.GetMessageCount(999)
	if count != 0 {
		t.Errorf("Expected count 0 for non-existent chat, got %d", count)
	}

	// Add messages
	manager.IncrementMessageCount(123, "Test", "group")
	manager.IncrementMessageCount(123, "Test", "group")
	manager.IncrementMessageCount(123, "Test", "group")

	count = manager.GetMessageCount(123)
	if count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}
}

func TestGetAllChats(t *testing.T) {
	manager := NewManager()

	// Add multiple chats
	manager.IncrementMessageCount(1, "Chat 1", "group")
	manager.IncrementMessageCount(2, "Chat 2", "supergroup")
	manager.IncrementMessageCount(3, "Chat 3", "group")

	allChats := manager.GetAllChats()
	if len(allChats) != 3 {
		t.Errorf("Expected 3 chats, got %d", len(allChats))
	}

	// Verify each chat
	for id := int64(1); id <= 3; id++ {
		chat, exists := allChats[id]
		if !exists {
			t.Errorf("Chat with ID %d not found in GetAllChats", id)
		}
		if chat.ID != id {
			t.Errorf("Chat ID mismatch: expected %d, got %d", id, chat.ID)
		}
	}
}

func TestChatTitleUpdate(t *testing.T) {
	manager := NewManager()

	// Create chat with initial title
	manager.IncrementMessageCount(123, "Old Title", "group")

	// Update with new title
	manager.IncrementMessageCount(123, "New Title", "group")

	chat := manager.GetChat(123)
	if chat.Title != "New Title" {
		t.Errorf("Expected title to update to 'New Title', got '%s'", chat.Title)
	}
}
