package site

import (
	_ "embed"

	"github.com/ocuroot/ui/css"
)

//go:embed style.css
var styleCSS []byte

//go:embed mermaid.css
var mermaidCSS []byte

//go:embed syntax-highlighting.css
var syntaxHighlightingCSS []byte

func init() {
	css.Default().Add(styleCSS)
	css.Default().Add(mermaidCSS)
	css.Default().Add(syntaxHighlightingCSS)
}
