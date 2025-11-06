package users

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestNewManager(t *testing.T) {
	manager := NewManager()
	if manager == nil {
		t.Fatal("NewManager returned nil")
	}
	if manager.users == nil {
		t.Fatal("Manager users map is nil")
	}
}

func TestUpdateUser(t *testing.T) {
	manager := NewManager()
	
	// Create a mock user
	tgUser := &tgbotapi.User{
		ID:        12345,
		UserName:  "testuser",
		FirstName: "Test",
		LastName:  "User",
	}
	
	// Update user for the first time
	manager.UpdateUser(tgUser)
	
	// Retrieve the user
	user := manager.GetUser(12345)
	if user == nil {
		t.Fatal("GetUser returned nil")
	}
	
	if user.ID != 12345 {
		t.Errorf("Expected ID 12345, got %d", user.ID)
	}
	if user.UserName != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", user.UserName)
	}
	if user.FirstName != "Test" {
		t.Errorf("Expected first name 'Test', got '%s'", user.FirstName)
	}
	if user.MessageCount != 1 {
		t.Errorf("Expected message count 1, got %d", user.MessageCount)
	}
	
	// Update user again
	manager.UpdateUser(tgUser)
	user = manager.GetUser(12345)
	if user.MessageCount != 2 {
		t.Errorf("Expected message count 2 after second update, got %d", user.MessageCount)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	manager := NewManager()
	
	user := manager.GetUser(99999)
	if user != nil {
		t.Errorf("Expected nil for non-existent user, got %+v", user)
	}
}

func TestGetAllUsers(t *testing.T) {
	manager := NewManager()
	
	// Add multiple users
	users := []*tgbotapi.User{
		{ID: 1, UserName: "user1", FirstName: "First1", LastName: "Last1"},
		{ID: 2, UserName: "user2", FirstName: "First2", LastName: "Last2"},
		{ID: 3, UserName: "user3", FirstName: "First3", LastName: "Last3"},
	}
	
	for _, u := range users {
		manager.UpdateUser(u)
	}
	
	allUsers := manager.GetAllUsers()
	if len(allUsers) != 3 {
		t.Errorf("Expected 3 users, got %d", len(allUsers))
	}
	
	// Verify each user
	for _, u := range users {
		user, exists := allUsers[u.ID]
		if !exists {
			t.Errorf("User with ID %d not found in GetAllUsers", u.ID)
		}
		if user.UserName != u.UserName {
			t.Errorf("Username mismatch for user %d: expected %s, got %s", u.ID, u.UserName, user.UserName)
		}
	}
}
