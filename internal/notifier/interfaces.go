package notifier

import (
	"context"

	"github.com/anlukk/faceit-tracker/internal/notifier/event_handlers"
)

type EventHandlers interface {
	EventType() string
	GetEvents(ctx context.Context) ([]event_handlers.Event, error)
}

type Messenger interface {
	SendMessage(chatID int64, text string) error
}
