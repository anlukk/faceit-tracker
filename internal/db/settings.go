package db

import "context"

type SettingsDB interface {
	GetNotificationsEnabled(ctx context.Context, chatID int64) (bool, error)
	SetNotificationsEnabled(ctx context.Context, chatID int64, enabled bool) error

	GetLanguage(ctx context.Context, chatID int64) string
	SetLanguage(ctx context.Context, chatID int64, language string) error
}
