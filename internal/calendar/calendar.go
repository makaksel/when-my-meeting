package calendar

import (
	"makaksel/when-my-meeting/internal/config"
	"makaksel/when-my-meeting/internal/storage"
)

type Service struct {
	Config  *config.Service
	Storage *storage.Service
}

func New(
	config *config.Service,
	storage *storage.Service,
) *Service {
	return &Service{
		Config:  config,
		Storage: storage,
	}
}
