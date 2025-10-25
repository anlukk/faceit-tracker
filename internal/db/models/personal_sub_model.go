package models

import "time"

type PersonalSub struct {
	ChatID    int64  `gorm:"uniqueIndex"`
	PlayerID  string `gorm:"not null"`
	LastCheck time.Time
}
