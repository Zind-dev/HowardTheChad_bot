package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteStorage implements the Storage interface using SQLite
type SQLiteStorage struct {
	db *sql.DB
}

// NewSQLiteStorage creates a new SQLite storage instance
func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	storage := &SQLiteStorage{db: db}
	return storage, nil
}

// Initialize creates all necessary tables
func (s *SQLiteStorage) Initialize() error {
	schema := `
	CREATE TABLE IF NOT EXISTS chats (
		id INTEGER PRIMARY KEY,
		title TEXT,
		type TEXT,
		message_count INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		username TEXT,
		first_name TEXT,
		last_name TEXT,
		message_count INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS chat_settings (
		chat_id INTEGER PRIMARY KEY,
		response_frequency INTEGER DEFAULT 10,
		always_respond_to_mentions BOOLEAN DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chat_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		text TEXT,
		is_bot BOOLEAN DEFAULT 0,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (chat_id) REFERENCES chats(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);

	CREATE INDEX IF NOT EXISTS idx_messages_chat_id ON messages(chat_id);
	CREATE INDEX IF NOT EXISTS idx_messages_timestamp ON messages(timestamp);

	CREATE TABLE IF NOT EXISTS user_profiles (
		chat_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		interests TEXT,
		topics TEXT,
		personality TEXT,
		last_interaction DATETIME,
		interaction_count INTEGER DEFAULT 0,
		notes TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (chat_id, user_id),
		FOREIGN KEY (chat_id) REFERENCES chats(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`

	_, err := s.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	return nil
}

// Close closes the database connection
func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}

// SaveChat saves or updates a chat
func (s *SQLiteStorage) SaveChat(chat *Chat) error {
	query := `
	INSERT INTO chats (id, title, type, message_count, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		title = excluded.title,
		type = excluded.type,
		message_count = excluded.message_count,
		updated_at = excluded.updated_at
	`

	_, err := s.db.Exec(query,
		chat.ID, chat.Title, chat.Type, chat.MessageCount,
		chat.CreatedAt, time.Now())

	return err
}

