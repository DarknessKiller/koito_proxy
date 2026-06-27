//go:build embed

package admin

import (
	_ "embed"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed dist/index.html
var adminPageHTML string

var adminPageFragment string

func init() {
	var b strings.Builder

	for _, tag := range []string{"script", "style"} {
		pos := 0
		open := "<" + tag
		close := "</" + tag + ">"
		for {
			s := strings.Index(adminPageHTML[pos:], open)
			if s == -1 {
				break
			}
			e := strings.Index(adminPageHTML[pos+s:], close)
			if e == -1 {
				break
			}
			b.WriteString(adminPageHTML[pos+s : pos+s+e+len(close)])
			b.WriteByte('\n')
			pos += s + e + len(close)
		}
	}

	// Extract the mount div from <body>
	const marker = `<div id="koito-admin-mount">`
	if idx := strings.Index(adminPageHTML, marker); idx != -1 {
		end := strings.Index(adminPageHTML[idx:], "</div>")
		if end != -1 {
			b.WriteString(adminPageHTML[idx : idx+end+6])
			b.WriteByte('\n')
		}
	}

	adminPageFragment = b.String()
}

func UIHandler(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(adminPageFragment))
}
