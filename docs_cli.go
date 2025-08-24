package main

import (
	"bytes"
	"context"
	"fmt"

	"github.com/a-h/templ"
	"github.com/ocuroot/ocuroot/client/commands"
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
