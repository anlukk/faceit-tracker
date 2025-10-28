package models

import "time"

type PersonalSub struct {
	ChatID   int64  `gorm:"uniqueIndex"`
	Nickname string `gorm:"not null"`
	//DetailedInfo DetailedInfo
	LastCheck time.Time
}

//type DetailedInfo struct {
//	MatchID string `gorm:"uniqueIndex"`
//	Elo     int
//	KD      float64
//	WinRate float64
//}
