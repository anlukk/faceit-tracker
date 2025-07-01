package commands

import (
	"context"
	"fmt"
	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/anlukk/faceit-tracker/internal/db/models"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"strings"
	"time"
)

type Subscription struct {
	deps                       *core.Dependencies
	waitingForUsernameToAdd    map[telego.ChatID]bool
	waitingForUsernameToRemove map[telego.ChatID]bool
}

func NewSubscription(deps *core.Dependencies) *Subscription {
	return &Subscription{
		deps:                       deps,
		waitingForUsernameToAdd:    make(map[telego.ChatID]bool),
		waitingForUsernameToRemove: make(map[telego.ChatID]bool),
	}
}

func (s *Subscription) HandleButtonAdd(bot *telego.Bot, update telego.Update) {
	if update.CallbackQuery == nil || update.CallbackQuery.Message == nil {
		s.deps.Logger.Errorw("invalid callback query", "update", update)
		return
	}

	chatID := update.CallbackQuery.Message.GetChat().ID
	messageID := update.CallbackQuery.Message.GetMessageID()

	s.deps.Logger.Debugw("handle button add",
		"chat_id", chatID,
		"message_id", messageID,
		"callback_data", update.CallbackQuery.Data)

	userId := tu.ID(chatID)
	s.waitingForUsernameToAdd[userId] = true
	s.waitingForUsernameToRemove[userId] = false

	msg, err := bot.SendMessage(
		tu.Message(userId, s.deps.Messages.NicknameForSubs).
			WithReplyMarkup(tu.ForceReply()),
	)
	if err != nil {
		s.deps.Logger.Errorw("failed to send message", "error", err)
		return
	}

	s.deps.Logger.Debugw("sent message with force reply",
		"message_id", msg.MessageID)

	if err := bot.AnswerCallbackQuery(&telego.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	}); err != nil {
		s.deps.Logger.Errorw("failed to answer callback", "error", err)
	}
}

func (s *Subscription) HandleUserMessageFromAdd(bot *telego.Bot, update telego.Update) {
	if update.Message == nil || update.Message.From == nil {
		s.deps.Logger.Errorw("nil message or sender", "update", update)
		return
	}

	userId := tu.ID(update.Message.From.ID)
	if !s.waitingForUsernameToAdd[userId] {
		return
	}

	if update.Message.ReplyToMessage == nil {
		s.deps.Logger.Debugw("message is not a reply", "update", update)
		return
	}

	defer func() {
		s.waitingForUsernameToAdd[userId] = false
	}()

	userMessage := strings.TrimSpace(update.Message.Text)
	if userMessage == "" {
		_, err := bot.SendMessage(
			tu.Message(userId, "Please enter a valid username.").
				WithParseMode(telego.ModeHTML),
		)
		if err != nil {
			s.deps.Logger.Errorw("send message error", "error", err)
		}
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	playerId, err := s.deps.Faceit.GetUserIDByUsername(ctx, userMessage)
	if err != nil {
		s.deps.Logger.Errorw("failed to get user", "error", err)
		_, sendErr := bot.SendMessage(
			tu.Message(userId, "Error fetching data from FACEIT API.").
				WithParseMode(telego.ModeHTML),
		)
		if sendErr != nil {
			s.deps.Logger.Errorw("send message error", "error", sendErr)
		}
		return
	}

	s.deps.Logger.Debugw("get user", "id", userId)
	s.deps.Logger.Infof("Player ID: %s", playerId)

	chatID := update.Message.Chat.ID
	err = s.deps.Services.Sub.Subscribe(ctx, chatID, playerId, userMessage)
	if err != nil {
		s.deps.Logger.Errorw("failed to subscribe", "error", err)
		_, sendErr := bot.SendMessage(
			tu.Message(userId, "Failed to subscribe. Please try again.").
				WithParseMode(telego.ModeHTML),
		)
		if sendErr != nil {
			s.deps.Logger.Errorw("send message error", "error", sendErr)
		}
		return
	}

	_, botErr := bot.SendMessage(
		tu.Message(userId, fmt.Sprintf(s.deps.Messages.SuccessSubs, userMessage)).
			WithParseMode(telego.ModeHTML),
	)
	if botErr != nil {
		s.deps.Logger.Errorw("bot error", "error", botErr)
	}
}

func (s *Subscription) HandleButtonRemove(bot *telego.Bot, update telego.Update) {
	if update.CallbackQuery == nil || update.CallbackQuery.Message == nil {
		s.deps.Logger.Errorw("invalid callback query", "update", update)
		return
	}

	chatID := update.CallbackQuery.Message.GetChat().ID
	messageID := update.CallbackQuery.Message.GetMessageID()

	s.deps.Logger.Debugw("handle button remove", "chat_id", chatID, "message_id", messageID)
	userId := tu.ID(chatID)

	s.waitingForUsernameToRemove[userId] = true
	s.waitingForUsernameToAdd[userId] = false

	msg, err := bot.SendMessage(
		tu.Message(userId, s.deps.Messages.NicknameForUnsubs).WithReplyMarkup(tu.ForceReply()))
	if err != nil {
		s.deps.Logger.Errorw("failed to send message", "error", err)
		return
	}

	s.deps.Logger.Debugw("sent message with force reply", "message_id", msg.MessageID)

	if err := bot.AnswerCallbackQuery(&telego.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	}); err != nil {
		s.deps.Logger.Errorw("failed to answer callback", "error", err)
	}
}

