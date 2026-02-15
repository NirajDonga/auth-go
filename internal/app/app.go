package app

import (
	"context"
	"fmt"
	"go-auth/internal/config"
	"go-auth/internal/db"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Config      config.Config
	MongoClient *mongo.Client
	DB          *mongo.Database
}

func New(ctx context.Context) (*App, error) {
	// load env
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// connect to db
	mongocli, err := db.Connect(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		Config:      cfg,
		MongoClient: mongocli.Client,
		DB:          mongocli.DB,
	}, nil
}

func (a *App) Close(ctx context.Context) error {
	if a.MongoClient == nil {
		return nil
	}

	closeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := a.MongoClient.Disconnect(closeCtx); err != nil {
		return fmt.Errorf("mongo Disconnect failed: %w", err)
	}
	return nil
}
