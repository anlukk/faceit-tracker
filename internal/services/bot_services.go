package services

import (
	"github.com/sirupsen/logrus"

	"github.com/anlukk/faceit-tracker/internal/faceit"
	"github.com/anlukk/faceit-tracker/internal/config"
)

type Services struct {
	Logger 				*logrus.Logger
	FaceitService *faceit.FaceitService
	config 				*config.Config
}

func NewServices() Services {
	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	faceit, err := faceit.NewFaceit(config.FaceitAPIToken)
	if err != nil {
		panic(err)
	}

	return Services{
		Logger: logrus.New(),
		FaceitService: faceit,
	}
}