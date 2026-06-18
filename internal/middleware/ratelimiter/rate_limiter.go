package ratelimiter

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"koito_proxy/internal/config"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter interface {
	Allow(ip string) bool
}

type TokenBucket struct {
	capacity   float64
	rate       float64
	tokens     float64
	lastRefill time.Time
	mu         sync.Mutex
}

func NewTokenBucket(capacity float64, rate float64) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		rate:       rate,
		tokens:     capacity,
		lastRefill: time.Now(),
	}
}

func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)
	tokensToAdd := elapsed.Seconds() * tb.rate
	tb.tokens = min(tb.capacity, tb.tokens+tokensToAdd)
	tb.lastRefill = now
}

func (tb *TokenBucket) Allow(ip string) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.refill()
	if tb.tokens >= 1.0 {
		tb.tokens -= 1.0
		return true
	}
	return false
}

func GlobalLimiterFactory(config *config.Config) (RateLimiter, error) {
	if config.RateLimitEnabled != "true" {
		return nil, fmt.Errorf("rate limiting is disabled")
	}

	rate := parseRate(config.RequestsPerSecond)
	burst := parseBurst(config.Burst)

	if rate <= 0 || burst < 1 {
		return nil, fmt.Errorf("invalid rate limit configuration: requests_per_second=%s, burst=%s", config.RequestsPerSecond, config.Burst)
	}

	return NewTokenBucket(burst, rate), nil
}

func parseRate(s string) float64 {
	if s == "" {
		return 0
	}
	rate, err := strconv.ParseFloat(s, 64)
	if err != nil || rate <= 0 {
		return 0
	}
	return rate
}

func parseBurst(s string) float64 {
	if s == "" {
		return 0
	}
	burst, err := strconv.Atoi(s)
	if err != nil || burst < 1 {
		return 0
	}
	return float64(burst)
}

func RateLimitMiddleware(config *config.Config) gin.HandlerFunc {
	limiter, err := GlobalLimiterFactory(config)
	if err != nil {
		slog.Warn("rate limiting disabled", "error", err)
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if ip == "" {
			c.Next()
			return
		}

		if !limiter.Allow(ip) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}
