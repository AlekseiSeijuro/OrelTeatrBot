package scraper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type Fetcher struct {
	client *http.Client
	url    string
}

func NewFetcher(url string) *Fetcher {
	jar, _ := cookiejar.New(nil)
	return &Fetcher{
		client: &http.Client{
			Jar:     jar,
			Timeout: 15 * time.Second,
		},
		url: url}
}

func (f *Fetcher) Fetch(ctx context.Context) ([]byte, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.url, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 "+
		"(KHTML, like Gecko) Chrome/124.0 Safari/537.36")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en;q=0.8")
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}
