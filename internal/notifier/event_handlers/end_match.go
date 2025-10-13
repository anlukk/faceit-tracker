package event_handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/anlukk/faceit-tracker/internal/faceit"
	"github.com/anlukk/faceit-tracker/internal/notifier/cache"
)

type EndMatch struct {
	deps  *core.Dependencies
	cache *cache.MatchCache
}

func NewEndMatch(deps *core.Dependencies, cache *cache.MatchCache) *EndMatch {
	return &EndMatch{
		deps:  deps,
		cache: cache,
	}
}

func (m *EndMatch) GetEvents(ctx context.Context) ([]Event, error) {
	var events []Event

	subs, err := m.deps.SubscriptionRepo.GetAllSubscription(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all subscription")
	}

	for _, sub := range subs {
		lastMatch, err := m.deps.Faceit.GetLastMatch(ctx, sub.Nickname)
		if err != nil {
			return nil, fmt.Errorf("failed to get last match")
		}

		if lastMatch.Status != "FINISHED" ||
			m.cache.AlreadyNotified(sub.Nickname, lastMatch.MatchId) {
			continue
		}

		finishedAt := time.Unix(lastMatch.FinishedAt, 0)

		if time.Since(finishedAt) > 10*time.Minute {
			continue
		}

		finishMatchResult, err := m.deps.Faceit.GetFinishMatchResult(ctx, sub.Nickname)
		if err != nil {
			return nil, fmt.Errorf("failed to get finish match result")
		}

		events = append(events, Event{
			Type:      m.EventType(),
			ChatID:    sub.ChatID,
			Message:   m.formatMessage(finishMatchResult),
			Nickname:  sub.Nickname,
			Timestamp: time.Now(),
		})

		m.cache.MarkNotified(sub.Nickname, lastMatch.MatchId)
	}

	return events, nil
}

func (m *EndMatch) EventType() string {
	return "match_end"
}

func (m *EndMatch) formatMessage(
	matchInfo *faceit.FinishMatchResult) string {
	isWin := map[bool]string{
		true:  m.deps.Messages.MatchWin,
		false: m.deps.Messages.MatchLoose,
	}

	result := fmt.Sprintf(
		"%s\n%s%s\n%s\n%s%s\n",
		m.deps.Messages.MatchFinish,
		m.deps.Messages.Nickname,
		matchInfo.Nickname,
		isWin[matchInfo.Win],
		m.deps.Messages.MatchScore,
		matchInfo.Score,
	)

	return result
}
