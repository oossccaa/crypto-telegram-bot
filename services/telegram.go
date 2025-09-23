package services

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"crypto-telegram-bot/models"
)

type TelegramService struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

func NewTelegramService(token string, chatID int64) *TelegramService {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(fmt.Sprintf("Failed to create Telegram bot: %v", err))
	}

	return &TelegramService{
		bot:    bot,
		chatID: chatID,
	}
}

func (t *TelegramService) SendPriceAlert(change *models.PriceChange) error {
	emoji := "📈"
	if change.ChangePercent < 0 {
		emoji = "📉"
	}

	message := fmt.Sprintf(
		"%s *%s Price Alert*\n\n"+
			"💰 Current Price: $%.6f\n"+
			"📊 Previous Price: $%.6f\n"+
			"📈 Change: %.2f%% (%.6f)\n"+
			"⏰ Time: %s",
		emoji,
		change.Coin,
		change.CurrentPrice,
		change.PreviousPrice,
		change.ChangePercent,
		change.Change,
		getCurrentTime(),
	)

	msg := tgbotapi.NewMessage(t.chatID, message)
	msg.ParseMode = "Markdown"

	_, err := t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (t *TelegramService) SendWelcomeMessage(currentPrices []models.CoinPrice) error {
	message := "🚀 *Crypto Telegram Bot Started*\n\n"
	message += "✅ System is now online and monitoring crypto prices\n\n"
	message += "*Current Prices:*\n"

	for _, price := range currentPrices {
		emoji := "🟢"
		if price.Change24h < 0 {
			emoji = "🔴"
		}

		message += fmt.Sprintf(
			"%s *%s*: $%.6f (%.2f%%)\n",
			emoji,
			price.Symbol,
			price.Price,
			price.Change24h,
		)
	}

	message += fmt.Sprintf("\n⏰ Started at: %s", getCurrentTime())

	msg := tgbotapi.NewMessage(t.chatID, message)
	msg.ParseMode = "Markdown"

	_, err := t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send welcome message: %w", err)
	}

	return nil
}

func (t *TelegramService) StartMessageListener(onNewClient func(int64) error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			// Check if this is a /start command or new client
			if update.Message.IsCommand() && update.Message.Command() == "start" {
				log.Printf("New client connected: %d", update.Message.Chat.ID)

				// Send welcome message to this specific client
				if err := onNewClient(update.Message.Chat.ID); err != nil {
					log.Printf("Error sending welcome message to new client %d: %v", update.Message.Chat.ID, err)
				}
			}
		}
	}
}

func (t *TelegramService) SendWelcomeMessageToChat(chatID int64, currentPrices []models.CoinPrice) error {
	message := "🚀 *Welcome to Crypto Price Monitor*\n\n"
	message += "✅ You are now subscribed to crypto price alerts\n\n"
	message += "*Current Prices:*\n"

	for _, price := range currentPrices {
		emoji := "🟢"
		if price.Change24h < 0 {
			emoji = "🔴"
		}

		message += fmt.Sprintf(
			"%s *%s*: $%.6f (%.2f%%)\n",
			emoji,
			price.Symbol,
			price.Price,
			price.Change24h,
		)
	}

	message += fmt.Sprintf("\n⏰ Connected at: %s", getCurrentTime())

	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "Markdown"

	_, err := t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send welcome message to chat %d: %w", chatID, err)
	}

	return nil
}

func getCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}