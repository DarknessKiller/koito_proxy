package auth

import (
	"context"
	"koito_proxy/internal/config"
	"koito_proxy/internal/response"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

type KoitoAuth struct {
	config     *config.Config
	cache      *Cache
	httpClient *http.Client
}

func NewKoitoAuth(cfg *config.Config, cache *Cache, httpClient *http.Client) *KoitoAuth {
	return &KoitoAuth{config: cfg, cache: cache, httpClient: httpClient}
}

func (a *KoitoAuth) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := c.Cookie("koito_session")
		if err != nil || session == "" {
			slog.Warn("missing koito_session cookie", "error", err, "path", c.Request.URL.Path)
			response.RespondUnauthorized(c, response.ErrMissingKoitoSession)
			c.Abort()
			return
		}

		if ok, found := a.cache.Get(session); found && ok {
			c.Next()
			return
		}

		if !a.validateUpstream(c.Request.Context(), session) {
			slog.Warn("koito session validation failed", "path", c.Request.URL.Path)
			response.RespondUnauthorized(c, response.ErrInvalidKoitoSession)
			c.Abort()
			return
		}

		a.cache.Set(session, 15*time.Minute)
		c.Next()
	}
}

func (a *KoitoAuth) validateUpstream(ctx context.Context, token string) bool {
	pathBuilder := newPathBuilder()
	targetURL, err := a.targetURL(pathBuilder.KoitoAuthorization())
	if err != nil {
		slog.Error("failed to build koito auth URL", "error", err)
		return false
	}

	req, err := http.NewRequestWithContext(ctx, "GET", targetURL.String(), nil)
	if err != nil {
		slog.Error("failed to create koito auth request", "error", err)
		return false
	}

	req.AddCookie(&http.Cookie{
		Name:  "koito_session",
		Value: token,
	})

	resp, err := a.httpClient.Do(req)
	if err != nil {
		slog.Error("failed to validate koito session upstream", "error", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true
	}

	slog.Warn("koito validation failed with upstream", "status", resp.StatusCode)
	return false
}

func (a *KoitoAuth) HasValidSession(r *http.Request) bool {
	if r == nil {
		return false
	}

	session, err := r.Cookie("koito_session")
	if err != nil || session.Value == "" {
		return false
	}

	token := session.Value
	if ok, found := a.cache.Get(token); found {
		return ok
	}

	if !a.validateUpstream(r.Context(), token) {
		return false
	}

	a.cache.Set(token, 15*time.Minute)
	return true
}

func (a *KoitoAuth) targetURL(apiPath APIPath) (*url.URL, error) {
	return apiPath.URL(a.config.UpstreamURL)
}
