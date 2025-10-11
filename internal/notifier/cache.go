package notifier

import "sync"

type MatchCache struct {
	notifiedMatches sync.Map
}

func NewMatchCache() *MatchCache {
	return &MatchCache{
		notifiedMatches: sync.Map{},
	}
}

func (m *MatchCache) alreadyNotified(nickname, matchID string) bool {
	if nickname == "" || matchID == "" {
		return false
	}

	val, ok := m.notifiedMatches.Load(nickname)
	if !ok {
		return false
	}

	return val == matchID
}

func (m *MatchCache) markNotified(nickname, matchID string) {
	if nickname == "" || matchID == "" {
		return
	}

	m.notifiedMatches.Store(nickname, matchID)
}
