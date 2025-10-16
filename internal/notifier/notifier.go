package notifier

import (
	"context"
	"time"

	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/anlukk/faceit-tracker/internal/events"
)

const notificationTicker = 3 * time.Minute

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

			for _, event := range collectEvents {
				err := n.messenger.SendMessage(event.ChatID, event.Message)
				if err != nil {
					n.deps.Logger.Errorw("failed to send message", "error", err)
					continue
				}
			}
		}
	}
}
