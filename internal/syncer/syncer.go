package syncer

import (
	"OrelTeatrBot/internal/domain"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"
)

type Fetcher interface {
	Fetch(ctx context.Context) ([]byte, error)
}

type Storage interface {
	SaveEvents(ctx context.Context, events []domain.Event) error
}

type Syncer struct {
	fetcher       Fetcher
	store         Storage
	fetchInterval time.Duration
	parse         func(io.Reader) ([]domain.Event, error)
}

func NewSyncer(fetcher Fetcher,
	store Storage,
	fetchInterval time.Duration,
	parseFunc func(io.Reader) ([]domain.Event, error)) *Syncer {
	return &Syncer{
		fetcher:       fetcher,
		store:         store,
		fetchInterval: fetchInterval,
		parse:         parseFunc}
}

func (s *Syncer) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.fetchInterval)
	defer ticker.Stop()
	for {
		err := s.doIteration(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return err
			}
			slog.Error("failed to do iteration", "err", err)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func (s *Syncer) doIteration(ctx context.Context) error {
	data, err := s.fetcher.Fetch(ctx)
	if err != nil {
		return fmt.Errorf("fetching: %w", err)
	}
	events, err := s.parse(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("parsing: %w", err)
	}
	if len(events) == 0 {
		return fmt.Errorf("no events found")
	}
	err = s.store.SaveEvents(ctx, events)
	if err != nil {
		return fmt.Errorf("saving events: %w", err)
	}
	slog.Info("iteration done",
		"total events", len(events))
	return nil
}
