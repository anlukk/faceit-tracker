package commands

import (
	"context"
	"fmt"

	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func BuildMainKeyboard(
	deps *core.Dependencies) *telego.InlineKeyboardMarkup {
	inlineKeyBoard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(
				deps.Messages.StartCommand.
					InlineKeyboard.KeyboardRow1.Options,
			).
				WithCallbackData("subscription"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(
				deps.Messages.StartCommand.
					InlineKeyboard.KeyboardRow3.Settings,
			).
				WithCallbackData("settings"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(
				deps.Messages.StartCommand.
					InlineKeyboard.KeyboardRow4.GitHub,
			).
				WithCallbackData("github").WithURL(
				"https://github.com/anlukk/faceit-tracker"),
		),
	)
	return inlineKeyBoard
}

func BuildSubscriptionKeyboard(
	deps *core.Dependencies,
	chatID int64) *telego.InlineKeyboardMarkup {
	ctx := context.Background()
	subs, _ := deps.SubscriptionRepo.GetSubscriptionByChatID(ctx, chatID)

	rows := make([][]telego.InlineKeyboardButton, 0)

	rows = append(rows, tu.InlineKeyboardRow(
		tu.InlineKeyboardButton(
			deps.Messages.SubscriptionsCommand.
				InlineKeyboard.KeyboardRow1.AddPlayer,
		).WithCallbackData("add_player"),
		tu.InlineKeyboardButton(
			deps.Messages.SubscriptionsCommand.
				InlineKeyboard.KeyboardRow2.RemovePlayer,
		).WithCallbackData("remove_player"),
	))

	if len(subs) == 0 {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(
				deps.Messages.SubscriptionsCommand.
					InlineKeyboard.KeyboardRow5.Back,
			).WithCallbackData("back")))
	}

	for _, sub := range subs {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(fmt.Sprintf("%s", sub.Nickname)).
				WithCallbackData(sub.Nickname),
		))
	}

	if len(subs) >= 1 {
		rows = append(rows,
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(
					deps.Messages.SubscriptionsCommand.
						InlineKeyboard.KeyboardRow5.Back,
				).WithCallbackData("back"),
			),
		)
	}

	return &telego.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func BuildSettingsKeyboard(
	deps *core.Dependencies,
	notificationsEnabled bool) *telego.InlineKeyboardMarkup {

	notificationsText := deps.Messages.SettingsCommand.
		InlineKeyboard.KeyboardRow2.Notifications
	if notificationsEnabled {
		notificationsText += " 🔔"
	} else {
		notificationsText += " 🔕"
	}

	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(
				deps.Messages.SettingsCommand.
					InlineKeyboard.KeyboardRow1.Language).
				WithCallbackData("language:1"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(notificationsText).
				WithCallbackData("notification"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(
				deps.Messages.SettingsCommand.
					InlineKeyboard.KeyboardRow3.Back).
				WithCallbackData("back"),
		),
	)
}

func BuildLanguageKeyboard() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("English").
				WithCallbackData("language:en"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Russian").
				WithCallbackData("language:ru"),
		),
	)
}
