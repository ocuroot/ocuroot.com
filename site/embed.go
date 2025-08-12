package site

import (
	_ "embed"

	"github.com/ocuroot/ui/css"
)

//go:embed style.css
var styleCSS []byte

func init() {
	css.Default().Add(styleCSS)
}
