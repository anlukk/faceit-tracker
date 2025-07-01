package models

import (
	"gorm.io/gorm"
	"time"
)

type Player struct {
	gorm.Model
	ActivatedAt        time.Time    `gorm:"index" json:"activated_at,omitempty"`
	Avatar             string       `gorm:"size:255" json:"avatar,omitempty"`
	Country            string       `gorm:"size:2" json:"country,omitempty"`
	CoverFeaturedImage string       `gorm:"size:255" json:"cover_featured_image,omitempty"` // Deprecated: no more in use
	CoverImage         string       `gorm:"size:255" json:"cover_image,omitempty"`
	FaceitUrl          string       `gorm:"size:255" json:"faceit_url,omitempty"`
	FriendsIds         []string     `gorm:"type:text[]" json:"friends_ids,omitempty"`
	Games              []Game       `gorm:"foreignKey:PlayerID" json:"games,omitempty"`
	MembershipType     string       `gorm:"size:50" json:"membership_type,omitempty"` // Deprecated: use memberships instead
	Memberships        []string     `gorm:"type:text[]" json:"memberships,omitempty"`
	NewSteamId         string       `gorm:"size:50;uniqueIndex" json:"new_steam_id,omitempty"`
	Nickname           string       `gorm:"size:100;index" json:"nickname,omitempty"`
	Platforms          []Platform   `gorm:"foreignKey:PlayerID" json:"platforms,omitempty"`
	PlayerId           string       `gorm:"size:50;uniqueIndex" json:"player_id,omitempty"`
	Settings           UserSettings `gorm:"embedded" json:"settings,omitempty"`
	SteamId64          string       `gorm:"size:50;uniqueIndex" json:"steam_id_64,omitempty"`
	SteamNickname      string       `gorm:"size:100" json:"steam_nickname,omitempty"`
	Verified           bool         `gorm:"default:false" json:"verified,omitempty"`
}

type Subscription struct {
	gorm.Model
	ChatID    int64  `gorm:"index"` // ID чата в Telegram
	PlayerID  string `gorm:"index"` // ID игрока на Faceit
	Nickname  string `gorm:"size:100"`
	LastStats string `gorm:"type:text"`
	LastCheck time.Time
}

type Game struct {
	gorm.Model
	PlayerID   string     `gorm:"size:50;index"` // Связь с Player
	GameName   string     `gorm:"size:100"`      // Ключ в map[string]GameDetail
	GameDetail GameDetail `gorm:"embedded"`
}

type GameDetail struct {
	FaceitElo  int    `json:"faceit_elo,omitempty"`
	GameId     string `gorm:"size:50" json:"game_id,omitempty"`
	GameName   string `gorm:"size:100" json:"game_name,omitempty"`
	SkillLevel int    `json:"skill_level,omitempty"`
}

type Platform struct {
	gorm.Model
	PlayerID  string `gorm:"size:50;index"` // Связь с Player
	Name      string `gorm:"size:50"`       // Ключ в map[string]string
	AccountID string `gorm:"size:100"`      // Значение в map[string]string
}

type UserSettings struct {
	Language string `gorm:"size:10" json:"language,omitempty"`
}
