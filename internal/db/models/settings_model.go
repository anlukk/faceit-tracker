package models

type UserSettings struct {
	ChatID               int64  `gorm:"primaryKey"`
	NotificationsEnabled bool   `gorm:"column:notifications_enabled"`
	Language             string `gorm:"size:10" json:"language,omitempty"`
}
