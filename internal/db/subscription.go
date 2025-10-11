package db

import (
	"context"

	"github.com/anlukk/faceit-tracker/internal/db/models"
)

type SubscriptionDB interface {
	Subscribe(ctx context.Context, chatID int64, playerID, nickname string) error
	Unsubscribe(ctx context.Context, chatID int64, playerID string) error

	IsSubscribed(ctx context.Context, chatID int64, playerID string) (bool, error)
	GetSubscriptionByChatID(ctx context.Context, chatID int64) ([]models.Subscription, error)

	GetAllSubscription(ctx context.Context) ([]models.Subscription, error)

	//GetPlayerIDs(ctx context.Context, chatID int64) ([]string, error)
}
