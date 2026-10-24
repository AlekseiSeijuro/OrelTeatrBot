package main

import (
	"OrelTeatrBot/internal/scraper"
	syncer "OrelTeatrBot/internal/syncer"
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	f := scraper.NewFetcher("https://quicktickets.ru/orel-teatr-svobodnoe-prostranstvo")
	s := syncer.NewSyncer(f, 60*time.Second, scraper.Parse)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	err := s.Run(ctx)
	if !errors.Is(err, context.Canceled) {
		slog.Error("Failure!", "error:", err)
	}
}
