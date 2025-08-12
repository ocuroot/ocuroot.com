package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/ocuroot/templbuildr/site"
	"github.com/ocuroot/ui/assets"
	"github.com/ocuroot/ui/css"
	"github.com/ocuroot/ui/js"
)

func main() {
	// Define flags
	dev := flag.Bool("dev", false, "run in development mode with hot reload server")
	devPort := flag.Int("dev-port", 3000, "port for development server (default 3000)")
	flag.Parse()

	// Create renderer
	r := &ConcreteRenderer{
		Paths: make(map[string]templ.Component),
	}
	r.Register("favicon.ico", StaticComponent(assets.Favicon))
	r.Register("static/anon_user.svg", StaticComponent(assets.AnonUser))
	r.Register("static/logo.svg", StaticComponent(assets.Logo))
	r.Register("static/images/why-ocuroot-cover.jpg", StaticFileComponent("static/images/why-ocuroot-cover.jpg"))
	r.Register("static/images/social.png", StaticFileComponent("static/images/social.png"))
	r.Register("images/see-production.svg", StaticFileComponent("static/images/see-production.svg"))
	r.Register("images/software-demo.jpg", StaticFileComponent("static/images/software-demo.jpg"))
	r.Register("images/server_racks.jpg", StaticFileComponent("static/images/server_racks.jpg"))
	r.Register("images/code.png", StaticFileComponent("static/images/code.png"))
	r.Register("images/export-history.gif", StaticFileComponent("static/images/export-history.gif"))
	r.Register("images/add_environment.gif", StaticFileComponent("static/images/add_environment.gif"))
	r.Register("index.html", site.Index())
	r.Register("solutions/cost-management/index.html", site.CostManagementPage())
	r.Register("solutions/productivity-satisfaction/index.html", site.ProductivitySatisfactionPage())
	r.Register("solutions/security-compliance/index.html", site.SecurityCompliancePage())
	r.Register("solutions/technical-agility/index.html", site.TechnicalAgilityPage())
	r.Register("demo/index.html", site.DemoPage())

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

	if *dev {
		// Run development server
		devServer := &DevServer{
			Renderer: r,
		}

		for path := range r.Paths {
			fmt.Printf("Registered path %s\n", path)
		}

		addr := fmt.Sprintf(":%d", *devPort)
		fmt.Printf("Starting development server on http://localhost%s\n", addr)
		log.Fatal(http.ListenAndServe(addr, devServer))
		return
	}

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
