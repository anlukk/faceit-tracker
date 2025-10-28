package faceit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	faceit3 "github.com/anlukk/faceit-tracker/internal/faceit/pkg/go-faceit"
	"github.com/antihax/optional"
)

var (
	ErrEmptyNickname = errors.New("empty nickname")
)

type Client struct {
	client *faceit3.APIClient
	token  string
}

func NewClient(apiToken string) (*Client, error) {
	if apiToken == "" {
		return nil, errors.New("faceit api token is empty")
	}

	cfg := faceit3.NewConfiguration()
	cfg.AddDefaultHeader("Authorization", "Bearer "+apiToken)

	return &Client{
		client: faceit3.NewAPIClient(cfg),
		token:  apiToken,
	}, nil
}

func (f *Client) GetPlayerByNickname(
	ctx context.Context,
	nickname string) (faceit3.Player, error) {
	if nickname == "" {
		return faceit3.Player{}, ErrEmptyNickname
	}

	playerID, err := f.GetPlayerIDByNickname(ctx, nickname)
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

func (f *Client) GetPlayerIDByNickname(
	ctx context.Context,
	nickname string) (string, error) {
	if nickname == "" {
		return "", ErrEmptyNickname
	}

	res, _, err := f.client.SearchApi.
		SearchPlayers(ctx, nickname, &faceit3.SearchApiSearchPlayersOpts{
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
	playerID string) (faceit3.Match, error) {
	if playerID == "" {
		return faceit3.Match{}, errors.New("player ID is empty")
	}

	history, _, err := f.client.PlayersApi.
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
	fmt.Println(matchID)
	fmt.Println(matchID)
	fmt.Println(matchID)
	fmt.Println(matchID)
	fmt.Println(matchID)
	fmt.Println(matchID)
	fmt.Println(matchID)
	fmt.Println(matchID)
	fmt.Println(matchID)
	fmt.Println(matchID)
	fmt.Println(matchID)
	fmt.Println(matchID)
	fmt.Println(matchID)
	fmt.Println(matchID)
	fmt.Println(matchID)
	fmt.Println(matchID)
	fmt.Println(matchID)
	fmt.Println(matchID)
	fmt.Println(matchID)
	fmt.Println(matchID)
	if matchID == "" {
		return faceit3.Match{}, fmt.Errorf("match id is empty")
	}

	match, _, err := f.client.MatchesApi.GetMatch(ctx, matchID)
	if err != nil {
		return faceit3.Match{}, fmt.Errorf("get match: %w", err)
	}

	return match, nil
}

func (f *Client) GetFinishMatchResult(
	ctx context.Context,
	nickname string) (*FinishMatchResult, error) {
	if nickname == "" {
		return nil, ErrEmptyNickname
	}

	playerID, err := f.GetPlayerIDByNickname(ctx, nickname)
	if err != nil {
		return nil, fmt.Errorf("get player id: %w", err)
	}

	lastMatch, err := f.GetLastMatch(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("get last match: %w", err)
	}

	if lastMatch.Status != "FINISHED" ||
		lastMatch.Results == nil ||
		lastMatch.Results.Score == nil {
		return nil, fmt.Errorf("match not finished or missing score")
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

	win := lastMatch.Results.Winner == playerTeamKey

	mapName := ""
	if lastMatch.Voting != nil {
		raw, _ := json.Marshal(lastMatch.Voting)
		var v Voting
		_ = json.Unmarshal(raw, &v)
		mapName = v.Map.Voted
	}

	t1, t2, err := f.getRoundScore(ctx, lastMatch.MatchId)
	if err != nil {
		return nil, fmt.Errorf("get round score: %w", err)
	}

	return &FinishMatchResult{
		Nickname:   nickname,
		MatchId:    lastMatch.MatchId,
		Win:        win,
		FinishedAt: lastMatch.FinishedAt,
		Score:      fmt.Sprintf("%d - %d", t1, t2),
		Teams: fmt.Sprintf(
			"%s - %s",
			lastMatch.Teams[opponentTeamKey].Name, lastMatch.Teams[playerTeamKey].Name,
		),
		Map: mapName,
	}, nil
}

func (f *Client) getRoundScore(
	ctx context.Context,
	matchID string) (int, int, error) {
	req, _ := http.NewRequestWithContext(ctx,
		http.MethodGet,
		fmt.Sprintf("https://open.faceit.com/data/v4/matches/%s/stats", matchID),
		nil,
	)
	req.Header.Set("Authorization", "Bearer "+f.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var stats MatchStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return 0, 0, err
	}

	if len(stats.Rounds) == 0 {
		return 0, 0, fmt.Errorf("no rounds data")
	}

	team1 := stats.Rounds[0].Teams[0].TeamStats["Final Score"]
	team2 := stats.Rounds[0].Teams[1].TeamStats["Final Score"]

	t1, _ := strconv.Atoi(team1)
	t2, _ := strconv.Atoi(team2)

	return t1, t2, nil
}

func (f *Client) GetStatForLastTenMatches(
	ctx context.Context,
	nickname string) ([]faceit3.MatchStats, error) {
	if nickname == "" {
		return nil, ErrEmptyNickname
	}

	playerID, err := f.GetPlayerIDByNickname(ctx, nickname)
	if err != nil {
		return nil, fmt.Errorf("get player id: %w", err)
	}

	history, _, err := f.client.
		PlayersApi.
		GetPlayerHistory(ctx, playerID, "cs2", &faceit3.PlayersApiGetPlayerHistoryOpts{
			Limit: optional.NewInt32(10),
		})
	if err != nil {
		return nil, fmt.Errorf("get player history: %w", err)
	}

	if len(history.Items) == 0 {
		return nil, fmt.Errorf("player has no matches")
	}

	matches := make([]faceit3.MatchStats, 0, len(history.Items))

	for _, item := range history.Items[:10] {
		stat, _, err := f.client.MatchesApi.GetMatchStats(ctx, item.MatchId)
		if err != nil {
			return nil, fmt.Errorf("get match stats: %w", err)
		}

		matches = append(matches, stat)

	}

	return matches, nil
}
