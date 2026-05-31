package auth

import (
	"fmt"
	"koito_proxy/internal/config"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

type KoitoAuth struct {
	config *config.Config
	cache  *Cache
}

func NewKoitoAuth(cfg *config.Config, cache *Cache) *KoitoAuth {
	return &KoitoAuth{config: cfg, cache: cache}
}

func (a *KoitoAuth) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := c.Cookie("koito_session")
		if err != nil || session == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing koito_session"})
			c.Abort()
			return
		}

		if ok, found := a.cache.Get(session); found && ok {
			c.Next()
			return
		}

		if !a.validateUpstream(c, session) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
			c.Abort()
			return
		}

		a.cache.Set(session, 5*time.Minute)
		c.Next()
	}
}

// validateUpstream calls ListenBrainz API to verify the token
func (a *KoitoAuth) validateUpstream(c *gin.Context, token string) bool {
	pathBuilder := newPathBuilder()
	targetURL, err := a.targetURL(pathBuilder.KoitoAuthorization())
	if err != nil {
		return false
	}

	req, err := http.NewRequest("GET", targetURL.String(), nil)
	if err != nil {
		return false
	}

	req.AddCookie(&http.Cookie{
		Name:  "koito_session",
		Value: token,
	})

	client := httpClient()
	resp, err := client.Do(req.WithContext(c.Request.Context()))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// 200 OK means token is valid
	if resp.StatusCode == http.StatusOK {
		return true
	}

	fmt.Printf("Koito validation failed with status: %d\n", resp.StatusCode)
	return false
}

func (a *KoitoAuth) targetURL(apiPath APIPath) (*url.URL, error) {
	return apiPath.URL(a.config.UpstreamURL)
}
