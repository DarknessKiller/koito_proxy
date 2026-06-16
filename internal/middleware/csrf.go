package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost ||
			c.Request.Method == http.MethodPut ||
			c.Request.Method == http.MethodDelete ||
			c.Request.Method == http.MethodPatch {
			if c.GetHeader("X-Requested-With") != "KOITO_PROXY_API" {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "CSRF protection: missing X-Requested-With header",
				})
				return
			}
		}
		c.Next()
	}
}
