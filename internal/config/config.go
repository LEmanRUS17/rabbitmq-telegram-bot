package config

import (
	"errors"
	"os"
)

// settings bot
type Config struct {
	RabbitMQURL   string
	QueueName     string
	TelegramToken string
}

func Load() (*Config, error) {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	token := os.Getenv("TELEGRAM_TOKEN")
	queueName := os.Getenv("QUEUE_NAME")

	if token == "" {
		return nil, errors.New("TELEGRAM_TOKEN not found")
	}

	cfg := &Config{
		RabbitMQURL:   rabbitURL,
		QueueName:     queueName,
		TelegramToken: token,
	}

	return cfg, nil
}
