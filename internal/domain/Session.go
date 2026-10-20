package domain

import "time"

type Session struct {
	id        int
	startsAt  time.Time
	actualNow bool
}
