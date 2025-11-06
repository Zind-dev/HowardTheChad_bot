package storage

import (
	"time"
)

// Storage defines the interface for data persistence
// This allows switching between SQLite, PostgreSQL, or other backends
type Storage interface {
	// Initialize sets up the storage (creates tables, connections, etc.)
	Initialize() error

	// Close closes the storage connection
	Close() error

	// Chat operations
	SaveChat(chat *Chat) error
	GetChat(chatID int64) (*Chat, error)
	GetAllChats() ([]*Chat, error)
	UpdateChatMessageCount(chatID int64, count int) error

	// User operations
	SaveUser(user *User) error
	GetUser(userID int64) (*User, error)
	GetAllUsers() ([]*User, error)
	GetChatUsers(chatID int64) ([]*User, error)
	UpdateUserMessageCount(userID int64, count int) error

	// Settings operations
	SaveChatSettings(chatID int64, settings *ChatSettings) error
	GetChatSettings(chatID int64) (*ChatSettings, error)
	DeleteChatSettings(chatID int64) error

	// Message history operations (for AI context)
	SaveMessage(msg *Message) error
	GetRecentMessages(chatID int64, limit int) ([]*Message, error)
	GetUserMessagesInChat(chatID int64, userID int64, limit int) ([]*Message, error)
	GetMessagesByTimeRange(chatID int64, start, end time.Time) ([]*Message, error)

	// User profile operations (for AI personalization)
	SaveUserProfile(profile *UserProfile) error
	GetUserProfile(chatID int64, userID int64) (*UserProfile, error)
	UpdateUserProfile(chatID int64, userID int64, updates map[string]interface{}) error
}

// Chat represents a Telegram chat
type Chat struct {
	ID           int64
	Title        string
	Type         string
	MessageCount int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// User represents a Telegram user
type User struct {
	ID           int64
	UserName     string
	FirstName    string
	LastName     string
	MessageCount int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ChatSettings represents bot settings for a specific chat
type ChatSettings struct {
	ChatID                  int64
	ResponseFrequency       int
	AlwaysRespondToMentions bool
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

// Message represents a chat message (for AI context)
type Message struct {
	ID        int64
	ChatID    int64
	UserID    int64
	Text      string
	IsBot     bool
	Timestamp time.Time
}

// UserProfile stores AI-relevant information about a user in a specific chat
type UserProfile struct {
	ChatID int64
	UserID int64
	// AI context data
	Interests        string // JSON or comma-separated
	Topics           string // Topics user talks about
	Personality      string // Observed personality traits
	LastInteraction  time.Time
	InteractionCount int
	Notes            string // Any other AI-relevant notes
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
