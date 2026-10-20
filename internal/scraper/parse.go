package scraper

import (
	"OrelTeatrBot/internal/domain"
	"io"
	"time"
)

func Parse(reader io.Reader, now time.Time) ([]domain.Event, error) {
	return make([]domain.Event, 0), nil
}
