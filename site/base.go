package site

import (
	"net/url"
	"path"
	"strings"
)

var BaseURL, _ = url.Parse("https://www.ocuroot.com/")

func Canonical(p ...string) string {
	u := *BaseURL
	u.Path = path.Join(u.Path, path.Join(p...))

	// Add a trailing slash to directories to avoid redirects
	if !strings.HasSuffix(u.Path, "/") && !strings.Contains(path.Base(u.Path), ".") {
		u.Path += "/"
	}
	return u.String()
}
