package domain

type Event struct {
	id       int
	title    string
	sessions []Session
}

func NewEvent(id int, title string, sessions []Session) *Event {
	return &Event{id: id, title: title, sessions: sessions}
}
