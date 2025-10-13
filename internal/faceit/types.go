package faceit

type Voting struct {
	Map `json:"map"`
}

type Map struct {
	Voted string `json:"voted"`
}

type MatchStats struct {
	Rounds []struct {
		Teams []struct {
			TeamID    string            `json:"team_id"`
			TeamStats map[string]string `json:"team_stats"`
		} `json:"teams"`
	} `json:"rounds"`
}

type OngoingMatchInfo struct {
	Nickname string
	MatchID  string
	Map      string
	Elo      bool
	Team1    string
	Team2    string
	StartAt  int64
}

type FinishMatchResult struct {
	Nickname string
	MatchId  string
	Win      bool
	Elo      bool
	Score    string
	Teams    string
	Map      string
}
