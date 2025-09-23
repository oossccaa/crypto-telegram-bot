package services

import (
	"fmt"
	"math"
	"time"

	"crypto-telegram-bot/models"
)

type MonitorService struct {
	coingecko    *CoinGeckoService
	telegram     *TelegramService
	previousPrices map[string]float64
}

func NewMonitorService(coingecko *CoinGeckoService, telegram *TelegramService) *MonitorService {
	return &MonitorService{
		coingecko:      coingecko,
		telegram:       telegram,
		previousPrices: make(map[string]float64),
	}
}

func (m *MonitorService) SendStartupMessage() error {
	coinIDs := []string{"cardano", "ethereum"}

	prices, err := m.coingecko.GetPrices(coinIDs)
	if err != nil {
		return fmt.Errorf("failed to get startup prices: %w", err)
	}

	if err := m.telegram.SendWelcomeMessage(prices); err != nil {
		return fmt.Errorf("failed to send startup message: %w", err)
	}

	// Initialize previous prices for monitoring (no alerts on startup)
	for _, price := range prices {
		m.previousPrices[price.ID] = price.Price
	}

	return nil
}

func (m *MonitorService) CheckPrices() error {
	coinIDs := []string{"cardano", "ethereum"}

	fmt.Printf("[DEBUG] Starting price check at %v\n", time.Now().Format("15:04:05"))

	prices, err := m.coingecko.GetPrices(coinIDs)
	if err != nil {
		fmt.Printf("[ERROR] Failed to get prices: %v\n", err)
		return fmt.Errorf("failed to get prices: %w", err)
	}

	fmt.Printf("[DEBUG] Got %d prices\n", len(prices))

	for _, price := range prices {
		if err := m.checkPriceChange(price); err != nil {
			fmt.Printf("Error checking price change for %s: %v\n", price.Symbol, err)
		}
	}

	return nil
}

func (m *MonitorService) checkPriceChange(price models.CoinPrice) error {
	previousPrice, exists := m.previousPrices[price.ID]

	if !exists {
		m.previousPrices[price.ID] = price.Price
		return nil
	}

	change := price.Price - previousPrice
	changePercent := (change / previousPrice) * 100

	// 調試日志：每次檢查都打印
	fmt.Printf("[DEBUG] %s: $%.6f -> $%.6f (%.3f%%)\n",
		price.Symbol, previousPrice, price.Price, changePercent)

	if math.Abs(changePercent) >= 0.2 { // 0.2% 觸發警報
		priceChange := &models.PriceChange{
			Coin:          price.Symbol,
			CurrentPrice:  price.Price,
			PreviousPrice: previousPrice,
			Change:        change,
			ChangePercent: changePercent,
		}

		if err := m.telegram.SendPriceAlert(priceChange); err != nil {
			return fmt.Errorf("failed to send alert: %w", err)
		}
	}

	m.previousPrices[price.ID] = price.Price
	return nil
}

func (m *MonitorService) HandleNewClient(chatID int64) error {
	coinIDs := []string{"cardano", "ethereum"}

	prices, err := m.coingecko.GetPrices(coinIDs)
	if err != nil {
		return fmt.Errorf("failed to get prices for new client: %w", err)
	}

	if err := m.telegram.SendWelcomeMessageToChat(chatID, prices); err != nil {
		return fmt.Errorf("failed to send welcome message to new client: %w", err)
	}

	return nil
}

func (m *MonitorService) StartListening() {
	m.telegram.StartMessageListener(m.HandleNewClient)
}