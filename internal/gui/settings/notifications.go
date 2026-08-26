package settings

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (s *Service) buildNotificationsSection() fyne.CanvasObject {
	s.notifyActive = widget.NewCheck(
		"Включить уведомления",
		func(value bool) {
			s.draft.Notifications.Active = value
		},
	)

	s.notifyActive.SetChecked(
		s.draft.Notifications.Active,
	)

	s.notifyBefore = widget.NewEntry()
	s.notifyBefore.SetText(
		fmt.Sprintf("%d", s.draft.Notifications.Before),
	)

	return container.NewVBox(
		widget.NewLabelWithStyle(
			"Уведомления",
			fyne.TextAlignLeading,
			fyne.TextStyle{Bold: true},
		),

		s.notifyActive,

		container.NewBorder(
			nil,
			nil,
			widget.NewLabel("За сколько минут до встречи:"),
			nil,
			s.notifyBefore,
		),
	)
}
