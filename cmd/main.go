package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/anlukk/faceit-tracker/internal/config"
	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/anlukk/faceit-tracker/internal/db"
	"github.com/anlukk/faceit-tracker/internal/events"
	"github.com/anlukk/faceit-tracker/internal/faceit"
	"github.com/anlukk/faceit-tracker/internal/notifier"
	"github.com/anlukk/faceit-tracker/internal/telegram"
	"github.com/anlukk/faceit-tracker/internal/telegram/adapters"
	"github.com/anlukk/faceit-tracker/pkg/logger"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	fmt.Println(cfg.LoggerLevel)

	_, sugar, err := logger.BuildLogger(cfg.LoggerLevel)
	if err != nil {
		panic(fmt.Sprintf("failed to build logger: %v", err))
	}
	defer sugar.Sync()

	messages, err := config.LoadMessages()
	if err != nil {
		sugar.Fatalw("failed to load messages",
			"error", err)
	}

	faceitClient, err := faceit.NewClient(cfg.FaceitAPIToken)
	if err != nil {
		sugar.Fatalw("failed to create faceit client",
			"error", err)
	}

	dbConn, err := db.New(cfg)
	if err != nil {
		sugar.Fatalw("failed to initialize postgres",
			"error", err)
	}
	defer func() {
		if err := db.Close(dbConn); err != nil {
			sugar.Errorw("failed to close postgres connection",
				"error", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dependencies := &core.Dependencies{
		Config:           cfg,
		Messages:         &messages,
		Logger:           sugar,
		Faceit:           faceitClient,
		SettingsRepo:     db.NewSettingsDBImpl(dbConn),
		SubscriptionRepo: db.NewSubscriptionDBImpl(dbConn),
		PersonalSubRepo:  db.NewPersonalSubDBImpl(dbConn),
		Ctx:              ctx,
	}

	telegramService, err := telegram.NewTelegram(dependencies)
	if err != nil {
		sugar.Fatalw("failed to initialize telegram service", "error", err)
	}

	if err := telegramService.Start(); err != nil {
		sugar.Fatalw("failed to start telegram service", "error", err)
	}

	messenger := adapters.NewMessengerAdapter(telegramService.Bot())
	eventRegistry := events.Registry(dependencies)
	n := notifier.New(dependencies, messenger, *eventRegistry)

	sugar.Info("Starting notifier...")
	go n.Run(ctx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(
		sigChan,
		os.Interrupt,
		syscall.SIGTERM)
	<-sigChan

	sugar.Info("Shutting down...")
	telegramService.Stop()
	sugar.Info("Application shutdown complete")
}
