package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gin-gonic/gin"

	"koito_proxy/internal/middleware"
)

func TestMiddleware(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Middleware Suite")
}

var _ = Describe("CSRFMiddleware", func() {
	var (
		router *gin.Engine
		rec    *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		router = gin.New()
		router.Use(middleware.CSRFMiddleware())
		router.Any("/*path", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})
		rec = httptest.NewRecorder()
	})

	doRequest := func(method, path string, headers map[string]string) {
		req := httptest.NewRequest(method, path, nil)
		for k, v := range headers {
			if k == "Host" {
				req.Host = v
			} else {
				req.Header.Set(k, v)
			}
		}
		router.ServeHTTP(rec, req)
	}

	Describe("safe methods", func() {
		It("allows GET without any CSRF headers", func() {
			doRequest("GET", "/apis/admin/rules", nil)
			Expect(rec.Code).To(Equal(http.StatusOK))
		})

		It("allows HEAD without any CSRF headers", func() {
			doRequest("HEAD", "/apis/admin/rules", nil)
			Expect(rec.Code).To(Equal(http.StatusOK))
		})

		It("allows OPTIONS without any CSRF headers", func() {
			doRequest("OPTIONS", "/apis/admin/rules", nil)
			Expect(rec.Code).To(Equal(http.StatusOK))
		})
	})

	Describe("mutating methods without Origin or Referer", func() {
		It("rejects POST without Origin or Referer", func() {
			doRequest("POST", "/apis/admin/rules", nil)
			Expect(rec.Code).To(Equal(http.StatusForbidden))
		})

		It("rejects PUT without Origin or Referer", func() {
			doRequest("PUT", "/apis/admin/rules/123", nil)
			Expect(rec.Code).To(Equal(http.StatusForbidden))
		})

		It("rejects DELETE without Origin or Referer", func() {
			doRequest("DELETE", "/apis/admin/rules/123", nil)
			Expect(rec.Code).To(Equal(http.StatusForbidden))
		})

		It("rejects PATCH without Origin or Referer", func() {
			doRequest("PATCH", "/apis/admin/rules/123", nil)
			Expect(rec.Code).To(Equal(http.StatusForbidden))
		})
	})

	Describe("mutating methods with Origin header", func() {
		It("allows POST with matching Origin", func() {
			doRequest("POST", "/apis/admin/rules", map[string]string{
				"Host":   "localhost:4112",
				"Origin": "http://localhost:4112",
			})
			Expect(rec.Code).To(Equal(http.StatusOK))
		})

		It("rejects POST with mismatched Origin", func() {
			doRequest("POST", "/apis/admin/rules", map[string]string{
				"Host":   "localhost:4112",
				"Origin": "http://evil.example.com",
			})
			Expect(rec.Code).To(Equal(http.StatusForbidden))
		})
	})

	Describe("mutating methods with Referer header", func() {
		It("allows POST with matching Referer", func() {
			doRequest("POST", "/apis/admin/rules", map[string]string{
				"Host":    "localhost:4112",
				"Referer": "http://localhost:4112/apis/admin/rules",
			})
			Expect(rec.Code).To(Equal(http.StatusOK))
		})

		It("rejects POST with mismatched Referer", func() {
			doRequest("POST", "/apis/admin/rules", map[string]string{
				"Host":    "localhost:4112",
				"Referer": "http://evil.example.com/page",
			})
			Expect(rec.Code).To(Equal(http.StatusForbidden))
		})
	})
})
