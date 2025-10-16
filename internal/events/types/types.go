package types

import "time"

type Event struct {
	Type      string
	ChatID    int64
	Message   string
	Nickname  string
	Payload   any
	Timestamp time.Time
}
