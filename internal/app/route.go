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
	"koito_proxy/internal/middleware/ratelimiter"
	"koito_proxy/internal/proxy"
	"koito_proxy/internal/proxy/koito"
	"koito_proxy/internal/proxy/listenbrainz"

	"github.com/gin-gonic/gin"
)

func (a *App) SetupRoute(ctx context.Context) {

	r := a.engine

	// Base Router
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

	rl := r.Group("", ratelimiter.RateLimitMiddleware(a.config))

	// Koito Proxy Routes
	{

		rl.GET("/apis/health", func(c *gin.Context) {
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

		rl.POST(
			"/apis/listenbrainz/1/submit-listens",
			lbAuth.Middleware(),
			lbHandler.InterceptSubmitListen,
		)

		rl.POST(
			"/apis/web/v1/:entity/:id/merge",
			koitoAuth.Middleware(),
			koitoHandler.InterceptMerge,
		)

		admin.RegisterRoutes(
			rl.Group("/apis/admin"),
			a.ruleService,
			koitoAuth.Middleware(),
		)
	}

	// Koito Upstream Routes
	{

		fb := r.Group("/apis/web/v1")
		fb.GET("/config", fallbackProxy)
		fb.POST("/login", fallbackProxy)
		fb.POST("/logout", fallbackProxy)
		fb.GET("/top/tracks", fallbackProxy)
		fb.GET("/top/albums", fallbackProxy)
		fb.GET("/top/artists", fallbackProxy)
		fb.GET("/listens", fallbackProxy)
		fb.POST("/listens", fallbackProxy)
		fb.DELETE("/listens", fallbackProxy)
		fb.GET("/listen-activity", fallbackProxy)
		fb.GET("/first-activity", fallbackProxy)
		fb.GET("/now-playing", fallbackProxy)
		fb.GET("/stats", fallbackProxy)
		fb.GET("/search", fallbackProxy)
		fb.GET("/summary", fallbackProxy)
		fb.GET("/user", fallbackProxy)
		fb.PATCH("/user", fallbackProxy)
		fb.GET("/user/apikeys", fallbackProxy)
		fb.POST("/user/apikeys", fallbackProxy)
		fb.PATCH("/user/apikeys/:id", fallbackProxy)
		fb.DELETE("/user/apikeys/:id", fallbackProxy)
		fb.GET("/export", fallbackProxy)
		fb.DELETE("/data", fallbackProxy)

		// Entity routes
		for _, entity := range []string{"artist", "album", "track"} {
			fb.GET("/"+entity+"/:id", fallbackProxy)
			fb.DELETE("/"+entity+"/:id", fallbackProxy)
			fb.PATCH("/"+entity+"/:id", fallbackProxy)
			fb.GET("/"+entity+"/:id/aliases", fallbackProxy)
			fb.POST("/"+entity+"/:id/aliases", fallbackProxy)
			fb.DELETE("/"+entity+"/:id/aliases", fallbackProxy)
			fb.PATCH("/"+entity+"/:id/aliases/primary", fallbackProxy)
			fb.GET("/"+entity+"/:id/interest", fallbackProxy)
		}

		// Album & track artist subroutes
		for _, entity := range []string{"album", "track"} {
			fb.GET("/"+entity+"/:id/artists", fallbackProxy)
			fb.PATCH("/"+entity+"/:id/artists/:artist_id", fallbackProxy)
		}

		// Track-only artist management
		fb.POST("/track/:id/artists", fallbackProxy)
		fb.DELETE("/track/:id/artists/:artist_id", fallbackProxy)

		// Non-web fallback routes
		rl.GET("/apis/listenbrainz/1/validate-token", fallbackProxy)
		r.GET("/image/:image_id/:filename", fallbackProxy)

		r.NoRoute(func(c *gin.Context) {
			if c.Request.Method != http.MethodGet {
				c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
				return
			}
			fallbackProxy(c)
		})
	}
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
