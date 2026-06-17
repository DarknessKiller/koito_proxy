package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodPost &&
			c.Request.Method != http.MethodPut &&
			c.Request.Method != http.MethodDelete &&
			c.Request.Method != http.MethodPatch {
			c.Next()
			return
		}

		origin := c.GetHeader("Origin")
		if origin != "" {
			if !strings.HasPrefix(origin, "//") && !strings.Contains(origin, "://") {
				origin = "http://" + origin
			}
			parsed, err := http.NewRequest(http.MethodGet, origin, nil)
			if err != nil || parsed.Host != c.Request.Host {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "CSRF protection: request origin mismatch",
				})
				return
			}
			c.Next()
			return
		}

		referer := c.GetHeader("Referer")
		if referer != "" {
			parsed, err := http.NewRequest(http.MethodGet, referer, nil)
			if err != nil || parsed.Host != c.Request.Host {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "CSRF protection: request origin mismatch",
				})
				return
			}
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "CSRF protection: request origin mismatch",
		})
	}
}
