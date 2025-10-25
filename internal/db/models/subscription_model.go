package models

import (
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model
	ChatID       int64        // FK on UserSettings.ChatID
	UserSettings UserSettings `gorm:"foreignKey:ChatID;references:ChatID"`

	//IsPersonal bool   `gorm:"default:false"`
	//PersonalSub *PersonalSub `gorm:"foreignKey:ChatID,PlayerID;references:ChatID,PlayerID"`

	PlayerID  string `gorm:"index"`
	Nickname  string `gorm:"size:100"`
	LastStats string `gorm:"type:text"`
	LastCheck time.Time
}
