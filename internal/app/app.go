package app

import (
	"context"
	"fmt"
	"makaksel/when-my-meeting/internal/calendar"
	"makaksel/when-my-meeting/internal/config"
	"makaksel/when-my-meeting/internal/notification"
	"makaksel/when-my-meeting/internal/scheduler"
	"makaksel/when-my-meeting/internal/state"
	"makaksel/when-my-meeting/internal/storage"
	"makaksel/when-my-meeting/internal/tray"
)

type App struct {
	config       *config.Service
	state        *state.Service
	storage      *storage.Service
	calendar     *calendar.Service
	notification *notification.Service
	tray         *tray.Service
	scheduler    *scheduler.Service
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

	tr := tray.New(cfg, s, cancel, c.SyncRemote)

	sch := scheduler.New(cfg, s, c, n)

	return &App{
		config:       cfg,
		storage:      strg,
		calendar:     c,
		notification: n,
		tray:         tr,
		state:        s,
		scheduler:    sch,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.calendar.SyncLocalOnly()
	a.calendar.SyncRemote()

	a.notification.Check()

	go a.scheduler.Start(ctx)

	a.tray.Run(ctx)

	return nil
}
