package models

import (
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model
	ChatID       int64        // FK на UserSettings.ChatID
	UserSettings UserSettings `gorm:"foreignKey:ChatID;references:ChatID"`

	PlayerID  string `gorm:"index"`
	Nickname  string `gorm:"size:100"`
	LastStats string `gorm:"type:text"`
	LastCheck time.Time
}
