package match

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/anlukk/faceit-tracker/internal/db/models"
	cache2 "github.com/anlukk/faceit-tracker/internal/events/cache"
	"github.com/anlukk/faceit-tracker/internal/events/types"
)

const (
	checkCooldown = 50 * time.Minute
	numWorkers    = 20
)

type End struct {
	deps  *core.Dependencies
	cache *cache2.NotifyCache

	timeCheck map[string]time.Time
	mu        sync.RWMutex
}

func NewMatchEnd(deps *core.Dependencies, cache *cache2.NotifyCache) *End {
	return &End{
		deps:      deps,
		cache:     cache,
		timeCheck: make(map[string]time.Time),
	}
}

// TODO: add cache invalidation to the production version
func (m *End) GetEvents(ctx context.Context) ([]types.Event, error) {
	var events []types.Event

	subs, err := m.deps.SubscriptionRepo.GetAllSubscription(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all subscription: %w", err)
	}

	subsChannel := make(chan models.Subscription, numWorkers)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for sub := range subsChannel {
				if sub.UserSettings.NotificationsEnabled == false {
					continue
				}

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
					m.deps.Logger.Debugw("already notified", "nickname", sub.Nickname)
					continue
				}

				finishedAt := time.Unix(finishMatchResult.FinishedAt, 0)
				if time.Since(finishedAt) > 10*time.Minute {
					m.deps.Logger.Debugw("match is too old", "nickname", sub.Nickname)
					continue
				}

				m.cache.MarkNotified(sub.Nickname, finishMatchResult.MatchId)

				m.mu.Lock()
				events = append(events, types.Event{
					Type:   m.EventType(),
					ChatID: sub.ChatID,
					Message: FormatMatchEndMessage(
						m.deps.Messages,
						finishMatchResult),
					Nickname:  sub.Nickname,
					Timestamp: time.Now(),
				})
				m.mu.Unlock()
			}
		}()
	}

	for _, sub := range subs {
		subsChannel <- sub
	}

	close(subsChannel)
	wg.Wait()

	return events, nil
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
