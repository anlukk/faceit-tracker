package events

import (
	"context"

	"github.com/anlukk/faceit-tracker/internal/events/types"
)

type EventService interface {
	EventType() string
	GetEvents(ctx context.Context) ([]types.Event, error)
}
