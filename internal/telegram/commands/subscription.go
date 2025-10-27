package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

const (
	customCtxTimeOutForSubscription = 10 * time.Second
)

type Subscription struct {
	deps *core.Dependencies
}

func NewSubscription(deps *core.Dependencies) *Subscription {
	return &Subscription{
		deps: deps,
	}
}

func (s *Subscription) HandleSubscribeButton(bot *telego.Bot, update telego.Update) {
	callbackQueryID := update.CallbackQuery.Message.GetChat().ID
	sendForceReply(bot, callbackQueryID, s.deps.Messages.NicknameForSubs, s.deps.Logger)
	answerCallback(bot, update.CallbackQuery.ID, s.deps.Logger)
}

func (s *Subscription) HandleSubscriptionNicknameReply(bot *telego.Bot, update telego.Update) {
	ctx, cancel := context.WithTimeout(s.deps.Ctx, customCtxTimeOutForSubscription)
	defer cancel()

	nickname := strings.TrimSpace(update.Message.Text)
	playerId, err := s.deps.Faceit.GetPlayerIDByNickname(ctx, nickname)
	if err != nil {
		s.deps.Logger.Errorw("failed to get user", "error", err)
	}

	telegramChatID := update.Message.Chat.ID
	err = s.deps.SubscriptionRepo.Subscribe(ctx, telegramChatID, playerId, nickname)
	if err != nil {
		s.deps.Logger.Errorw("failed to subscribe", "error", err)
		reply(bot, telegramChatID, fmt.Sprintf(s.deps.Messages.FailedSubs+" "+err.Error()), s.deps.Logger)
		return
	}

	reply(bot, telegramChatID, s.deps.Messages.SuccessSubs, s.deps.Logger)
}

func IsSubscriptionReplyMessage() th.Predicate {
	return func(update telego.Update) bool {
		return update.Message != nil &&
			update.Message.ReplyToMessage != nil &&
			strings.Contains(update.Message.ReplyToMessage.Text, "add")
	}
}

func (s *Subscription) HandleUnsubscribeButton(bot *telego.Bot, update telego.Update) {
	callbackQueryID := update.CallbackQuery.Message.GetChat().ID
	sendForceReply(bot, callbackQueryID, s.deps.Messages.NicknameForUnsubs, s.deps.Logger)
	answerCallback(bot, update.CallbackQuery.ID, s.deps.Logger)
}

func (s *Subscription) HandleUnsubscriptionNicknameReply(bot *telego.Bot, update telego.Update) {
	ctx, cancel := context.WithTimeout(s.deps.Ctx, customCtxTimeOutForSubscription)
	defer cancel()

	nickname := strings.TrimSpace(update.Message.Text)
	playerId, err := s.deps.Faceit.GetPlayerIDByNickname(ctx, nickname)
	if err != nil {
		s.deps.Logger.Errorw("failed to get user", "error", err)
	}

	telegramChatID := update.Message.Chat.ID
	err = s.deps.SubscriptionRepo.Unsubscribe(ctx, telegramChatID, playerId)
	if err != nil {
		s.deps.Logger.Errorw("failed to unsubscribe", "error", err)
		reply(bot, telegramChatID, fmt.Sprintf(s.deps.Messages.FailedSubs+" "+err.Error()), s.deps.Logger)
	}

	reply(bot, telegramChatID, s.deps.Messages.SuccessUnsubs, s.deps.Logger)
}

func IsUnsubscriptionReplyMessage() th.Predicate {
	return func(update telego.Update) bool {
		return update.Message != nil &&
			update.Message.ReplyToMessage != nil &&
			strings.Contains(update.Message.ReplyToMessage.Text, "delete")
	}
}

func (s *Subscription) HandleNewPersonalSubButton(bot *telego.Bot, update telego.Update) {
	callbackQueryID := update.CallbackQuery.Message.GetChat().ID
	sendForceReply(bot, callbackQueryID, "enter the new main player", s.deps.Logger)
	answerCallback(bot, update.CallbackQuery.ID, s.deps.Logger)
}

func (s *Subscription) HandleNewPersonalSubReply(bot *telego.Bot, update telego.Update) {
	ctx, cancel := context.WithTimeout(s.deps.Ctx, customCtxTimeOutForSubscription)
	defer cancel()

	nickname := strings.TrimSpace(update.Message.Text)
	telegramChatID := update.Message.Chat.ID
	err := s.deps.PersonalSubRepo.SetPersonalSub(ctx, telegramChatID, nickname)
	if err != nil {
		s.deps.Logger.Errorw("failed to set personal sub", "error", err)
		reply(bot, telegramChatID, "Error setting personal sub. Please try again.", s.deps.Logger)
		return
	}

	reply(bot, telegramChatID, "Successfully set personal sub.", s.deps.Logger)
}

func IsNewPersonalSubReplyMessage() th.Predicate {
	return func(update telego.Update) bool {
		return update.Message != nil &&
			update.Message.ReplyToMessage != nil &&
			update.Message.ReplyToMessage.Text == "enter the new main player"
	}
}
