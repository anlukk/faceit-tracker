package commands

import (
	"github.com/anlukk/faceit-tracker/internal/core"
)

type Notification struct {
	deps *core.Dependencies
}

func NewNotification(deps *core.Dependencies) *Notification {
	return &Notification{
		deps: deps,
	}
}

//func (n *Notification) HandlePreMatchNotifications(bot *telego.Bot, update telego.Update) {
//	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//	defer cancel()
//
//	enable, err := n.deps.
//		Services.
//		Settings.
//		GetNotificationsEnabled(ctx, update.CallbackQuery.From.ID)
//	if err != nil {
//		n.deps.Logger.Errorw("failed to get notifications enabled", "error", err)
//		return
//	}
//
//	if !enable {
//		return
//	}
//
//}
//
//func (n *Notification) HandleMatchResultNotifications(bot *telego.Bot, update telego.Update) {
//	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//	defer cancel()
//
//	n.deps.Services.Notifications.GetFinishMatchResult(ctx, update.CallbackQuery.From.Username)
//}
