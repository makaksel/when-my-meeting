package config

import (
	"fmt"
	"makaksel/when-my-meeting/internal/domain"
	"makaksel/when-my-meeting/internal/paths"
	"os"

	"gopkg.in/yaml.v3"
)

type Service struct {
	Paths *paths.Paths
}

func New(p *paths.Paths) (*Service, error) {
	s := &Service{
		Paths: p,
	}

	err := os.MkdirAll(s.Paths.ConfigDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Если конфиг есть, закгружаем, если нет, создаем базовый
	if _, err := os.Stat(s.Paths.ConfigPath); os.IsNotExist(err) {

		if err := s.Save(domain.DefaultConfig); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Service) Get() (*domain.Config, error) {
	data, err := os.ReadFile(s.Paths.ConfigPath)
	if err != nil {
		return nil, err
	}

	var cfg domain.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (s *Service) Save(cfg *domain.Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(s.Paths.ConfigPath, data, 0644)
}
