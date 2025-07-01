package service

import (
	"github.com/anlukk/faceit-tracker/internal/db"
	"github.com/anlukk/faceit-tracker/internal/service/sub"
)

type Services struct {
	Sub *sub.Service
}

func NewServices(repos *db.Repositories) *Services {
	return &Services{
		Sub: sub.NewService(repos.Sub),
	}
}
