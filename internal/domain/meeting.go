package domain

import "time"

type Meeting struct {
	CalendarID string
	Calendar   string

	Title       string
	Description string

	Start    *time.Time
	End      *time.Time
	Location string
}
