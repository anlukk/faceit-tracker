package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type CommandsText struct {
	Description            string `yaml:"description"`
	About                  string `yaml:"about"`
	StartTrackingCommand   string `yaml:"start_tracking_message"`

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
}

func InitCommandsText(path string) (CommandsText, error) {
	var command CommandsText

	err := cleanenv.ReadConfig(path, &command)
	if err != nil {
		return command, fmt.Errorf("ConfigErr: %v", err)
	}

	return command, nil
}