package commands

import (
	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/anlukk/faceit-tracker/internal/telegram/menu"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type Start struct {
	deps *core.Dependencies
	menu *menu.MenuManager
}

func NewStart(deps *core.Dependencies, menu *menu.MenuManager) *Start {
	return &Start{
		menu: menu,
		deps: deps,
	}
}

func (s *Start) StartCommand(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID
	keyboard := BuildMainKeyboard(s.deps)
	msg, err := bot.SendMessage(tu.Message(tu.ID(chatID), s.deps.Messages.Description).
		WithReplyMarkup(keyboard).
		WithParseMode(telego.ModeHTML))
	if err != nil {
		s.deps.Logger.Errorw("bot error", "error", err)
		return
	}

	if msg.MessageID == 0 {
		s.deps.Logger.Errorw("message id is 0", "chat_id", chatID)
		return
	}

	s.menu.SetActive(chatID, "main", msg.MessageID)
}

func (s *Start) HandleOptionsCallback(bot *telego.Bot, update telego.Update) {
	if update.CallbackQuery == nil {
		s.deps.Logger.Errorw("callback query is nil")
		return
	}

	msg := update.CallbackQuery.Message
	if msg == nil {
		s.deps.Logger.Errorw("message is nil")
		return
	}

	chatID := msg.GetChat().ID
	messageID := msg.GetMessageID()

	s.deps.Logger.Debugw("handle options callback", "chat_id", chatID, "message_id", messageID)

	s.menu.SetActive(chatID, "options", messageID)

	_, err := bot.EditMessageText(&telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        s.deps.Messages.SubsCommandMessage,
		ReplyMarkup: BuildSubscriptionKeyboard(s.deps),
		ParseMode:   telego.ModeHTML,
	})

	if err != nil {
		s.deps.Logger.Errorw("bot error", "error", err)
	}

	err = bot.AnswerCallbackQuery(tu.CallbackQuery(update.CallbackQuery.ID))
	if err != nil {
		s.deps.Logger.Errorw("bot error", "error", err)
	}
}

func (s *Start) HandleSettingsCallback(bot *telego.Bot, update telego.Update) {
	if update.CallbackQuery == nil {
		s.deps.Logger.Errorw("callback query is nil")
		return
	}

	msg := update.CallbackQuery.Message
	if msg == nil {
		s.deps.Logger.Errorw("message is nil")
		return
	}

	chatID := msg.GetChat().ID
	messageID := msg.GetMessageID()

	s.deps.Logger.Debugw("handle settings callback", "chat_id", chatID, "message_id", messageID)

	s.menu.SetActive(chatID, "options", messageID)

	_, err := bot.EditMessageText(&telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        "Settings Menu",
		ReplyMarkup: BuildSettingsKeyboard(s.deps),
		ParseMode:   telego.ModeHTML,
	})

	if err != nil {
		s.deps.Logger.Errorw("bot error", "error", err)
	}

	err = bot.AnswerCallbackQuery(tu.CallbackQuery(update.CallbackQuery.ID))
	if err != nil {
		s.deps.Logger.Errorw("bot error", "error", err)
	}
}

func (s *Start) HandleBackCallback(bot *telego.Bot, update telego.Update) {
	if update.CallbackQuery == nil {
		s.deps.Logger.Errorw("callback query is nil")
		return
	}

	chatID := update.CallbackQuery.Message.GetChat().ID
	messageID := update.CallbackQuery.Message.GetMessageID()

	_, err := bot.EditMessageText(&telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        s.deps.Messages.Description,
		ReplyMarkup: BuildMainKeyboard(s.deps),
	})
	if err != nil {
		s.deps.Logger.Errorw("bot error", "error", err)
	}

	err = bot.AnswerCallbackQuery(tu.CallbackQuery(update.CallbackQuery.ID))
	if err != nil {
		s.deps.Logger.Errorw("bot error", "error", err)
	}
}
