package config

import "testing"

func TestLoad_MissingToken(t *testing.T) {
	t.Setenv("TELEGRAM_TOKEN", "")

	_, err := Load()

	if err == nil {
		t.Errorf("expected error when TELEGRAM_TOKEN is missing? got nil")
	}
}

func TestLoad_Success(t *testing.T) {
	t.Setenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	t.Setenv("QUEUE_NAME", "test_queue")
	t.Setenv("UPDATES_QUEUE_NAME", "test_updates_queue")
	t.Setenv("TELEGRAM_TOKEN", "fake-token")

	cfg, err := Load()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cfg.RabbitMQURL != "amqp://guest:guest@localhost:5672/" {
		t.Errorf("RabbitMQURL = %q, want %q", cfg.RabbitMQURL, "amqp://guest:guest@localhost:5672/")
	}

	if cfg.QueueName != "test_queue" {
		t.Errorf("QueueName = %q, want %q", cfg.QueueName, "test_queue")
	}

	if cfg.UpdatesQueueName != "test_updates_queue" {
		t.Errorf("UpdatesQueueName = %q, want %q", cfg.UpdatesQueueName, "test_updates_queue")
	}

	if cfg.TelegramToken != "fake-token" {
		t.Errorf("TelegramToken = %q, want %q", cfg.TelegramToken, "fake-token")
	}
}
