package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/anlukk/faceit-tracker/internal/faceit/pkg/go-faceit"

	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/mymmrac/telego"
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

func (s *SearchPlayer) PromptPlayerSearch(bot *telego.Bot, update telego.Update) {
	userId := tu.ID(update.Message.Chat.ID)

	s.deps.Logger.Debugw(
		"prompt player search",
		"user_id", userId,
	)

	_, botErr := bot.SendMessage(tu.Message(userId, "Enter the player you want to find").
		WithReplyMarkup(tu.ForceReply()))
	if botErr != nil {
		s.deps.Logger.Errorw("bot error", "error", botErr)
	}
}

func (s *SearchPlayer) HandleUserMessage(bot *telego.Bot, update telego.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	userId := tu.ID(update.Message.From.ID)
	userMessage := update.Message.Text
	if strings.TrimSpace(userMessage) == "" {
		_, err := bot.SendMessage(tu.Message(userId, "Please enter a valid nickname.").
			WithParseMode(telego.ModeHTML))
		if err != nil {
			s.deps.Logger.Errorw("send message error", "error", err)
		}
		return
	}

	response, err := s.deps.Faceit.GetPlayer(ctx, userMessage)
	if err != nil {
		s.deps.Logger.Errorw("failed to get user", "error", err)
		_, sendErr := bot.SendMessage(tu.Message(userId, "Error retrieving player data.").
			WithParseMode(telego.ModeHTML))
		if sendErr != nil {
			s.deps.Logger.Errorw("send message error", "error", sendErr)
		}
		return
	}

	formattedResponse := formatResponse(&response)
	_, err = bot.SendMessage(tu.Message(userId, formattedResponse).
		WithParseMode(telego.ModeHTML))
	if err != nil {
		s.deps.Logger.Errorw("send message error", "error", err)
		return
	}

	if userMessage == "cancel" {
		_, err := bot.SendMessage(tu.Message(userId, "Canceled").
			WithParseMode(telego.ModeHTML))
		if err != nil {
			s.deps.Logger.Errorw("send message error", "error", err)
		}
		return
	}
}

func formatResponse(response *faceit.Player) string {
	gamesStr := ""
	for game, gameInfo := range response.Games {
		gamesStr += fmt.Sprintf("Game: %s, FaceitElo: %d, SkillLevel: %d\n",
			game, gameInfo.FaceitElo, gameInfo.SkillLevel)
	}

	return fmt.Sprintf("Nickname: %s\n"+"Country: %s\n"+"Games: %s\n"+"Steam nickname: %s\n",
		response.Nickname, response.Country, gamesStr, response.SteamNickname)
}
