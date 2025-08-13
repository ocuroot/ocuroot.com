package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark-highlighting/v2"
	"gopkg.in/yaml.v3"
)

type Content[F any] struct {
	FrontMatter F
	Date        time.Time `yaml:"date"`
	Content     string    // HTML content from markdown
	Raw         string    // Original markdown content
}

// Parser handles parsing markdown files with YAML frontmatter
type Parser[F any] struct {
	markdown goldmark.Markdown
}

// NewParser creates a new parser
func NewParser[F any]() *Parser[F] {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
				highlighting.WithFormatOptions(
					chromahtml.WithClasses(true),
				),
			),
			NewTemplInjector(),
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

	return &Parser[F]{
		markdown: md,
	}
}

// ParseFile parses a markdown file with YAML frontmatter
func (bp *Parser[F]) ParseFile(filename string) (*Content[F], error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	return bp.Parse(content)
}

// Parse parses markdown content with YAML frontmatter
func (bp *Parser[F]) Parse(content []byte) (*Content[F], error) {
	// Split frontmatter and content
	parts := bytes.SplitN(content, []byte("---"), 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid frontmatter format")
	}

	frontmatter := parts[1]
	markdown := parts[2]

	// Parse frontmatter
	var post Content[F]
	if err := yaml.Unmarshal(frontmatter, &post); err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}
	if err := yaml.Unmarshal(frontmatter, &post.FrontMatter); err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Convert markdown to HTML
	var buf bytes.Buffer
	if err := bp.markdown.Convert(markdown, &buf); err != nil {
		return nil, fmt.Errorf("failed to convert markdown: %w", err)
	}

	post.Content = buf.String()
	post.Raw = string(markdown)

	// Parse date if it's a string
	if post.Date.IsZero() {
		// Try to parse date from various formats
		dateStr := strings.TrimSpace(string(frontmatter))
		if strings.Contains(dateStr, "date:") {
			lines := strings.Split(dateStr, "\n")
			for _, line := range lines {
				if strings.HasPrefix(strings.TrimSpace(line), "date:") {
					dateValue := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "date:"))
					dateValue = strings.Trim(dateValue, `"'`)

					// Try different date formats
					formats := []string{
						"2006-01-02",
						"2006-01-02T15:04:05Z07:00",
						"2006-01-02 15:04:05",
					}

					for _, format := range formats {
						if parsedDate, err := time.Parse(format, dateValue); err == nil {
							post.Date = parsedDate
							break
						}
					}
					break
				}
			}
		}
	}

	return &post, nil
}
