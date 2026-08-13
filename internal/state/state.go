package state

import (
	"makaksel/when-my-meeting/internal/domain"
	"sync"
)

type Service struct {
	mu sync.RWMutex
	domain.State
}

func New() *Service {

	return &Service{}
}
