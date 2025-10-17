package commands

import (
	"context"
	"strings"
	"time"

	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type SearchPlayer struct {
	deps *core.Dependencies
}

func NewSearchPlayer(deps *core.Dependencies) *SearchPlayer {
	return &SearchPlayer{
		deps: deps,
	}
}

// TODO: add i18n support
func (s *SearchPlayer) PromptPlayerSearch(bot *telego.Bot, update telego.Update) {
	chatID := tu.ID(update.Message.Chat.ID)

	_, botErr := bot.SendMessage(tu.Message(chatID, "Enter the player you want to find").
		WithReplyMarkup(tu.ForceReply()))
	if botErr != nil {
		s.deps.Logger.Errorw("bot error", "error", botErr)
	}

	if err := bot.AnswerCallbackQuery(&telego.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	}); err != nil {
		s.deps.Logger.Errorw("failed to answer callback", "error", err)
	}
}

func FindPlayerReplyMessage() th.Predicate {
	return func(update telego.Update) bool {
		return update.Message != nil &&
			update.Message.ReplyToMessage != nil &&
			update.Message.ReplyToMessage.Text == "Enter the player you want to find"
	}
}

// TODO: add i18n support
func (s *SearchPlayer) HandleUserMessage(bot *telego.Bot, update telego.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	chatID := tu.ID(update.Message.Chat.ID)

	userMessage := update.Message.Text
	if strings.TrimSpace(userMessage) == "" {
		_, err := bot.SendMessage(tu.Message(chatID, "Please enter a valid nickname.").
			WithParseMode(telego.ModeHTML))
		if err != nil {
			s.deps.Logger.Errorw("send message error", "error", err)
		}
		return
	}

	response, err := s.deps.Faceit.GetPlayer(ctx, userMessage)
	if err != nil {
		s.deps.Logger.Errorw("failed to get user", "error", err)
		_, sendErr := bot.SendMessage(tu.Message(chatID, "Error retrieving player data.").
			WithParseMode(telego.ModeHTML))
		if sendErr != nil {
			s.deps.Logger.Errorw("send message error", "error", sendErr)
		}
		return
	}

	formattedResponse := formatSearchCommandResponse(&response)
	_, err = bot.SendMessage(tu.Message(chatID, formattedResponse).WithParseMode(telego.ModeHTML))
	if err != nil {
		s.deps.Logger.Errorw("send message error", "error", err)
		return
	}
}
