# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Commands

### Development
- `go mod tidy` - Install and clean up dependencies
- `go run main.go` - Run the bot directly
- `go build -o crypto-bot` - Build executable binary
- `./crypto-bot` - Run the compiled binary

### Docker Deployment
- `docker-compose up -d` - Start bot in background
- `docker-compose down` - Stop bot
- `docker-compose logs -f` - View live logs
- `docker-compose restart` - Restart bot

### Testing
- `go test ./...` - Run all tests in the project
- `go test -v ./services/...` - Run tests for services package with verbose output

## Architecture Overview

This is a Go-based cryptocurrency price monitoring Telegram bot with a clean modular architecture:

### Core Components
- **main.go**: Entry point that orchestrates services and runs the monitoring loop (10-minute intervals)
- **config/**: YAML-based configuration management for Telegram credentials
- **models/**: Data structures for CoinPrice and PriceChange
- **services/**: Business logic layer with three key services

### Service Layer Architecture
The application follows a service-oriented pattern:

1. **CoinGeckoService** (`services/coingecko.go`): Handles external API calls to CoinGecko for price data
2. **TelegramService** (`services/telegram.go`): Manages Telegram bot messaging and formatting
3. **MonitorService** (`services/monitor.go`): Core business logic that coordinates price checking, comparison, and alerting

### Data Flow
1. MonitorService triggers price checks every 10 minutes
2. Fetches current prices for ADA and ETH via CoinGeckoService
3. Compares against stored previous prices
4. If price change ≥ 0.5%, sends formatted alert via TelegramService
5. Updates stored prices for next comparison

### Configuration
- Requires `config.yaml` file with `bot_token` and `chat_id`
- Copy from `config.yaml.example` and fill in credentials
- Bot token from @BotFather, Chat ID from `/getUpdates` API endpoint
- For Docker: Place `config.yaml` in project root (mounted as read-only)

### Dependencies
- `github.com/go-telegram-bot-api/telegram-bot-api/v5`: Telegram Bot API client
- `gopkg.in/yaml.v3`: YAML configuration parsing

### Key Behavior
- Monitors Cardano (ADA) and Ethereum (ETH) only
- Triggers alerts on ±0.5% price changes
- Uses markdown formatting for Telegram messages
- Stores previous prices in memory (resets on restart)