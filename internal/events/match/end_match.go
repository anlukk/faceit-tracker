package match

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/anlukk/faceit-tracker/internal/events/cache"
	"github.com/anlukk/faceit-tracker/internal/events/helper"
	"github.com/anlukk/faceit-tracker/internal/events/types"
)

const checkCooldown = 50 * time.Minute

type End struct {
	deps  *core.Dependencies
	cache *cache.NotifyCache

	timeCheck map[string]time.Time
	mu        sync.RWMutex
}

func NewMatchEnd(deps *core.Dependencies, cache *cache.NotifyCache) *End {
	return &End{
		deps:      deps,
		cache:     cache,
		timeCheck: make(map[string]time.Time),
	}
}

func (m *End) GetEvents(ctx context.Context) ([]types.Event, error) {
	var e []types.Event

	subs, err := m.deps.SubscriptionRepo.GetAllSubscription(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all subscription")
	}

	for _, sub := range subs {
		lastChecked := m.getLastChecked(sub.Nickname)
		if time.Since(lastChecked) < checkCooldown {
			continue
		}

		finishMatchResult, err := m.deps.Faceit.GetFinishMatchResult(ctx, sub.Nickname)
		if err != nil {
			m.deps.Logger.Errorw("failed to get finish match result", "error", err)
			continue
		}

		m.setLastChecked(sub.Nickname, time.Now())

		if m.cache.AlreadyNotified(sub.Nickname, finishMatchResult.MatchId) {
			continue
		}

		finishedAt := time.Unix(finishMatchResult.FinishedAt, 0)
		if time.Since(finishedAt) > 10*time.Minute {
			continue
		}

		formatMessage := helper.FormatMatchEndMessage(m.deps.Messages, finishMatchResult)

		e = append(e, types.Event{
			Type:      m.EventType(),
			ChatID:    sub.ChatID,
			Message:   formatMessage,
			Nickname:  sub.Nickname,
			Timestamp: time.Now(),
		})

		m.cache.MarkNotified(sub.Nickname, finishMatchResult.MatchId)
	}

	return e, nil
}

func (m *End) EventType() string {
	return "match_end"
}

func (m *End) getLastChecked(nickname string) time.Time {
	m.mu.RLock()
	t, ok := m.timeCheck[nickname]
	m.mu.RUnlock()

	if !ok {
		return time.Time{}
	}

	return t
}

func (m *End) setLastChecked(nickname string, t time.Time) {
	m.mu.Lock()
	m.timeCheck[nickname] = t
	m.mu.Unlock()
}
