package app

import (
	"context"
	"fmt"
	"koito_proxy/internal/bootstrap"
	"koito_proxy/internal/config"
	"koito_proxy/internal/model"
	"koito_proxy/internal/repository"
	"koito_proxy/internal/rules"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
	engine     *gin.Engine
	config     *config.Config
	repository repository.Repository[model.Rule]
	ruleEngine *rules.RuleEngine
	httpClient *http.Client
}

func New() (*App, error) {
	bs, err := bootstrap.New()
	if err != nil {
		return nil, err
	}

	gin.SetMode(gin.ReleaseMode)

	app := &App{
		config:     bs.Config,
		repository: bs.Repository,
		engine:     gin.New(),
		ruleEngine: bs.RuleEngine,
		httpClient: bs.HttpClient,
	}

	app.SetupRoute()

	return app, nil
}

func (a *App) Run() error {
	server := &http.Server{
		Addr:         ":" + a.config.Port,
		Handler:      a.engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	rules, err := a.repository.GetAll(context.Background())
	if err != nil {
		slog.Error("failed to load rules from database", "error", err)
		return err
	}
	a.ruleEngine.UpdateRules(rules)

	go func() {
		slog.Info("Starting Koito Proxy server", "port", a.config.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Could not start HTTP server", "error", err)
		}
	}()

	<-quit
	slog.Warn("Received termination signal, initiating graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		return fmt.Errorf("server shutdown error: %w", err)
	}

	slog.Info("Server shut down gracefully")
	return nil
}
