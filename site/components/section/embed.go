package section

import (
	_ "embed"

	"github.com/ocuroot/ui/css"
)

//go:embed section.css
var CSS []byte

func init() {
	css.Default().Add(CSS)
}
