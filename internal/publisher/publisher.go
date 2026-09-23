package publisher

import (
	"context"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	conn      *amqp.Connection
	ch        *amqp.Channel
	queueName string
}

func NewPublisher(url string, queueName string) (*Publisher, error) {
	conn, err := amqp.Dial(url)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to broker: %w", err)
	}

	ch, err := conn.Channel()

	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)

	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare the queue: %w", err)
	}

	return &Publisher{conn: conn, ch: ch, queueName: q.Name}, nil
}

func (p *Publisher) Publish(ctx context.Context, body []byte) error {
	err := p.ch.PublishWithContext(
		ctx,
		"",
		p.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	return nil
}

func (p *Publisher) Close() error {
	return errors.Join(p.ch.Close(), p.conn.Close())
}
