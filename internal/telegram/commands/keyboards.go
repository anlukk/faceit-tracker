package commands

import (
	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/anlukk/faceit-tracker/internal/db/models"
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
	subs []models.Subscription,
	mainNickname string) *telego.InlineKeyboardMarkup {
	rows := make([][]telego.InlineKeyboardButton, 0)

	for _, sub := range subs {
		nickname := sub.Nickname
		if sub.Nickname == mainNickname {
			nickname += " ⭐"
		}

		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(nickname).
				WithCallbackData("player:"+sub.Nickname),
		))
	}

	rows = append(rows, tu.InlineKeyboardRow(
		tu.InlineKeyboardButton(
			deps.Messages.SubscriptionsCommand.
				InlineKeyboard.KeyboardRow1.AddPlayer,
		).WithCallbackData("add_player"),
		tu.InlineKeyboardButton(
			deps.Messages.SubscriptionsCommand.
				InlineKeyboard.KeyboardRow2.RemovePlayer,
		).WithCallbackData("remove_player"),
	),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(
				deps.Messages.SubscriptionsCommand.
					InlineKeyboard.KeyboardRow3.NewPersonalSub,
			).WithCallbackData("create_new_personal_sub"),

			tu.InlineKeyboardButton(
				deps.Messages.SubscriptionsCommand.
					InlineKeyboard.KeyboardRow5.Back,
			).WithCallbackData("back"),
		))

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

//func BuildLanguageKeyboard() *telego.InlineKeyboardMarkup {
//	return tu.InlineKeyboard(
//		tu.InlineKeyboardRow(
//			tu.InlineKeyboardButton("English").
//				WithCallbackData("language:en"),
//		),
//		tu.InlineKeyboardRow(
//			tu.InlineKeyboardButton("Russian").
//				WithCallbackData("language:ru"),
//		),
//	)
//}
