package app

import (
	"net/http"

	"koito_proxy/internal/middleware/auth"
	"koito_proxy/internal/proxy"
	"koito_proxy/internal/proxy/koito"
	"koito_proxy/internal/proxy/listenbrainz"

	"github.com/gin-gonic/gin"
)

func (a *App) SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	cache := auth.NewCache()

	lbAuth := auth.NewListenBrainzAuth(a.bs.Config, cache)

	lbHandler := listenbrainz.NewHandler(a.bs.Engine, a.bs.Config)
	koitoHandler := koito.NewHandler(a.bs.Engine, a.bs.Store, a.bs.Config)

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
		koitoHandler.InterceptMerge,
	)

	r.NoRoute(proxy.New(a.bs.Config).Handler())

	return r
}
