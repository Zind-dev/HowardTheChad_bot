# Storage Layer Documentation

## Overview

The HowardTheChad bot uses a storage abstraction layer to persist data for AI context and user knowledge. The storage system is designed to support future AI integration by tracking:

- Message history for context
- User profiles with interests and personality
- Chat information and settings
- Conversation patterns and topics

## Architecture

### Storage Interface
Located in `storage/storage.go`, defines the contract for all storage implementations:
- **Chat operations**: Save, retrieve, and update chat information
- **User operations**: Manage user data and message counts
- **Settings operations**: Per-chat configuration persistence
- **Message operations**: Store and retrieve conversation history
- **Profile operations**: AI-focused user profiling (interests, personality, topics)

### Implementations

#### SQLiteStorage (`storage/sqlite.go`)
Production implementation using SQLite3 database:
- **Database file**: `bot_data.db` (created automatically)
- **Requirements**: CGO enabled, MinGW GCC (C:\mingw64\bin)
- **Features**:
  - ACID transactions
  - Indexed queries for performance
  - Foreign key constraints
  - Automatic timestamps

#### MockStorage (`storage/mock.go`)
In-memory implementation for testing:
- No persistence (data lost on restart)
- Thread-safe with mutex locks
- Perfect for unit tests
- No external dependencies

## Database Schema

### Tables

#### `chats`
- `id` (INTEGER PRIMARY KEY): Telegram chat ID
- `title` (TEXT): Chat name
- `type` (TEXT): Chat type (group, supergroup, private)
- `message_count` (INTEGER): Total messages in chat
- `created_at`, `updated_at` (DATETIME): Timestamps

#### `users`
- `id` (INTEGER PRIMARY KEY): Telegram user ID
- `username` (TEXT): @username
- `first_name`, `last_name` (TEXT): User display names
- `message_count` (INTEGER): Total user messages
- `created_at`, `updated_at` (DATETIME): Timestamps

#### `chat_settings`
- `chat_id` (INTEGER PRIMARY KEY): Links to chats.id
- `response_frequency` (INTEGER): Respond every N messages
- `always_respond_to_mentions` (BOOLEAN): Mention behavior
- `created_at`, `updated_at` (DATETIME): Timestamps

#### `messages`
- `id` (INTEGER AUTOINCREMENT): Unique message ID
- `chat_id` (INTEGER): Links to chats.id
- `user_id` (INTEGER): Links to users.id
- `text` (TEXT): Message content
- `is_bot` (BOOLEAN): Bot response flag
- `timestamp` (DATETIME): Message time

Indexes:
- `idx_messages_chat_id`: Fast chat message lookup
- `idx_messages_timestamp`: Time-based queries

#### `user_profiles`
- `chat_id`, `user_id` (COMPOSITE PRIMARY KEY): Per-chat user profile
- `interests` (TEXT): User interests (for AI context)
- `topics` (TEXT): Discussion topics
- `personality` (TEXT): Personality traits
- `last_interaction` (DATETIME): Last message timestamp
- `interaction_count` (INTEGER): Number of interactions
- `notes` (TEXT): AI notes about user
- `created_at`, `updated_at` (DATETIME): Timestamps

## Usage

### Initialization (main.go)
```go
// Initialize storage
store, err := storage.NewSQLiteStorage("bot_data.db")
if err != nil {
    log.Fatalf("Failed to create storage: %v", err)
}
defer store.Close()

if err := store.Initialize(); err != nil {
    log.Fatalf("Failed to initialize database: %v", err)
}

// Pass to bot
b, err := bot.New(cfg, store)
```

### Saving Messages (bot.go)
Messages are automatically saved when received:
```go
msg := &storage.Message{
    ChatID:    message.Chat.ID,
    UserID:    message.From.ID,
    Text:      message.Text,
    IsBot:     message.From.IsBot,
    Timestamp: time.Now(),
}
storage.SaveMessage(msg)
```

### Retrieving Context for AI
```go
// Get last 20 messages for context
messages, err := storage.GetRecentMessages(chatID, 20)

// Get user profile
profile, err := storage.GetUserProfile(chatID, userID)

// Get user's recent messages
userMsgs, err := storage.GetUserMessagesInChat(chatID, userID, 10)
```

### User Profiles (Future AI Integration)
```go
profile := &storage.UserProfile{
    ChatID:           chatID,
    UserID:           userID,
    Interests:        "gaming, music, technology",
    Topics:           "video games, concerts, AI",
    Personality:      "friendly, enthusiastic, helpful",
    LastInteraction:  time.Now(),
    InteractionCount: 42,
    Notes:            "Active participant, likes detailed explanations",
}
storage.SaveUserProfile(profile)

// Update specific fields
updates := map[string]interface{}{
    "interaction_count": 43,
    "notes": "Prefers casual conversation",
}
storage.UpdateUserProfile(chatID, userID, updates)
```

## Building with CGO

SQLite requires CGO (C compiler):

### Windows Setup
1. Install MinGW: https://www.mingw-w64.org/
2. Add to PATH: `C:\mingw64\bin`
3. Build with CGO:
   ```powershell
   $env:PATH = "C:\mingw64\bin;$env:PATH"
   $env:CGO_ENABLED = "1"
   go build -o HowardTheChad_bot.exe
   ```

The `start_bot.ps1` script automatically sets these variables.

## Testing

### Run Storage Tests
```powershell
cd storage
$env:PATH = "C:\mingw64\bin;$env:PATH"
$env:CGO_ENABLED = "1"
go test -v
```

### Using Mock Storage in Tests
```go
mockStore := storage.NewMockStorage()
bot, err := bot.New(cfg, mockStore)
```

## Data Persistence

- **Location**: `bot_data.db` in the bot's directory
- **Backup**: Copy `bot_data.db` file
- **Reset**: Delete `bot_data.db` (will recreate on next run)
- **Migration**: Database schema auto-creates on Initialize()

## Future Enhancements

### AI Integration Points
1. **Context Building**: Use GetRecentMessages() to build AI context
2. **User Understanding**: Load UserProfile for personalized responses
3. **Topic Tracking**: Analyze message history for conversation topics
4. **Personality Adaptation**: Adjust bot responses based on user profiles
5. **Memory**: Maintain long-term conversation context across sessions

### Planned Features
- [ ] Conversation summarization
- [ ] Sentiment analysis storage
- [ ] Topic categorization
- [ ] User preference learning
- [ ] Multi-language support
- [ ] Message embeddings for semantic search

## Performance

### Optimization Tips
- Messages table grows quickly - consider cleanup policies
- Use GetRecentMessages() with reasonable limits (10-50 messages)
- Index frequently queried columns
- Batch inserts when possible
- Regular VACUUM for SQLite maintenance

### Monitoring
```sql
-- Check database size
SELECT page_count * page_size as size FROM pragma_page_count(), pragma_page_size();

-- Count messages per chat
SELECT chat_id, COUNT(*) FROM messages GROUP BY chat_id;

-- Most active users
SELECT user_id, COUNT(*) as msg_count FROM messages GROUP BY user_id ORDER BY msg_count DESC LIMIT 10;
```

## Troubleshooting

### "cgo: C compiler not found"
- Ensure MinGW is installed and in PATH
- Check: `where gcc` should show `C:\mingw64\bin\gcc.exe`
- Set CGO_ENABLED=1 before building

### Database Locked Errors
- Only one process can write at a time
- Ensure bot is not running multiple instances
- Check for zombie processes

### Performance Issues
- Check database size (messages table)
- Run VACUUM to optimize
- Add indexes on frequently queried columns
- Consider message retention policies
