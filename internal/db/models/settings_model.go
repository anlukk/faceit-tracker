package models

type UserSettings struct {
	ChatID               int64 `gorm:"primaryKey"`
	NotificationsEnabled bool
	Language             string `gorm:"size:10"`
}
