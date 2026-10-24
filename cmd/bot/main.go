package main

import (
	config "OrelTeatrBot/internal/config"
	"OrelTeatrBot/internal/scraper"
	"OrelTeatrBot/internal/storage"
	syncer "OrelTeatrBot/internal/syncer"
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	err := godotenv.Load()
	if err != nil {
		slog.Info("Error loading .env file")
		return
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

	err = s.Run(ctx)
	if !errors.Is(err, context.Canceled) {
		slog.Error("Failure!", "error:", err)
	}
}
