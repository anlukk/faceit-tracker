package service

import (
	"github.com/anlukk/faceit-tracker/internal/db"
	"github.com/anlukk/faceit-tracker/internal/faceit"
	"github.com/anlukk/faceit-tracker/internal/service/match"
	"github.com/anlukk/faceit-tracker/internal/service/settings"
	"github.com/anlukk/faceit-tracker/internal/service/subscription"
	"gorm.io/gorm"
)

type Services struct {
	Subscription *subscription.Service
	Settings     *settings.Service
	Match        *match.Service
}

func NewServices(gormDb *gorm.DB, client faceit.FaceitClient) *Services {
	subRepo := db.NewSubscriptionDBImpl(gormDb)
	settingsRepo := db.NewSettngsDBImpl(gormDb)
	return &Services{
		Subscription: subscription.NewService(subRepo),
		Settings:     settings.NewService(settingsRepo),
		Match:        match.NewService(client),
	}
}
