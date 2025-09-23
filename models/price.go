package models

type CoinPrice struct {
	ID       string  `json:"id"`
	Symbol   string  `json:"symbol"`
	Name     string  `json:"name"`
	Price    float64 `json:"current_price"`
	Change24h float64 `json:"price_change_percentage_24h"`
}

type PriceChange struct {
	Coin          string
	CurrentPrice  float64
	PreviousPrice float64
	Change        float64
	ChangePercent float64
}