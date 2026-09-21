package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Receiver struct {
	bot *tgbotapi.BotAPI
}

func NewReceiver(token string) (*Receiver, error) {
	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		return nil, fmt.Errorf("failed to init telegram bot: %w", err)
	}

	return &Receiver{bot: bot}, nil
}

func (r *Receiver) Updates() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return r.bot.GetUpdatesChan(u)
}
