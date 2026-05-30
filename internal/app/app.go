package app

import (
	"koito_proxy/internal/bootstrap"
	"log/slog"
)

type App struct {
	bs *bootstrap.Bootstrap
}

func New() (*App, error) {
	bs, err := bootstrap.New()
	if err != nil {
		return nil, err
	}

	return &App{bs: bs}, nil
}

func (a *App) Run() error {
	r := a.SetupRouter()
	if err := r.Run(":" + a.bs.Config.Port); err != nil {
		slog.Error("failed to run server", "error", err)
		return err
	}
	return nil
}
