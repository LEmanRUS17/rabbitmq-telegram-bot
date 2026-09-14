package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
)

type payload struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func Consume(ctx context.Context, url string, queueName string, handler func(chatID int64, text string) error) error {

	conn, err := amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to connect to broker: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare the queue: %w", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to subscribe to the queue: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			slog.Info("shutdown signal received")
			return nil
		case msg, ok := <-msgs:
			if !ok {
				slog.Info("stop consume")
				return nil
			}

			var p payload

			if err := json.Unmarshal(msg.Body, &p); err != nil {
				slog.Warn("skip malformed message", "error", err)
				continue
			}

			if err := handler(p.ChatID, p.Text); err != nil {
				slog.Error("handler failed", "error", err)
			}
		}
	}
}
