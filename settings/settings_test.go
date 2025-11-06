package settings

import "testing"

func TestNewDefaultSettings(t *testing.T) {
	settings := NewDefaultSettings()

	if settings == nil {
		t.Fatal("NewDefaultSettings returned nil")
	}

	if settings.ResponseFrequency != 10 {
		t.Errorf("Expected ResponseFrequency 10, got %d", settings.ResponseFrequency)
	}

	if !settings.AlwaysRespondToMentions {
		t.Error("Expected AlwaysRespondToMentions to be true")
	}
}

func TestNewCustomSettings(t *testing.T) {
	settings := NewCustomSettings(5, false)

	if settings == nil {
		t.Fatal("NewCustomSettings returned nil")
	}

	if settings.ResponseFrequency != 5 {
		t.Errorf("Expected ResponseFrequency 5, got %d", settings.ResponseFrequency)
	}

	if settings.AlwaysRespondToMentions {
		t.Error("Expected AlwaysRespondToMentions to be false")
	}
}

func TestShouldRespondToRegularMessage(t *testing.T) {
	tests := []struct {
		name         string
		frequency    int
		messageCount int
		expected     bool
	}{
		{
			name:         "10th message with frequency 10",
			frequency:    10,
			messageCount: 10,
			expected:     true,
		},
		{
			name:         "20th message with frequency 10",
			frequency:    10,
			messageCount: 20,
			expected:     true,
		},
		{
			name:         "5th message with frequency 10",
			frequency:    10,
			messageCount: 5,
			expected:     false,
		},
		{
			name:         "1st message with frequency 10",
			frequency:    10,
			messageCount: 1,
			expected:     false,
		},
		{
			name:         "5th message with frequency 5",
			frequency:    5,
			messageCount: 5,
			expected:     true,
		},
		{
			name:         "Frequency 0 never responds",
			frequency:    0,
			messageCount: 10,
			expected:     false,
		},
		{
			name:         "Negative frequency never responds",
			frequency:    -1,
			messageCount: 10,
			expected:     false,
		},
		{
			name:         "Every message (frequency 1)",
			frequency:    1,
			messageCount: 7,
			expected:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			settings := NewCustomSettings(tt.frequency, true)
			result := settings.ShouldRespondToRegularMessage(tt.messageCount)
			if result != tt.expected {
				t.Errorf("ShouldRespondToRegularMessage(%d) = %v, expected %v",
					tt.messageCount, result, tt.expected)
			}
		})
	}
}

