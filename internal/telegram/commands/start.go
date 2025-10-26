package commands

import (
	"context"
	"errors"
	"time"

	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/anlukk/faceit-tracker/internal/telegram/menu"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"gorm.io/gorm"
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

// TODO: Move everything related to subtitles to a separate file (separate the logic)
func (s *Start) HandleSubscriptionMenuCallback(bot *telego.Bot, update telego.Update) {
	msg := update.CallbackQuery.Message
	chatID := msg.GetChat().ID
	messageID := msg.GetMessageID()

	s.menu.SetActive(chatID, "options", messageID)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) //TODO: remove context
	defer cancel()

	personalSub, err := s.deps.PersonalSubRepo.GetPersonalSub(ctx, chatID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.deps.Logger.Errorw("failed to get personal sub", "error", err)
		return
	}

	subs, err := s.deps.SubscriptionRepo.GetSubscriptionsByChatID(ctx, chatID)
	if err != nil {
		s.deps.Logger.Errorw("failed to get subscriptions", "error", err)
		return
	}

	var mainNickname string
	if personalSub != nil {
		mainNickname = personalSub.Nickname
	} else {
		mainNickname = ""
	}

	_, err = bot.EditMessageText(&telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        s.deps.Messages.SubsCommandMessage,
		ReplyMarkup: BuildSubscriptionKeyboard(s.deps, subs, mainNickname),
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

func (s *Start) HandleSubscriptionToggleCallback(bot *telego.Bot, update telego.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) //TODO: remove context
	defer cancel()

	msg := update.CallbackQuery.Message
	chatID := msg.GetChat().ID
	personalSub, err := s.deps.PersonalSubRepo.GetPersonalSub(ctx, chatID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.deps.Logger.Errorw("failed to get personal sub", "error", err)
		return
	}

	subs, err := s.deps.SubscriptionRepo.GetSubscriptionsByChatID(ctx, chatID)
	if err != nil {
		s.deps.Logger.Errorw("failed to get subscriptions", "error", err)
		return
	}

	var mainNickname string
	if personalSub != nil {
		mainNickname = personalSub.Nickname
	} else {
		mainNickname = ""
	}

	messageID := msg.GetMessageID()
	s.menu.SetActive(chatID, "options", messageID)

	_, err = bot.EditMessageText(&telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        s.deps.Messages.SettingsCommandMessage,
		ReplyMarkup: BuildSubscriptionKeyboard(s.deps, subs, mainNickname),
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

func (s *Start) HandleSettingsMenuCallback(bot *telego.Bot, update telego.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) //TODO: remove context
	defer cancel()

	msg := update.CallbackQuery.Message

	chatID := msg.GetChat().ID
	current, err := s.deps.SettingsRepo.GetNotificationsEnabled(ctx, chatID)
	if err != nil {
		s.deps.Logger.Errorw("failed to get notifications status", "error", err)
		return
	}

	messageID := msg.GetMessageID()
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) //TODO: remove context
	defer cancel()

	msg := update.CallbackQuery.Message
	chatID := msg.GetChat().ID
	current, err := s.deps.SettingsRepo.GetNotificationsEnabled(ctx, chatID)
	if err != nil {
		s.deps.Logger.Errorw("failed to get notifications status", "error", err)
		return
	}

	newState := !current
	err = s.deps.SettingsRepo.SetNotificationsEnabled(ctx, chatID, newState)
	if err != nil {
		s.deps.Logger.Errorw("failed to set notifications status", "error", err)
		return
	}

	messageID := msg.GetMessageID()
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

	_, err := bot.EditMessageText(&telego.EditMessageTextParams{
		ChatID:      tu.ID(msg.GetChat().ID),
		MessageID:   msg.GetMessageID(),
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
