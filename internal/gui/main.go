package gui

import (
	"context"

	"makaksel/when-my-meeting/internal/config"
	"makaksel/when-my-meeting/internal/gui/settings"
	"makaksel/when-my-meeting/internal/gui/tray"
	"makaksel/when-my-meeting/internal/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type Service struct {
	app      fyne.App
	settings *settings.Service
	tray     *tray.Service
}

func New(
	cfg *config.Service,
	state *state.Service,
	cancel context.CancelFunc,
	handleSync func() error,
	deleteCalendar func(string),
) *Service {
	a := app.NewWithID("com.makaksel.when-my-meeting")
	st := settings.New(cfg, a, handleSync, deleteCalendar)

	tr := tray.New(
		cfg,
		state,
		cancel,
		handleSync,
		st.Open,
	)

	return &Service{
		app:      a,
		settings: st,
		tray:     tr,
	}
}

func (g *Service) Run(ctx context.Context) {

	startTray, endTray := g.tray.Run(ctx)
	g.settings.Init()

	startTray()

	go func() {
		<-ctx.Done()
		fyne.Do(func() {
			g.app.Quit()
		})
	}()

	g.app.Run()

	endTray()
}