// GetChat retrieves a chat by ID
func (s *SQLiteStorage) GetChat(chatID int64) (*Chat, error) {
	query := `SELECT id, title, type, message_count, created_at, updated_at 
	          FROM chats WHERE id = ?`

	chat := &Chat{}
	err := s.db.QueryRow(query, chatID).Scan(
		&chat.ID, &chat.Title, &chat.Type, &chat.MessageCount,
		&chat.CreatedAt, &chat.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return chat, nil
}

// GetAllChats retrieves all chats
func (s *SQLiteStorage) GetAllChats() ([]*Chat, error) {
	query := `SELECT id, title, type, message_count, created_at, updated_at FROM chats`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []*Chat
	for rows.Next() {
		chat := &Chat{}
		err := rows.Scan(&chat.ID, &chat.Title, &chat.Type, &chat.MessageCount,
			&chat.CreatedAt, &chat.UpdatedAt)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	return chats, nil
}

// UpdateChatMessageCount updates the message count for a chat
func (s *SQLiteStorage) UpdateChatMessageCount(chatID int64, count int) error {
	query := `UPDATE chats SET message_count = ?, updated_at = ? WHERE id = ?`
	_, err := s.db.Exec(query, count, time.Now(), chatID)
	return err
}

// SaveUser saves or updates a user
func (s *SQLiteStorage) SaveUser(user *User) error {
	query := `
	INSERT INTO users (id, username, first_name, last_name, message_count, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		username = excluded.username,
		first_name = excluded.first_name,
		last_name = excluded.last_name,
		message_count = excluded.message_count,
		updated_at = excluded.updated_at
	`

	_, err := s.db.Exec(query,
		user.ID, user.UserName, user.FirstName, user.LastName,
		user.MessageCount, user.CreatedAt, time.Now())

	return err
}

// GetUser retrieves a user by ID
func (s *SQLiteStorage) GetUser(userID int64) (*User, error) {
	query := `SELECT id, username, first_name, last_name, message_count, created_at, updated_at 
	          FROM users WHERE id = ?`

	user := &User{}
	err := s.db.QueryRow(query, userID).Scan(
		&user.ID, &user.UserName, &user.FirstName, &user.LastName,
		&user.MessageCount, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetAllUsers retrieves all users
func (s *SQLiteStorage) GetAllUsers() ([]*User, error) {
	query := `SELECT id, username, first_name, last_name, message_count, created_at, updated_at FROM users`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.UserName, &user.FirstName, &user.LastName,
			&user.MessageCount, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetChatUsers retrieves all users who have participated in a chat
func (s *SQLiteStorage) GetChatUsers(chatID int64) ([]*User, error) {
	query := `
	SELECT DISTINCT u.id, u.username, u.first_name, u.last_name, 
	       u.message_count, u.created_at, u.updated_at
	FROM users u
	INNER JOIN messages m ON u.id = m.user_id
	WHERE m.chat_id = ?
	ORDER BY u.username
	`

	rows, err := s.db.Query(query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.UserName, &user.FirstName, &user.LastName,
			&user.MessageCount, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// UpdateUserMessageCount updates the message count for a user
func (s *SQLiteStorage) UpdateUserMessageCount(userID int64, count int) error {
	query := `UPDATE users SET message_count = ?, updated_at = ? WHERE id = ?`
	_, err := s.db.Exec(query, count, time.Now(), userID)
	return err
}

// SaveChatSettings saves or updates chat settings
func (s *SQLiteStorage) SaveChatSettings(chatID int64, settings *ChatSettings) error {
	query := `
	INSERT INTO chat_settings (chat_id, response_frequency, always_respond_to_mentions, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(chat_id) DO UPDATE SET
		response_frequency = excluded.response_frequency,
		always_respond_to_mentions = excluded.always_respond_to_mentions,
		updated_at = excluded.updated_at
	`

	_, err := s.db.Exec(query,
		chatID, settings.ResponseFrequency, settings.AlwaysRespondToMentions,
		settings.CreatedAt, time.Now())

	return err
}

// GetChatSettings retrieves settings for a chat
func (s *SQLiteStorage) GetChatSettings(chatID int64) (*ChatSettings, error) {
	query := `SELECT chat_id, response_frequency, always_respond_to_mentions, created_at, updated_at 
	          FROM chat_settings WHERE chat_id = ?`

	settings := &ChatSettings{}
	err := s.db.QueryRow(query, chatID).Scan(
		&settings.ChatID, &settings.ResponseFrequency, &settings.AlwaysRespondToMentions,
		&settings.CreatedAt, &settings.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return settings, nil
}

// DeleteChatSettings deletes settings for a chat
func (s *SQLiteStorage) DeleteChatSettings(chatID int64) error {
	query := `DELETE FROM chat_settings WHERE chat_id = ?`
	_, err := s.db.Exec(query, chatID)
	return err
}

// SaveMessage saves a message to the database
func (s *SQLiteStorage) SaveMessage(msg *Message) error {
	query := `INSERT INTO messages (chat_id, user_id, text, is_bot, timestamp) 
	          VALUES (?, ?, ?, ?, ?)`

	result, err := s.db.Exec(query, msg.ChatID, msg.UserID, msg.Text, msg.IsBot, msg.Timestamp)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err == nil {
		msg.ID = id
	}

	return nil
}

// GetRecentMessages retrieves recent messages from a chat (for AI context)
func (s *SQLiteStorage) GetRecentMessages(chatID int64, limit int) ([]*Message, error) {
	query := `
	SELECT id, chat_id, user_id, text, is_bot, timestamp 
	FROM messages 
	WHERE chat_id = ? 
	ORDER BY timestamp DESC 
	LIMIT ?
	`

	rows, err := s.db.Query(query, chatID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		msg := &Message{}
		err := rows.Scan(&msg.ID, &msg.ChatID, &msg.UserID, &msg.Text, &msg.IsBot, &msg.Timestamp)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	// Reverse to get chronological order
	for i := 0; i < len(messages)/2; i++ {
		j := len(messages) - 1 - i
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// GetUserMessagesInChat retrieves recent messages from a specific user in a chat
func (s *SQLiteStorage) GetUserMessagesInChat(chatID int64, userID int64, limit int) ([]*Message, error) {
	query := `
	SELECT id, chat_id, user_id, text, is_bot, timestamp 
	FROM messages 
	WHERE chat_id = ? AND user_id = ?
	ORDER BY timestamp DESC 
	LIMIT ?
	`

	rows, err := s.db.Query(query, chatID, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		msg := &Message{}
		err := rows.Scan(&msg.ID, &msg.ChatID, &msg.UserID, &msg.Text, &msg.IsBot, &msg.Timestamp)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// GetMessagesByTimeRange retrieves messages within a time range
func (s *SQLiteStorage) GetMessagesByTimeRange(chatID int64, start, end time.Time) ([]*Message, error) {
	query := `
	SELECT id, chat_id, user_id, text, is_bot, timestamp 
	FROM messages 
	WHERE chat_id = ? AND timestamp BETWEEN ? AND ?
	ORDER BY timestamp ASC
	`

	rows, err := s.db.Query(query, chatID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		msg := &Message{}
		err := rows.Scan(&msg.ID, &msg.ChatID, &msg.UserID, &msg.Text, &msg.IsBot, &msg.Timestamp)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// SaveUserProfile saves or updates a user profile
func (s *SQLiteStorage) SaveUserProfile(profile *UserProfile) error {
	query := `
	INSERT INTO user_profiles (chat_id, user_id, interests, topics, personality, 
	                           last_interaction, interaction_count, notes, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(chat_id, user_id) DO UPDATE SET
		interests = excluded.interests,
		topics = excluded.topics,
		personality = excluded.personality,
		last_interaction = excluded.last_interaction,
		interaction_count = excluded.interaction_count,
		notes = excluded.notes,
		updated_at = excluded.updated_at
	`

	_, err := s.db.Exec(query,
		profile.ChatID, profile.UserID, profile.Interests, profile.Topics,
		profile.Personality, profile.LastInteraction, profile.InteractionCount,
		profile.Notes, profile.CreatedAt, time.Now())

	return err
}

// GetUserProfile retrieves a user profile for a specific chat
func (s *SQLiteStorage) GetUserProfile(chatID int64, userID int64) (*UserProfile, error) {
	query := `
	SELECT chat_id, user_id, interests, topics, personality, last_interaction, 
	       interaction_count, notes, created_at, updated_at
	FROM user_profiles 
	WHERE chat_id = ? AND user_id = ?
	`

	profile := &UserProfile{}
	err := s.db.QueryRow(query, chatID, userID).Scan(
		&profile.ChatID, &profile.UserID, &profile.Interests, &profile.Topics,
		&profile.Personality, &profile.LastInteraction, &profile.InteractionCount,
		&profile.Notes, &profile.CreatedAt, &profile.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// UpdateUserProfile updates specific fields of a user profile
func (s *SQLiteStorage) UpdateUserProfile(chatID int64, userID int64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	// Build dynamic query
	query := "UPDATE user_profiles SET "
	args := []interface{}{}
	first := true

	for key, value := range updates {
		if !first {
			query += ", "
		}
		query += key + " = ?"
		args = append(args, value)
		first = false
	}

	query += ", updated_at = ? WHERE chat_id = ? AND user_id = ?"
	args = append(args, time.Now(), chatID, userID)

	_, err := s.db.Exec(query, args...)
	return err
}
