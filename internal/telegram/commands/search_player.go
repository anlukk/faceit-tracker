package commands

import (
	"context"
	"fmt"
	"github.com/anlukk/faceit-tracker/internal/faceit/pkg/go-faceit"
	"strings"

	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type SearchPlayer struct {
	deps               *core.Dependencies
	waitingForUsername map[telego.ChatID]bool
}

func NewSearchPlayer(deps *core.Dependencies) *SearchPlayer {
	return &SearchPlayer{
		deps:               deps,
		waitingForUsername: make(map[telego.ChatID]bool),
	}
}

func (s *SearchPlayer) PromptPlayerSearch(bot *telego.Bot, update telego.Update) {
	if update.Message == nil ||
		update.Message.Chat == (telego.Chat{}) ||
		update.Message.Chat.ID == 0 {
		s.deps.Logger.Errorw("invalid message", "update", update)
		return
	}

	userId := tu.ID(update.Message.Chat.ID)

	s.deps.Logger.Debugw(
		"prompt player search",
		"user_id", userId,
	)
	s.waitingForUsername[userId] = true

	_, botErr := bot.SendMessage(
		tu.Message(userId, "Enter the player you want to find").
			WithReplyMarkup(tu.ForceReply()),
	)
	if botErr != nil {
		s.deps.Logger.Errorw("bot error", "error", botErr)
	}
}

func (s *SearchPlayer) HandleUserMessage(bot *telego.Bot, update telego.Update) {
	if update.Message == nil || update.Message.From == nil {
		s.deps.Logger.Errorw("nil message or sender", "update", update)
		return
	}

	userId := tu.ID(update.Message.From.ID)
	if !s.waitingForUsername[userId] {
		return
	}

	if s.deps.Faceit == nil {
		s.deps.Logger.Errorw("faceit is nil")
		return
	}

	userMessage := update.Message.Text
	if strings.TrimSpace(userMessage) == "" {
		_, err := bot.SendMessage(
			tu.Message(userId, "Please enter a valid username.").
				WithParseMode(telego.ModeHTML))
		if err != nil {
			s.deps.Logger.Errorw("send message error", "error", err)
		}
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := s.deps.Faceit.GetPlayer(ctx, userMessage)
	if err != nil {
		s.deps.Logger.Errorw("failed to get user", "error", err)
		_, sendErr := bot.SendMessage(
			tu.Message(userId, "Error fetching data from FACEIT API.").
				WithParseMode(telego.ModeHTML))
		if sendErr != nil {
			s.deps.Logger.Errorw("send message error", "error", sendErr)
		}
		return
	}

	formattedResponse := formatResponse(&response)
	_, err = bot.SendMessage(tu.Message(userId, formattedResponse).WithParseMode(telego.ModeHTML))
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

	s.waitingForUsername[userId] = false
}

func formatResponse(response *faceit.Player) string {
	gamesStr := ""
	for game, gameInfo := range response.Games {
		gamesStr += fmt.Sprintf(
			"Game: %s, FaceitElo: %d, SkillLevel: %d\n",
			game, gameInfo.FaceitElo, gameInfo.SkillLevel,
		)
	}

	return fmt.Sprintf(
		"Nickname: %s\n"+
			"Country: %s\n"+
			"Games: %s\n"+
			"Steam nickname: %s\n",
		response.Nickname,
		response.Country,
		gamesStr,
		response.SteamNickname,
	)
}
