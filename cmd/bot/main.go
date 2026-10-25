package main

import (
	config "OrelTeatrBot/internal/config"
	"OrelTeatrBot/internal/notifier"
	"OrelTeatrBot/internal/scraper"
	"OrelTeatrBot/internal/storage"
	syncer "OrelTeatrBot/internal/syncer"
	"OrelTeatrBot/internal/telegram"
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	err := godotenv.Load()
	if err != nil {
		slog.Info("Error loading .env file")
	}
	cnf, err := config.ReadConfigFromEnv()
	if err != nil {
		slog.Error("Can't read config", "err", err)
		return
	}

	pool, err := pgxpool.New(ctx, cnf.DbConnectionString)
	if err != nil {
		slog.Error("Can't connect to database", "err", err)
		return
	}
	defer pool.Close()

	stg := storage.NewStorage(pool)
	f := scraper.NewFetcher("https://quicktickets.ru/orel-teatr-svobodnoe-prostranstvo")
	s := syncer.NewSyncer(f, stg, 60*time.Second, scraper.Parse)

	tbot, err := telegram.NewBot(cnf.BotToken, stg)
	if err != nil {
		slog.Error("Can't connect to telegram", "err", err)
		return
	}

	ntf := notifier.NewNotifier("https://quicktickets.ru/orel-teatr-svobodnoe-prostranstvo", 30*time.Second,
		stg, tbot)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error { return s.Run(ctx) })
	g.Go(func() error { return ntf.Run(ctx) })
	g.Go(func() error {
		tbot.Start(ctx)
		return nil
	})
	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("stopped with error", "err", err)
	}
}
