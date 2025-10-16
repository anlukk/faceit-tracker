package events

import (
	"context"

	"github.com/anlukk/faceit-tracker/internal/events/types"
)

type Controller struct {
	services []EventService
}

func NewController(services ...EventService) *Controller {
	return &Controller{services: services}
}

func (d *Controller) CollectEvents(ctx context.Context) ([]types.Event, error) {
	var events []types.Event

	for _, s := range d.services {
		evs, err := s.GetEvents(ctx)
		if err != nil {
			return nil, err
		}

		events = append(events, evs...)
	}

	return events, nil
}
