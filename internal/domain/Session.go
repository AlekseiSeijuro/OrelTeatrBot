package domain

type Session struct {
	id        int
	startsAt  string
	actualNow bool
}

func NewSession(id int, startsAt string, actualNow bool) *Session {
	return &Session{id: id, startsAt: startsAt, actualNow: actualNow}
}
