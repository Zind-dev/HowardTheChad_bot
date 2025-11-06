# Usage Guide

## Getting Started

### 1. Create Your Bot

1. Open Telegram and search for [@BotFather](https://t.me/botfather)
2. Send `/newbot` and follow the instructions
3. Choose a name for your bot (e.g., "HowardTheChad")
4. Choose a username for your bot (e.g., "howardthechad_bot")
5. BotFather will give you a token - save it securely!

### 2. Configure the Bot

Create a `.env` file or set environment variables:

```bash
export TELEGRAM_BOT_TOKEN="123456789:ABCdefGHIjklMNOpqrsTUVwxyz"
export BOT_USERNAME="howardthechad_bot"
```

### 3. Run the Bot

```bash
go run main.go
```

Or build and run:

```bash
go build -o howardthechad_bot
./howardthechad_bot
```

You should see:
```
Authorized on account howardthechad_bot
Bot is starting...
```

## Using the Bot

### Private Conversations

1. Find your bot in Telegram
2. Send any message
3. The bot will respond with a greeting

Example:
```
You: Hello!
Bot: Hello! I'm HowardTheChad bot. I can see you, John! Add me to a group and mention me with @ to chat!
```

### Group Conversations

#### Adding the Bot to a Group

1. Create a group or use an existing one
2. Click on the group name → Add Members
3. Search for your bot username
4. Add the bot to the group

#### Interacting with the Bot

The bot will respond when:

1. **You mention it with @**
   ```
   You: @howardthechad_bot what do you think?
   Bot: Hey John! What's up?
   ```

2. **You reply to one of its messages**
   ```
   You: [Reply to bot's message] Tell me more
   Bot: Hello John! I'm here to help.
   ```

## Bot Behavior

### Current Functionality

- **User Tracking**: The bot remembers all users it interacts with
- **Context Awareness**: Uses user's first name when responding
- **Mention Detection**: Responds when mentioned with @
- **Reply Detection**: Responds when you reply to its messages

### Response Patterns

The bot currently uses simple response patterns. It selects from:
- "Hey {name}! What's up?"
- "Hello {name}! I'm here to help."
- "Hi {name}! What can I do for you?"
- "{name}, I'm listening!"
- "Yo {name}! How can I contribute?"

### Future AI Integration

The bot is designed to integrate with AI models. The response generation happens in `bot.generateResponse()`, which can be easily replaced with:

- OpenAI GPT API
- Anthropic Claude API
- Local LLM models (Llama, Mistral, etc.)
- Custom fine-tuned models

## Troubleshooting

### Bot Not Responding

1. **Check the bot is running**
   - Look for "Bot is starting..." in the console
   
2. **Verify environment variables**
   - `TELEGRAM_BOT_TOKEN` must be set correctly
   - `BOT_USERNAME` must match your bot's username (without @)

3. **In groups, make sure to mention the bot**
   - Use `@yourbotusername` in your message
   - Or reply to one of the bot's messages

### Permission Issues

If the bot can't read messages in a group:
1. Go to BotFather
2. Send `/setprivacy`
3. Choose your bot
4. Select "Disable" to allow the bot to see all messages
   - Note: This is optional, the bot works fine with privacy enabled as long as you mention it

### Common Errors

```
Failed to load configuration: TELEGRAM_BOT_TOKEN environment variable is required
```
→ Set the `TELEGRAM_BOT_TOKEN` environment variable

```
Failed to create bot: unauthorized
```
→ Check that your bot token is correct

## Advanced Usage

### Monitoring User Activity

The bot tracks:
- User ID
- Username
- First and Last name
- Message count

Access this data programmatically:
```go
// Get specific user info
userInfo := bot.GetUserInfo(userID)

// Get all users
allUsers := bot.GetAllUsers()
```

### Customizing Responses

Edit `bot/bot.go`, find the `generateResponse` method:

```go
func (b *Bot) generateResponse(message *tgbotapi.Message, userInfo *users.User) string {
    // Your custom logic here
    // This is where you'd integrate an AI model
}
```

## Next Steps

- Integrate an AI model for intelligent responses
- Add persistent storage for user data
- Implement conversation context tracking
- Add commands (e.g., /help, /stats)
- Build personality profiling system
