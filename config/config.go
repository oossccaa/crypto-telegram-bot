package config

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	BotToken string `yaml:"bot_token"`
	ChatID   int64  `yaml:"chat_id"`
}

func Load() (*Config, error) {
	// Try to load from environment variables first (for Railway deployment)
	botToken := os.Getenv("BOT_TOKEN")
	chatIDStr := os.Getenv("CHAT_ID")

	if botToken != "" {
		if chatIDStr == "" {
			return nil, os.ErrNotExist
		}

		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil {
			return nil, err
		}

		return &Config{
			BotToken: botToken,
			ChatID:   chatID,
		}, nil
	}

	// Fallback to config.yaml file
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}