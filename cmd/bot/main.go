package main

import (
	"OrelTeatrBot/internal/scraper"
	"os"
	"time"
)

func main() {
	f, err := os.Open("C:\\Users\\User\\GolandProjects\\OrelTeatrBot\\internal\\scraper\\testdata\\page.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scraper.Parse(f, time.Now())
}
