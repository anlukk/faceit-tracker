package models

import (
	"gorm.io/gorm"
	"time"
)

type Subscription struct {
	gorm.Model
	ChatID    int64  `gorm:"index"`
	PlayerID  string `gorm:"index"`
	Nickname  string `gorm:"size:100"`
	LastStats string `gorm:"type:text"`
	LastCheck time.Time
}
