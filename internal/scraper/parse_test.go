package scraper

import (
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	f, err := os.Open("testdata/page.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	events, err := Parse(f)

	if err != nil {
		t.Errorf("Error parsing events: %v", err)
	}
	if events == nil {
		t.Errorf("No events parsed")
	}
	if len(events) != 16 {
		t.Errorf("Wrong number of events: %d", len(events))
	}
}
