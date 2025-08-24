package main

import (
	"context"
	"fmt"
	"io"
	"path"
	"path/filepath"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/ocuroot/templbuildr/site"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// DocPage represents a parsed documentation page
type DocPage struct {
	Title string `yaml:"title"`
	Path  string `yaml:"path"`
}

// DocsManager handles loading and managing documentation
type DocsManager struct {
	parser *Parser[DocPage]
	pages  map[string]*Content[DocPage] // path -> page
}

// NewDocsManager creates a new docs manager
func NewDocsManager() *DocsManager {
	return &DocsManager{
		parser: NewParser[DocPage](),
		pages:  make(map[string]*Content[DocPage]),
	}
}

// DocsParser handles parsing markdown files with YAML frontmatter
type DocsParser struct {
	markdown goldmark.Markdown
}

// NewDocsParser creates a new docs parser
func NewDocsParser() *DocsParser {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
				highlighting.WithFormatOptions(
					chromahtml.WithClasses(true),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
			// html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	return &DocsParser{
		markdown: md,
	}
}

func (dm *DocsManager) LoadAndRegister(r *ConcreteRenderer) error {
	if err := dm.LoadPages(); err != nil {
		return err
	}
	if err := dm.LoadCLIPages(); err != nil {
		return err
	}
	if err := dm.LoadSDKPages(); err != nil {
		return err
	}

	dm.RegisterWithRenderer(r)
	return nil
}

// LoadPages loads all documentation pages from the _pages directory
func (dm *DocsManager) LoadPages() error {
	// Find all markdown files in _pages directory
	files, err := filepath.Glob("_docs/*.md")
	if err != nil {
		return fmt.Errorf("failed to find documentation pages: %w", err)
	}

	for _, file := range files {
		fmt.Printf("Loading documentation page %s\n", file)
		page, err := dm.parser.ParseFile(file)
		if err != nil {
			fmt.Printf("Warning: failed to parse %s: %v\n", file, err)
			continue
		}

		dm.pages[page.FrontMatter.Path] = page
	}

	return nil
}

type DocComponent struct {
	Post *Content[DocPage]
}

func (dc *DocComponent) Render(ctx context.Context, w io.Writer) error {
	// Convert main package DocPage to site package DocPage
	sitePost := &site.DocPage{
		Title:   dc.Post.FrontMatter.Title,
		Path:    dc.Post.FrontMatter.Path,
		Content: dc.Post.Content,

		CLINav: CLINav(),
		SDKNav: SDKNav(),
	}

	// Render template
	return site.DocsPage(sitePost).Render(ctx, w)
}

// RegisterWithRenderer registers all blog routes with the ConcreteRenderer
func (dm *DocsManager) RegisterWithRenderer(r *ConcreteRenderer) {
	// Register individual blog post pages
	for subPath, post := range dm.pages {
		path := path.Join("docs", subPath, "index.html")
		r.Register(path, &DocComponent{Post: post})
	}
}
