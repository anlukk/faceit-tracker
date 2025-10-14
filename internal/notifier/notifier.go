package notifier

import (
	"context"
	"time"

	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/anlukk/faceit-tracker/internal/notifier/cache"
	"github.com/anlukk/faceit-tracker/internal/notifier/event_handlers"
)

const notificationTicker = 3 * time.Minute

type Notifier struct {
	deps          *core.Dependencies
	messenger     Messenger
	eventHandlers []EventHandlers
}

func New(deps *core.Dependencies, messenger Messenger) *Notifier {
	n := &Notifier{
		deps:          deps,
		messenger:     messenger,
		eventHandlers: []EventHandlers{},
	}

	matchCache := cache.NewMatchCache()
	n.RegisterEventHandlers(event_handlers.NewMatchEnd(deps, matchCache))

	return n
}

func (n *Notifier) RegisterEventHandlers(eventHandlers EventHandlers) {
	if eventHandlers == nil {
		n.deps.Logger.Errorw("event handlers is nil")
	}

	n.eventHandlers = append(n.eventHandlers, eventHandlers)
}

func (n *Notifier) Run(ctx context.Context) {
	ticker := time.NewTicker(notificationTicker)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for _, eventHandler := range n.eventHandlers {
				events, err := eventHandler.GetEvents(ctx)
				if err != nil {
					n.deps.Logger.Errorw("failed to get events", "error", err)
					continue
				}

				for _, event := range events {
					err := n.messenger.SendMessage(event.ChatID, event.Message)
					if err != nil {
						n.deps.Logger.Errorw("failed to send message", "error", err)
						continue
					}
				}

			}
		}
	}
}
