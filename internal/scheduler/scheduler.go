package scheduler

import (
	"context"
	"fmt"
	"log"
	"makaksel/when-my-meeting/internal/calendar"
	"makaksel/when-my-meeting/internal/config"
	"makaksel/when-my-meeting/internal/notification"
	"makaksel/when-my-meeting/internal/state"
	"makaksel/when-my-meeting/internal/tray"
	"time"
)

type Service struct {
	Config       *config.Service
	State        *state.Service
	Calendar     *calendar.Service
	Notification *notification.Service
	Tray         *tray.Service
}

func New(
	cfg *config.Service,
	s *state.Service,
	c *calendar.Service,
	n *notification.Service,
	t *tray.Service,
) *Service {
	return &Service{Config: cfg, State: s, Calendar: c, Notification: n, Tray: t}
}

func (s *Service) Start(ctx context.Context) {
	cfg, err := s.Config.Get()
	if err != nil {
		log.Printf("scheduler: config is not available: %v", err)
		return
	}

	// Воркер для синхронизации календарей
	if cfg.RefreshInterval != 0 {
		go s.worker(ctx, time.Duration(cfg.RefreshInterval)*time.Minute, "syncCalendars", func() {
			s.Calendar.SyncRemote(ctx)
		})
	}

	// Воркер для проверки уведомлений
	go s.worker(ctx, time.Minute, "checkNotifications", func() {
		s.Notification.Check()
	})

	// Воркер для обновления трея
	go s.worker(ctx, time.Second*30, "refreshTray", func() {
		s.Tray.Refresh(ctx)
	})

}

func (s *Service) worker(ctx context.Context, dur time.Duration, name string, task func()) {
	ticker := time.NewTicker(dur)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			task()
		case <-ctx.Done():
			fmt.Printf("scheduler: %s stopped\n", name)
			return
		}
	}
}

func (s *Service) StopAll(ctx context.Context) error {
	ctx.Done()

	return nil
}
