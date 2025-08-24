package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/a-h/templ"
	"github.com/ocuroot/templbuildr/site"
	"github.com/ocuroot/ui/assets"
	"github.com/ocuroot/ui/css"
	"github.com/ocuroot/ui/js"
)

func main() {
	// Create renderer
	r := &ConcreteRenderer{
		Paths: make(map[string]templ.Component),
	}
	r.Register("favicon.ico", StaticComponent(assets.Favicon))
	r.Register("static/anon_user.svg", StaticComponent(assets.AnonUser))
	r.Register("static/logo.svg", StaticComponent(assets.Logo))
	r.Register("index.html", site.Index())
	r.Register("demo/index.html", site.DemoPage())
	r.Register("privacy/index.html", site.PrivacyPolicyPage())

	// Initialize and register blog posts
	blogManager := NewBlogManager()
	if err := blogManager.LoadAndRegister(r); err != nil {
		log.Printf("Warning: failed to load blog posts: %v", err)
	}

	// Initialize and register documentation
	docsManager := NewDocsManager()
	if err := docsManager.LoadAndRegister(r); err != nil {
		log.Printf("Warning: failed to load documentation: %v", err)
	}

	r.Register(css.Default().GetVersionedURL(), AsComponent(css.Default().GetCombined()))
	r.Register(js.Default().GetVersionedURL(), AsComponent(js.Default().GetCombined()))

	RegisterStatic(r)

	// Build static files
	err := r.RenderAll(context.Background(), "dist")
	if err != nil {
		log.Fatalf("failed to render: %v", err)
	}
	fmt.Println("Templates rendered to dist/")

}

type AsComponent string

func (ac AsComponent) Render(ctx context.Context, w io.Writer) error {
	w.Write([]byte(ac))
	return nil
}
