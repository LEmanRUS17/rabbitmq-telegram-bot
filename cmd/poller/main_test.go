package main

import (
	"encoding/json"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestUpdateJSONContract(t *testing.T) {
	update := tgbotapi.Update{
		UpdateID: 1,
		Message: &tgbotapi.Message{
			MessageID: 2,
			Text:      "hello",
			Chat:      &tgbotapi.Chat{ID: 42},
		},
	}

	body, err := json.Marshal(update)

	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var decoded map[string]any

	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if decoded["update_id"] != float64(1) {
		t.Errorf("update_id = %v, want 1", decoded["update_id"])
	}

	msg, ok := decoded["message"].(map[string]any)

	if !ok {
		t.Fatalf("message field missing or wrong type")
	}

	if msg["text"] != "hello" {
		t.Errorf("message.text = %v, want %q", msg["text"], "hello")
	}

	chat, ok := msg["chat"].(map[string]any)

	if !ok {
		t.Fatalf("chat field missing or wrong type")
	}

	if chat["id"] != float64(42) {
		t.Errorf("message.chat.id = %v, want 42", chat["id"])
	}
}
