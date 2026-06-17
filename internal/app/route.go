package app

import (
	"context"
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

func (a *App) SetupRoute(ctx context.Context) {

	r := a.engine

	r.Use(
		limit.BodyLimitMiddleware(5),
		middleware.RequestIDMiddleware(),
		GinSlogLogger(),
		gin.Recovery(),
	)

	cache := auth.NewCache()
	cache.StartCleanup(ctx, 5*time.Minute)

	lbAuth := auth.NewListenBrainzAuth(a.config, cache, a.httpClient)
	koitoAuth := auth.NewKoitoAuth(a.config, cache, a.httpClient)

	lbHandler := listenbrainz.NewHandler(a.ruleEngine, a.config, a.httpClient)
	koitoHandler := koito.NewHandler(a.koitoService, a.config, a.httpClient)

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
		a.ruleService,
		koitoAuth.Middleware(),
		middleware.CSRFMiddleware(),
	)

	// Upstream Koito routes — only valid paths are registered.
	r.GET("/apis/web/v1/config", fallbackProxy)
	r.POST("/apis/web/v1/login", fallbackProxy)
	r.POST("/apis/web/v1/logout", fallbackProxy)
	r.GET("/apis/web/v1/top/tracks", fallbackProxy)
	r.GET("/apis/web/v1/top/albums", fallbackProxy)
	r.GET("/apis/web/v1/top/artists", fallbackProxy)
	r.GET("/apis/web/v1/listens", fallbackProxy)
	r.POST("/apis/web/v1/listens", fallbackProxy)
	r.DELETE("/apis/web/v1/listens", fallbackProxy)
	r.GET("/apis/web/v1/listen-activity", fallbackProxy)
	r.GET("/apis/web/v1/first-activity", fallbackProxy)
	r.GET("/apis/web/v1/now-playing", fallbackProxy)
	r.GET("/apis/web/v1/stats", fallbackProxy)
	r.GET("/apis/web/v1/search", fallbackProxy)
	r.GET("/apis/web/v1/summary", fallbackProxy)
	r.GET("/apis/web/v1/user", fallbackProxy)
	r.PATCH("/apis/web/v1/user", fallbackProxy)
	r.GET("/apis/web/v1/user/apikeys", fallbackProxy)
	r.POST("/apis/web/v1/user/apikeys", fallbackProxy)
	r.PATCH("/apis/web/v1/user/apikeys/:id", fallbackProxy)
	r.DELETE("/apis/web/v1/user/apikeys/:id", fallbackProxy)
	r.GET("/apis/web/v1/export", fallbackProxy)
	r.DELETE("/apis/web/v1/data", fallbackProxy)
	r.GET("/apis/listenbrainz/1/validate-token", fallbackProxy)
	r.GET("/image/:image_id/:filename", fallbackProxy)

	// Entity routes
	for _, entity := range []string{"artist", "album", "track"} {
		r.GET("/apis/web/v1/"+entity+"/:id", fallbackProxy)
		r.DELETE("/apis/web/v1/"+entity+"/:id", fallbackProxy)
		r.PATCH("/apis/web/v1/"+entity+"/:id", fallbackProxy)
		r.GET("/apis/web/v1/"+entity+"/:id/aliases", fallbackProxy)
		r.POST("/apis/web/v1/"+entity+"/:id/aliases", fallbackProxy)
		r.DELETE("/apis/web/v1/"+entity+"/:id/aliases", fallbackProxy)
		r.PATCH("/apis/web/v1/"+entity+"/:id/aliases/primary", fallbackProxy)
		r.GET("/apis/web/v1/"+entity+"/:id/interest", fallbackProxy)
	}

	// Album & track artist subroutes
	for _, entity := range []string{"album", "track"} {
		r.GET("/apis/web/v1/"+entity+"/:id/artists", fallbackProxy)
		r.PATCH("/apis/web/v1/"+entity+"/:id/artists/:artist_id", fallbackProxy)
	}

	// Track-only artist management
	r.POST("/apis/web/v1/track/:id/artists", fallbackProxy)
	r.DELETE("/apis/web/v1/track/:id/artists/:artist_id", fallbackProxy)

	// Frontend
	r.NoRoute(func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}
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
