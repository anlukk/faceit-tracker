package models

import "time"

type PersonalSub struct {
	ChatID    int64  `gorm:"uniqueIndex"`
	Nickname  string `gorm:"not null"`
	LastCheck time.Time
}
