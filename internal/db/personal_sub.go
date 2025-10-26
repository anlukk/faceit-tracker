package db

import (
	"context"

	"github.com/anlukk/faceit-tracker/internal/db/models"
)

type PersonalSubDB interface {
	GetPersonalSub(ctx context.Context, chatID int64) (*models.PersonalSub, error)
	SetPersonalSub(ctx context.Context, chatID int64, nickname string) error
}
