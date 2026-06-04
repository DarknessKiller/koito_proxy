package auth

import (
	"koito_proxy/internal/config"
	"koito_proxy/internal/response"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

type ListenBrainzAuth struct {
	config     *config.Config
	cache      *Cache
	httpClient *http.Client
}

func NewListenBrainzAuth(cfg *config.Config, cache *Cache, httpClient *http.Client) *ListenBrainzAuth {
	return &ListenBrainzAuth{config: cfg, cache: cache, httpClient: httpClient}
}

func (a *ListenBrainzAuth) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			slog.Warn("missing authorization header", "path", c.Request.URL.Path)
			response.RespondUnauthorized(c, response.ErrMissingAuthHeader)
			c.Abort()
			return
		}

		if ok, found := a.cache.Get(token); found && ok {
			c.Next()
			return
		}

		if !a.validateUpstream(c, token) {
			slog.Warn("listenbrainz token validation failed", "path", c.Request.URL.Path)
			response.RespondUnauthorized(c, response.ErrInvalidToken)
			c.Abort()
			return
		}

		a.cache.Set(token, 5*time.Minute)
		c.Next()
	}
}

func (a *ListenBrainzAuth) validateUpstream(c *gin.Context, token string) bool {
	pathBuilder := newPathBuilder()
	targetURL, err := a.targetURL(pathBuilder.LBAuthorization())
	if err != nil {
		slog.Error("failed to build listenbrainz auth URL", "error", err)
		return false
	}

	req, err := http.NewRequest("GET", targetURL.String(), nil)
	if err != nil {
		slog.Error("failed to create listenbrainz auth request", "error", err)
		return false
	}

	req.Header.Set("Authorization", token)

	resp, err := a.httpClient.Do(req.WithContext(c.Request.Context()))
	if err != nil {
		slog.Error("failed to validate listenbrainz token upstream", "error", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true
	}

	slog.Warn("listenbrainz validation failed with upstream", "status", resp.StatusCode)
	return false
}

func (a *ListenBrainzAuth) targetURL(apiPath APIPath) (*url.URL, error) {
	return apiPath.URL(a.config.UpstreamURL)
}
