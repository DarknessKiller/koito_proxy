//go:build !embed

package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UIHandler(c *gin.Context) {
	c.Data(http.StatusNotFound, "text/plain; charset=utf-8", []byte("admin UI not available"))
}
