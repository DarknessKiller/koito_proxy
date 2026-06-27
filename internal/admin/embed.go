//go:build embed

package admin

import "embed"

//go:embed all:dist
var FS embed.FS
