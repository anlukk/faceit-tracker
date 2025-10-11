package faceit

import (
	"context"

	faceit3 "github.com/anlukk/faceit-tracker/internal/faceit/pkg/go-faceit"
)

type FaceitClient interface {
	GetPlayer(
		ctx context.Context,
		nickname string) (faceit3.Player, error)
	GetPlayerIDByNickname(
		ctx context.Context,
		nickname string) (string, error)
	GetLastMatch(
		ctx context.Context,
		nickname string) (faceit3.Match, error)
	GetOngoingMatchInfo(
		ctx context.Context,
		nickname string) (*OngoingMatchInfo, error)
	GetFinishMatchResult(
		ctx context.Context,
		nickname string) (*FinishMatchResult, error)
}
