package main

import (
	"OrelTeatrBot/internal/scraper"
	"bytes"
	"context"
	"fmt"
)

func main() {

	f := scraper.NewFetcher("https://quicktickets.ru/orel-teatr-svobodnoe-prostranstvo")
	data, _ := f.Fetch(context.Background())
	events, _ := scraper.Parse(bytes.NewReader(data))
	fmt.Print(events)
}
