# Quick Start Guide

## Step 1: Register Commands with BotFather (IMPORTANT!)

For commands like `/help` and `/settings` to work, register them with BotFather:

1. Open Telegram and message `@BotFather`
2. Send `/mybots`
3. Select: `HowardTheChad_bot`
4. Click `Edit Bot` → `Edit Commands`
5. Paste this command list:

```
settings - Show current bot settings for this group
setfrequency - Change response frequency (admin only)
togglementions - Toggle mention responses on/off (admin only)
resetsettings - Reset settings to defaults (admin only)
help - Show available commands and usage
```

## Step 2: Your Current Setup

- ✅ Bot is added to a group
- ✅ You are the only member (you're the admin)
- ✅ Default: Respond every 10th message and when mentioned

## Running with Default Settings (Recommended)

The bot is already configured with your desired defaults:
- Responds every **10th message**
- Responds **every time** it's mentioned

**Just run:**
```powershell
cd "c:\Users\Zaslon\Desktop\New folder\test_bot\HowardTheChad_bot"
$env:TELEGRAM_BOT_TOKEN = "your_token_here"
$env:BOT_USERNAME = "HowardTheChad_bot"
go run main.go
```

## Testing in Your Group

### Test Mention Response
1. In your group, type: `@HowardTheChad_bot hello`
2. Bot should respond immediately

### Test Every 10th Message
1. Send 9 regular messages in the group (without mentioning the bot)
2. On the 10th message, the bot will automatically respond
3. It will respond again on the 20th, 30th, 40th message, etc.

## Changing Settings

### More Frequent Responses (Every 5th Message)
```powershell
$env:TELEGRAM_BOT_TOKEN = "your_token_here"
$env:BOT_USERNAME = "HowardTheChad_bot"
$env:BOT_RESPONSE_FREQUENCY = "5"
go run main.go
```

### Only Respond to Mentions (No Automatic Messages)
```powershell
$env:TELEGRAM_BOT_TOKEN = "your_token_here"
$env:BOT_USERNAME = "HowardTheChad_bot"
$env:BOT_RESPONSE_FREQUENCY = "0"
go run main.go
```

### Respond to Every Message
```powershell
$env:TELEGRAM_BOT_TOKEN = "your_token_here"
$env:BOT_USERNAME = "HowardTheChad_bot"
$env:BOT_RESPONSE_FREQUENCY = "1"
go run main.go
```

## What Gets Tracked

The bot now tracks:
- **User Information**: ID, username, first name, last name, message count per user
- **Chat Information**: Chat ID, title, type, total message count per chat
- **Each chat has its own message counter** (independent from other chats)

## Verifying It Works

### Check Message Counter
The bot logs will show activity. Watch for:
- "Authorized on account HowardTheChad_bot" - Bot is connected
- Each message increments the internal counter for that chat
- When the counter hits a multiple of 10, the bot responds

### Private Messages
If you message the bot privately, it **always responds** regardless of settings.

## Troubleshooting

### Bot Doesn't Respond in Group
1. **Check bot has permissions**: Make sure the bot can read and send messages
2. **Count your messages**: Remember it only responds every 10th message by default
3. **Try mentioning**: Type `@HowardTheChad_bot test` - this should always work

### Bot Responds Too Often/Too Little
Adjust `BOT_RESPONSE_FREQUENCY`:
- Higher number = less frequent responses
- Lower number = more frequent responses
- 0 = never respond automatically (only mentions)

### Bot Doesn't See Messages
- Make sure the bot has **privacy mode disabled** in BotFather
- Send `/setprivacy` to @BotFather
- Choose your bot
- Select "Disable" to allow bot to see all messages

## Testing the Bot

Run the test suite:
```powershell
cd "c:\Users\Zaslon\Desktop\New folder\test_bot\HowardTheChad_bot"
go test ./... -v
```

Run with your real token for integration testing:
```powershell
.\test_with_token.ps1
```

## Current Test Coverage

- ✅ **config**: 100% coverage
- ✅ **users**: 100% coverage  
- ✅ **chats**: 100% coverage
- ✅ **settings**: 100% coverage
- ⚠️ **bot**: 28.8% coverage (some parts need API mocking)

## Next Steps

1. Start the bot with default settings
2. Test in your group by sending messages
3. Adjust `BOT_RESPONSE_FREQUENCY` if needed
4. Monitor the behavior and fine-tune

For more details:
- **Settings Configuration**: See [SETTINGS.md](SETTINGS.md)
- **Testing Guide**: See [TESTING.md](TESTING.md)
- **Project Overview**: See [README.md](README.md)
