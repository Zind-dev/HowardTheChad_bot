# Bot Settings Configuration

The HowardTheChad bot now supports **per-group configurable behavior settings** to control when and how it responds to messages.

## 🌟 Key Features

- ✅ **Per-Group Settings** - Each group can have independent settings
- ✅ **Admin-Only Control** - Only group administrators can change settings
- ✅ **Easy Commands** - Simple slash commands to configure the bot
- ✅ **Persistent Settings** - Settings remain active for each group independently

## Default Behavior

By default, the bot is configured with:
- **Response Frequency**: Responds every **10th message** in group chats
- **Always Respond to Mentions**: **Enabled** (responds every time it's mentioned with @)

These defaults can be overridden globally via environment variables or per-group using bot commands.

## 🤖 Bot Commands (In-Group Configuration)

Group administrators can use these commands to configure the bot's behavior for their specific group:

### View Current Settings
```
/settings
```
Shows the current settings for the group.

### Change Response Frequency
```
/setfrequency <number>
```
Set how often the bot responds to regular messages.

**Examples:**
- `/setfrequency 10` - Respond every 10th message (default)
- `/setfrequency 5` - Respond every 5th message (more active)
- `/setfrequency 20` - Respond every 20th message (less active)
- `/setfrequency 0` - Only respond to mentions, never to regular messages
- `/setfrequency 1` - Respond to every message

### Toggle Mention Responses
```
/togglementions
```
Turn automatic responses to mentions on or off. Each use toggles the setting.

### Reset to Defaults
```
/resetsettings
```
Reset the group's settings back to the global defaults.

### Help
```
/help
```
Show all available commands.

**Note:** Only group administrators can use configuration commands. All members can view settings with `/settings` and `/help`.

## 🔧 Global Configuration Options (Environment Variables)

You can set global defaults using environment variables (these apply to all groups that haven't customized their settings):

### BOT_RESPONSE_FREQUENCY

Controls how often the bot responds to regular (non-mention) messages in group chats.

```powershell
# Respond every 10th message (default)
$env:BOT_RESPONSE_FREQUENCY = "10"

# Respond every 5th message
$env:BOT_RESPONSE_FREQUENCY = "5"

# Never respond to regular messages (only mentions)
$env:BOT_RESPONSE_FREQUENCY = "0"

# Respond to every message
$env:BOT_RESPONSE_FREQUENCY = "1"
```

### BOT_RESPOND_TO_MENTIONS

Controls whether the bot responds when mentioned.

```powershell
# Always respond to mentions (default)
$env:BOT_RESPOND_TO_MENTIONS = "true"

# Disable automatic response to mentions
$env:BOT_RESPOND_TO_MENTIONS = "false"
```

## Usage Examples

### Example 1: Default Settings
```powershell
$env:TELEGRAM_BOT_TOKEN = "your_token_here"
$env:BOT_USERNAME = "your_bot_username"
# Bot will respond every 10th message and always to mentions
go run main.go
```

### Example 2: More Frequent Responses
```powershell
$env:TELEGRAM_BOT_TOKEN = "your_token_here"
$env:BOT_USERNAME = "your_bot_username"
$env:BOT_RESPONSE_FREQUENCY = "5"
# Bot will respond every 5th message and always to mentions
go run main.go
```

### Example 3: Mention-Only Mode
```powershell
$env:TELEGRAM_BOT_TOKEN = "your_token_here"
$env:BOT_USERNAME = "your_bot_username"
$env:BOT_RESPONSE_FREQUENCY = "0"
# Bot will only respond when mentioned
go run main.go
```

### Example 4: Silent Mode
```powershell
$env:TELEGRAM_BOT_TOKEN = "your_token_here"
$env:BOT_USERNAME = "your_bot_username"
$env:BOT_RESPONSE_FREQUENCY = "0"
$env:BOT_RESPOND_TO_MENTIONS = "false"
# Bot will track messages but won't respond (useful for data collection)
go run main.go
```

## How It Works

### Per-Group Settings
- Each group has **independent settings**
- Settings are stored per group and persist during the bot's runtime
- Groups without custom settings use the global defaults
- Admin changes apply immediately to their group only

### Message Counting
- The bot tracks message counts **per chat**
- Each chat has its own independent counter
- The counter increments for every message in the group
- When the counter reaches a multiple of the group's configured frequency, the bot responds

### Mention Detection
The bot detects mentions in several ways:
1. **Text mentions**: Messages containing `@your_bot_username`
2. **Entity mentions**: Telegram mention entities
3. **Reply mentions**: When someone replies to the bot's message

### Response Logic
For each message in a group chat:
1. **If the bot is mentioned** AND `BOT_RESPOND_TO_MENTIONS` is true → Respond
2. **If message count is a multiple of `BOT_RESPONSE_FREQUENCY`** → Respond
3. Otherwise → Stay silent (but still track the message)

### Private Messages
In private (direct) messages, the bot **always responds** regardless of settings.

### Admin Verification
- The bot verifies that the user issuing a configuration command is a group administrator
- Only users with "creator" or "administrator" status can change settings
- Non-admins attempting to use configuration commands receive an error message

## 💡 Usage Scenarios

### Scenario 1: Different Activity Levels
```
# Active group - respond frequently
/setfrequency 5

# Quiet group - respond occasionally  
/setfrequency 20

# Very active group - mentions only
/setfrequency 0
```

### Scenario 2: Testing Bot Behavior
```
# Make bot very responsive for testing
/setfrequency 1
/togglementions  # Enable if needed

# After testing, reset
/resetsettings
```

### Scenario 3: Temporary Silence
```
# Temporarily disable automatic responses
/setfrequency 0
/togglementions  # Disable mention responses

# Re-enable later
/resetsettings
```

## 💻 Programmatic Settings Management

You can also manage settings programmatically:

```go
import "github.com/Zind-dev/HowardTheChad_bot/settings"

// Create custom settings for a specific chat
customSettings := settings.NewCustomSettings(
    3,     // Respond every 3rd message
    false, // Don't respond to mentions
)

// Update settings for a specific chat
bot.UpdateSettings(chatID, customSettings)

// Get settings for a specific chat
chatSettings := bot.GetSettings(chatID)

// Use the settings manager directly
bot.settingsManager.SetFrequency(chatID, 5)
bot.settingsManager.ToggleMentionResponse(chatID)
bot.settingsManager.ResetSettings(chatID)
```

## Monitoring

### Get Chat Information
```go
// Get info about a specific chat
chatInfo := bot.GetChatInfo(chatID)
fmt.Printf("Chat: %s, Messages: %d\n", chatInfo.Title, chatInfo.MessageCount)

// Get all tracked chats
allChats := bot.GetAllChats()
for id, chat := range allChats {
    fmt.Printf("Chat ID: %d, Title: %s, Messages: %d\n", 
        id, chat.Title, chat.MessageCount)
}
```

### Get User Information
```go
// Get info about a specific user
userInfo := bot.GetUserInfo(userID)
fmt.Printf("User: %s, Messages: %d\n", userInfo.FirstName, userInfo.MessageCount)

// Get all tracked users
allUsers := bot.GetAllUsers()
for id, user := range allUsers {
    fmt.Printf("User ID: %d, Name: %s, Messages: %d\n", 
        id, user.FirstName, user.MessageCount)
}
```

## Testing

Run tests to verify settings behavior:

```powershell
# Test settings package
go test ./settings -v

# Test chat tracking
go test ./chats -v

# Test configuration loading
go test ./config -v

# Test all packages
go test ./... -v
```

## Troubleshooting

### Bot isn't responding in group
1. Check that the bot is added to the group
2. Check current settings with `/settings`
3. Verify frequency is not set to 0 (unless you only want mention responses)
4. Ensure enough messages have been sent (needs to reach the frequency threshold)
5. Check logs for errors

### Bot responds too often
- Use `/setfrequency <higher number>` to reduce frequency
- Example: `/setfrequency 20` instead of `/setfrequency 5`

### Bot doesn't respond to mentions
- Check with `/settings` if mention responses are enabled
- Use `/togglementions` to enable them
- Ensure you're using the correct bot username with @ symbol
- Verify the bot has permission to read and send messages

### Can't change settings
- Verify you're a group administrator
- Check that you're running commands in the group chat (not private message)
- Make sure command syntax is correct (use `/help` to see examples)

### Settings not applying
- Settings apply immediately after command execution
- Use `/settings` to verify the change was saved
- Each group has independent settings - make sure you're configuring the right group

## Best Practices

1. **Start conservative**: Use higher frequency values (10-20) to avoid spamming
2. **Customize per group**: Each group has different needs - adjust settings accordingly
3. **Test first**: Try `/setfrequency 1` temporarily to test bot responses
4. **Monitor reactions**: Adjust based on group members' feedback
5. **Use `/settings` regularly**: Check current configuration when troubleshooting
6. **Communicate changes**: Let group members know when you change bot behavior
7. **Reset when unsure**: Use `/resetsettings` to return to safe defaults

## Examples of Real Usage

### Active Discussion Group
Members want frequent bot participation:
```
/setfrequency 5
```

### Casual Chat Group  
Less frequent responses preferred:
```
/setfrequency 15
```

### Announcement Group
Bot should only respond when directly mentioned:
```
/setfrequency 0
```

### Testing/Development Group
Maximum responsiveness for testing:
```
/setfrequency 1
```
