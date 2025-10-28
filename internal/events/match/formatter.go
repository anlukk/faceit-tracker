package match

import (
	"fmt"

	"github.com/anlukk/faceit-tracker/internal/config"
	"github.com/anlukk/faceit-tracker/internal/faceit"
)

func FormatMatchEndMessage(messages *config.BotMessages, info *faceit.FinishMatchResult) string {
	isWin := map[bool]string{
		true:  "🏆 *" + messages.MatchWin + "*",
		false: "💀 *" + messages.MatchLoose + "*",
	}

	return fmt.Sprintf(
		"*%s*\n━━━━━━━━━━━━━━━\n👤 *%s*: `%s`\n%s\n📊 *%s*: `%s`\n━━━━━━━━━━━━━━━",
		messages.MatchFinish,
		messages.Nickname,
		info.Nickname,
		isWin[info.Win],
		messages.MatchScore,
		info.Score,
	)
}
