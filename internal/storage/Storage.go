package storage

import (
	"OrelTeatrBot/internal/domain"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{pool: pool}
}

func (s *Storage) SaveEvents(ctx context.Context, events []domain.Event) error {
	if len(events) == 0 {
		return nil
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	var foundSessionIds []int64
	for _, event := range events {
		_, err = tx.Exec(ctx,
			"INSERT INTO events(id, title, updated_at) VALUES($1, $2, now()) "+
				"ON CONFLICT (id) DO UPDATE SET title = $2, updated_at = now()",
			event.ID, event.Title)
		if err != nil {
			return fmt.Errorf("could not insert event: %w", err)
		}
		if len(event.Sessions) == 0 {
			return fmt.Errorf("could not insert event: no sessions")
		}
		for _, session := range event.Sessions {
			foundSessionIds = append(foundSessionIds, session.ID)
			_, err = tx.Exec(ctx,
				"INSERT INTO sessions(id, event_id, date_text, actual_now, last_seen_at) VALUES($1, $2, $3, $4, now()) "+
					"ON CONFLICT (id) DO UPDATE SET date_text = $3, actual_now = $4, last_seen_at = now(), "+
					"avail_epoch = sessions.avail_epoch + (NOT sessions.actual_now AND EXCLUDED.actual_now)::int",
				session.ID, event.ID, session.DateText, session.ActualNow)
			if err != nil {
				return fmt.Errorf("could not insert session: %w", err)
			}
		}
	}
	_, err = tx.Exec(ctx, "UPDATE sessions SET actual_now = false WHERE id <> ALL($1) AND actual_now", foundSessionIds)
	if err != nil {
		return fmt.Errorf("could not update sessions: %w", err)
	}

	return tx.Commit(ctx)
}
