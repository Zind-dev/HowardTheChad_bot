package chats

import (
	"sync"
)

// Chat represents a chat with message tracking
type Chat struct {
	ID           int64
	Title        string
	Type         string
	MessageCount int
}

// Manager manages chat information and message counts
type Manager struct {
	chats map[int64]*Chat
	mu    sync.RWMutex
}

// NewManager creates a new chat manager
func NewManager() *Manager {
	return &Manager{
		chats: make(map[int64]*Chat),
	}
}

// IncrementMessageCount increments the message count for a chat
func (m *Manager) IncrementMessageCount(chatID int64, title, chatType string) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	if chat, exists := m.chats[chatID]; exists {
		chat.MessageCount++
		chat.Title = title
		chat.Type = chatType
		return chat.MessageCount
	}

	// Create new chat
	m.chats[chatID] = &Chat{
		ID:           chatID,
		Title:        title,
		Type:         chatType,
		MessageCount: 1,
	}
	return 1
}

// GetChat retrieves chat information by ID
func (m *Manager) GetChat(chatID int64) *Chat {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.chats[chatID]
}

// GetMessageCount returns the message count for a chat
func (m *Manager) GetMessageCount(chatID int64) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if chat, exists := m.chats[chatID]; exists {
		return chat.MessageCount
	}
	return 0
}

// GetAllChats returns all tracked chats
func (m *Manager) GetAllChats() map[int64]*Chat {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy to prevent concurrent access issues
	chats := make(map[int64]*Chat, len(m.chats))
	for k, v := range m.chats {
		chatCopy := *v
		chats[k] = &chatCopy
	}
	return chats
}
