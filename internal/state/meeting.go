package state

import (
	"makaksel/when-my-meeting/internal/domain"
	"slices"
)

func (s *Service) GetMeetings() []domain.Meeting {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.State.Meetings
}
func (s *Service) GetNextMeeting() domain.Meeting {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.State.Meetings) == 0 {
		return domain.Meeting{}
	}
	return s.State.Meetings[0]
}

func (s *Service) UpdateMeetings(meetings []domain.Meeting) {
	s.mu.Lock()
	defer s.mu.Unlock()

	slices.SortFunc(meetings, func(a, b domain.Meeting) int {
		return a.Start.Local().Compare(b.Start.Local())
	})

	s.State.Meetings = meetings
}
