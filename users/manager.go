package users

import (
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// User represents a Telegram user with their information
type User struct {
	ID           int64
	UserName     string
	FirstName    string
	LastName     string
	MessageCount int
}

// Manager manages user information
type Manager struct {
	users map[int64]*User
	mu    sync.RWMutex
}

// NewManager creates a new user manager
func NewManager() *Manager {
	return &Manager{
		users: make(map[int64]*User),
	}
}

// UpdateUser updates or creates user information
func (m *Manager) UpdateUser(from *tgbotapi.User) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if user, exists := m.users[from.ID]; exists {
		// Update existing user
		user.UserName = from.UserName
		user.FirstName = from.FirstName
		user.LastName = from.LastName
		user.MessageCount++
	} else {
		// Create new user
		m.users[from.ID] = &User{
			ID:           from.ID,
			UserName:     from.UserName,
			FirstName:    from.FirstName,
			LastName:     from.LastName,
			MessageCount: 1,
		}
	}
}

// GetUser retrieves user information by ID
func (m *Manager) GetUser(userID int64) *User {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.users[userID]
}

// GetAllUsers returns all stored users
func (m *Manager) GetAllUsers() map[int64]*User {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy to prevent concurrent access issues
	users := make(map[int64]*User, len(m.users))
	for k, v := range m.users {
		userCopy := *v
		users[k] = &userCopy
	}
	return users
}
