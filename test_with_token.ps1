# Script to run tests with a real bot token
# Usage: .\test_with_token.ps1

Write-Host "HowardTheChad Bot - Integration Testing" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Check if token is already set
if (-not $env:TEST_BOT_TOKEN) {
    Write-Host "Please enter your test bot token:" -ForegroundColor Yellow
    $token = Read-Host -AsSecureString
    $BSTR = [System.Runtime.InteropServices.Marshal]::SecureStringToBSTR($token)
    $env:TEST_BOT_TOKEN = [System.Runtime.InteropServices.Marshal]::PtrToStringAuto($BSTR)
}

if (-not $env:TEST_BOT_USERNAME) {
    Write-Host "Please enter your bot username (without @):" -ForegroundColor Yellow
    $env:TEST_BOT_USERNAME = Read-Host
}

Write-Host ""
Write-Host "Running all tests..." -ForegroundColor Green
go test ./... -v

Write-Host ""
Write-Host "Running integration tests with real token..." -ForegroundColor Green
go test ./bot -v -run TestNew_WithRealToken

Write-Host ""
Write-Host "Test Coverage:" -ForegroundColor Cyan
go test ./... -cover

# Clean up (optional - comment out if you want to keep the token in session)
# $env:TEST_BOT_TOKEN = $null
# $env:TEST_BOT_USERNAME = $null
