package app

import (
	"koito_proxy/internal/bootstrap"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type App struct {
	bs     *bootstrap.Bootstrap
	engine *gin.Engine
}

func New() (*App, error) {
	bs, err := bootstrap.New()
	if err != nil {
		return nil, err
	}

	gin.SetMode(gin.ReleaseMode)

	app := &App{
		bs:     bs,
		engine: gin.New(),
	}

	app.SetupRoute()

	return app, nil
}

func (a *App) Run() error {
	if err := a.engine.Run(":" + a.bs.Config.Port); err != nil {
		slog.Error("failed to run server", "error", err)
		return err
	}

	return nil
}
