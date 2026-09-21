package main

import (
	"rabbitmq-telegram-bot/internal/config"
	"rabbitmq-telegram-bot/internal/consumer"
	"rabbitmq-telegram-bot/internal/telegram"

	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config error", "error", err)
		return
	}

	sender, err := telegram.NewSender(cfg.TelegramToken)

	if err != nil {
		slog.Error("error create sender", "error", err)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	handler := func(chatID int64, text string) error {
		return sender.Send(chatID, text)
	}

	slog.Info("starting consumer", "queue", cfg.QueueName)

	if err := consumer.Consume(ctx, cfg.RabbitMQURL, cfg.QueueName, handler); err != nil {
		slog.Error("rabbitMQ error", "error", err)
	}
}
