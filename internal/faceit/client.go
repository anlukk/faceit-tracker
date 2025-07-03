package faceit

import (
	"context"
	"errors"
	"fmt"
	faceit3 "github.com/anlukk/faceit-tracker/internal/faceit/pkg/go-faceit"
	"github.com/antihax/optional"
)

var (
	ErrEmptyUsername = errors.New("empty username")
)

type Client struct {
	client *faceit3.APIClient
}

func NewClient(apiToken string) (*Client, error) {
	if apiToken == "" {
		return nil, errors.New("faceit api token is empty")
	}

	cfg := faceit3.NewConfiguration()
	cfg.AddDefaultHeader("Authorization", "Bearer "+apiToken)

	client := faceit3.NewAPIClient(cfg)
	return &Client{
		client: client,
	}, nil
}

func (f *Client) GetPlayer(
	ctx context.Context,
	username string) (faceit3.Player, error) {
	if username == "" {
		return faceit3.Player{}, ErrEmptyUsername
	}

	playerID, err := f.GetPlayerIDByUsername(ctx, username)
	if err != nil {
		return faceit3.Player{}, fmt.Errorf("get user id: %w", err)
	}

	player, _, err := f.client.PlayersApi.GetPlayer(ctx, playerID)
	if err != nil {
		return player, fmt.Errorf("get player: %w", err)
	}

	if player.PlayerId == "" {
		return player, fmt.Errorf("player not found")
	}

	return player, nil
}

func (f *Client) GetPlayerIDByUsername(
	ctx context.Context,
	username string) (string, error) {
	if username == "" {
		return "", ErrEmptyUsername
	}

	res, _, err := f.client.
		SearchApi.
		SearchPlayers(ctx, username, &faceit3.SearchApiSearchPlayersOpts{
			Limit: optional.NewInt32(1),
		})
	if err != nil {
		return "", fmt.Errorf("search player: %w", err)
	}

	if len(res.Items) == 0 {
		return "", fmt.Errorf("player not found")
	}

	return res.Items[0].PlayerId, nil
}

func (f *Client) GetLastMatch(
	ctx context.Context,
	username string) (faceit3.Match, error) {
	if username == "" {
		return faceit3.Match{}, ErrEmptyUsername
	}

	playerID, err := f.GetPlayerIDByUsername(ctx, username)
	if err != nil {
		return faceit3.Match{}, fmt.Errorf("get player id: %w", err)
	}

	history, _, err := f.client.
		PlayersApi.
		GetPlayerHistory(ctx, playerID, "cs2", &faceit3.PlayersApiGetPlayerHistoryOpts{
			Limit: optional.NewInt32(1),
		})
	if err != nil {
		return faceit3.Match{}, fmt.Errorf("get player history: %w", err)
	}

	if len(history.Items) == 0 {
		return faceit3.Match{}, fmt.Errorf("player has no matches")
	}

	matchID := history.Items[0].MatchId
	if matchID == "" {
		return faceit3.Match{}, fmt.Errorf("match id is empty")
	}

	match, _, err := f.client.MatchesApi.GetMatch(ctx, matchID)
	if err != nil {
		return faceit3.Match{}, fmt.Errorf("get match: %w", err)
	}

	return match, nil
}
