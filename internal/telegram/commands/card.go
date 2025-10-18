package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type PlayerCard struct {
	deps *core.Dependencies
}

func NewPlayerCard(deps *core.Dependencies) *PlayerCard {
	return &PlayerCard{
		deps: deps,
	}
}

func (s *PlayerCard) HandlePlayerButton(bot *telego.Bot, update telego.Update) {
	msg := update.CallbackQuery.Message

	chatID := msg.GetChat().ID
	userId := tu.ID(chatID)

	subs, err := s.deps.SubscriptionRepo.GetSubscriptionsByChatID(context.Background(), chatID)
	if err != nil {
		s.deps.Logger.Errorw("failed to get subscriber", "error", err)
		_, err = bot.SendMessage(tu.Message(tu.ID(chatID), fmt.Sprintf(s.deps.Messages.NotSubscribed+err.Error())).
			WithParseMode(telego.ModeHTML))
	}

	nickname := strings.TrimPrefix(update.CallbackQuery.Data, "player:")

	var found bool
	for _, v := range subs {
		if v.Nickname == nickname {
			found = true
			break
		}
	}

	player, err := s.deps.Faceit.GetPlayerByNickname(context.Background(), nickname)
	if err != nil {
		s.deps.Logger.Errorw("failed to get player", "error", err)
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "❌ Player not found").WithParseMode(telego.ModeHTML))
		return
	}

	lastTenMatches, err := s.deps.Faceit.GetStatForLastTenMatches(context.Background(), nickname)
	if err != nil {
		s.deps.Logger.Errorw("failed to get last ten matches", "error", err)
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), "❌ Error retrieving last ten matches").
			WithParseMode(telego.ModeHTML))
		return
	}

	if found {
		msg, err = bot.SendMessage(tu.Message(userId, formatPlayerCard(&player, lastTenMatches)).
			WithReplyMarkup(tu.ForceReply()).WithParseMode(telego.ModeHTML))
		if err != nil {
			s.deps.Logger.Errorw("failed to send message", "error", err)
			return
		}
	}

	if err := bot.AnswerCallbackQuery(&telego.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	}); err != nil {
		s.deps.Logger.Errorw("failed to answer callback", "error", err)
	}
}
