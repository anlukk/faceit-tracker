package handlers

import (
	"log"

	"github.com/anlukk/faceit-tracker/internal/telegram/types"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/sirupsen/logrus"

	"github.com/anlukk/faceit-tracker/internal/config"

)

type CommandsHandler struct {
	StartCommand  *StartCommand
	UserLogin     *UserLoginCommand
}

type StartCommand struct {
	*types.CommandsOptions
	logger *logrus.Logger
}

func NewStartCommand() *StartCommand {
	return &StartCommand{
		logger: logrus.New(),
	}
}

func MessageError(
	userId telego.ChatID,
	replyToMessageID int,
	message string,
	isReply bool,
) *telego.SendMessageParams {
	inlineKeyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(
				"Cancel",
			).WithCallbackData("cancel"),
		),
	)

	msg := tu.Message(
		userId,
		message,
	).WithReplyMarkup(inlineKeyboard).
		WithParseMode(telego.ModeHTML)

		if isReply {
			msg = msg.WithReplyMarkup(inlineKeyboard).
				WithReplyParameters(
					&telego.ReplyParameters{
						MessageID: replyToMessageID,
					})
		}

	return msg
}

func BuildKeyboard() *telego.InlineKeyboardMarkup {
	messages, err := config.InitCommandsText("../locales/en.yaml")
	if err != nil {
		log.Fatal(err)
	}

	inlineKeyBoard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(
				messages.StartCommand.InlineKeyboard.
					KeyboardRow1.Options,
			).
				WithCallbackData("options"),
		),

		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(
				messages.StartCommand.InlineKeyboard.
				KeyboardRow3.Settings,
			).
				WithCallbackData("settings"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(
				messages.StartCommand.InlineKeyboard.
					KeyboardRow4.GitHub,
			).
				WithCallbackData("github").WithURL(
					"https://github.com/anlukk/faceit-tracker",),
		),
	)
	return inlineKeyBoard
}

func (start *StartCommand) NewStartCommand(bot *telego.Bot, update telego.Update) {
	inlineKeyBoard := BuildKeyboard()

	user_id := tu.ID(update.Message.From.ID)

	messages, err := config.InitCommandsText("../locales/en.yaml")
	if err != nil {
		log.Fatal(err)
	}

	msgTextHello := messages.Description
	if msgTextHello == "" {
		msgTextHello = "Hello, please choose an option"
	}

	message := tu.Message(
		user_id,
		msgTextHello,
	).WithReplyMarkup(inlineKeyBoard).
		WithParseMode(telego.ModeHTML)

	_, sendMsgErr := bot.SendMessage(message)
	if sendMsgErr != nil {
		log.Printf("send message error: %v\n", sendMsgErr)
	}
}

//TODO:
func (start *StartCommand) HandleStartCallback(bot *telego.Bot, update telego.Update) {
	callbackId := update.CallbackQuery.ID
	userId := tu.ID(update.CallbackQuery.From.ID)

	messages, err := config.InitCommandsText("../locales/en.yaml")
	if err != nil {
		log.Fatal(err)
	}

	_, Boterr := bot.SendMessage(tu.Message(
		userId, messages.StartTrackingCommand).
		WithParseMode(telego.ModeHTML))
	if err != nil {
		start.logger.Errorf("send message error: %v\n", Boterr)
	}

	callback := tu.CallbackQuery(callbackId)
	err = bot.AnswerCallbackQuery(callback)
	if err != nil {
		start.logger.Errorf("answer callback error: %v\n", err)
	}
}

func (start *StartCommand) HandleBackCallback(
	bot *telego.Bot, update telego.Update) {
	start.HandleStartCallback(bot, update)
}

func BuildBackKeyboard() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("back").
			WithCallbackData("back"),
		),
	)
}

func (start *StartCommand) HandleAboutCallback(bot *telego.Bot, update telego.Update) {
	messages, err := config.InitCommandsText("../locales/en.yaml")
	if err != nil {
		log.Fatal(err)
	}

	callbackId := update.CallbackQuery.ID
	userId := tu.ID(update.CallbackQuery.From.ID)

	backKeyboard := BuildBackKeyboard()

	_, sendErr := bot.SendMessage(tu.Message(
		userId, messages.About,
	).WithParseMode(telego.ModeHTML).WithReplyMarkup(backKeyboard))
	if sendErr != nil {
		start.logger.Errorf("send message error: %v\n", sendErr)
	}

	if err != nil {
		start.logger.Errorf("Error sending message: %v\n", err)
	}

	callback := tu.CallbackQuery(callbackId)
	err = bot.AnswerCallbackQuery(callback)
	if err != nil {
		start.logger.Errorf("Ошибка при ответе на колбэк: %v\n", err)
	}
}


func (start *StartCommand) HandleExitCallback(bot *telego.Bot, update telego.Update) {
	callbackId := update.CallbackQuery.ID
	userId := tu.ID(update.CallbackQuery.From.ID)

	_, err := bot.SendMessage(tu.Message(
		userId, "Good bye!",
	).WithParseMode(telego.ModeHTML))

	if err != nil {
		start.logger.Errorf("Error sending message: %v\n", err)
	}

	callback := tu.CallbackQuery(callbackId)
	err = bot.AnswerCallbackQuery(callback)
	if err != nil {
		start.logger.Errorf("answer callback error: %v\n", err)
	}
}