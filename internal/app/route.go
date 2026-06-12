package app

import (
	"log/slog"
	"net/http"
	"time"

	"koito_proxy/internal/admin"
	"koito_proxy/internal/middleware"
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
		middleware.RequestIDMiddleware(),
		GinSlogLogger(),
		gin.Recovery(),
	)

	cache := auth.NewCache()

	lbAuth := auth.NewListenBrainzAuth(a.config, cache, a.httpClient)
	koitoAuth := auth.NewKoitoAuth(a.config, cache, a.httpClient)

	lbHandler := listenbrainz.NewHandler(a.ruleEngine, a.config)
	koitoHandler := koito.NewHandler(a.ruleEngine, a.repository, a.config)

	fallbackProxy := proxy.New(a.config).Handler()

	r.GET("/apis/health", func(c *gin.Context) {

		ruleEngineStatus := "inactive"
		if a.ruleEngine != nil {
			ruleEngineStatus = "active"
		}

		c.JSON(http.StatusOK, gin.H{
			"ok":                 true,
			"rule_engine_status": ruleEngineStatus,
			"timestamp":          time.Now().Format(time.RFC3339),
		})
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

	admin.RegisterRoutes(
		r.Group("/apis/admin"),
		a.repository,
		a.ruleEngine,
		koitoAuth.Middleware(),
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
			slog.String("request_id", getRequestID(c)),
			slog.Duration("latency", time.Since(start)),
		)
	}
}

func getRequestID(c *gin.Context) string {
	if id, exists := c.Get("request_id"); exists {
		return id.(string)
	}
	return "unknown"
}
