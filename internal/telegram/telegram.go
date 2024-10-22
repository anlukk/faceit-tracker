package telegram

import (
	"os/signal"
	"syscall"
	"os"

	"sync"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/sirupsen/logrus"

	"github.com/anlukk/faceit-tracker/internal/services"
)

type TelegramService struct {
	Bot               *telego.Bot
	Handlers 		  		*th.BotHandler
	Logger            *logrus.Logger
	Services          services.Services

	updates           chan telego.Update
	stop              chan struct{}
	done              chan struct{}
	wg                *sync.WaitGroup

}


func NewTelegram(token string, logg *logrus.Logger) (*TelegramService, error) {
	bot, err := telego.NewBot(
		token,
		telego.WithDefaultDebugLogger(),
		)
	if err != nil {
		return nil, err
	}

	updates, err := bot.UpdatesViaLongPolling(nil)
	if err != nil {
		return nil, err
	}


	botHandlers, err := th.NewBotHandler(bot, updates)
	if err != nil {
		return nil, err
	}

	botHandlers.Group()

	done := make(chan struct{}, 1)
	stop := make(chan struct{}, 1)

	// botHandlers.Use(
	// 	func(bot *telego.Bot, update telego.Update, next th.Handler) {
	// 		go func() {
	// 			defer func() {
	// 				if r := recover(); r != nil {
	// 					logg.Error(r)
	// 				}
	// 			}()
	// 			next(bot, update)
	// 		}()
	// 	},
	// )

	// botServices = types.BotServices{
	// 	Config: &cfg,
	// }

	TelegramService := &TelegramService{
		Bot:         bot,
		Handlers:    botHandlers,
		Logger:      logg,
		stop:        stop,
		done:        done,
		wg:          &sync.WaitGroup{},
	}

	return TelegramService, nil
}

func (service *TelegramService) StartService() error {
	service.handleStopSignal()

	service.handlersInit()

	signals := make(chan os.Signal, 1)
	signal.Notify(
		signals,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	go func() {
		service.wg.Add(1)
		defer service.wg.Done()

		service.Handlers.Start()
	} ()

	service.Logger.Info("Telegram service started")

	select {
	case sig := <-signals:
		service.Logger.Infof("Received signal: %v", sig)
		service.stop <- struct{}{}
	case <-service.done:
	}

	service.wg.Wait()
	service.Logger.Info("Telegram service stopped")

	if service.Handlers.IsRunning() {
	service.
		Logger.
		Fatal("Telegram service not stopped")
		return nil
	}

	return nil
}


func (service *TelegramService) handleStopSignal() {
	go func() {
		<-service.stop

		service.Logger.Info("Stopping the bot...")

		service.Bot.StopLongPolling()
		service.Logger.Info("Long polling stopped")

		service.Handlers.Stop()
		service.Logger.Info("Handler stopped")

		service.done <- struct{}{}
	}()
}