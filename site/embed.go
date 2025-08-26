package site

import (
	_ "embed"

	"github.com/ocuroot/ui/css"
	"github.com/ocuroot/ui/js"
)

//go:embed style.css
var styleCSS []byte

//go:embed mermaid.css
var mermaidCSS []byte

//go:embed syntax-highlighting.css
var syntaxHighlightingCSS []byte

//go:embed gitcard.css
var gitcardCSS []byte

//go:embed gitcard.js
var gitcardJS []byte

//go:embed components/code/code.css
var codeCSS []byte

//go:embed components/code/code.js
var codeJS []byte

func init() {
	css.Default().Add(styleCSS)
	css.Default().Add(mermaidCSS)
	css.Default().Add(syntaxHighlightingCSS)
	css.Default().Add(gitcardCSS)
	css.Default().Add(codeCSS)

	js.Default().Add(gitcardJS)
	js.Default().Add(codeJS)
}
