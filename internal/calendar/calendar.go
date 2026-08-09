package calendar

import (
	"context"
	"log"
	"makaksel/when-my-meeting/internal/config"
	"makaksel/when-my-meeting/internal/state"
	"makaksel/when-my-meeting/internal/storage"
	"makaksel/when-my-meeting/internal/tray"
	"os"

	"github.com/apognu/gocal"
)

type Service struct {
	Config  *config.Service
	Storage *storage.Service
	State   *state.Service
	Tray    *tray.Service
}

func New(
	cfg *config.Service,
	strg *storage.Service,
	s *state.Service,
	t *tray.Service,
) *Service {
	return &Service{
		Config:  cfg,
		Storage: strg,
		State:   s,
		Tray:    t,
	}
}

func (s *Service) Sync(ctx context.Context) error {
	cfg, err := s.Config.Get()
	if err != nil {
		return err
	}

	for _, calendar := range cfg.Calendars {

		if !calendar.Enabled {
			continue
		}

		filePath := cfg.TemporaryFilesPath + calendar.ID + ".ics"
		_, err := s.ParseICV(ctx, filePath)
		if err != nil {
			log.Printf("read calendar %s: %v", calendar.ID, err)
			continue
		}

	}

	// берем календари из конфига
	// проверяем скачанные, удаляем те что нет
	// парсим те что остались
	// обновляем стейт

	// отправляем загружаться те что есть в конфиге из ремута
	// после загрузки обновляем стейт

	return nil
}

func (s *Service) ParseICV(ctx context.Context, path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	c := gocal.NewParser(f)
	c.Parse()
	for _, e := range c.Events {
		log.Printf("%s at %s in %s", e.Summary, e.Start, e.Location)
	}

	return "", nil
}
