package calendar

import (
	"fmt"
	"log"
	"makaksel/when-my-meeting/internal/config"
	"makaksel/when-my-meeting/internal/domain"
	"makaksel/when-my-meeting/internal/state"
	"makaksel/when-my-meeting/internal/storage"
	"makaksel/when-my-meeting/internal/utils"
	"os"
	"strings"
	"time"

	"github.com/apognu/gocal"
)

type Service struct {
	Config  *config.Service
	State   *state.Service
	Storage *storage.Service
}

func New(
	cfg *config.Service,
	s *state.Service,
	strg *storage.Service,
) *Service {
	return &Service{
		Config:  cfg,
		Storage: strg,
		State:   s,
	}
}

func (s *Service) SyncLocalOnly() error {
	cfg, err := s.Config.Get()
	if err != nil {
		return err
	}

	newMeetings := s.parseCalendars(s.onlyEnabled(cfg.Calendars), cfg.TemporaryFilesPath)
	s.State.UpdateMeetings(newMeetings)

	return nil
}

func (s *Service) SyncRemote() error {
	cfg, err := s.Config.Get()
	if err != nil {
		return err
	}

	s.loadCalendars(s.onlyEnabled(cfg.Calendars))

	return s.SyncLocalOnly()
}

func (s *Service) parseCalendars(calendars []domain.Calendar, filesPath string) []domain.Meeting {
	meetings := make([]domain.Meeting, 0, 25)

	for _, calendar := range calendars {
		newMeetings, err := s.parseICS(filesPath+calendar.Name+".ics", calendar.Name, calendar.Name)
		if err != nil {
			log.Printf("read calendar %s: %v", calendar.Name, err)
			continue
		}

		meetings = append(meetings, newMeetings...)
	}

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

func (s *Service) loadCalendars(calendars []domain.Calendar) {
	for _, calendar := range calendars {
		err := s.Storage.LoadCalendar(calendar.URL, calendar.Name, calendar.User, calendar.Password)
		if err != nil {
			log.Printf("load calendar %s: %v", calendar.Name, err)
			continue
		}
	}
}

func (s *Service) onlyEnabled(calendars []domain.Calendar) []domain.Calendar {
	filtered := make([]domain.Calendar, 0, len(calendars))
	for _, calendar := range calendars {
		if calendar.Enabled {
			filtered = append(filtered, calendar)
		}
	}

	return filtered
}
