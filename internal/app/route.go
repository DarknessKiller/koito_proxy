package app

import (
	"log/slog"
	"net/http"
	"time"

	"koito_proxy/internal/middleware/auth"
	"koito_proxy/internal/middleware/limit"
	"koito_proxy/internal/proxy"
	"koito_proxy/internal/proxy/koito"
	"koito_proxy/internal/proxy/listenbrainz"

	"github.com/gin-gonic/gin"
)

func (a *App) SetupRoute() {

	r := a.engine

	r.Use(
		limit.BodyLimitMiddleware(5),
		limit.RateLimiterMiddleware(50, 100),
		GinSlogLogger(),
		gin.Recovery(),
	)

	cache := auth.NewCache()

	lbAuth := auth.NewListenBrainzAuth(a.bs.Config, cache)
	koitoAuth := auth.NewKoitoAuth(a.bs.Config, cache)

	lbHandler := listenbrainz.NewHandler(a.bs.RuleEngine, a.bs.Config)
	koitoHandler := koito.NewHandler(a.bs.RuleEngine, a.bs.Store, a.bs.Config)

	fallbackProxy := proxy.New(a.bs.Config).Handler()

	r.GET("/apis/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	r.POST(
		"/apis/listenbrainz/1/submit-listens",
		lbAuth.Middleware(),
		lbHandler.InterceptSubmitListen,
	)

	r.POST(
		"/apis/web/v1/:entity/:id/merge",
		koitoAuth.Middleware(),
		koitoHandler.InterceptMerge,
	)

	r.NoRoute(func(c *gin.Context) {
		fallbackProxy(c)
	})
}

func GinSlogLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery

		c.Next()

		status := c.Writer.Status()

		level := slog.LevelInfo
		if status >= 500 {
			level = slog.LevelError
		} else if status >= 400 {
			level = slog.LevelWarn
		}

		slog.LogAttrs(c.Request.Context(), level, "http_request",
			slog.Int("status", status),
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("query", rawQuery),
			slog.String("ip", c.ClientIP()),
			slog.Duration("latency", time.Since(start)),
		)
	}
}
