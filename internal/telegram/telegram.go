package telegram

import (
	"fmt"
	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/anlukk/faceit-tracker/internal/telegram/commands"
	"github.com/anlukk/faceit-tracker/internal/telegram/menu"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"go.uber.org/zap"
	"sync"
)

type Telegram struct {
	bot         *telego.Bot
	handlers    *th.BotHandler
	commands    *commands.BotCommands
	deps        *core.Dependencies
	menuManager *menu.MenuManager

	wg       sync.WaitGroup
	stopChan chan struct{}
}

func NewTelegram(deps *core.Dependencies) (*Telegram, error) {
	bot, err := telego.NewBot(
		deps.Config.TelegramToken,
		telego.WithLogger(deps.Logger),
		telego.WithDefaultLogger(true, true),
	)
	if err != nil {
		return nil, fmt.Errorf("telegram service init: %v", err)
	}

	updates, err := bot.UpdatesViaLongPolling(nil)
	if err != nil {
		return nil, fmt.Errorf("updates error: %v", err)
	}

	botHandlers, err := th.NewBotHandler(bot, updates)
	if err != nil {
		return nil, fmt.Errorf("bot commands error: %v", err)
	}

	newMenu := menu.NewMenuManager(deps.Logger)
	newCommands := commands.NewBotCommands(deps, newMenu)

	service := &Telegram{
		bot:         bot,
		handlers:    botHandlers,
		menuManager: newMenu,
		commands:    newCommands,
		deps:        deps,
		stopChan:    make(chan struct{}),
	}

	botHandlers.Use(
		func(bot *telego.Bot, update telego.Update, next th.Handler) {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						service.
							deps.
							Logger.
							Error("panic in handler", zap.Any("panic", r))
					}
				}()
				next(bot, update)
			}()
		})

	err = service.registerCommands()
	if err != nil {
		return nil, fmt.Errorf("register commands: %v", err)
	}

	return service, nil
}

func (t *Telegram) registerCommands() error {
	t.handlers.Handle(
		t.commands.StartCommand.StartCommand,
		th.CommandEqual("start"),
	)

	t.handlers.Handle(
		t.commands.StartCommand.HandleSubscriptionMenuCallback,
		th.CallbackDataPrefix("menu:"),
	)

	t.handlers.Handle(
		t.commands.StartCommand.HandleSubscriptionMenuCallback,
		th.CallbackDataEqual("subscription"),
	)

	t.handlers.Handle(
		t.commands.StartCommand.HandleBackCallback,
		th.CallbackDataEqual("back"),
	)

	t.handlers.Handle(
		t.commands.StartCommand.HandleSettingsMenuCallback,
		th.CallbackDataEqual("settings"),
	)

	t.handlers.Handle(
		t.commands.StartCommand.HandleNotificationToggleCallback,
		th.CallbackDataEqual("notification"),
	)

	t.handlers.Handle(
		t.commands.Subscription.HandleSubscribeButton,
		th.CallbackDataEqual("add_player"),
	)
	t.handlers.Handle(
		t.commands.Subscription.HandleSubscriptionUsernameReply,
		commands.IsSubscriptionReplyMessage(),
	)

	t.handlers.Handle(
		t.commands.Subscription.HandleUnsubscribeButton,
		th.CallbackDataEqual("remove_player"),
	)
	t.handlers.Handle(
		t.commands.Subscription.HandleUnsubscriptionUsernameReply,
		commands.IsUnsubscriptionReplyMessage(),
	)

	t.handlers.Handle(
		t.commands.Subscription.HandleSubscriptionsListButton,
		th.CallbackDataEqual("list"),
	)

	t.handlers.Handle(
		t.commands.SearchPlayerCommand.PromptPlayerSearch,
		th.TextEqual("/searchplayer"),
	)
	t.handlers.Handle(
		t.commands.SearchPlayerCommand.HandleUserMessage,
		th.AnyMessage(),
	)

	return nil
}

func (t *Telegram) Start() error {
	t.deps.Logger.Info("Telegram bot started")

	t.wg.Add(1)
	go func() {
		defer t.wg.Done()

		t.handlers.Start()
		<-t.stopChan
		t.handlers.Stop()
	}()

	return nil
}

func (t *Telegram) Stop() {
	close(t.stopChan)
	t.wg.Wait()
	t.deps.Logger.Info("Telegram bot stopped")
}
