package cache

import "sync"

type NotifyCache struct {
	mu       sync.RWMutex
	notified map[string]string
}

func NewNotifyCache() *NotifyCache {
	return &NotifyCache{
		notified: make(map[string]string),
	}
}

func (c *NotifyCache) AlreadyNotified(nickname, matchID string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.notified[nickname] == matchID
}

func (c *NotifyCache) MarkNotified(nickname, matchID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.notified[nickname] = matchID
}
