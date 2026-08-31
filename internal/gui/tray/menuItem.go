package tray

import (
	"fmt"
	"makaksel/when-my-meeting/internal/domain"
	"makaksel/when-my-meeting/internal/utils"

	"fyne.io/fyne/v2"
	"fyne.io/systray"
)

func (s *Service) addBtn(title, tooltip string, callback func()) *systray.MenuItem {
	btn := systray.AddMenuItem(title, tooltip)

	go func() {
		for range btn.ClickedCh {
			fyne.Do(func() {
				callback()
			})
		}
	}()

	return btn
}

func (s *Service) addMeetingItem(m *domain.Meeting) {
	title := fmt.Sprintf("%s-%s %s", m.Start.Local().Format("15:04"), m.End.Local().Format("15:04"), m.Title)

	s.addBtn(title, title, func() {
		if m.Location != "" {
			utils.OpenBrowser(m.Location)
		}
	})
}
