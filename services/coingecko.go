package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"crypto-telegram-bot/models"
)

type CoinGeckoService struct {
	baseURL string
}

func NewCoinGeckoService() *CoinGeckoService {
	return &CoinGeckoService{
		baseURL: "https://api.coingecko.com/api/v3",
	}
}

func (c *CoinGeckoService) GetPrices(coinIDs []string) ([]models.CoinPrice, error) {
	url := fmt.Sprintf("%s/coins/markets?vs_currency=usd&ids=%s&order=market_cap_desc&per_page=250&page=1&sparkline=false&price_change_percentage=24h",
		c.baseURL, joinStrings(coinIDs, ","))

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get prices: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var prices []models.CoinPrice
	if err := json.NewDecoder(resp.Body).Decode(&prices); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return prices, nil
}

func joinStrings(strings []string, sep string) string {
	if len(strings) == 0 {
		return ""
	}

	result := strings[0]
	for i := 1; i < len(strings); i++ {
		result += sep + strings[i]
	}
	return result
}