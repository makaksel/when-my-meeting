package state

import (
	"makaksel/when-my-meeting/internal/domain"
	"makaksel/when-my-meeting/internal/utils"
	"slices"
)

func (s *Service) GetMeetings() []domain.Meeting {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.State.Meetings
}

func (s *Service) GetFolowingMeetings() []domain.Meeting {
	s.mu.RLock()
	defer s.mu.RUnlock()

	meetings := s.State.Meetings

	folowing := make([]domain.Meeting, 0, len(meetings))

	for i := range meetings {
		if utils.IsFolowing(meetings[i].Start) {
			folowing = append(folowing, meetings[i])
		}
	}

	if len(folowing) == 0 {
		return nil
	}
	return folowing
}

func (s *Service) GetNextMeeting() *domain.Meeting {
	s.mu.RLock()
	defer s.mu.RUnlock()

	m := s.GetFolowingMeetings()

	if len(m) == 0 {
		return &domain.Meeting{}
	}
	return &m[0]
}

func (s *Service) UpdateMeetings(meetings []domain.Meeting) {
	s.mu.Lock()

	for i := range meetings {
		meetings[i] = normalizeMeeting(meetings[i])
	}

	slices.SortFunc(meetings, func(a, b domain.Meeting) int {
		return a.Start.Local().Compare(b.Start.Local())
	})

	s.State.Meetings = meetings

	s.mu.Unlock()
	select {
	case s.changed <- struct{}{}:
	default:
	}
}

func normalizeMeeting(m domain.Meeting) domain.Meeting {
	if u := utils.ParseURL(m.Location); u != nil {
		m.Location = u.String()
	}

	return m
}
