package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gin-gonic/gin"

	"koito_proxy/internal/config"
	"koito_proxy/internal/middleware/auth"
)

func TestAuth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth Suite")
}

var _ = Describe("KoitoAuth CSRF Cookie", func() {
	var (
		rec     *httptest.ResponseRecorder
		cache   *auth.Cache
		cfg     *config.Config
		handler gin.HandlerFunc
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		cache = auth.NewCache()
		cfg = &config.Config{UpstreamURL: "http://localhost:9999"}

		handler = auth.NewKoitoAuth(cfg, cache, &http.Client{Timeout: 5 * time.Second}).Middleware()
	})

	doRequest := func(cookie string) *httptest.ResponseRecorder {
		rec = httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)
		req := httptest.NewRequest("POST", "/apis/admin/rules", nil)
		if cookie != "" {
			req.Header.Set("Cookie", cookie)
		}
		ctx.Request = req
		handler(ctx)
		return rec
	}

	It("sets koito_csrf cookie on cache hit", func() {
		cache.Set("valid-session", 15*time.Minute)

		rec = doRequest("koito_session=valid-session")

		cookies := rec.Result().Cookies()
		var csrfCookie *http.Cookie
		for _, c := range cookies {
			if c.Name == "koito_csrf" {
				csrfCookie = c
				break
			}
		}

		Expect(csrfCookie).NotTo(BeNil())
		Expect(csrfCookie.Value).To(Equal("1"))
		Expect(csrfCookie.SameSite).To(Equal(http.SameSiteStrictMode))
		Expect(csrfCookie.HttpOnly).To(BeTrue())
		Expect(csrfCookie.Path).To(Equal("/"))
	})

	It("does not set koito_csrf cookie on auth failure", func() {
		rec = doRequest("koito_session=invalid-session")

		cookies := rec.Result().Cookies()
		for _, c := range cookies {
			Expect(c.Name).NotTo(Equal("koito_csrf"))
		}
	})

	It("does not set koito_csrf cookie without session", func() {
		rec = doRequest("")

		cookies := rec.Result().Cookies()
		for _, c := range cookies {
			Expect(c.Name).NotTo(Equal("koito_csrf"))
		}
	})
})
