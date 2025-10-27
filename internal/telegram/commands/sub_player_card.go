package commands

import (
	"fmt"
	"strings"

	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/mymmrac/telego"
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
	chatID := update.CallbackQuery.Message.GetChat().ID
	subs, err := s.deps.SubscriptionRepo.GetSubscriptionsByChatID(s.deps.Ctx, chatID)
	if err != nil {
		s.deps.Logger.Errorw("failed to get subscriber", "error", err)
		reply(bot, chatID, fmt.Sprintf(s.deps.Messages.NotSubscribed+err.Error()), s.deps.Logger)
	}

	nickname := strings.TrimPrefix(update.CallbackQuery.Data, "player:")

	var found bool
	for _, sub := range subs {
		if sub.Nickname == nickname {
			found = true
			break
		}
	}

	player, err := s.deps.Faceit.GetPlayerByNickname(s.deps.Ctx, nickname)
	if err != nil {
		s.deps.Logger.Errorw("failed to get player", "error", err)
		reply(bot, chatID, "❌ Player not found", s.deps.Logger)
	}

	lastTenMatches, err := s.deps.Faceit.GetStatForLastTenMatches(s.deps.Ctx, nickname)
	if err != nil {
		s.deps.Logger.Errorw("failed to get last ten matches", "error", err)
		reply(bot, chatID, "❌ Error retrieving last ten matches", s.deps.Logger)
	}

	form := formatPlayerCard(&player, lastTenMatches)
	if found {
		sendForceReply(bot, chatID, form, s.deps.Logger)
	}

	answerCallback(bot, update.CallbackQuery.ID, s.deps.Logger)
}
