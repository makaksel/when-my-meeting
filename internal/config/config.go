package config

import (
	"makaksel/when-my-meeting/internal/domain"
	"os"

	"gopkg.in/yaml.v3"
)

type Service struct {
	Path string
}

func New(path string) (*Service, error) {
	s := &Service{
		Path: path,
	}

	if _, err := s.Init(path); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Service) Init(path string) (*Service, error) {
	// Если конфиг есть, закгружаем, если нет, создаем базовый
	if _, err := os.Stat(path); os.IsNotExist(err) {
		defaultConfig := &domain.Config{
			RefreshInterval: 10,
			Notifications: domain.Notify{
				Before: 5,
				Active: true,
			},
		}

		if err := s.Save(defaultConfig); err != nil {
			return nil, err
		}
	}

	return &Service{
		Path: path,
	}, nil
}

func (c *Service) Get() (*domain.Config, error) {
	data, err := os.ReadFile(c.Path)
	if err != nil {
		return nil, err
	}

	var cfg domain.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Service) Save(cfg *domain.Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(c.Path, data, 0644)
}
