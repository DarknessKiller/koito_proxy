package ratelimiter

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateLimiterMiddleware(requestsPerSec rate.Limit, burst int) gin.HandlerFunc {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	var clients = make(map[string]*client)

	go func() {
		for {
			time.Sleep(1 * time.Minute)
			for ip, item := range clients {
				if time.Since(item.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if _, exists := clients[ip]; !exists {
			clients[ip] = &client{
				limiter: rate.NewLimiter(requestsPerSec, burst),
			}
		}
		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			slog.Warn("rate_limit_exceeded", "ip", ip, "path", c.Request.URL.Path)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please slow down.",
			})
			return
		}

		c.Next()
	}
}
