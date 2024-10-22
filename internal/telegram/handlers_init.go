package telegram

import (
	"github.com/anlukk/faceit-tracker/internal/telegram/handlers"
	th "github.com/mymmrac/telego/telegohandler"
)

func (service *TelegramService) handlersInit() error {
	handlersInit := handlers.CommandsHandler{
		StartCommand: handlers.NewStartCommand(),
		UserLogin:    handlers.NewUserLoginCommand(),
	}

	service.Handlers.Handle(
		handlersInit.StartCommand.NewStartCommand,
		th.CommandEqual("start"),
	)

	service.Handlers.Handle(
		handlersInit.StartCommand.HandleAboutCallback,
		th.CallbackDataEqual("about"),
	)

	//TODO:
	// service.Handlers.Handle(
	// 	handlersInit.StartCommand.HandleBackCallback,
	// 	th.CallbackDataEqual("back"),
	// )

	// service.Handlers.Handle(
		// 	handlersInit.StartCommand.HandleStartCallback,
		// 	th.CallbackDataEqual("cancel"),
		// )

	service.Handlers.Handle(
		handlersInit.StartCommand.HandleExitCallback,
		th.CallbackDataEqual("exit"),
	)

	// service.Handlers.Handle(
	// 	handlersInit.StartCommand.HandleStartCallback,
	// 	th.CallbackDataEqual("start_tracking"),
	// )

	// Login commands handler
	service.Handlers.Handle(
		handlersInit.UserLogin.HandleLoginRequest,
		th.TextEqual("/login"),
	)

	service.Handlers.Handle(
		handlersInit.UserLogin.HandleUserMessage,
		th.AnyMessage(),
	)

	return nil
}