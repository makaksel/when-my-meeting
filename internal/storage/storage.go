package storage

import (
	"fmt"
	"io"
	"makaksel/when-my-meeting/internal/config"
	"net/http"
	"os"
	"path/filepath"
)

type Service struct {
	Config *config.Service
}

func New(config *config.Service) *Service {
	return &Service{Config: config}
}

func (s *Service) LoadCalendar(fileUrl, fileName, username, password string) error {
	cfg, err := s.Config.Get()
	if err != nil {
		return err
	}

	err = os.MkdirAll(cfg.TemporaryFilesPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	req, err := http.NewRequest("GET", fileUrl, nil)
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

	filePath := filepath.Join(cfg.TemporaryFilesPath, fileName+".ics")
	fmt.Printf("Saving file as: %s\n", filePath)

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

func (s *Service) DeleteCalendar(id string) {

}

func (s *Service) ListCalendars() {

}
