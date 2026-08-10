package domain

import "time"

type State struct {
	Meetings    []Meeting
	NextMeeting Meeting
	LastSync    time.Time
	SyncStatus  string
	SyncError   string
}
