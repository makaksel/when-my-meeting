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

	// Если конфиг есть, закгружаем, если нет, создаем базовый
	if _, err := os.Stat(path); os.IsNotExist(err) {

		if err := s.Save(domain.DefaultConfig); err != nil {
			return nil, err
		}
	}

	return s, nil
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
