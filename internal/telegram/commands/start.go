package commands

import (
	"context"
	"errors"
	"time"

	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/anlukk/faceit-tracker/internal/telegram/menu"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	customCtxTimeOutForMenu = 5 * time.Second
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

func sendForceReply(bot *telego.Bot, chatID int64, text string, logger *zap.SugaredLogger) {
	if _, err := bot.SendMessage(tu.Message(tu.ID(chatID), text).
		WithReplyMarkup(tu.ForceReply()).
		WithParseMode(telego.ModeHTML)); err != nil {
		logger.Errorw("failed to send force reply", "error", err)
	}
}

func reply(bot *telego.Bot, chatID int64, text string, logger *zap.SugaredLogger) {
	if _, err := bot.SendMessage(tu.Message(tu.ID(chatID), text).
		WithParseMode(telego.ModeHTML)); err != nil {
		logger.Errorw("failed to send message", "error", err)
	}
}

func answerCallback(bot *telego.Bot, callbackID string, logger *zap.SugaredLogger) {
	if err := bot.AnswerCallbackQuery(tu.CallbackQuery(callbackID)); err != nil {
		logger.Errorw("failed to answer callback query", "callbackID", callbackID, "error", err)
	}
}

func getCallbackQueryChatAndMessageID(update telego.Update) (chatID int64, messageID int) {
	msg := update.CallbackQuery.Message
	return msg.GetChat().ID, msg.GetMessageID()
}

func editMessageText(
	bot *telego.Bot,
	keyboard *telego.InlineKeyboardMarkup,
	logger *zap.SugaredLogger,
	chatID int64,
	messageID int,
	text string) {

	if chatID == 0 || messageID == 0 || text == "" {
		return
	}

	if _, err := bot.EditMessageText(&telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        text,
		ReplyMarkup: keyboard,
		ParseMode:   telego.ModeHTML,
	}); err != nil {
		logger.Errorw("failed to edit message",
			"chatID", chatID,
			"messageID", messageID,
			"error", err)
	}
}

func (s *Start) StartCommand(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID
	msg, err := bot.SendMessage(tu.Message(tu.ID(chatID), s.deps.Messages.Description).
		WithReplyMarkup(BuildMainKeyboard(s.deps)).
		WithParseMode(telego.ModeHTML))
	if err != nil {
		s.deps.Logger.Errorw("failed to send start command message", "chatID", chatID, "error", err)
		return
	}

	s.menu.SetActive(chatID, "main", msg.MessageID)
}

func (s *Start) HandleSubscriptionMenuCallback(bot *telego.Bot, update telego.Update) {
	ctx, cancel := context.WithTimeout(s.deps.Ctx, customCtxTimeOutForMenu)
	defer cancel()

	chatID, messageID := getCallbackQueryChatAndMessageID(update)

	s.menu.SetActive(chatID, "options", messageID)

	personalSub, err := s.deps.PersonalSubRepo.GetPersonalSub(ctx, chatID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) { //TODO:
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

	editMessageText(bot,
		BuildSubscriptionKeyboard(s.deps, subs, mainNickname),
		s.deps.Logger,
		chatID,
		messageID,
		s.deps.Messages.SubsCommandMessage)

	answerCallback(bot, update.CallbackQuery.ID, s.deps.Logger)
}

func (s *Start) HandleSubscriptionToggleCallback(bot *telego.Bot, update telego.Update) {
	ctx, cancel := context.WithTimeout(s.deps.Ctx, customCtxTimeOutForMenu)
	defer cancel()

	chatID, messageID := getCallbackQueryChatAndMessageID(update)
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

	s.menu.SetActive(chatID, "options", messageID)

	editMessageText(bot,
		BuildSubscriptionKeyboard(s.deps, subs, mainNickname),
		s.deps.Logger,
		chatID,
		messageID,
		s.deps.Messages.SubsCommandMessage)

	answerCallback(bot, update.CallbackQuery.ID, s.deps.Logger)
}

func (s *Start) HandleSettingsMenuCallback(bot *telego.Bot, update telego.Update) {
	ctx, cancel := context.WithTimeout(s.deps.Ctx, customCtxTimeOutForMenu)
	defer cancel()

	chatID, messageID := getCallbackQueryChatAndMessageID(update)
	current, err := s.deps.SettingsRepo.GetNotificationsEnabled(ctx, chatID)
	if err != nil {
		s.deps.Logger.Errorw("failed to get notifications status", "error", err)
		return
	}

	s.menu.SetActive(chatID, "options", messageID)

	editMessageText(bot,
		BuildSettingsKeyboard(s.deps, current),
		s.deps.Logger,
		chatID,
		messageID,
		s.deps.Messages.SettingsCommandMessage)

	answerCallback(bot, update.CallbackQuery.ID, s.deps.Logger)
}

func (s *Start) HandleNotificationToggleCallback(bot *telego.Bot, update telego.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), customCtxTimeOutForMenu)
	defer cancel()

	chatID, messageID := getCallbackQueryChatAndMessageID(update)

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

	s.menu.SetActive(chatID, "options", messageID)

	editMessageText(bot,
		BuildSettingsKeyboard(s.deps, newState),
		s.deps.Logger,
		chatID,
		messageID,
		s.deps.Messages.SettingsCommandMessage)

	answerCallback(bot, update.CallbackQuery.ID, s.deps.Logger)
}

func (s *Start) HandleBackCallback(bot *telego.Bot, update telego.Update) {
	editMessageText(bot,
		BuildMainKeyboard(s.deps),
		s.deps.Logger,
		update.CallbackQuery.Message.GetChat().ID,
		update.CallbackQuery.Message.GetMessageID(),
		s.deps.Messages.Description)

	answerCallback(bot, update.CallbackQuery.ID, s.deps.Logger)
}
