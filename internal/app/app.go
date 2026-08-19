package app

import (
	"context"
	"fmt"
	"makaksel/when-my-meeting/internal/calendar"
	"makaksel/when-my-meeting/internal/config"
	"makaksel/when-my-meeting/internal/gui"
	"makaksel/when-my-meeting/internal/gui/tray"
	"makaksel/when-my-meeting/internal/notification"
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
	cfg, err := config.New("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("config is not initialized: %s", err)
	}

	s := state.New()

	strg := storage.New(cfg)

	n := notification.New(cfg, s)

	c := calendar.New(cfg, s, strg)

	sch := scheduler.New(cfg, s, c, n)

	gui := gui.New(cfg, s, cancel, c.SyncRemote)

	return &App{
		config:       cfg,
		storage:      strg,
		calendar:     c,
		notification: n,
		state:        s,
		scheduler:    sch,
		gui:          gui,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.calendar.SyncLocalOnly()
	a.calendar.SyncRemote()

	a.notification.Check()

	go a.scheduler.Start(ctx)

	a.gui.Run(ctx)

	return nil
}
