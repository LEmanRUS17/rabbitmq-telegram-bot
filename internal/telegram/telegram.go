package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Sender struct {
	bot *tgbotapi.BotAPI
}

func NewSender(token string) (*Sender, error) {
	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		return nil, fmt.Errorf("failed to init telegram bot: %w", err)
	}
	return &Sender{bot: bot}, nil
}

func (s *Sender) Send(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := s.bot.Send(msg)
	return err
}
