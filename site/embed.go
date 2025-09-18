package site

import (
	"embed"
	_ "embed"
	"io/fs"
	"log"
	"sort"
	"strings"

	"github.com/ocuroot/ui/css"
	"github.com/ocuroot/ui/js"
)

//go:embed **
var allFiles embed.FS

func init() {
	// Walk the embedded filesystem to collect all CSS and JS files
	cssFiles := make(map[string]string) // name -> path
	jsFiles := make(map[string]string)  // name -> path

	err := fs.WalkDir(allFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			name := d.Name()
			if strings.HasSuffix(name, ".css") {
				cssFiles[name] = path
			} else if strings.HasSuffix(name, ".js") {
				jsFiles[name] = path
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to walk embedded filesystem: %v", err)
	}

	// Get sorted file names for stable output
	var cssNames []string
	for name := range cssFiles {
		cssNames = append(cssNames, name)
	}
	sort.Strings(cssNames)

	var jsNames []string
	for name := range jsFiles {
		jsNames = append(jsNames, name)
	}
	sort.Strings(jsNames)

	// Add each CSS file to the default CSS registry
	for _, name := range cssNames {
		path := cssFiles[name]
		content, err := allFiles.ReadFile(path)
		if err != nil {
			log.Fatalf("Failed to read CSS file %s: %v", path, err)
		}
		css.Default().Add(content)
	}

	// Add each JS file to the default JS registry
	for _, name := range jsNames {
		path := jsFiles[name]
		content, err := allFiles.ReadFile(path)
		if err != nil {
			log.Fatalf("Failed to read JS file %s: %v", path, err)
		}
		js.Default().Add(content)
	}
}
