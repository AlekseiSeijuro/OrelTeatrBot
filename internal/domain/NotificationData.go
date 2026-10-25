package domain

type NotificationData struct {
	ChatID     int64
	EventID    int64
	EventTitle string
	SessionID  int64
	DateText   string
	AvailEpoch int
}
