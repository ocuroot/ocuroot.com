package main

import (
	"context"
	"flag"
	"fmt"
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
	r.Register("style.css", css.NewService())
	r.Register("script.js", js.NewService())
	r.Register("static/anon_user.svg", StaticComponent(assets.AnonUser))
	r.Register("static/logo.svg", StaticComponent(assets.Logo))
	r.Register("index.html", site.Index())
	r.Register("hello.html", hello("Tom"))

	if *dev {
		// Run development server
		devServer := &DevServer{
			Renderer: r,
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
