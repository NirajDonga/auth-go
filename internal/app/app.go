package app

import (
	"context"
	"go-auth/internal/config"
	"go-auth/internal/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Config config.Config
	DB     *pgxpool.Pool
}

func New(ctx context.Context) (*App, error) {
	// load env
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// connect to db
	pg, err := db.Connect(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		Config: cfg,
		DB:     pg.Pool,
	}, nil
}

func (a *App) Close(ctx context.Context) error {
	if a.DB == nil {
		return nil
	}

	a.DB.Close()
	return nil
}
