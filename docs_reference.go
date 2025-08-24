package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"

	"github.com/a-h/templ"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/ocuroot/ocuroot/client/commands"
	"github.com/ocuroot/ocuroot/sdk"
	"github.com/ocuroot/templbuildr/site"
	"github.com/ocuroot/ui/components/docnav"
)

func CLIPath(name string) string {
	return fmt.Sprintf("/reference/cli/%s", name)
}

func CLINav() []docnav.NavLink {
	root := commands.RootCmd
	allCommands := root.Commands()
	var navLinks []docnav.NavLink
	for _, cmd := range allCommands {
		navLinks = append(navLinks, docnav.NavLink{
			Title: cmd.Name(),
			URL:   templ.SafeURL(fmt.Sprintf("/docs/%v", CLIPath(cmd.Name()))),
		})
	}
	return navLinks
}

func (dm *DocsManager) LoadCLIPages() error {
	root := commands.RootCmd
	allCommands := root.Commands()
	for _, cmd := range allCommands {
		page := &Content[DocPage]{
			FrontMatter: DocPage{
				Title: cmd.Name(),
				Path:  CLIPath(cmd.Name()),
			},
		}

		render := site.CLIContent(cmd)

		var buf bytes.Buffer
		if err := render.Render(context.Background(), &buf); err != nil {
			return err
		}
		page.Content = buf.String()

		dm.pages[CLIPath(cmd.Name())] = page
	}

	return nil
}

func SDKPath(name string) string {
	return fmt.Sprintf("/reference/sdk/%s", name)
}

func sdkFiles() []string {
	var names []string
	err := fs.WalkDir(sdk.Builtins, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() &&
			strings.HasSuffix(p, ".star") &&
			!strings.Contains(p, "readme") &&
			!strings.Contains(p, "after") {
			names = append(names, p)
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return names
}

func sdkName(p string) string {
	name := strings.TrimSuffix(p, ".star")
	name = path.Base(name)
	return name
}

func SDKNav() []docnav.NavLink {
	var navLinks []docnav.NavLink
	for _, p := range sdkFiles() {
		name := sdkName(p)
		navLinks = append(navLinks, docnav.NavLink{
			Title: name,
			URL:   templ.SafeURL(fmt.Sprintf("/docs/%v", SDKPath(name))),
		})
	}
	return navLinks
}

func (dm *DocsManager) LoadSDKPages() error {
	for _, f := range sdkFiles() {
		page := &Content[DocPage]{
			FrontMatter: DocPage{
				Title: sdkName(f),
				Path:  SDKPath(sdkName(f)),
			},
		}

		data, err := sdk.Builtins.ReadFile(fmt.Sprintf("%v", f))
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		err = Highlight(&buf, string(data), "python", "html", "monokai")
		if err != nil {
			return err
		}

		page.Content = buf.String()

		dm.pages[SDKPath(sdkName(f))] = page
	}

	return nil
}

func Highlight(w io.Writer, source, lexer, formatter, style string) error {
	// Determine lexer.
	l := lexers.Get(lexer)
	if l == nil {
		l = lexers.Analyse(source)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	// Determine formatter.
	f := html.New(html.WithClasses(true))

	// Determine style.
	s := styles.Get(style)
	if s == nil {
		s = styles.Fallback
	}

	it, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}
	return f.Format(w, s, it)
}
