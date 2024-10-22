package utils

import (
	"github.com/mconnat/go-faceit/pkg/models"
	"fmt"
)

func FormatResponseMessage(response *models.Player) string {
	gamesStr := ""

	for game, gameInfo := range response.Games {
		gamesStr += fmt.Sprintf(
			"Game: %s, FaceitElo: %d, SkillLevel: %d\n",
			game, gameInfo.FaceitElo, gameInfo.SkillLevel,
		)
	}

	return fmt.Sprintf(
		"Nickname: %s\n" +
			"Country: %s\n" +
			"Games: %s\n" +
			"Steam nickname: %s\n",
		response.Nickname,
		response.Country,
		gamesStr,
		response.SteamNickname,
	)
}