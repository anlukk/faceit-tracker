package main

import (
	"github.com/anlukk/faceit-tracker/internal/config"
	//"github.com/anlukk/faceit-tracker/internal/faceit"
	"github.com/anlukk/faceit-tracker/internal/logger"
	"github.com/anlukk/faceit-tracker/internal/telegram"
)

func main() {
	logg := logger.InitLogger()

	cfg, err := config.NewConfig()
	if err != nil {
		logg.Panicf("config init: %v", err)
	}

	_, err = config.InitCommandsText("../locales/en.yaml")
	if err != nil {
		logg.Panicf("commands init: %v", err)
	}

	//TODO: remove config from services
	// services := &services.Services{
	// 	Logger:           logg,
	// }

	// newFaceit, err := faceit.NewFaceit(cfg.FaceitAPIToken)
	// if err != nil {
	// 	logg.Panicf("faceit service init: %v", err)
	// }


	bot, err := telegram.NewTelegram(cfg.TelegramToken, logg)
	if err != nil {
		logg.Panicf("telegram service init: %v", err)
	}


	bot.StartService()
}