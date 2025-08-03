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
	"gopkg.in/yaml.v3"
)

// BlogParser handles parsing markdown files with YAML frontmatter
type BlogParser struct {
	markdown goldmark.Markdown
}

// NewBlogParser creates a new blog parser
func NewBlogParser() *BlogParser {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	return &BlogParser{
		markdown: md,
	}
}

// ParseFile parses a markdown file with YAML frontmatter
func (bp *BlogParser) ParseFile(filename string) (*BlogPost, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	return bp.Parse(content)
}

// Parse parses markdown content with YAML frontmatter
func (bp *BlogParser) Parse(content []byte) (*BlogPost, error) {
	// Split frontmatter and content
	parts := bytes.SplitN(content, []byte("---"), 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid frontmatter format")
	}

	frontmatter := parts[1]
	markdown := parts[2]

	// Parse frontmatter
	var post BlogPost
	if err := yaml.Unmarshal(frontmatter, &post); err != nil {
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
