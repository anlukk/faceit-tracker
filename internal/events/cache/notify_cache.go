package cache

import "sync"

type NotifyCache struct {
	notifiedMatches sync.Map
}

func NewNotifyCache() *NotifyCache {
	return &NotifyCache{
		notifiedMatches: sync.Map{},
	}
}

func (m *NotifyCache) AlreadyNotified(nickname, matchID string) bool {
	if nickname == "" || matchID == "" {
		return false
	}

	val, ok := m.notifiedMatches.Load(nickname)
	if !ok {
		return false
	}

	return val == matchID
}

func (m *NotifyCache) MarkNotified(nickname, matchID string) {
	if nickname == "" || matchID == "" {
		return
	}

	m.notifiedMatches.Store(nickname, matchID)
}
