package auth

import (
	"fmt"
	"koito_proxy/internal/config"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

type ListenBrainzAuth struct {
	config *config.Config
	cache  *Cache
}

func NewListenBrainzAuth(cfg *config.Config, cache *Cache) *ListenBrainzAuth {
	return &ListenBrainzAuth{config: cfg, cache: cache}
}

func (a *ListenBrainzAuth) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			c.Abort()
			return
		}

		if ok, found := a.cache.Get(token); found && ok {
			c.Next()
			return
		}

		if !a.validateUpstream(c, token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		a.cache.Set(token, 5*time.Minute)
		c.Next()
	}
}

// validateUpstream calls ListenBrainz API to verify the token
func (a *ListenBrainzAuth) validateUpstream(c *gin.Context, token string) bool {
	pathBuilder := newPathBuilder()
	targetURL, err := a.targetURL(pathBuilder.LBAuthorization())
	if err != nil {
		return false
	}

	req, err := http.NewRequest("GET", targetURL.String(), nil)
	if err != nil {
		return false
	}

	req.Header.Set("Authorization", token)

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

	fmt.Printf("ListenBrainz validation failed with status: %d\n", resp.StatusCode)
	return false
}

func (a *ListenBrainzAuth) targetURL(apiPath APIPath) (*url.URL, error) {
	return apiPath.URL(a.config.UpstreamURL)
}
