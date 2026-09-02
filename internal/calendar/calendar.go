package calendar

import (
	"fmt"
	"log"
	"makaksel/when-my-meeting/internal/config"
	"makaksel/when-my-meeting/internal/domain"
	"makaksel/when-my-meeting/internal/paths"
	"makaksel/when-my-meeting/internal/state"
	"makaksel/when-my-meeting/internal/storage"
	"makaksel/when-my-meeting/internal/utils"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/apognu/gocal"
)

type Service struct {
	Config  *config.Service
	State   *state.Service
	Storage *storage.Service
	Paths   *paths.Paths
}

func New(
	cfg *config.Service,
	s *state.Service,
	strg *storage.Service,
	p *paths.Paths,
) *Service {
	return &Service{
		Config:  cfg,
		Storage: strg,
		State:   s,
		Paths:   p,
	}
}

func (s *Service) SyncLocalOnly() error {
	cfg, err := s.Config.Get()
	if err != nil {
		return err
	}

	newMeetings := s.parseCalendars(cfg.Calendars)

	s.State.UpdateMeetings(newMeetings)

	return nil
}

func (s *Service) SyncRemote() error {
	cfg, err := s.Config.Get()
	if err != nil {
		return err
	}

	s.loadCalendars(&cfg.Calendars)

	s.Config.Save(cfg) // TODO подумать и вынести в память

	return s.SyncLocalOnly()
}

func (s *Service) parseCalendars(calendars []domain.Calendar) []domain.Meeting {
	meetings := make([]domain.Meeting, 0, 25)

	for _, calendar := range calendars {
		newMeetings, err := s.parseICS(filepath.Join(s.Paths.DataDir, calendar.Name+".ics"), calendar)
		if err != nil {
			log.Printf("read calendar %s: %v", calendar.Name, err)
			continue
		}

		meetings = append(meetings, newMeetings...)
	}

	return meetings
}

func (s *Service) parseICS(path string, calendar domain.Calendar) ([]domain.Meeting, error) {
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
				Calendar:    calendar.Name,
				Title:       e.Summary,
				Description: e.Description,
				Start:       e.Start,
				End:         e.End,
				Location:    e.Location,
				Disabled:    !calendar.Enabled,
			})
		}
	}

	return meetings, nil
}

func (s *Service) loadCalendars(calendars *[]domain.Calendar) {
	for i, calendar := range *calendars {
		err := s.Storage.LoadCalendar(calendar.URL, calendar.Name, calendar.User, calendar.Password)
		if err != nil {
			log.Printf("load calendar %s: %v", calendar.Name, err)
			(*calendars)[i].Error = fmt.Sprintf("%v", err)
			continue
		}
		(*calendars)[i].Error = ""
		(*calendars)[i].LastSync = time.Now().Unix()
	}

}
