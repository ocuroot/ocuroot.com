package main

import (
	"fmt"
	"net/http"
	"path"
	"strings"
)

type DevServer struct {
	Renderer Renderer
}

func (d *DevServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestPath := r.URL.Path
	requestPath = strings.TrimPrefix(requestPath, "/")

	if !d.Renderer.HasPath(requestPath) {
		if d.Renderer.HasPath(path.Join(requestPath, "index.html")) {
			requestPath = path.Join(requestPath, "index.html")
		}
		if d.Renderer.HasPath(path.Join(requestPath, "index.htm")) {
			requestPath = path.Join(requestPath, "index.htm")
		}
	}

	// Render the confirmed target file
	data, err := d.Renderer.RenderPath(r.Context(), requestPath)
	if err != nil {
		fmt.Printf("error rendering path %s: %v\n", requestPath, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	contentType := http.DetectContentType(data)
	if strings.HasSuffix(requestPath, ".svg") {
		contentType = "image/svg+xml"
	}
	if strings.HasSuffix(requestPath, ".css") {
		contentType = "text/css"
	}
	w.Header().Set("Content-Type", contentType)
	w.Write(data)
}
