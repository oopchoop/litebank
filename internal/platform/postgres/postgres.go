package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)

	if err != nil {
		return nil, fmt.Errorf("создание пула: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()

		return nil, fmt.Errorf("пинг к бд: %w", err)
	}

	return pool, nil
}
