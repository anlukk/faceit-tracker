package faceit

import (
	"fmt"

	// "github.com/mitchellh/mapstructure"
	"github.com/mconnat/go-faceit/pkg/client"
	"github.com/mconnat/go-faceit/pkg/models"
)

type FaceitService struct {
	client client.FaceITClient
}

func NewFaceit(apiToken string) (*FaceitService, error) {
	if apiToken == "" {
		return nil, fmt.Errorf("empty Bearer token")
	}

	client, err := client.New(apiToken)
	if err != nil {
		return nil, err
	}

	return &FaceitService{
		client: client,
		}, nil
}

// func (c *Client) GetUserViaSearch(username string) (models.Player, error) {
// 	user, err := c.client.SearchPlayers(username, nil)
// 	if err != nil {
// 		return models.Player{}, err
// 	}

// 	if len(user.Items) == 0 {
// 		return models.Player{}, fmt.Errorf("user not found")
// 	}

// 	return user.Items[0], nil
// }

func (f *FaceitService) GetUser(username string) (models.Player, error) {
	player, err := f.client.GetPlayer(map[string]interface{}{
		"nickname": username,
	})
	if err != nil {
		return models.Player{}, err
	}

	if player.PlayerId == "" {
		return models.Player{}, fmt.Errorf("player id is empty")
	}

	return player, nil
}

func (f *FaceitService) GetUserMatches(playerID string) (models.PlayerHistory, error) {
	if playerID == "" {
		return models.PlayerHistory{}, fmt.Errorf("empty player id")
	}

	matchesHistory, err := f.client.GetPlayerHistory(
		playerID,
		"game",
		map[string]interface{}{
    "game": "10",
	})

	if err != nil {
		return models.PlayerHistory{}, err
	}

	if len(matchesHistory.Items) == 0 {
		return models.PlayerHistory{}, fmt.Errorf("no matches found")
	}


	return matchesHistory, nil
}

