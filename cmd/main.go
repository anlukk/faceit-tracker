package main

import (
	"fmt"
	"github.com/anlukk/faceit-tracker/internal/config"
	"github.com/anlukk/faceit-tracker/internal/core"
	"github.com/anlukk/faceit-tracker/internal/db/postgres"
	"github.com/anlukk/faceit-tracker/internal/faceit"
	"github.com/anlukk/faceit-tracker/internal/service"
	"github.com/anlukk/faceit-tracker/internal/telegram"
	"github.com/anlukk/faceit-tracker/pkg/logger"
	"os"
	"os/signal"
	"syscall"
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

	dbConn, err := postgres.NewDb(cfg)
	if err != nil {
		sugar.Fatalw("failed to initialize postgres",
			"error", err)
	}
	defer func() {
		if err := postgres.Close(dbConn); err != nil {
			sugar.Errorw("failed to close postgres connection",
				"error", err)
		}
	}()

	services := service.NewServices(dbConn, faceitClient)
	dependencies := &core.Dependencies{
		Config:   cfg,
		Messages: &messages,
		Logger:   sugar,
		Db:       dbConn,
		Faceit:   faceitClient,
		Services: services,
	}

	bot, err := telegram.NewTelegram(dependencies)
	if err != nil {
		sugar.Fatalw("failed to initialize telegram bot",
			"error", err)
	}

	if err := bot.Start(); err != nil {
		sugar.Fatalw("failed to start bot",
			"error", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(
		sigChan,
		os.Interrupt,
		syscall.SIGTERM)
	<-sigChan

	sugar.Info("Shutting down...")
	bot.Stop()
	sugar.Info("Application shutdown complete")
}
