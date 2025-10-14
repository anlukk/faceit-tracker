package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

const (
	timeout = 10 * time.Second
)

type Subscription struct {
	deps *core.Dependencies
}

func NewSubscription(deps *core.Dependencies) *Subscription {
	return &Subscription{
		deps: deps,
	}
}

func (s *Subscription) sendForceReply(bot *telego.Bot, chatID int64, text string) error {
	_, err := bot.SendMessage(tu.Message(tu.ID(chatID), text).WithReplyMarkup(tu.ForceReply()))
	return err
}

func (s *Subscription) HandleSubscribeButton(bot *telego.Bot, update telego.Update) {
	msg := update.CallbackQuery.Message
	if msg == nil {
		s.deps.Logger.Errorw("message is nil")
		return
	}

	chatID := msg.GetChat().ID
	messageID := msg.GetMessageID()

	s.deps.Logger.Debugw("handle button add",
		"chat_id", chatID,
		"message_id", messageID,
		"callback_data", update.CallbackQuery.Data)

	if err := s.sendForceReply(bot, chatID, s.deps.Messages.NicknameForSubs); err != nil {
		s.deps.Logger.Errorw("failed to send message", "error", err)
		return
	}

	s.deps.Logger.Debugw("sent message with force reply",
		"message_id", messageID)

	if err := bot.AnswerCallbackQuery(&telego.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	}); err != nil {
		s.deps.Logger.Errorw("failed to answer callback", "error", err)
	}
}

