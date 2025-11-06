# Testing Guide

This document describes how to run tests for the HowardTheChad bot.

## Running Basic Tests

Run all tests without requiring a real bot token:

```powershell
go test ./...
```

Run tests with verbose output:

```powershell
go test ./... -v
```

Run tests with coverage:

```powershell
go test ./... -cover
```

## Running Integration Tests with Real Bot Token

Some tests can optionally use a real Telegram bot token for more comprehensive integration testing.

### Method 1: Using the PowerShell Script

```powershell
.\test_with_token.ps1
```

The script will prompt you for your test bot token and username.

### Method 2: Setting Environment Variables Manually

Set the environment variables:

```powershell
$env:TEST_BOT_TOKEN = "your_bot_token_here"
$env:TEST_BOT_USERNAME = "your_bot_username"
```

Then run the integration tests:

```powershell
go test ./bot -v -run TestNew_WithRealToken
```

### Creating a Test Bot

1. Open Telegram and search for [@BotFather](https://t.me/botfather)
2. Send `/newbot` command
3. Follow the instructions to create a new bot
4. Copy the bot token provided by BotFather
5. Use this token for testing (keep your production bot token separate!)

## Test Coverage

Current test coverage by package:

- `config`: 100% - All configuration loading scenarios tested
- `users`: 100% - All user management operations tested
- `bot`: ~34% - Core bot logic tested (some functions require API mocking)

### Generating Coverage Report

Generate an HTML coverage report:

```powershell
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

Open `coverage.html` in your browser to see detailed coverage.

## Test Structure

### Unit Tests

- `config/config_test.go` - Configuration loading tests
- `users/manager_test.go` - User management tests
- `bot/bot_test.go` - Bot logic tests

### What's Tested

✅ Configuration loading from environment variables
✅ User creation and updates
✅ User retrieval and listing
✅ Bot mention detection in messages
✅ Response generation
✅ Bot initialization with valid token (integration test)

### What's Not Fully Tested

⚠️ Actual message sending (requires API mocking or integration tests)
⚠️ Message handling workflows (skipped - require mocked Telegram API)

## Best Practices

1. **Never commit tokens** - Always use environment variables
2. **Use a separate test bot** - Don't use your production bot for testing
3. **Run tests before committing** - Ensure all tests pass
4. **Keep test tokens secure** - Don't share test bot tokens

## Running Specific Tests

Run tests for a specific package:

```powershell
go test ./bot -v
go test ./config -v
go test ./users -v
```

Run a specific test function:

```powershell
go test ./bot -v -run TestIsBotMentioned
```

## Continuous Integration

For CI/CD pipelines, set the following secrets:

- `TEST_BOT_TOKEN` - Test bot token
- `TEST_BOT_USERNAME` - Test bot username

The integration tests will automatically run if these are set, otherwise they'll be skipped.
