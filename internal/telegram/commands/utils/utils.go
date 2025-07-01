package utils

import (
	"fmt"
	"github.com/anlukk/faceit-tracker/internal/faceit/pkg/go-faceit"
)

func FormatResponseMessage(response *faceit.Player) string {
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
