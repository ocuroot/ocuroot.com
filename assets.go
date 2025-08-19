package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RegisterStatic adds static content to the renderer.
// Used for images and other media.
func RegisterStatic(r *ConcreteRenderer) {
	err := filepath.Walk("static", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories, only register files
		if info.IsDir() {
			return nil
		}

		webPath := strings.TrimPrefix(path, "static/")
		r.Register(webPath, StaticFileComponent(path))
		return nil
	})

	if err != nil {
		fmt.Printf("Warning: failed to walk images directory: %v\n", err)
	}
}
