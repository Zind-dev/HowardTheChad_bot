# HowardTheChad Bot Launcher
# This script sets environment variables and launches the bot

Write-Host "HowardTheChad Bot Launcher" -ForegroundColor Cyan
Write-Host "==============================" -ForegroundColor Cyan
Write-Host ""

# Add MinGW to PATH for CGO support (required for SQLite)
$env:PATH = "C:\mingw64\bin;$env:PATH"
$env:CGO_ENABLED = "1"

# Check if .env file exists and load from it first (as defaults)
if (Test-Path ".env") {
    Write-Host "Loading default values from .env file..." -ForegroundColor Yellow
    Get-Content .env | ForEach-Object {
        if ($_ -match '^([^#][^=]+)=(.+)$') {
            $name = $matches[1].Trim()
            $value = $matches[2].Trim()
            Set-Item -Path "env:$name" -Value $value
        }
    }
}

# Then check Windows User environment variables (these override .env)
Write-Host "Checking Windows environment variables..." -ForegroundColor Yellow
$userToken = [System.Environment]::GetEnvironmentVariable("TELEGRAM_BOT_TOKEN", "User")
$userUsername = [System.Environment]::GetEnvironmentVariable("BOT_USERNAME", "User")

if ($userToken) {
    $env:TELEGRAM_BOT_TOKEN = $userToken
    Write-Host "[OK] Using TELEGRAM_BOT_TOKEN from Windows User variables (overriding .env)" -ForegroundColor Green
}

if ($userUsername) {
    $env:BOT_USERNAME = $userUsername
    Write-Host "[OK] Using BOT_USERNAME from Windows User variables (overriding .env)" -ForegroundColor Green
}

# Check if variables are set, if not prompt for them
if (-not $env:TELEGRAM_BOT_TOKEN) {
    Write-Host "TELEGRAM_BOT_TOKEN not found." -ForegroundColor Yellow
    $token = Read-Host "Enter your Telegram Bot Token"
    $env:TELEGRAM_BOT_TOKEN = $token
}

if (-not $env:BOT_USERNAME) {
    Write-Host "BOT_USERNAME not found." -ForegroundColor Yellow
    $username = Read-Host "Enter your Bot Username (without @)"
    $env:BOT_USERNAME = $username
}

# Optional: Set response frequency if not set (default will be used otherwise)
if (-not $env:BOT_RESPONSE_FREQUENCY) {
    $env:BOT_RESPONSE_FREQUENCY = "10"
}

# Optional: Set respond to mentions if not set (default will be used otherwise)
if (-not $env:BOT_RESPOND_TO_MENTIONS) {
    $env:BOT_RESPOND_TO_MENTIONS = "true"
}

Write-Host ""
Write-Host "Starting bot with settings:" -ForegroundColor Green
Write-Host "  Bot Username: $env:BOT_USERNAME"
Write-Host "  Response Frequency: $env:BOT_RESPONSE_FREQUENCY"
Write-Host "  Respond to Mentions: $env:BOT_RESPOND_TO_MENTIONS"
Write-Host ""

# Launch the bot
.\HowardTheChad_bot.exe
