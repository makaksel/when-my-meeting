package domain

import "time"

type Meeting struct {
	CalendarID string
	Calendar   string

	Title string

	Start    *time.Time
	Location string
}
