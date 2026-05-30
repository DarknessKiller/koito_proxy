package bootstrap

import (
	"context"
	"database/sql"
	"log/slog"

	"koito_proxy/internal/config"
	"koito_proxy/internal/db"
	"koito_proxy/internal/rules"
)

type Bootstrap struct {
	Config *config.Config
	DB     *sql.DB
	Store  *rules.Store
	Engine *rules.Engine
}

func New() (*Bootstrap, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	sqlDB := db.OpenAndMigrate(cfg.DBPath)

	store := rules.NewStore(sqlDB)
	rulesList, err := store.Load(context.Background())
	if err != nil {
		slog.Error("failed to load rules", "error", err)
		return nil, err
	}
	engine := rules.NewEngine(rulesList)

	return &Bootstrap{
		Config: cfg,
		DB:     sqlDB,
		Store:  store,
		Engine: engine,
	}, nil
}
