package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// EventStore persists event records to PostgreSQL.
type EventStore struct {
	pool *pgxpool.Pool
}

func NewEventStore(pool *pgxpool.Pool) *EventStore {
	return &EventStore{pool: pool}
}

func (s *EventStore) BatchInsert(ctx context.Context, sourceID string, events []map[string]interface{}) error {
	// Direct batch insert — called synchronously on the request path.
	// Under high concurrency this creates write bottlenecks.
	return nil
}
