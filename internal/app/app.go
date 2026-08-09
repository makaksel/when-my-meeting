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

func New() (*App, error) {
	cfg, err := config.New("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("config is not initialized: %s", err)
	}

	s := state.New()

	strg := storage.New(cfg)

	tr := tray.New(cfg, s)

	n := notification.New(cfg)

	c := calendar.New(cfg, strg, s, tr)

	return &App{
		config:       cfg,
		state:        s,
		storage:      strg,
		calendar:     c,
		notification: n,
		tray:         tr,
		scheduler:    scheduler.New(cfg, c, n),
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.calendar.Sync(ctx)

	a.scheduler.SyncCalendar(ctx)
	a.scheduler.CheckNotifications(ctx)

	// Обработка кнопок
	// событие -> на локешн если он есть
	// рефреш, немедленно перезапускает цикл
	// общая форма для настроек и редактирования календарей
	// выход

	return nil
}
