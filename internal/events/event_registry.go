package events

import (
	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/anlukk/faceit-tracker/internal/events/cache"
	"github.com/anlukk/faceit-tracker/internal/events/match"
)

func Registry(deps *core.Dependencies) *Controller {
	sources := []EventService{
		match.NewMatchEnd(deps, cache.NewNotifyCache()),
	}

	return NewController(sources...)
}
