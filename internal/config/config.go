package config

import (
	"fmt"
	"makaksel/when-my-meeting/internal/domain"
	"makaksel/when-my-meeting/internal/paths"
	"os"

	"gopkg.in/yaml.v3"
)

type Service struct {
	Paths  *paths.Paths
	Config *domain.Config
}

func New(p *paths.Paths) *Service {
	return &Service{
		Paths: p,
	}
}

func (s *Service) Get() (*domain.Config, error) {
	if s.Config == nil {
		return nil, fmt.Errorf("config is not loaded")
	}
	return s.Config, nil
}

func (s *Service) Load() error {
	err := os.MkdirAll(s.Paths.ConfigDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Если конфиг есть, закгружаем, если нет, создаем базовый
	if _, err := os.Stat(s.Paths.ConfigPath); os.IsNotExist(err) {
		if err := s.Save(domain.DefaultConfig); err != nil {
			return err
		}
	}

	data, err := os.ReadFile(s.Paths.ConfigPath)
	if err != nil {
		return err
	}

	var cfg domain.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}

	s.Config = &cfg

	return nil
}

func (s *Service) Update(cfg *domain.Config) {
	s.Config = cfg
}

func (s *Service) Save(cfg *domain.Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	s.Config = cfg
	return os.WriteFile(s.Paths.ConfigPath, data, 0644)
}
