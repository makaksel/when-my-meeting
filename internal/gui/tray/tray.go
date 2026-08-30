package tray

import (
	_ "embed"
	"sync"
	"time"

	"context"
	"fmt"
	"makaksel/when-my-meeting/internal/assets"
	"makaksel/when-my-meeting/internal/config"
	"makaksel/when-my-meeting/internal/state"

	"fyne.io/systray"
)

type Service struct {
	Config *config.Service
	State  *state.Service

	cancel       context.CancelFunc
	handleSync   func() error
	openSettings func()

	refreshMu sync.Mutex
}

func New(
	cfg *config.Service,
	s *state.Service,
	cancel context.CancelFunc,
	handleSync func() error,
	openSettings func(),
) *Service {

	return &Service{
		Config:       cfg,
		State:        s,
		cancel:       cancel,
		handleSync:   handleSync,
		refreshMu:    sync.Mutex{},
		openSettings: openSettings,
	}
}

func (s *Service) Run(ctx context.Context) (
	start func(),
	end func(),
) {


	return systray.RunWithExternalLoop(func() {
		s.onReady(ctx)
	}, s.onExit)
}

func (s *Service) onReady(ctx context.Context) {
	systray.SetIcon(assets.Icon)

	go s.listenState(ctx)

	s.refreshMenu()
}

func (s *Service) refreshMenu() {
	s.refreshMu.Lock()
	defer s.refreshMu.Unlock()

	s.updateTitle()
	systray.ResetMenu()
	s.updateMenu()
}

func (s *Service) updateTitle() {
	n := s.State.GetNextMeeting()

	if n == nil || n.Start == nil {
		systray.SetTitle("")
		systray.SetTooltip("Сегодня нет встреч")
		return
	}

	t := fmt.Sprintf("%s", n.Start.Local().Format("15:04"))

	systray.SetTitle(t)
	systray.SetTooltip(fmt.Sprintf("%s %s", t, n.Title))
}

func (s *Service) updateMenu() {
	n := s.State.GetNextMeeting()
	if n != nil && n.Start != nil {
		s.addMeetingItem(n)
		systray.AddSeparator() // ------------
	}

	meetings := s.State.GetFolowingMeetings()
	if len(meetings) > 1 {
		for _, m := range meetings[1:] {
			s.addMeetingItem(&m)
		}
		systray.AddSeparator() // ------------
	}

	s.addBtn("Обновить 🔄", "Обновить", func() {
		s.handleSync()

	})

	s.addBtn("Настройки", "Настройки", func() {
		s.openSettings()
	})

	s.addBtn("Выйти", "Выйти", func() {
		s.onExit()
	})
}

func (s *Service) onExit() {
	now := time.Now()
	fmt.Println("Exit at", now.String())

	systray.Quit()
	s.cancel()
}

func (s *Service) listenState(ctx context.Context) {
	changes := s.State.Changes()

	for {
		select {
		case <-changes:
			s.refreshMenu()

		case <-ctx.Done():
			return
		}
	}
}
