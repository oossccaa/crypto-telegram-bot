package main

import (
	"fmt"
	"log"
	"time"

	"crypto-telegram-bot/config"
	"crypto-telegram-bot/services"
)

func main() {
	fmt.Println("Starting Crypto Telegram Bot...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	coingeckoService := services.NewCoinGeckoService()
	telegramService := services.NewTelegramService(cfg.BotToken, cfg.ChatID)
	monitorService := services.NewMonitorService(coingeckoService, telegramService)

	// Send startup message with current prices
	fmt.Println("Sending startup message with current prices...")
	if err := monitorService.SendStartupMessage(); err != nil {
		log.Printf("Error sending startup message: %v", err)
	}

	// Start listening for new client connections in a separate goroutine
	go func() {
		fmt.Println("Starting message listener for new clients...")
		monitorService.StartListening()
	}()

	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	fmt.Println("Bot started successfully. Monitoring every 10 minutes...")

	for range ticker.C {
		if err := monitorService.CheckPrices(); err != nil {
			log.Printf("Error checking prices: %v", err)
		}
	}
}