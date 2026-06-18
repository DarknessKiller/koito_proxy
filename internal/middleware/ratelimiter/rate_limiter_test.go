package ratelimiter_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gin-gonic/gin"

	"koito_proxy/internal/config"
	"koito_proxy/internal/middleware/ratelimiter"
)

const (
	enabled  = "true"
	disabled = "false"

	defaultRPS  = "10"
	defaultBurst = "5"

	singleTokenRPS  = "1"
	singleTokenBurst = "1"

	invalidRPS  = "0"
	invalidBurst = "0"

	testPath = "/test"
	testIP   = "127.0.0.1"
)

func TestRateLimiter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RateLimiter Suite")
}

func newConfig(rateLimitEnabled, rps, burst string) *config.Config {
	return &config.Config{
		RateLimitEnabled:  rateLimitEnabled,
		RequestsPerSecond: rps,
		Burst:             burst,
	}
}

var _ = Describe("RateLimitMiddleware", func() {
	var (
		router *gin.Engine
		rec    *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		rec = httptest.NewRecorder()
	})

	setupRouter := func(cfg *config.Config) {
		router = gin.New()
		router.Use(ratelimiter.RateLimitMiddleware(cfg))
		router.Any("/*path", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})
	}

	doRequest := func(method, path string) {
		req := httptest.NewRequest(method, path, nil)
		router.ServeHTTP(rec, req)
	}

	Describe("rate limiting enabled", func() {
		Context("with sufficient burst capacity", func() {
			BeforeEach(func() {
				setupRouter(newConfig(enabled, defaultRPS, defaultBurst))
			})

			It("allows requests within the limit", func() {
				doRequest(http.MethodGet, testPath)
				Expect(rec.Code).To(Equal(http.StatusOK))
			})
		})

		Context("with exhausted burst", func() {
			BeforeEach(func() {
				setupRouter(newConfig(enabled, singleTokenRPS, singleTokenBurst))
			})

			It("rejects requests exceeding the limit with 429", func() {
				doRequest(http.MethodGet, testPath)
				Expect(rec.Code).To(Equal(http.StatusOK))

				rec = httptest.NewRecorder()
				doRequest(http.MethodGet, testPath)
				Expect(rec.Code).To(Equal(http.StatusTooManyRequests))
			})
		})
	})

	Describe("rate limiting disabled", func() {
		BeforeEach(func() {
			setupRouter(newConfig(disabled, defaultRPS, defaultBurst))
		})

		It("allows all requests", func() {
			doRequest(http.MethodGet, testPath)
			Expect(rec.Code).To(Equal(http.StatusOK))

			rec = httptest.NewRecorder()
			doRequest(http.MethodGet, testPath)
			Expect(rec.Code).To(Equal(http.StatusOK))
		})
	})

	Describe("token bucket", func() {
		It("allows requests up to capacity", func() {
			tb := ratelimiter.NewTokenBucket(3, 10)
			Expect(tb.Allow(testIP)).To(BeTrue())
			Expect(tb.Allow(testIP)).To(BeTrue())
			Expect(tb.Allow(testIP)).To(BeTrue())
			Expect(tb.Allow(testIP)).To(BeFalse())
		})

		It("rejects when capacity is exhausted", func() {
			tb := ratelimiter.NewTokenBucket(1, 0)
			Expect(tb.Allow(testIP)).To(BeTrue())
			Expect(tb.Allow(testIP)).To(BeFalse())
		})
	})

	Describe("GlobalLimiterFactory", func() {
		It("returns error when disabled", func() {
			_, err := ratelimiter.GlobalLimiterFactory(newConfig(disabled, defaultRPS, defaultBurst))
			Expect(err).To(HaveOccurred())
		})

		It("returns error for invalid config", func() {
			_, err := ratelimiter.GlobalLimiterFactory(newConfig(enabled, invalidRPS, invalidBurst))
			Expect(err).To(HaveOccurred())
		})

		It("creates limiter with valid config", func() {
			limiter, err := ratelimiter.GlobalLimiterFactory(newConfig(enabled, defaultRPS, defaultBurst))
			Expect(err).NotTo(HaveOccurred())
			Expect(limiter).NotTo(BeNil())
		})
	})
})
