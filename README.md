# HowardTheChad_bot

A Telegram bot built in Go with AI capabilities for group conversations.

## Features

- 🤖 Connect to Telegram groups and channels
- 👥 Retrieve and track user information
- 💬 Respond when mentioned with @ symbol
- ⚙️ **Per-Group Configurable Settings** - Each group has independent settings
- 👑 **Admin-Only Controls** - Group admins can configure the bot via commands
- 📊 User and chat activity tracking (per-group message counts)
- 🧠 Context-aware message handling (ready for AI model integration)
- � Easy configuration with slash commands (`/settings`, `/setfrequency`, etc.)

## Prerequisites

- Go 1.24 or higher
- Telegram Bot Token (obtain from [@BotFather](https://t.me/botfather))

## Installation

1. Clone the repository:
```bash
git clone https://github.com/Zind-dev/HowardTheChad_bot.git
cd HowardTheChad_bot
```

2. Install dependencies:
```bash
go mod download
```

3. Build the bot:
```bash
go build -o howardthechad_bot
```

## Configuration

### Required Environment Variables

- `TELEGRAM_BOT_TOKEN` - Your Telegram bot token from BotFather
- `BOT_USERNAME` - Your bot's username (without @)

### Optional Settings (with defaults)

- `BOT_RESPONSE_FREQUENCY` - How often to respond to regular messages (default: `10` = every 10th message)
- `BOT_RESPOND_TO_MENTIONS` - Always respond when mentioned (default: `true`)

**PowerShell Example:**
```powershell
$env:TELEGRAM_BOT_TOKEN = "your_bot_token_here"
$env:BOT_USERNAME = "your_bot_username"
$env:BOT_RESPONSE_FREQUENCY = "10"  # Optional
$env:BOT_RESPOND_TO_MENTIONS = "true"  # Optional
```

**Linux/Mac Example:**
```bash
export TELEGRAM_BOT_TOKEN="your_bot_token_here"
export BOT_USERNAME="your_bot_username"
export BOT_RESPONSE_FREQUENCY="10"  # Optional
export BOT_RESPOND_TO_MENTIONS="true"  # Optional
```

See [SETTINGS.md](SETTINGS.md) for detailed configuration options.

## Usage

Run the bot:
```bash
./howardthechad_bot
```

Or run directly with Go:
```bash
go run main.go
```

### Using the Bot

1. **Private Messages**: Send a direct message to the bot to receive a greeting

2. **Group Chats**: 
   - Add the bot to a group
   - Mention the bot with `@your_bot_username` to get a response
   - Reply to bot's messages to continue the conversation
   - Bot will automatically respond every Nth message (default: every 10th message)

3. **Admin Configuration** (in groups):
   - `/settings` - View current bot settings for your group
   - `/setfrequency <number>` - Change response frequency (admin only)
   - `/togglementions` - Toggle mention responses on/off (admin only)
   - `/resetsettings` - Reset to default settings (admin only)
   - `/help` - Show available commands

Each group can have independent settings configured by its administrators!

For detailed configuration guide, see [SETTINGS.md](SETTINGS.md).

## Project Structure

```
.
├── main.go           # Entry point
├── config/           # Configuration management
│   ├── config.go
│   └── config_test.go
├── bot/              # Bot logic and message handling
│   ├── bot.go
│   └── bot_test.go
├── users/            # User information management
│   ├── manager.go
│   └── manager_test.go
├── chats/            # Chat tracking and message counting
│   ├── manager.go
│   └── manager_test.go
├── settings/         # Bot behavior settings
│   ├── settings.go
│   └── settings_test.go
├── SETTINGS.md       # Settings configuration guide
└── TESTING.md        # Testing guide
```

## Future Enhancements

- [ ] AI model integration for intelligent responses
- [ ] Message context analysis
- [ ] User personality profiling
- [ ] Advanced conversation handling
- [ ] Persistent storage for user data
- [ ] Analytics and insights

## Development

### Building
```bash
go build -o howardthechad_bot
```

### Running Tests
```bash
go test ./...
```

See [TESTING.md](TESTING.md) for detailed testing instructions including integration tests.

### Code Formatting
```bash
go fmt ./...
```

## License

This project is open source and available under the MIT License.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.