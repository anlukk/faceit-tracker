package service

import (
	"github.com/anlukk/faceit-tracker/internal/db"
	"github.com/anlukk/faceit-tracker/internal/faceit"
	"github.com/anlukk/faceit-tracker/internal/service/notification"
	"github.com/anlukk/faceit-tracker/internal/service/settings"
	"github.com/anlukk/faceit-tracker/internal/service/subscription"
	"gorm.io/gorm"
)

type Services struct {
	Subscriptions subscription.Subscription
	Notifications notification.Notifications
	Settings      settings.Settings
}

func NewServices(gormDb *gorm.DB, faceit faceit.Faceit) *Services {
	subRepo := db.NewSubscriptionDBImpl(gormDb)
	settingsRepo := db.NewSettngsDBImpl(gormDb)
	return &Services{
		Subscriptions: subscription.NewService(subRepo),
		Notifications: notification.NewService(faceit),
		Settings:      settings.NewService(settingsRepo),
	}
}
