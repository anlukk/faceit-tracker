package types

import (
	"github.com/anlukk/faceit-tracker/internal/config"
	"github.com/anlukk/faceit-tracker/internal/services"
)

type CommandsOptions struct {
	Commands *config.CommandsText
	Services *services.Services
}

func NewCommandsOptions() *CommandsOptions {
	services := services.NewServices()
	return &CommandsOptions{
		Commands: &config.CommandsText{},
		Services: &services,
	}
}

type MainMenu struct {
	CurrentPage 		int
	TotalPages  		int
}