package limit

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BodyLimitMiddleware(limitInMB int64) gin.HandlerFunc {

	limit := limitInMB * 1024 * 1024

	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(
			c.Writer,
			c.Request.Body,
			limit,
		)

		c.Next()
	}
}
