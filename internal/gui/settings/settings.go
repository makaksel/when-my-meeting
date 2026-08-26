package settings

import (
	"makaksel/when-my-meeting/internal/assets"
	"makaksel/when-my-meeting/internal/config"
	"makaksel/when-my-meeting/internal/domain"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Service struct {
	Config *config.Service

	app    fyne.App
	window fyne.Window

	handleSync     func() error
	deleteCalendar func(string)

	draft *domain.Config

	refreshInterval *widget.Entry
	notifyBefore    *widget.Entry
	notifyActive    *widget.Check
}

func New(
	cfg *config.Service,
	app fyne.App,
	handleSync func() error,
	deleteCalendar func(string),
) *Service {

	return &Service{
		Config:         cfg,
		app:            app,
		handleSync:     handleSync,
		deleteCalendar: deleteCalendar,
	}
}

func (s *Service) Init() {
	s.createWindow()

	s.window.SetCloseIntercept(func() {
		s.draft = nil
		s.window.Hide()
	})
}

func (s *Service) Open() {
	cfg, err := s.Config.Get()
	if err != nil {
		return
	}
	s.draft = cfg

	if s.window == nil {
		return
	}

	fyne.Do(func() {
		s.createMenu()
		s.window.Show()
		s.window.RequestFocus()
	})
}

func (s *Service) createWindow() {
	s.window = s.app.NewWindow("When My Meeting")

	s.window.SetIcon(
		fyne.NewStaticResource(
			"icon.png",
			assets.Icon,
		),
	)

	s.window.CenterOnScreen()
	s.window.Resize(fyne.NewSize(500, 500))
}

func (s *Service) createMenu() {
	s.window.SetContent(
		container.NewBorder(
			nil,
			container.NewPadded(
				widget.NewButton("Сохранить", func() {
					s.save()
				}),
			),
			nil,
			nil,
			container.NewPadded(
				container.NewVBox(
					s.buildCalendarsSection(),
					widget.NewSeparator(),
					s.buildRefreshSection(),
					widget.NewSeparator(),
					s.buildNotificationsSection(),
				),
			),
		),
	)
}

func (s *Service) rebuild() {
	s.createMenu()
}

func (s *Service) save() {
	refresh, err := strconv.Atoi(s.refreshInterval.Text)
	if err != nil || refresh <= 0 {
		return
	}

	before, err := strconv.Atoi(s.notifyBefore.Text)
	if err != nil || before < 0 {
		return
	}

	s.draft.RefreshInterval = refresh
	s.draft.Notifications.Before = before
	s.draft.Notifications.Active = s.notifyActive.Checked

	if err := s.Config.Save(s.draft); err != nil {
		// показать ошибку
		return
	}

	s.window.Hide()

	s.handleSync()
}
