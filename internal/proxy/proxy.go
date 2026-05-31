package proxy

import (
	"koito_proxy/internal/config"
	"log/slog"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

type Proxy struct {
	target *url.URL
	rp     *httputil.ReverseProxy
}

func New(cfg *config.Config) *Proxy {
	parsed, err := url.Parse(cfg.UpstreamURL)
	if err != nil {
		slog.Error("invalid upstream URL", "error", err)
		panic(err)
	}
	rp := &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetURL(parsed)
			pr.Out.Host = parsed.Host
		},
	}
	return &Proxy{target: parsed, rp: rp}
}

func (p *Proxy) Handler() func(c *gin.Context) {
	return func(c *gin.Context) {
		p.rp.ServeHTTP(c.Writer, c.Request)
	}
}
