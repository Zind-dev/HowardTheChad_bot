package settings

import (
	"sync"
)

// Settings holds bot behavior configuration
type Settings struct {
	// ResponseFrequency determines how often bot responds to regular messages (e.g., every 10th message)
	// 0 means never respond to regular messages, only when mentioned
	ResponseFrequency int

	// AlwaysRespondToMentions when true, bot always responds when mentioned
	AlwaysRespondToMentions bool
}

// Manager manages settings per chat
type Manager struct {
	chatSettings map[int64]*Settings
	defaults     *Settings
	mu           sync.RWMutex
}

// NewManager creates a new settings manager with default settings
func NewManager(defaults *Settings) *Manager {
	if defaults == nil {
		defaults = NewDefaultSettings()
	}
	return &Manager{
		chatSettings: make(map[int64]*Settings),
		defaults:     defaults,
	}
}

// NewDefaultSettings creates settings with default values
func NewDefaultSettings() *Settings {
	return &Settings{
		ResponseFrequency:       10, // Respond every 10th message
		AlwaysRespondToMentions: true,
	}
}

// NewCustomSettings creates settings with custom values
func NewCustomSettings(frequency int, alwaysRespondToMentions bool) *Settings {
	return &Settings{
		ResponseFrequency:       frequency,
		AlwaysRespondToMentions: alwaysRespondToMentions,
	}
}

// ShouldRespondToRegularMessage determines if bot should respond based on message count
func (s *Settings) ShouldRespondToRegularMessage(messageCount int) bool {
	if s.ResponseFrequency <= 0 {
		return false
	}
	return messageCount%s.ResponseFrequency == 0
}

// ShouldRespondToMention determines if bot should respond to mentions
func (s *Settings) ShouldRespondToMention() bool {
	return s.AlwaysRespondToMentions
}

// GetSettings returns settings for a specific chat, or defaults if not set
func (m *Manager) GetSettings(chatID int64) *Settings {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if settings, exists := m.chatSettings[chatID]; exists {
		return settings
	}
	return m.defaults
}

// SetSettings sets custom settings for a specific chat
func (m *Manager) SetSettings(chatID int64, settings *Settings) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.chatSettings[chatID] = settings
}

// SetFrequency sets the response frequency for a specific chat
func (m *Manager) SetFrequency(chatID int64, frequency int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if settings, exists := m.chatSettings[chatID]; exists {
		settings.ResponseFrequency = frequency
	} else {
		// Create new settings based on defaults
		m.chatSettings[chatID] = &Settings{
			ResponseFrequency:       frequency,
			AlwaysRespondToMentions: m.defaults.AlwaysRespondToMentions,
		}
	}
}

// ToggleMentionResponse toggles the mention response setting for a specific chat
func (m *Manager) ToggleMentionResponse(chatID int64) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if settings, exists := m.chatSettings[chatID]; exists {
		settings.AlwaysRespondToMentions = !settings.AlwaysRespondToMentions
		return settings.AlwaysRespondToMentions
	} else {
		// Create new settings based on defaults
		newValue := !m.defaults.AlwaysRespondToMentions
		m.chatSettings[chatID] = &Settings{
			ResponseFrequency:       m.defaults.ResponseFrequency,
			AlwaysRespondToMentions: newValue,
		}
		return newValue
	}
}

// ResetSettings resets a chat to default settings
func (m *Manager) ResetSettings(chatID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.chatSettings, chatID)
}

// GetAllChatSettings returns all custom chat settings
func (m *Manager) GetAllChatSettings() map[int64]*Settings {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy to prevent concurrent access issues
	settings := make(map[int64]*Settings, len(m.chatSettings))
	for k, v := range m.chatSettings {
		settingsCopy := *v
		settings[k] = &settingsCopy
	}
	return settings
}
