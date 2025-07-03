package faceit

import (
	"context"
	faceit3 "github.com/anlukk/faceit-tracker/internal/faceit/pkg/go-faceit"
)

type FaceitClient interface {
	GetPlayer(
		ctx context.Context,
		username string) (faceit3.Player, error)
	GetPlayerIDByUsername(
		ctx context.Context,
		username string) (string, error)
	GetLastMatch(
		ctx context.Context,
		username string) (faceit3.Match, error)
}
