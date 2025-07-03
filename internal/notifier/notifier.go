package notifier

import (
	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/mymmrac/telego"
)

type Notifier struct {
	deps *core.Dependencies
	bot  *telego.Bot
}

func NewNotifier(deps *core.Dependencies, bot *telego.Bot) *Notifier {
	return &Notifier{
		deps: deps,
		bot:  bot,
	}
}
