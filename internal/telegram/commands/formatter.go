package commands

import (
	"fmt"

	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/anlukk/faceit-tracker/internal/db/models"
	"github.com/anlukk/faceit-tracker/internal/faceit/pkg/go-faceit"
)

// TODO: change to i18n
func formatSearchCommandResponse(response *faceit.Player) string {
	gamesStr := ""
	for game, gameInfo := range response.Games {
		gamesStr += fmt.Sprintf("Game: %s, FaceitElo: %d, SkillLevel: %d\n",
			game, gameInfo.FaceitElo, gameInfo.SkillLevel)
	}

	return fmt.Sprintf("Nickname: %s\n"+"Country: %s\n"+"Games: %s\n"+"Steam nickname: %s\n",
		response.Nickname, response.Country, gamesStr, response.SteamNickname)
}

// TODO: change to i18n
func formatSubscriptionsList(deps *core.Dependencies, subs []models.Subscription) string {
	if len(subs) == 0 {
		return deps.Messages.NoSubscriptions
	}

	sb := "Your subscription:\n"
	for i, sub := range subs {
		sb += fmt.Sprintf("%d. %s\n", i+1, sub.Nickname)
	}

	return sb
}
