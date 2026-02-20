package db

import (
	"context"
	"fmt"
	"time"

	"go-auth/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

func Connect(ctx context.Context, cfg config.Config) (*Postgres, error) {
	connectCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(connectCtx, cfg.PostgresURL)
	if err != nil {
		return nil, fmt.Errorf("postgres connect failed: %w", err)
	}

	// verify connection
	if err := pool.Ping(connectCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("postgres ping failed: %w", err)
	}

	return &Postgres{Pool: pool}, nil
}

func (p *Postgres) Disconnect(ctx context.Context) {
	p.Pool.Close()
}
