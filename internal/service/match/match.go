package match

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anlukk/faceit-tracker/internal/faceit"
)

var ErrEmptyUsername = errors.New("empty username")

type Service struct {
	client faceit.FaceitClient
}

func NewService(client faceit.FaceitClient) *Service {
	return &Service{client: client}
}

func (s *Service) GetOngoingMatchInfo(
	ctx context.Context,
	username string) (*OngoingMatchInfo, error) {
	if username == "" {
		return nil, ErrEmptyUsername
	}

	lastMatch, err := s.client.GetLastMatch(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("get last match: %w", err)
	}

	if lastMatch.Status != "ONGOING" && lastMatch.Status != "READY" {
		return nil, nil
	}

	teams := make([]string, 0, 2)
	for _, team := range lastMatch.Teams {
		teams = append(teams, team.Name)
	}

	return &OngoingMatchInfo{
		MatchID: lastMatch.MatchId,
		Elo:     lastMatch.CalculateElo,
		StartAt: lastMatch.StartedAt,
		Team1:   teams[0],
		Team2:   teams[1],
	}, nil
}

func (s *Service) GetFinishMatchResult(
	ctx context.Context,
	username string) (*FinishMatchResult, error) {
	if username == "" {
		return nil, ErrEmptyUsername
	}

	lastMatch, err := s.client.GetLastMatch(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("get last match: %w", err)
	}

	if lastMatch.Status != "FINISHED" ||
		lastMatch.Results == nil ||
		lastMatch.Results.Score == nil {
		return nil, fmt.Errorf("match not finished or missing score")
	}

	playerID, err := s.client.GetPlayerIDByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("get player id: %w", err)
	}

	var playerTeamKey, opponentTeamKey string

	for teamKey, team := range lastMatch.Teams {
		for _, player := range team.Roster {
			if player.PlayerId == playerID {
				playerTeamKey = teamKey
			}
		}
	}

	for key := range lastMatch.Teams {
		if key != playerTeamKey {
			opponentTeamKey = key
		}
	}

	if playerTeamKey == "" {
		return nil, fmt.Errorf("cannot find player team")
	}

	if opponentTeamKey == "" {
		return nil, fmt.Errorf("cannot find opponent team")
	}

	playerScore := lastMatch.Results.Score[playerTeamKey]
	opponentScore := lastMatch.Results.Score[opponentTeamKey]
	win := lastMatch.Results.Winner == playerTeamKey

	mapName := ""
	if lastMatch.Voting != nil {
		raw, _ := json.Marshal(lastMatch.Voting)
		var v Voting
		_ = json.Unmarshal(raw, &v)
		mapName = v.Map.Voted
	}

	return &FinishMatchResult{
		MatchID:  lastMatch.MatchId,
		Win:      win,
		Score:    fmt.Sprintf("%d - %d", playerScore, opponentScore),
		Opponent: lastMatch.Teams[opponentTeamKey].Name,
		Map:      mapName,
	}, nil
}
