package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type BotMessages struct {
	Description            string `yaml:"description"`
	About                  string `yaml:"about"`
	StartTrackingCommand   string `yaml:"start_tracking_message"`
	SubsCommandMessage     string `yaml:"subscriptions_command_message"`
	NicknameForSubs        string `yaml:"nickname_for_subscription"`
	NicknameForUnsubs      string `yaml:"nickname_for_unsubscription"`
	SuccessSubs            string `yaml:"success_subscription"`
	SuccessUnsubs          string `yaml:"success_unsubscription"`
	NotSubscribed          string `yaml:"you_are_not_subscribed"`
	FailedSubs             string `yaml:"failed_subscription"`
	FailedToGetSubs        string `yaml:"failed_to_get_subscriptions"`
	SettingsCommandMessage string `yaml:"settings_command_message"`
	NoSubscriptions        string `yaml:"no_subscriptions"`

	StartCommand struct {
		InlineKeyboard struct {
			KeyboardRow1 struct {
				Options string `yaml:"options"`
			} `yaml:"keyboard_row_1"`

			KeyboardRow3 struct {
				Settings string `yaml:"settings"`
			} `yaml:"keyboard_row_2"`

			KeyboardRow4 struct {
				GitHub string `yaml:"github"`
			} `yaml:"keyboard_row_4"`
		} `yaml:"inline_keyboard"`
	} `yaml:"start_command"`

	SubscriptionsCommand struct {
		InlineKeyboard struct {
			KeyboardRow1 struct {
				AddPlayer string `yaml:"add_player"`
			} `yaml:"keyboard_row_1"`

			KeyboardRow2 struct {
				RemovePlayer string `yaml:"remove_player"`
			} `yaml:"keyboard_row_2"`

			KeyboardRow4 struct {
				List string `yaml:"list"`
			} `yaml:"keyboard_row_4"`

			KeyboardRow5 struct {
				Back string `yaml:"back"`
			} `yaml:"keyboard_row_5"`
		} `yaml:"inline_keyboard"`
	} `yaml:"subscriptions_command"`

	SettingsCommand struct {
		InlineKeyboard struct {
			KeyboardRow1 struct {
				Language string `yaml:"language"`
			} `yaml:"keyboard_row_1"`

			KeyboardRow2 struct {
				Notifications string `yaml:"notification"`
			} `yaml:"keyboard_row_2"`

			KeyboardRow3 struct {
				Back string `yaml:"back"`
			} `yaml:"keyboard_row_3"`
		} `yaml:"inline_keyboard"`
	} `yaml:"settings_command"`
}

func LoadMessages() (BotMessages, error) {
	var command BotMessages

	err := cleanenv.ReadConfig("./locales/en.yaml", &command)
	if err != nil {
		return command, fmt.Errorf("ConfigErr: %v", err)
	}

	return command, nil
}
