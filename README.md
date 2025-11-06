# HowardTheChad_bot

A Telegram bot built in Go with AI capabilities for group conversations.

## Features

- 🤖 Connect to Telegram groups and channels
- 👥 Retrieve and track user information
- 💬 Respond when mentioned with @ symbol
- 🧠 Context-aware message handling (ready for AI model integration)
- 📊 User activity tracking

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

Set the following environment variables:

- `TELEGRAM_BOT_TOKEN` - Your Telegram bot token from BotFather
- `BOT_USERNAME` - Your bot's username (without @)

Example:
```bash
export TELEGRAM_BOT_TOKEN="your_bot_token_here"
export BOT_USERNAME="your_bot_username"
```

Or create a `.env` file (not tracked in git):
```
TELEGRAM_BOT_TOKEN=your_bot_token_here
BOT_USERNAME=your_bot_username
```

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

## Project Structure

```
.
├── main.go           # Entry point
├── config/           # Configuration management
│   └── config.go
├── bot/              # Bot logic and message handling
│   └── bot.go
└── users/            # User information management
    └── manager.go
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

### Code Formatting
```bash
go fmt ./...
```

## License

This project is open source and available under the MIT License.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.