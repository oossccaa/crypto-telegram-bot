package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	BotToken string `yaml:"bot_token"`
	ChatID   int64  `yaml:"chat_id"`
}

func Load() (*Config, error) {
	fmt.Println("Loading configuration...")

	// Try to load from environment variables first (for Railway deployment)
	botToken := os.Getenv("BOT_TOKEN")
	chatIDStr := os.Getenv("CHAT_ID")

	fmt.Printf("Environment variables - BOT_TOKEN exists: %t, CHAT_ID exists: %t\n",
		botToken != "", chatIDStr != "")

	if botToken != "" {
		if chatIDStr == "" {
			fmt.Println("BOT_TOKEN found but CHAT_ID is missing")
			return nil, fmt.Errorf("CHAT_ID environment variable is required when BOT_TOKEN is set")
		}

		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil {
			fmt.Printf("Failed to parse CHAT_ID '%s': %v\n", chatIDStr, err)
			return nil, fmt.Errorf("invalid CHAT_ID format: %v", err)
		}

		fmt.Println("Successfully loaded configuration from environment variables")
		return &Config{
			BotToken: botToken,
			ChatID:   chatID,
		}, nil
	}

	// Fallback to config.yaml file
	fmt.Println("Environment variables not found, trying config.yaml file...")
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		fmt.Printf("Failed to read config.yaml: %v\n", err)
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Printf("Failed to parse config.yaml: %v\n", err)
		return nil, fmt.Errorf("failed to parse config.yaml: %v", err)
	}

	fmt.Println("Successfully loaded configuration from config.yaml")
	return &config, nil
}