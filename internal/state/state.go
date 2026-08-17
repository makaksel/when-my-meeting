package state

import (
	"makaksel/when-my-meeting/internal/domain"
	"sync"
)

type Service struct {
	mu sync.RWMutex
	domain.State

	changed chan struct{}
}

func New() *Service {

	return &Service{
		changed: make(chan struct{}, 1),
	}
}

func (s *Service) Changes() <-chan struct{} {
	return s.changed
}
