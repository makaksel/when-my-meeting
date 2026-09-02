package app

import (
	"context"
	"fmt"
	"makaksel/when-my-meeting/internal/calendar"
	"makaksel/when-my-meeting/internal/config"
	"makaksel/when-my-meeting/internal/gui"
	"makaksel/when-my-meeting/internal/gui/tray"
	"makaksel/when-my-meeting/internal/notification"
	"makaksel/when-my-meeting/internal/paths"
	"makaksel/when-my-meeting/internal/scheduler"
	"makaksel/when-my-meeting/internal/state"
	"makaksel/when-my-meeting/internal/storage"
)

type App struct {
	config       *config.Service
	state        *state.Service
	storage      *storage.Service
	calendar     *calendar.Service
	notification *notification.Service
	tray         *tray.Service
	scheduler    *scheduler.Service
	gui          *gui.Service
}

func New(cancel context.CancelFunc) (*App, error) {
	p, err := paths.New("when-my-meeting")
	if err != nil {
		return nil, fmt.Errorf("paths is not initialized: %s", err)
	}

	cfg, err := config.New(p)
	if err != nil {
		return nil, fmt.Errorf("config is not initialized: %s", err)
	}
	// Debug
	cfgD, err := cfg.Get()
	if err != nil {
		return nil, fmt.Errorf("config is not initialized: %s", err)
	}

	fmt.Println("config:", cfgD)
	fmt.Println("config initialized")
	// Debug 
	s := state.New()

	strg := storage.New(cfg, p)

	n := notification.New(cfg, s)

	c := calendar.New(cfg, s, strg, p)

	sch := scheduler.New(cfg, s, c, n)

	gui := gui.New(cfg, s, cancel, c.SyncRemote, strg.DeleteCalendar)

	return &App{
		storage:      strg,
		calendar:     c,
		notification: n,
		scheduler:    sch,
		gui:          gui,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.storage.Init()

	a.calendar.SyncLocalOnly()
	a.calendar.SyncRemote()

	a.notification.Check()

	go a.scheduler.Start(ctx)

	a.gui.Run(ctx)

	return nil
}
