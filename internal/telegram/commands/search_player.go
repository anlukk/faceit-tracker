package commands

import (
	"context"
	"strings"
	"time"

	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

const (
	customCtxTimeOutForPlayerSearchCommand = 3 * time.Second
)

type SearchPlayer struct {
	deps *core.Dependencies
}

func NewSearchPlayer(deps *core.Dependencies) *SearchPlayer {
	return &SearchPlayer{
		deps: deps,
	}
}

func (s *SearchPlayer) PromptPlayerSearch(bot *telego.Bot, update telego.Update) {
	sendForceReply(bot, update.Message.Chat.ID, "Enter the player you want to find", s.deps.Logger)
	answerCallback(bot, update.CallbackQuery.ID, s.deps.Logger)
}

func FindPlayerReplyMessage() th.Predicate {
	return func(update telego.Update) bool {
		return update.Message != nil &&
			update.Message.ReplyToMessage != nil &&
			update.Message.ReplyToMessage.Text == "Enter the player you want to find"
	}
}

func (s *SearchPlayer) HandleUserMessage(bot *telego.Bot, update telego.Update) {
	ctx, cancel := context.WithTimeout(s.deps.Ctx, customCtxTimeOutForPlayerSearchCommand)
	defer cancel()

	nickname := strings.TrimSpace(update.Message.Text)
	response, err := s.deps.Faceit.GetPlayerByNickname(ctx, nickname)
	if err != nil {
		s.deps.Logger.Errorw("failed to get user from search player command", "error", err)
		return
	}

	lastTenMatches, err := s.deps.Faceit.GetStatForLastTenMatches(s.deps.Ctx, nickname)
	if err != nil {
		s.deps.Logger.Errorw("failed to get last ten matches", "error", err)
	}

	formattedResponse := formatSearchCommandPlayerCard(&response, lastTenMatches)

	telegramChatID := update.Message.Chat.ID
	reply(bot, telegramChatID, formattedResponse, s.deps.Logger)
}
