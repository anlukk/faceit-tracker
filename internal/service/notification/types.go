package notification

type Voting struct {
	Map `json:"map"`
}

type Map struct {
	Voted string `json:"voted"`
}

type OngoingMatchInfo struct {
	MatchID string
	Map     string
	Elo     bool
	Team1   string
	Team2   string
	StartAt int64
}

type FinishMatchResult struct {
	MatchID  string
	Win      bool
	Score    string
	Opponent string
	Map      string
}
