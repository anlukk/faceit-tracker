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

type Client interface {
	GetUser(ctx context.Context, username string) (faceit3.Player, error)
	GetUserIDByUsername(ctx context.Context, username string) (string, error)
}

type ClientImpl struct {
	client *faceit3.APIClient
}

func NewClientImpl(apiToken string) (*ClientImpl, error) {
	if apiToken == "" {
		return nil, errors.New("faceit api token is empty")
	}

	cfg := faceit3.NewConfiguration()
	cfg.AddDefaultHeader("Authorization", "Bearer "+apiToken)
	client := faceit3.NewAPIClient(cfg)

	return &ClientImpl{
		client: client,
	}, nil
}

func (f *ClientImpl) GetUser(ctx context.Context, username string) (faceit3.Player, error) {
	if username == "" {
		return faceit3.Player{}, ErrEmptyUsername
	}

	playerID, err := f.GetUserIDByUsername(ctx, username)
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

func (f *ClientImpl) GetUserIDByUsername(ctx context.Context, username string) (string, error) {
	if username == "" {
		return "", ErrEmptyUsername
	}

	res, _, err := f.client.SearchApi.SearchPlayers(ctx, username, &faceit3.SearchApiSearchPlayersOpts{
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
