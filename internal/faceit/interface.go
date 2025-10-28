package faceit

import (
	"context"

	faceit3 "github.com/anlukk/faceit-tracker/internal/faceit/pkg/go-faceit"
)

type FaceitClient interface {
	GetPlayerByNickname(
		ctx context.Context,
		nickname string) (faceit3.Player, error)
	GetPlayerIDByNickname(
		ctx context.Context,
		nickname string) (string, error)
	GetLastMatch(
		ctx context.Context,
		playerID string) (faceit3.Match, error)
	GetFinishMatchResult(
		ctx context.Context,
		nickname string) (*FinishMatchResult, error)

	GetStatForLastTenMatches(
		ctx context.Context,
		nickname string) ([]faceit3.MatchStats, error)
}
