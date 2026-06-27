package proxy

import (
	"bytes"
	"io"
	"koito_proxy/internal/config"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

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
		ModifyResponse: injectOverlay,
	}
	return &Proxy{target: parsed, rp: rp}
}

func (p *Proxy) Handler() func(c *gin.Context) {
	return func(c *gin.Context) {
		p.rp.ServeHTTP(c.Writer, c.Request)
	}
}

func injectOverlay(r *http.Response) error {
	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "text/html") {
		return nil
	}

	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return nil
	}

	overlay := overlayScript

	modified := bytes.Replace(body, []byte("</body>"), []byte(overlay+"\n</body>"), 1)
	if len(modified) == len(body) {
		modified = append(body, []byte(overlay)...)
	}

	r.Body = io.NopCloser(bytes.NewReader(modified))
	r.ContentLength = int64(len(modified))
	r.Header.Set("Content-Length", strconv.Itoa(len(modified)))
	return nil
}

const overlayScript = `<style>
#kp-admin-btn {
  display: inline-block;
  padding: 8px;
  border-radius: 6px;
  color: var(--color-fg-secondary, #cfc3b7);
  cursor: pointer;
  transition: color .1s ease, background-color .1s ease;
  line-height: 0;
}

#kp-admin-btn:hover {
  background: var(--color-bg-tertiary, #3c2e2a);
  color: var(--color-fg, #f5ece3);
}

#kp-admin-btn svg {
  display: block;
}
</style>
<script>
(function () {
  fetch('/apis/admin/check').then(function (r) {
    if (!r.ok) return;

    var i = setInterval(function () {
      var s = document.querySelector('[class*="h-screen"][class*="justify-between"]');
      if (!s) return;
      clearInterval(i);

      var g = s.children[1];
      if (!g) return;

      var a = document.createElement('a');
      a.id = 'kp-admin-btn';
      a.innerHTML = [
        '<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24"',
        '  fill="none" stroke="currentColor" stroke-width="2"',
        '  stroke-linecap="round" stroke-linejoin="round">',
        '  <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10"/>',
        '</svg>',
      ].join('');
      a.title = 'Admin';

      var items = g.children;
      if (items.length > 0) {
        g.insertBefore(a, items[items.length - 1]);
      } else {
        g.appendChild(a);
      }

      a.onclick = function (e) {
        e.preventDefault();

        var root = document.getElementById('koito-admin-root');
        if (root) {
          root.style.display = '';
          return;
        }

        root = document.createElement('div');
        root.id = 'koito-admin-root';
        document.body.appendChild(root);

        fetch('/apis/admin/ui')
          .then(function (r) { return r.text() })
          .then(function (html) {
            root.innerHTML = html;
            root.querySelectorAll('script').forEach(function (s) {
              var ns = document.createElement('script');
              ns.textContent = s.textContent;
          s.replaceWith(ns);
            });
          });
      };
    }, 100);
  });
})();
</script>`
