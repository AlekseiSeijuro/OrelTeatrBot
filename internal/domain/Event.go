package domain

type Event struct {
	ID       int64
	Title    string
	Sessions []Session
}

func NewEvent(id int64, title string, sessions []Session) *Event {
	return &Event{ID: id, Title: title, Sessions: sessions}
}
