package notifier

import (
	"sync"
	"time"

	"github.com/anlukk/faceit-tracker/internal/core"
)

type SubsInMemoryStorage struct {
	deps       *core.Dependencies
	playerSubs map[string][]int64
	timeCheck  map[string]time.Time

	mu sync.RWMutex
}

func NewSubscriber(deps *core.Dependencies) *SubsInMemoryStorage {
	return &SubsInMemoryStorage{
		deps:       deps,
		playerSubs: make(map[string][]int64),
	}
}

//
//func (s *SubsInMemoryStorage) InitSubscribers(ctx context.Context) error {
//	subs, err := s.deps.SubscriptionRepo.GetAllSubscription(ctx)
//	if err != nil {
//		return fmt.Errorf("failed to get all subscription")
//	}
//
//	for _, sub := range subs {
//		s.playerSubs[sub.Nickname] = append(s.playerSubs[sub.Nickname], sub.ChatID)
//	}
//
//	return nil
//}
//
//func (s *SubsInMemoryStorage) StartPeriodicSync(ctx context.Context, interval time.Duration) {
//	ticker := time.NewTicker(interval)
//	go func() {
//		for {
//			select {
//			case <-ctx.Done():
//				return
//			case <-ticker.C:
//				_ = s.InitSubscribers(ctx)
//			}
//		}
//	}()
//}
//
//func (s *SubsInMemoryStorage) LastChecked(nickname string) (time.Time, error) {
//	s.mu.RLock()
//	_, ok := s.timeCheck[nickname]
//	s.mu.RUnlock()
//	if ok {
//		return s.timeCheck[nickname], nil
//	}
//
//	s.mu.Lock()
//	defer s.mu.Unlock()
//	if _, ok := s.timeCheck[nickname]; !ok {
//		s.timeCheck[nickname] = time.Now()
//	}
//
//	return s.timeCheck[nickname], nil
//}
