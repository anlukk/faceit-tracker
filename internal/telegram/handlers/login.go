package handlers

import (
	//"github.com/anlukk/faceit-tracker/internal/faceit"
	//"github.com/anlukk/faceit-tracker/internal/services"
	"github.com/anlukk/faceit-tracker/internal/telegram/types"
	"github.com/anlukk/faceit-tracker/internal/telegram/utils"
	"github.com/sirupsen/logrus"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type UserLoginCommand struct {
	*types.CommandsOptions
	logger *logrus.Logger

}

func NewUserLoginCommand() *UserLoginCommand {
	return &UserLoginCommand{
		logger: logrus.New(),
		CommandsOptions: types.NewCommandsOptions(),
	}
}

func (user *UserLoginCommand) HandleLoginRequest(
	bot *telego.Bot, update telego.Update) {
	if bot == nil {
			user.logger.Println("Bot instance is nil")
			return
	}

	if update.Message == nil ||
	update.Message.Chat == (telego.Chat{}) ||
	update.Message.Chat.ID == 0 {
			user.logger.Printf("Invalid update structure: Message or Chat is nil")
			return
	}

	userId := tu.ID(update.Message.Chat.ID)
	if user.logger == nil {
			user.logger.Printf("Logger is not initialized for user: %s", userId)
			return
	}

	user.logger.Infof("User %s: %s", userId, update.Message.Text)

	_, botErr := bot.SendMessage(
			tu.Message(userId, "Enter your FACEIT username: ").
					WithReplyMarkup(tu.ForceReply()),
	)
	if botErr != nil {
			user.logger.Errorf("send message error: %v\n", botErr)
	}
}


func (user *UserLoginCommand) HandleUserMessage(
	bot *telego.Bot, update telego.Update) {
	userId := tu.ID(update.Message.From.ID)
	userMessage := update.Message.Text

	user.logger.Infof("User %s: %s", userId, userMessage)

	if userMessage == "cancel" {
		return
	}

	if user.logger == nil {
		user.logger.Printf("Logger is not initialized for user: %s", userId)
		return
	}

	if user.CommandsOptions.Services.FaceitService == nil {
		user.logger.Println("Faceit service is not initialized")
		return
	}

	response, err := user.Services.FaceitService.GetUser(userMessage)

	if err != nil {
		user.logger.Errorf("Error sending API request: %v\n", err)
		_, sendErr := bot.SendMessage(tu.Message(
			userId, "Error fetching data from FACEIT API.",
		).WithParseMode(telego.ModeHTML))
		if sendErr != nil {
			user.logger.Errorf("send message error: %v\n", sendErr)
		}
		return
	}

	formattedResponse := utils.FormatResponseMessage(&response)

	_, err = bot.SendMessage(tu.Message(
		userId, formattedResponse,
	).WithParseMode(telego.ModeHTML))
	if err != nil {
		user.logger.Errorf("send message error: %v\n", err)
		return
	}

	if userMessage == "cancel" {
		user.logger.Infof("User %s canceled", userId)

		_, err := bot.SendMessage(tu.Message(
			userId, "Canceled",
		).WithParseMode(telego.ModeHTML))
		if err != nil {
			user.logger.Errorf("send message error: %v\n", err)
		}
		return
	}

}