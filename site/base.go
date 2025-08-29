package site

import (
	"net/url"
	"path"
)

var BaseURL, _ = url.Parse("https://www.ocuroot.com/")

func Canonical(p ...string) string {
	u := *BaseURL
	u.Path = path.Join(u.Path, path.Join(p...))
	return u.String()
}
