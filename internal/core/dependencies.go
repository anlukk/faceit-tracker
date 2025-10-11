package core

import (
	"github.com/anlukk/faceit-tracker/internal/config"
	"github.com/anlukk/faceit-tracker/internal/db"
	"github.com/anlukk/faceit-tracker/internal/faceit"
	"go.uber.org/zap"
)

type Dependencies struct {
	Config           *config.Config
	Messages         *config.BotMessages
	Logger           *zap.SugaredLogger
	Faceit           faceit.FaceitClient
	SettingsRepo     db.SettingsDB
	SubscriptionRepo db.SubscriptionDB
}
