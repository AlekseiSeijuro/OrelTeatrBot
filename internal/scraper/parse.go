package scraper

import (
	"OrelTeatrBot/internal/domain"
	"log"
	"strconv"
	"strings"

	"io"

	"github.com/PuerkitoBio/goquery"
)

func Parse(reader io.Reader) ([]domain.Event, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	logIfError(err, "Error parsing document")
	events := make([]domain.Event, 0)

	doc.Find(".elem.quick-tickets-event-card").
		Each(func(i int, eventCard *goquery.Selection) {
			eventRef, _ := eventCard.Find(".quick-tickets-heading-link").Attr("href")
			eventId, err := strconv.Atoi(eventRef[strings.LastIndex(eventRef, "e")+1:])
			logIfError(err, "Could not parse event ref")
			title := eventCard.Find("span.underline").Text()

			sessions := make([]domain.Session, 0)
			eventCard.Find(".quick-tickets-session-option").Each(func(i int, sessionCard *goquery.Selection) {
				sessionRef, _ := sessionCard.Find(".sessions__date-link").Attr("href")
				sessionId, err := strconv.Atoi(sessionRef[strings.LastIndex(sessionRef, "s")+1:])
				logIfError(err, "Could not parse session ref")

				sessionDateText := sessionCard.Find(".quick-tickets-session-option__label").Text()

				actualNow := sessionCard.Find(".session-sold-label").Length() == 0

				sessions = append(sessions, *domain.NewSession(sessionId, sessionDateText, actualNow))
			})
			events = append(events, *domain.NewEvent(eventId, title, sessions))
		})

	return events, nil
}

func logIfError(err error, msg string) {
	if err != nil {
		log.Print(msg, err)
	}
}
