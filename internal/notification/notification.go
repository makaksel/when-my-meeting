package notification

import (
	"fmt"
	"log"
	"makaksel/when-my-meeting/internal/assets"
	"makaksel/when-my-meeting/internal/config"
	"makaksel/when-my-meeting/internal/domain"
	"makaksel/when-my-meeting/internal/state"
	"makaksel/when-my-meeting/internal/utils"
	"time"

	"github.com/gen2brain/beeep"
)

type Service struct {
	Config *config.Service
	State  *state.Service
	cache  map[string]bool
}

func New(
	cfg *config.Service,
	state *state.Service,
) *Service {
	beeep.AppName = "When My Meeting"

	return &Service{
		Config: cfg,
		State:  state,
		cache:  make(map[string]bool),
	}
}

func (s *Service) Check() {
	cfg, err := s.Config.Get()
	if err != nil {
		log.Printf("checkNotifications: config is not available: %v", err)
		return
	}

	if !cfg.Notifications.Active || cfg.Notifications.Before == 0 {
		log.Printf("notifications are disabled")
		return
	}

	m := s.State.GetFolowingMeetings()
	for _, meeting := range m {
		if meeting.Start.Local().Before(time.Now().Add(time.Duration(cfg.Notifications.Before) * time.Minute)) {
			if s.cache[makeKey(&meeting)] {
				continue
			}
			s.sendNotification(&meeting)
			s.cache[makeKey(&meeting)] = true
		}
	}
	s.ClearCache()
}

func (s *Service) ClearCache() {
	meetings := s.State.GetFolowingMeetings()

	for k := range s.cache {
		exists := false

		for _, m := range meetings {
			if makeKey(&m) == k {
				exists = true
				break
			}
		}
		if !exists {
			delete(s.cache, k)
		}
	}
}

func makeKey(m *domain.Meeting) string {
	return fmt.Sprintf("%s/%s/%s", m.Start.Local().Format("15:04"), m.End.Local().Format("15:04"), utils.LimitString(m.Title, 5))
}

func (s *Service) sendNotification(m *domain.Meeting) {
	title := fmt.Sprintf("%s-%s %s", m.Start.Local().Format("15:04"), m.End.Local().Format("15:04"), m.Title)
	location := fmt.Sprintf("Место провердения: %s\n\n", m.Location)

	err := beeep.Notify(title, "\n"+location+m.Description, assets.Icon)
	if err != nil {
		log.Printf("send Notification err: %v", err)
	}
}
