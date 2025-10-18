package commands

import (
	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/anlukk/faceit-tracker/internal/telegram/menu"
)

type BotCommands struct {
	StartCommand        *Start
	SearchPlayerCommand *SearchPlayer
	Subscription        *Subscription
	PlayerCard          *PlayerCard
}

func NewBotCommands(
	deps *core.Dependencies,
	menu *menu.MenuManager) *BotCommands {
	return &BotCommands{
		StartCommand:        NewStart(deps, menu),
		SearchPlayerCommand: NewSearchPlayer(deps),
		Subscription:        NewSubscription(deps),
		PlayerCard:          NewPlayerCard(deps),
	}
}
