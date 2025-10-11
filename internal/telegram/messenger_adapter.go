package telegram

import (
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type MessengerAdapter struct {
	bot *telego.Bot
}

func NewMessengerAdapter(bot *telego.Bot) *MessengerAdapter {
	return &MessengerAdapter{bot: bot}
}

func (m *MessengerAdapter) SendMessage(chatID int64, text string) error {
	_, err := m.bot.SendMessage(&telego.SendMessageParams{
		ChatID: tu.ID(chatID),
		Text:   text,
	})
	return err
}
