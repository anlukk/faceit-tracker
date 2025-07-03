package commands

import (
	"context"
	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/anlukk/faceit-tracker/internal/telegram/menu"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"time"
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
	msg, err := bot.SendMessage(tu.Message(tu.ID(chatID), s.deps.Messages.Description).
		WithReplyMarkup(BuildMainKeyboard(s.deps)).
		WithParseMode(telego.ModeHTML))
	if err != nil {
		s.deps.Logger.Errorw("bot error", "error", err)
		return
	}

	s.menu.SetActive(chatID, "main", msg.MessageID)
}

func (s *Start) HandleSubscriptionMenuCallback(bot *telego.Bot, update telego.Update) {
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

	s.deps.Logger.Debugw(
		"handle options callback",
		"chat_id", chatID,
		"message_id", messageID,
	)

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

func (s *Start) HandleSettingsMenuCallback(bot *telego.Bot, update telego.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	msg := update.CallbackQuery.Message
	if msg == nil {
		s.deps.Logger.Errorw("message is nil")
		return
	}

	chatID := msg.GetChat().ID
	messageID := msg.GetMessageID()

	current, err := s.deps.Services.Settings.GetNotificationsEnabled(ctx, chatID)
	if err != nil {
		s.deps.Logger.Errorw("failed to get notifications status", "error", err)
		return
	}

	s.menu.SetActive(chatID, "options", messageID)

	_, err = bot.EditMessageText(&telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        s.deps.Messages.SettingsCommandMessage,
		ReplyMarkup: BuildSettingsKeyboard(s.deps, current),
		ParseMode:   telego.ModeHTML,
	})
	if err != nil {
		s.deps.Logger.Errorw("failed to edit message", "error", err)
	}

	err = bot.AnswerCallbackQuery(tu.CallbackQuery(update.CallbackQuery.ID))
	if err != nil {
		s.deps.Logger.Errorw("bot error", "error", err)
	}
}

func (s *Start) HandleNotificationToggleCallback(bot *telego.Bot, update telego.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	msg := update.CallbackQuery.Message
	if msg == nil {
		s.deps.Logger.Errorw("message is nil")
		return
	}

	chatID := msg.GetChat().ID
	messageID := msg.GetMessageID()

	current, err := s.deps.Services.Settings.GetNotificationsEnabled(ctx, chatID)
	if err != nil {
		s.deps.Logger.Errorw("failed to get notifications status", "error", err)
		return
	}

	newState := !current
	err = s.deps.Services.Settings.SetNotificationsEnabled(ctx, chatID, newState)
	if err != nil {
		s.deps.Logger.Errorw("failed to set notifications status", "error", err)
		return
	}

	s.menu.SetActive(chatID, "options", messageID)

	_, err = bot.EditMessageText(&telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        s.deps.Messages.SettingsCommandMessage,
		ReplyMarkup: BuildSettingsKeyboard(s.deps, newState),
		ParseMode:   telego.ModeHTML,
	})
	if err != nil {
		s.deps.Logger.Errorw("failed to edit message", "error", err)
	}

	err = bot.AnswerCallbackQuery(tu.CallbackQuery(update.CallbackQuery.ID))
	if err != nil {
		s.deps.Logger.Errorw("bot error", "error", err)
	}
}

func (s *Start) HandleBackCallback(bot *telego.Bot, update telego.Update) {
	msg := update.CallbackQuery.Message
	if msg == nil {
		s.deps.Logger.Errorw("message is nil")
		return
	}

	chatID := msg.GetChat().ID
	messageID := msg.GetMessageID()

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
