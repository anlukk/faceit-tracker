package notifier

import (
	"context"
	"fmt"
	"time"

	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/anlukk/faceit-tracker/internal/faceit"
)

type Notifier struct {
	deps       *core.Dependencies
	messenger  Messenger
	playerSubs map[string][]int64

	cache *MatchCache
}

func New(deps *core.Dependencies, messenger Messenger) *Notifier {
	cache := NewMatchCache()

	return &Notifier{
		deps:       deps,
		messenger:  messenger,
		playerSubs: make(map[string][]int64),
		cache:      cache,
	}
}

func (n *Notifier) NotifyEndMatch(ctx context.Context) {
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := n.initSubscribers(ctx); err != nil {
				n.deps.Logger.Error("failed to init subscribers", err)
			}
			for nickname := range n.playerSubs {
				n.deps.Logger.Debugw("start processing subscriber", "nickname", nickname)

				match, err := n.deps.Faceit.GetLastMatch(ctx, nickname)
				if err != nil {
					n.deps.Logger.Errorw(
						"failed to get last match",
						"nickname", nickname,
						"error", err)
					continue
				}
				n.deps.Logger.Debugw(
					"checking match status",
					"nickname", nickname,
					"status", match.Status,
					"matchID", match.MatchId,
					"finished time", match.FinishedAt,
					"alreadyNotified", n.cache.alreadyNotified(nickname, match.MatchId),
				)
				if match.Status == "FINISHED" && !n.cache.alreadyNotified(nickname, match.MatchId) {
					finishedAt := time.Unix(match.FinishedAt, 0)
					if time.Since(finishedAt) > 10*time.Minute {
						n.deps.Logger.Debugw(
							"skipping match notification, match too old",
							"nickname", nickname,
							"matchID", match.MatchId,
							"finishedAt", finishedAt,
						)
						continue
					}

					matchFinish, err := n.deps.Faceit.GetFinishMatchResult(ctx, nickname)
					if err != nil {
						n.deps.Logger.Errorw(
							"failed to get finish match result",
							"nickname", nickname,
							"error", err)
						continue
					}

					for _, chatID := range n.playerSubs[matchFinish.Nickname] {
						if err := n.messenger.SendMessage(chatID, n.formatInfoMatchFinish(matchFinish)); err != nil {
							n.deps.Logger.Errorw("send message error", "error", err)
						}
					}

					n.cache.markNotified(nickname, matchFinish.MatchId)
				}
			}
		}
	}
}

func (n *Notifier) initSubscribers(ctx context.Context) error {
	subs, err := n.deps.SubscriptionRepo.GetAllSubscription(ctx)
	if err != nil {
		return fmt.Errorf("failed to get all subscription")
	}

	for _, sub := range subs {
		n.playerSubs[sub.Nickname] = append(n.playerSubs[sub.Nickname], sub.ChatID)
	}

	return nil
}

func (n *Notifier) formatInfoMatchFinish(
	matchInfo *faceit.FinishMatchResult) string {
	isWin := map[bool]string{
		true:  n.deps.Messages.MatchWin,
		false: n.deps.Messages.MatchLoose,
	}

	result := fmt.Sprintf(
		"%s\n%s%s\n%s\n%s%s\n",
		n.deps.Messages.MatchFinish,
		n.deps.Messages.Nickname,
		matchInfo.Nickname,
		isWin[matchInfo.Win],
		n.deps.Messages.MatchScore,
		matchInfo.Score,
	)

	return result
}
