package main

import (
      "rabbitmq-telegram-bot/internal/config"
      "rabbitmq-telegram-bot/internal/publisher"
      "rabbitmq-telegram-bot/internal/telegram"

      "context"
      "encoding/json"
      "log/slog"
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

	pub, err := publisher.NewPublisher(cfg.RabbitMQURL, cfg.UpdatesQueueName)

	if err != nil {
		slog.Error("error create publisher", "error", err)
		return
	}

	defer pub.Close()

	slog.Info("starting poller", "queue", cfg.UpdatesQueueName)

	for update := range receiver.Updates() {
		body, err := json.Marshal(update)
		
		if err != nil {
			slog.Warn("skip update: marshal failed", "error", err)
			continue
		}
	
		if err := pub.Publish(context.Background(), body); err != nil {
			slog.Error("publish failed", "error", err)
		}
	}
}
