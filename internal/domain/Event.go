package domain

type Event struct {
	id        int
	title     string
	ageRating string
	sessions  []Session
}
