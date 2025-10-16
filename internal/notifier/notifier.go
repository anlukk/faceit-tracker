package notifier

import (
	"context"
	"sync"
	"time"

	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/anlukk/faceit-tracker/internal/events"
	"github.com/anlukk/faceit-tracker/internal/events/types"
)

const (
	notificationTicker = 3 * time.Minute
	numWorkers         = 20
)

type Notifier struct {
	deps      *core.Dependencies
	messenger Messenger

	controller events.Controller
}

func New(
	deps *core.Dependencies,
	messenger Messenger,
	controller events.Controller) *Notifier {

	n := &Notifier{
		deps:       deps,
		messenger:  messenger,
		controller: controller,
	}

	return n
}

func (n *Notifier) Run(ctx context.Context) {
	ticker := time.NewTicker(notificationTicker)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			collectEvents, err := n.controller.CollectEvents(ctx)
			if err != nil {
				n.deps.Logger.Errorw("failed to collect events", "error", err)
				continue
			}

			eventsChannel := make(chan types.Event, len(collectEvents))
			var wg sync.WaitGroup

			for i := 0; i < numWorkers; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for event := range eventsChannel {
						if err := n.messenger.SendMessage(event.ChatID, event.Message); err != nil {
							n.deps.Logger.Errorw("failed to send message", "error", err)
						}
					}
				}()
			}

			for _, collectEvent := range collectEvents {
				eventsChannel <- collectEvent
			}

			close(eventsChannel)
			wg.Wait()
		}
	}
}
