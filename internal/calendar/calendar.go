package calendar

import (
	"context"
	"fmt"
	"log"
	"makaksel/when-my-meeting/internal/config"
	"makaksel/when-my-meeting/internal/domain"
	"makaksel/when-my-meeting/internal/state"
	"makaksel/when-my-meeting/internal/storage"
	"makaksel/when-my-meeting/internal/tray"
	"makaksel/when-my-meeting/internal/utils"
	"os"
	"slices"
	"strings"
	"time"

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

	meetings := s.parseCalendars(cfg.Calendars, cfg.TemporaryFilesPath)

	if len(meetings) != 0 {
		nextMeeting := meetings[0]
		log.Printf("nextMeeting %+v", nextMeeting.Start.Local())
	}

	// log.Printf("parsed %+v meetings from all calendars", meetings)

	// берем календари из конфига
	// проверяем скачанные, удаляем те что нет
	// парсим те что остались
	// обновляем стейт

	// отправляем загружаться те что есть в конфиге из ремута
	// после загрузки обновляем стейт

	return nil
}

func (s *Service) parseCalendars(calendars []domain.Calendar, filesPath string) []domain.Meeting {
	meetings := make([]domain.Meeting, 0, 25)

	for _, calendar := range calendars {

		if !calendar.Enabled {
			continue
		}

		newMeetings, err := s.parseICS(filesPath+calendar.ID+".ics", calendar.ID, calendar.Name)
		if err != nil {
			log.Printf("read calendar %s: %v", calendar.ID, err)
			continue
		}

		meetings = append(meetings, newMeetings...)
	}

	slices.SortFunc(meetings, func(a, b domain.Meeting) int {
		return a.Start.Local().Compare(b.Start.Local())
	})

	return meetings
}

func (s *Service) parseICS(path string, calendarID, calendar string) ([]domain.Meeting, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	gocal.SetTZMapper(func(tzid string) (*time.Location, error) {
		cleanTZID := strings.TrimSpace(tzid)
		cleanTZID = strings.Trim(cleanTZID, `"'`)

		loc, err := time.LoadLocation(cleanTZID)

		if err != nil {
			return nil, fmt.Errorf("unknown timezone: %s", err)
		}

		return loc, nil
	})

	meetings := make([]domain.Meeting, 0, 20)
	c := gocal.NewParser(f)
	c.Parse()

	for _, e := range c.Events {
		if utils.IsToday(e.Start) {
			meetings = append(meetings, domain.Meeting{
				CalendarID: calendarID,
				Calendar:   calendar,
				Title:      e.Summary,
				Start:      e.Start,
				Location:   e.Location,
			})
		}
	}

	return meetings, nil
}
