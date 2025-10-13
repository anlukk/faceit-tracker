package notifier

import (
	"context"
	"fmt"

	"github.com/anlukk/faceit-tracker/internal/core"
)

type Subscriber struct {
	deps       *core.Dependencies
	playerSubs map[string][]int64
}

func NewSubscriber(deps *core.Dependencies) *Subscriber {
	return &Subscriber{
		deps:       deps,
		playerSubs: make(map[string][]int64),
	}
}

func (s *Subscriber) InitSubscribers(ctx context.Context) error {
	subs, err := s.deps.SubscriptionRepo.GetAllSubscription(ctx)
	if err != nil {
		return fmt.Errorf("failed to get all subscription")
	}

	for _, sub := range subs {
		s.playerSubs[sub.Nickname] = append(s.playerSubs[sub.Nickname], sub.ChatID)
	}

	return nil
}
