package main

import (
	"rabbitmq-telegram-bot/internal/config"
	"rabbitmq-telegram-bot/internal/publisher"
	"rabbitmq-telegram-bot/internal/telegram"

	"context"
	"encoding/json"
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

	receiver, err := telegram.NewReceiver(cfg.TelegramToken)

	if err != nil {
		slog.Error("error create receiver", "error", err)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pub, err := publisher.NewPublisher(cfg.RabbitMQURL, cfg.UpdatesQueueName)

	if err != nil {
		slog.Error("error create publisher", "error", err)
		return
	}

	defer pub.Close()

	slog.Info("starting poller", "queue", cfg.UpdatesQueueName)

	updates := receiver.Updates()

	for {
		select {
		case <-ctx.Done():
			slog.Info("shutdown signal received")
			receiver.Stop()
			return

		case update, ok := <-updates:
			if !ok {
				slog.Info("stop poller")
				return
			}

			body, err := json.Marshal(update)

			if err != nil {
				slog.Warn("skip update: marshal failed", "error", err)
				continue
			}

			if err := pub.Publish(ctx, body); err != nil {
				slog.Error("publish failed", "error", err)
				panic(err)
			}
		}
	}
}
