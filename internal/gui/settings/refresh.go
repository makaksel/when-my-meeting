package settings

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (s *Service) buildRefreshSection() fyne.CanvasObject {
	s.refreshInterval = widget.NewEntry()
	s.refreshInterval.SetText(
		fmt.Sprintf("%d", s.draft.RefreshInterval),
	)

	return container.NewVBox(
		widget.NewLabelWithStyle(
			"Обновление",
			fyne.TextAlignLeading,
			fyne.TextStyle{Bold: true},
		),

		container.NewBorder(
			nil,
			nil,
			widget.NewLabel("Интервал (минуты):"),
			nil,
			s.refreshInterval,
		),
	)
}