func TestShouldRespondToMention(t *testing.T) {
	tests := []struct {
		name                    string
		alwaysRespondToMentions bool
		expected                bool
	}{
		{
			name:                    "Always respond enabled",
			alwaysRespondToMentions: true,
			expected:                true,
		},
		{
			name:                    "Always respond disabled",
			alwaysRespondToMentions: false,
			expected:                false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			settings := NewCustomSettings(10, tt.alwaysRespondToMentions)
			result := settings.ShouldRespondToMention()
			if result != tt.expected {
				t.Errorf("ShouldRespondToMention() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestNewManager(t *testing.T) {
	defaults := NewDefaultSettings()
	manager := NewManager(defaults)

	if manager == nil {
		t.Fatal("NewManager returned nil")
	}

	if manager.chatSettings == nil {
		t.Fatal("Manager chatSettings map is nil")
	}
}

func TestManagerGetSettings_Default(t *testing.T) {
	defaults := NewCustomSettings(15, false)
	manager := NewManager(defaults)

	// Get settings for a chat that hasn't been configured
	settings := manager.GetSettings(12345)

	if settings.ResponseFrequency != 15 {
		t.Errorf("Expected default frequency 15, got %d", settings.ResponseFrequency)
	}

	if settings.AlwaysRespondToMentions {
		t.Error("Expected default AlwaysRespondToMentions to be false")
	}
}

func TestManagerSetSettings(t *testing.T) {
	manager := NewManager(NewDefaultSettings())

	customSettings := NewCustomSettings(5, false)
	manager.SetSettings(100, customSettings)

	retrieved := manager.GetSettings(100)
	if retrieved.ResponseFrequency != 5 {
		t.Errorf("Expected frequency 5, got %d", retrieved.ResponseFrequency)
	}
	if retrieved.AlwaysRespondToMentions {
		t.Error("Expected AlwaysRespondToMentions to be false")
	}
}

func TestManagerSetFrequency(t *testing.T) {
	manager := NewManager(NewDefaultSettings())

	// Set frequency for new chat
	manager.SetFrequency(100, 20)

	settings := manager.GetSettings(100)
	if settings.ResponseFrequency != 20 {
		t.Errorf("Expected frequency 20, got %d", settings.ResponseFrequency)
	}

	// Update existing chat
	manager.SetFrequency(100, 30)
	settings = manager.GetSettings(100)
	if settings.ResponseFrequency != 30 {
		t.Errorf("Expected frequency 30, got %d", settings.ResponseFrequency)
	}
}

func TestManagerToggleMentionResponse(t *testing.T) {
	manager := NewManager(NewDefaultSettings())

	// Toggle for new chat (default is true, so should become false)
	result := manager.ToggleMentionResponse(100)
	if result {
		t.Error("Expected toggle to return false (toggled from default true)")
	}

	settings := manager.GetSettings(100)
	if settings.AlwaysRespondToMentions {
		t.Error("Expected AlwaysRespondToMentions to be false after toggle")
	}

	// Toggle again
	result = manager.ToggleMentionResponse(100)
	if !result {
		t.Error("Expected toggle to return true")
	}

	settings = manager.GetSettings(100)
	if !settings.AlwaysRespondToMentions {
		t.Error("Expected AlwaysRespondToMentions to be true after second toggle")
	}
}

func TestManagerResetSettings(t *testing.T) {
	manager := NewManager(NewDefaultSettings())

	// Set custom settings
	manager.SetFrequency(100, 5)

	// Verify custom settings
	settings := manager.GetSettings(100)
	if settings.ResponseFrequency != 5 {
		t.Errorf("Expected frequency 5, got %d", settings.ResponseFrequency)
	}

	// Reset
	manager.ResetSettings(100)

	// Should now return defaults
	settings = manager.GetSettings(100)
	if settings.ResponseFrequency != 10 {
		t.Errorf("Expected default frequency 10 after reset, got %d", settings.ResponseFrequency)
	}
}

func TestManagerGetAllChatSettings(t *testing.T) {
	manager := NewManager(NewDefaultSettings())

	// Set settings for multiple chats
	manager.SetFrequency(100, 5)
	manager.SetFrequency(200, 15)
	manager.SetFrequency(300, 25)

	allSettings := manager.GetAllChatSettings()

	if len(allSettings) != 3 {
		t.Errorf("Expected 3 custom settings, got %d", len(allSettings))
	}

	// Verify each setting
	if allSettings[100].ResponseFrequency != 5 {
		t.Errorf("Chat 100: expected frequency 5, got %d", allSettings[100].ResponseFrequency)
	}
	if allSettings[200].ResponseFrequency != 15 {
		t.Errorf("Chat 200: expected frequency 15, got %d", allSettings[200].ResponseFrequency)
	}
	if allSettings[300].ResponseFrequency != 25 {
		t.Errorf("Chat 300: expected frequency 25, got %d", allSettings[300].ResponseFrequency)
	}
}

func TestManagerIndependentChats(t *testing.T) {
	manager := NewManager(NewDefaultSettings())

	// Set different settings for different chats
	manager.SetFrequency(100, 5)
	manager.SetFrequency(200, 10)

	// Verify independence
	settings100 := manager.GetSettings(100)
	settings200 := manager.GetSettings(200)

	if settings100.ResponseFrequency != 5 {
		t.Errorf("Chat 100: expected frequency 5, got %d", settings100.ResponseFrequency)
	}
	if settings200.ResponseFrequency != 10 {
		t.Errorf("Chat 200: expected frequency 10, got %d", settings200.ResponseFrequency)
	}

	// Toggle mentions for one chat shouldn't affect the other
	manager.ToggleMentionResponse(100)

	settings100 = manager.GetSettings(100)
	settings200 = manager.GetSettings(200)

	if settings100.AlwaysRespondToMentions {
		t.Error("Chat 100: expected AlwaysRespondToMentions to be false")
	}
	if !settings200.AlwaysRespondToMentions {
		t.Error("Chat 200: expected AlwaysRespondToMentions to be true (unchanged)")
	}
}
