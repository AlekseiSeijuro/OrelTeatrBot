package domain

type Session struct {
	ID        int64
	DateText  string
	ActualNow bool
}

func NewSession(id int64, startsAt string, actualNow bool) *Session {
	return &Session{ID: id, DateText: startsAt, ActualNow: actualNow}
}
