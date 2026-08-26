package settings

import (
	"errors"
	"fmt"
	"makaksel/when-my-meeting/internal/domain"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (s *Service) buildCalendarsSection() fyne.CanvasObject {
	items := container.NewVBox()

	if len(s.draft.Calendars) == 0 {
		items.Add(
			widget.NewLabel("Календари не добавлены"),
		)
	}

	for i := range s.draft.Calendars {
		items.Add(s.calendarItem(&s.draft.Calendars[i], &i))
	}

	items.Add(
		widget.NewButton(
			"+ Добавить",
			func() {
				s.openCalendarForm(nil)
			},
		),
	)

	return container.NewVBox(
		widget.NewLabelWithStyle(
			"Календари",
			fyne.TextAlignLeading,
			fyne.TextStyle{Bold: true},
		),
		items,
	)
}

func (s *Service) calendarItem(
	calendar *domain.Calendar,
	index *int,
) fyne.CanvasObject {
	enabled := widget.NewCheck("", func(value bool) {
		calendar.Enabled = value
	})

	enabled.SetChecked(calendar.Enabled)

	errorLabel := fmt.Sprintf(" (⚠️ %s)", calendar.Error)
	textLabel := calendar.Name

	if calendar.Error != "" {
		textLabel += errorLabel
	}

	label := widget.NewLabel(textLabel)
	label.Wrapping = fyne.TextWrapWord

	edit := widget.NewButton("Изменить", func() {
		s.openCalendarForm(index)
	})
	remove := widget.NewButton("Удалить", func() {
		s.removeCalendar(index)
	})

	return container.NewBorder(
		nil,
		nil,
		enabled,
		container.NewHBox(
			edit,
			remove,
		),
		container.NewVBox(
			label,
		),
	)
}

func (s *Service) openCalendarForm(index *int) {
	var calendar domain.Calendar

	// Добавление
	if index == nil {
		calendar.Enabled = true
	} else {
		// Редактирование — копируем из draft
		calendar = s.draft.Calendars[*index]
	}

	window := s.app.NewWindow(
		map[bool]string{true: "Редактировать календарь", false: "Добавить календарь"}[index != nil],
	)

	window.Resize(fyne.NewSize(450, 400))

	name := widget.NewEntry()
	name.SetText(calendar.Name)
	name.SetPlaceHolder("Например: Work")

	url := widget.NewEntry()
	url.SetText(calendar.URL)
	url.SetPlaceHolder("https://example.com/calendar.ics")

	user := widget.NewEntry()
	user.SetText(calendar.User)

	password := widget.NewPasswordEntry()
	password.SetText(calendar.Password)

	enabled := widget.NewCheck("Календарь включен", nil)
	enabled.SetChecked(calendar.Enabled)

	form := widget.NewForm(
		widget.NewFormItem("Название", name),
		widget.NewFormItem("URL", url),
		widget.NewFormItem("Пользователь", user),
		widget.NewFormItem("Пароль", password),
	)

	formContainer := container.NewVBox(
		form,
		enabled,
	)

	buttonText := "Добавить"

	if index != nil {
		buttonText = "Сохранить"
	}

	buttons := container.NewHBox(
		widget.NewButton("Отмена", func() {
			window.Close()
		}),

		widget.NewButton(buttonText, func() {
			calendar := domain.Calendar{
				Name:     name.Text,
				URL:      url.Text,
				User:     user.Text,
				Password: password.Text,
				Enabled:  enabled.Checked,
			}

			if calendar.Name == "" {
				name.SetValidationError(
					errors.New("поле обязательно"),
				)
				form.Refresh()
				return
			} else {
				name.SetValidationError(nil)
			}

			if calendar.URL == "" {
				url.SetValidationError(
					errors.New("поле обязательно"),
				)
				form.Refresh()
				return
			} else {
				url.SetValidationError(nil)
			}

			if index == nil {
				// ADD
				s.draft.Calendars = append(
					s.draft.Calendars,
					calendar,
				)
			} else {
				// EDIT
				s.draft.Calendars[*index] = calendar
			}

			window.Close()
			s.rebuild()
		}),
	)

	window.SetContent(
		container.NewBorder(
			nil,
			container.NewPadded(buttons),
			nil,
			nil,
			container.NewPadded(formContainer),
		),
	)

	window.Show()
}

func (s *Service) removeCalendar(index *int) {
	if *index < 0 || *index >= len(s.draft.Calendars) {
		return
	}

	calendar := s.draft.Calendars[*index]

	dialog.ShowConfirm(
		"Удалить календарь?",
		"Календарь «"+calendar.Name+"» будет удалён.",
		func(ok bool) {
			if !ok {
				return
			}

			s.draft.Calendars = append(
				s.draft.Calendars[:*index],
				s.draft.Calendars[*index+1:]...,
			)

			s.deleteCalendar(calendar.Name + ".ics")

			s.rebuild()
		},
		s.window,
	)

}
