package storage

import (
	"fmt"
	"sync"
	"time"
)

// MockStorage is an in-memory implementation for testing
type MockStorage struct {
	chats    map[int64]*Chat
	users    map[int64]*User
	settings map[int64]*ChatSettings
	messages []*Message
	profiles map[string]*UserProfile // key: "chatID:userID"
	mu       sync.RWMutex
}

// NewMockStorage creates a new mock storage
func NewMockStorage() *MockStorage {
	return &MockStorage{
		chats:    make(map[int64]*Chat),
		users:    make(map[int64]*User),
		settings: make(map[int64]*ChatSettings),
		messages: []*Message{},
		profiles: make(map[string]*UserProfile),
	}
}

func (m *MockStorage) Initialize() error { return nil }
func (m *MockStorage) Close() error      { return nil }

func (m *MockStorage) SaveChat(chat *Chat) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.chats[chat.ID] = chat
	return nil
}

func (m *MockStorage) GetChat(chatID int64) (*Chat, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.chats[chatID], nil
}

func (m *MockStorage) GetAllChats() ([]*Chat, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	chats := make([]*Chat, 0, len(m.chats))
	for _, chat := range m.chats {
		chats = append(chats, chat)
	}
	return chats, nil
}

func (m *MockStorage) UpdateChatMessageCount(chatID int64, count int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if chat, ok := m.chats[chatID]; ok {
		chat.MessageCount = count
		chat.UpdatedAt = time.Now()
	}
	return nil
}

func (m *MockStorage) SaveUser(user *User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.users[user.ID] = user
	return nil
}

func (m *MockStorage) GetUser(userID int64) (*User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.users[userID], nil
}

func (m *MockStorage) GetAllUsers() ([]*User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	users := make([]*User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func (m *MockStorage) GetChatUsers(chatID int64) ([]*User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	users := []*User{}
	for _, msg := range m.messages {
		if msg.ChatID == chatID {
			if user, ok := m.users[msg.UserID]; ok {
				users = append(users, user)
			}
		}
	}
	return users, nil
}

func (m *MockStorage) UpdateUserMessageCount(userID int64, count int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if user, ok := m.users[userID]; ok {
		user.MessageCount = count
		user.UpdatedAt = time.Now()
	}
	return nil
}

func (m *MockStorage) SaveChatSettings(chatID int64, settings *ChatSettings) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.settings[chatID] = settings
	return nil
}

func (m *MockStorage) GetChatSettings(chatID int64) (*ChatSettings, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.settings[chatID], nil
}

func (m *MockStorage) DeleteChatSettings(chatID int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.settings, chatID)
	return nil
}

func (m *MockStorage) SaveMessage(msg *Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	msg.ID = int64(len(m.messages) + 1)
	m.messages = append(m.messages, msg)
	return nil
}

func (m *MockStorage) GetRecentMessages(chatID int64, limit int) ([]*Message, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var messages []*Message
	for i := len(m.messages) - 1; i >= 0 && len(messages) < limit; i-- {
		if m.messages[i].ChatID == chatID {
			messages = append([]*Message{m.messages[i]}, messages...)
		}
	}
	return messages, nil
}

func (m *MockStorage) GetUserMessagesInChat(chatID int64, userID int64, limit int) ([]*Message, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var messages []*Message
	for i := len(m.messages) - 1; i >= 0 && len(messages) < limit; i-- {
		if m.messages[i].ChatID == chatID && m.messages[i].UserID == userID {
			messages = append(messages, m.messages[i])
		}
	}
	return messages, nil
}

func (m *MockStorage) GetMessagesByTimeRange(chatID int64, start, end time.Time) ([]*Message, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var messages []*Message
	for _, msg := range m.messages {
		if msg.ChatID == chatID && !msg.Timestamp.Before(start) && !msg.Timestamp.After(end) {
			messages = append(messages, msg)
		}
	}
	return messages, nil
}

func (m *MockStorage) SaveUserProfile(profile *UserProfile) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := profileKey(profile.ChatID, profile.UserID)
	m.profiles[key] = profile
	return nil
}

func (m *MockStorage) GetUserProfile(chatID int64, userID int64) (*UserProfile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	key := profileKey(chatID, userID)
	return m.profiles[key], nil
}

func (m *MockStorage) UpdateUserProfile(chatID int64, userID int64, updates map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := profileKey(chatID, userID)
	if profile, ok := m.profiles[key]; ok {
		for k, v := range updates {
			switch k {
			case "interests":
				profile.Interests = v.(string)
			case "topics":
				profile.Topics = v.(string)
			case "personality":
				profile.Personality = v.(string)
			case "notes":
				profile.Notes = v.(string)
			case "interaction_count":
				profile.InteractionCount = v.(int)
			}
		}
		profile.UpdatedAt = time.Now()
	}
	return nil
}

func profileKey(chatID, userID int64) string {
	return fmt.Sprintf("%d:%d", chatID, userID)
}
