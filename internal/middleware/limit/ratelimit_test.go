package limit

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRateLimiterConcurrentAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := RateLimiterMiddleware(50, 100)

	router := gin.New()
	router.Use(middleware)
	router.Any("/*path", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	var wg sync.WaitGroup
	iterations := 100

	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				rec := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/test", nil)
				req.RemoteAddr = "192.168.1.1:12345"
				router.ServeHTTP(rec, req)
			}
		}(i)
	}

	wg.Wait()
}
