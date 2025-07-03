package core

import (
	"github.com/anlukk/faceit-tracker/internal/config"
	"github.com/anlukk/faceit-tracker/internal/faceit"
	"github.com/anlukk/faceit-tracker/internal/service"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Dependencies struct {
	Config   *config.Config
	Messages *config.BotMessages
	Logger   *zap.SugaredLogger
	Faceit   faceit.FaceitClient
	Services *service.Services
	Db       *gorm.DB
}