func (s *Subscription) HandleSubscriptionNicknameReply(bot *telego.Bot, update telego.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	userId := tu.ID(update.Message.From.ID)
	userMessage := strings.TrimSpace(update.Message.Text)
	if userMessage == "" {
		_, err := bot.SendMessage(tu.Message(userId, "Please enter a valid nickname.").
			WithParseMode(telego.ModeHTML),
		)
		if err != nil {
			s.deps.Logger.Errorw("send message error", "error", err)
		}
		return
	}

	playerId, err := s.deps.Faceit.GetPlayerIDByNickname(ctx, userMessage)
	if err != nil {
		s.deps.Logger.Errorw("failed to get user", "error", err)
		_, sendErr := bot.SendMessage(tu.Message(userId, "Error fetching data from FACEIT API.").
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

	exists, err := s.deps.SubscriptionRepo.IsSubscribed(ctx, chatID, playerId)
	if err != nil {
		s.deps.Logger.Errorw("check subscription failed: %w", err)
	}

	if exists {
		//_, sendErr := bot.SendMessage(tu.Message(userId,
		//	fmt.Sprintf(s.deps.Messages.AlreadySubscribed, userMessage)).
		//	WithParseMode(telego.ModeHTML),
		//)
	}

	err = s.deps.SubscriptionRepo.Subscribe(ctx, chatID, playerId, userMessage)
	if err != nil {
		s.deps.Logger.Errorw("failed to subscribe", "error", err)
		_, sendErr := bot.SendMessage(tu.Message(userId,
			fmt.Sprintf(s.deps.Messages.FailedSubs+" "+err.Error())).
			WithParseMode(telego.ModeHTML),
		)
		if sendErr != nil {
			s.deps.Logger.Errorw("send message error", "error", sendErr)
		}
		return
	}

	_, botErr := bot.SendMessage(tu.Message(userId, s.deps.Messages.SuccessSubs).
		WithParseMode(telego.ModeHTML))
	if botErr != nil {
		s.deps.Logger.Errorw("bot error", "error", botErr)
	}
}

func IsSubscriptionReplyMessage() th.Predicate {
	return func(update telego.Update) bool {
		return update.Message != nil &&
			update.Message.ReplyToMessage != nil &&
			strings.Contains(update.Message.ReplyToMessage.Text, "add")
	}
}

func (s *Subscription) HandleUnsubscribeButton(bot *telego.Bot, update telego.Update) {
	if update.CallbackQuery == nil || update.CallbackQuery.Message == nil {
		s.deps.Logger.Errorw("invalid callback query", "update", update)
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
		"handle button remove",
		"chat_id", chatID,
		"message_id", messageID,
	)
	userId := tu.ID(chatID)

	msg, err := bot.SendMessage(tu.Message(userId, s.deps.Messages.NicknameForUnsubs).
		WithReplyMarkup(tu.ForceReply()))
	if err != nil {
		s.deps.Logger.Errorw("failed to send message", "error", err)
		return
	}

	s.deps.Logger.Debugw("sent message with force reply", "message_id", messageID)

	if err := bot.AnswerCallbackQuery(&telego.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	}); err != nil {
		s.deps.Logger.Errorw("failed to answer callback", "error", err)
	}
}

func (s *Subscription) HandleUnsubscriptionNicknameReply(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID
	userId := tu.ID(chatID)

	userMessage := strings.TrimSpace(update.Message.Text)
	if userMessage == "" {
		_, err := bot.SendMessage(tu.Message(userId, "Please enter a valid nickname.").
			WithParseMode(telego.ModeHTML),
		)
		if err != nil {
			s.deps.Logger.Errorw("send message error", "error", err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	playerId, err := s.deps.Faceit.GetPlayerIDByNickname(ctx, userMessage)
	if err != nil {
		s.deps.Logger.Errorw("failed to get user", "error", err)
	}

	isSubscribed, err := s.deps.SubscriptionRepo.IsSubscribed(ctx, chatID, playerId)
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

	if !isSubscribed {
		_, sendErr := bot.SendMessage(tu.Message(userId,
			fmt.Sprintf(s.deps.Messages.NotSubscribed, userMessage)))
		if sendErr != nil {
			s.deps.Logger.Errorw("send message error", "error", sendErr)
		}
		return
	}

	err = s.deps.SubscriptionRepo.Unsubscribe(ctx, chatID, playerId)
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

	_, botErr := bot.SendMessage(tu.Message(userId, s.deps.Messages.SuccessUnsubs))
	if botErr != nil {
		s.deps.Logger.Errorw("bot error", "error", botErr)
	}
}
func IsUnsubscriptionReplyMessage() th.Predicate {
	return func(update telego.Update) bool {
		return update.Message != nil &&
			update.Message.ReplyToMessage != nil &&
			strings.Contains(update.Message.ReplyToMessage.Text, "delete")
	}
}

func (s *Subscription) HandleSubscriptionsListButton(bot *telego.Bot, update telego.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if update.CallbackQuery == nil ||
		update.CallbackQuery.Message == nil {
		s.deps.Logger.Errorw(
			"invalid callback query",
			"update", update,
		)
		return
	}

	msg := update.CallbackQuery.Message
	if msg == nil {
		s.deps.Logger.Errorw("message is nil")
		return
	}

	chatID := msg.GetChat().ID
	subs, err := s.deps.SubscriptionRepo.GetSubscriptionByChatID(ctx, chatID)
	if err != nil {
		s.deps.Logger.Errorw("failed to get subscriber", "error", err)
		_, err = bot.SendMessage(tu.Message(tu.ID(chatID),
			fmt.Sprintf(s.deps.Messages.NotSubscribed+err.Error())).
			WithParseMode(telego.ModeHTML))
		if err != nil {
			s.deps.Logger.Errorw("send message error", "error", err)
		}
	}

	_, err = bot.SendMessage(tu.Message(tu.ID(chatID), formatSubscriptionsList(s.deps, subs)).
		WithParseMode(telego.ModeHTML))
	if err != nil {
		s.deps.Logger.Errorw("failed to send message", "error", err)
	}

	if err := bot.AnswerCallbackQuery(tu.CallbackQuery(update.CallbackQuery.ID)); err != nil {
		s.deps.Logger.Errorw("failed to answer callback", "error", err)
		return
	}
}