func (s *Subscription) HandleUserMessageFromRemove(bot *telego.Bot, update telego.Update) {
	if update.Message == nil || update.Message.From == nil {
		s.deps.Logger.Errorw("nil message or sender", "update", update)
		return
	}

	chatID := update.Message.Chat.ID
	userId := tu.ID(chatID)

	if !s.waitingForUsernameToRemove[userId] {
		return
	}
	s.waitingForUsernameToRemove[userId] = false

	userMessage := strings.TrimSpace(update.Message.Text)
	if userMessage == "" {
		_, err := bot.SendMessage(
			tu.Message(userId, "Please enter a valid username.").
				WithParseMode(telego.ModeHTML),
		)
		if err != nil {
			s.deps.Logger.Errorw("send message error", "error", err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	playerId, err := s.deps.Faceit.GetUserIDByUsername(ctx, userMessage)
	if err != nil {
		s.deps.Logger.Errorw("failed to get user", "error", err)
	}

	subscribed, err := s.deps.Services.Sub.IsSubscribed(ctx, chatID, playerId)
	if err != nil {
		s.deps.Logger.Errorw("failed to check subscription", "error", err)
		_, sendErr := bot.SendMessage(
			tu.Message(userId, "Error checking subscription status."),
		)
		if sendErr != nil {
			s.deps.Logger.Errorw("send message error", "error", sendErr)
		}
		return
	}

	if !subscribed {
		_, sendErr := bot.SendMessage(
			tu.Message(userId, fmt.Sprintf("You are not subscribed to %s", userMessage)),
		)
		if sendErr != nil {
			s.deps.Logger.Errorw("send message error", "error", sendErr)
		}
		return
	}

	err = s.deps.Services.Sub.Unsubscribe(ctx, chatID, playerId)
	if err != nil {
		s.deps.Logger.Errorw("failed to unsubscribe", "error", err)
		_, sendErr := bot.SendMessage(
			tu.Message(userId, "Failed to unsubscribe. Please try again.").
				WithParseMode(telego.ModeHTML),
		)
		if sendErr != nil {
			s.deps.Logger.Errorw("send message error", "error", sendErr)
		}
		return
	}

	_, botErr := bot.SendMessage(
		tu.Message(userId, fmt.Sprintf(s.deps.Messages.SuccessUnsubs, userMessage)),
	)
	if botErr != nil {
		s.deps.Logger.Errorw("bot error", "error", botErr)
	}
}

func (s *Subscription) HandleButtonList(bot *telego.Bot, update telego.Update) {
	if update.CallbackQuery == nil || update.CallbackQuery.Message == nil {
		s.deps.Logger.Errorw("invalid callback query", "update", update)
		return
	}

	chatID := update.CallbackQuery.Message.GetChat().ID
	messageID := update.CallbackQuery.Message.GetMessageID()

	userId := tu.ID(chatID)

	s.deps.Logger.Debugw("handle button list", "chat_id", chatID, "message_id", messageID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	subs, err := s.deps.Services.Sub.GetSubscribers(ctx, chatID)
	if err != nil {
		s.deps.Logger.Errorw("failed to get subscriber", "error", err)
	}

	_, err = bot.SendMessage(tu.Message(userId, subsOutput(subs)).WithParseMode(telego.ModeHTML))
	if err != nil {
		s.deps.Logger.Errorw("failed to send message", "error", err)
	}
}

func subsOutput(subs []models.Subscription) string {
	sb := ""
	for _, sub := range subs {
		sb += fmt.Sprintf("%s\n", sub.Nickname)
	}

	return sb
}
