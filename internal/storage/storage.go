package storage

import (
	"fmt"
	"io"
	"log"
	"makaksel/when-my-meeting/internal/config"
	"makaksel/when-my-meeting/internal/domain"
	"makaksel/when-my-meeting/internal/paths"
	"net/http"
	"os"
	"path/filepath"
	"slices"
)

type Service struct {
	Config *config.Service
	Paths  *paths.Paths
}

func New(c *config.Service, p *paths.Paths) *Service {
	return &Service{Config: c, Paths: p}
}

func (s *Service) Init() error {
	err := os.MkdirAll(s.Paths.DataDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	s.ClearUnused()

	return nil
}

func (s *Service) LoadCalendar(url, fileName, username, password string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(username, password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad HTTP status: %s", resp.Status)
	}

	filePath := filepath.Join(s.Paths.DataDir, fileName+".ics")

	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file on disk: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file contents: %w", err)
	}

	return nil
}

func (s *Service) DeleteCalendar(fileName string) {
	err := os.Remove(filepath.Join(s.Paths.DataDir, fileName))
	if err != nil {
		return
	}
}

func (s *Service) ClearUnused() {
	cfg, err := s.Config.Get()
	if err != nil {
		return
	}

	icsFiles, err := os.ReadDir(s.Paths.DataDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range icsFiles {
		if file.IsDir() {
			continue
		}

		existsInConf := slices.ContainsFunc(cfg.Calendars, func(c domain.Calendar) bool {
			return c.Name+".ics" == file.Name()
		})

		if !existsInConf {
			s.DeleteCalendar(file.Name())
		}
	}

}
