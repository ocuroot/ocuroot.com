package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ocuroot/templbuildr/site"
	"github.com/ocuroot/ui/css"

	_ "embed"
)

//go:embed static/css/blog.css
var blogCSS []byte

func init() {
	css.Register("blog", blogCSS)
}

// Author represents blog post author information
type Author struct {
	Name    string `yaml:"name"`
	Picture string `yaml:"picture"`
}

// CoverImage represents blog post cover image information
type CoverImage struct {
	Src       string `yaml:"src"`
	Alt       string `yaml:"alt"`
	Credit    string `yaml:"credit"`
	CreditURL string `yaml:"creditURL"`
}

// OGImage represents Open Graph image information
type OGImage struct {
	URL string `yaml:"url"`
}

// BlogPost represents a parsed blog post
type BlogPost struct {
	Title      string     `yaml:"title"`
	Slug       string     `yaml:"slug"`
	Excerpt    string     `yaml:"excerpt"`
	Date       time.Time  `yaml:"date"`
	Author     Author     `yaml:"author"`
	CoverImage CoverImage `yaml:"coverImage"`
	OGImage    OGImage    `yaml:"ogImage"`
	Tags       []string   `yaml:"tags"`
	Content    string     // HTML content from markdown
	Raw        string     // Original markdown content
}

// BlogComponent wraps a blog post for the ConcreteRenderer
type BlogComponent struct {
	Post *BlogPost
}

func (bc *BlogComponent) Render(ctx context.Context, w io.Writer) error {
	// Convert main package BlogPost to site package BlogPost
	sitePost := &site.BlogPost{
		Title:      bc.Post.Title,
		Slug:       bc.Post.Slug,
		Excerpt:    bc.Post.Excerpt,
		Date:       bc.Post.Date,
		Author:     site.Author{Name: bc.Post.Author.Name, Picture: bc.Post.Author.Picture},
		CoverImage: site.CoverImage{Src: bc.Post.CoverImage.Src, Alt: bc.Post.CoverImage.Alt, Credit: bc.Post.CoverImage.Credit, CreditURL: bc.Post.CoverImage.CreditURL},
		Tags:       bc.Post.Tags,
		Content:    bc.Post.Content,
	}
	return site.BlogPostPage(sitePost).Render(ctx, w)
}

// BlogListComponent wraps a list of blog posts for the ConcreteRenderer
type BlogListComponent struct {
	Posts []*BlogPost
}

func (blc *BlogListComponent) Render(ctx context.Context, w io.Writer) error {
	// Convert main package BlogPost slice to site package BlogPost slice
	sitePosts := make([]*site.BlogPost, len(blc.Posts))
	for i, post := range blc.Posts {
		sitePosts[i] = &site.BlogPost{
			Title:      post.Title,
			Slug:       post.Slug,
			Excerpt:    post.Excerpt,
			Date:       post.Date,
			Author:     site.Author{Name: post.Author.Name, Picture: post.Author.Picture},
			CoverImage: site.CoverImage{Src: post.CoverImage.Src, Alt: post.CoverImage.Alt, Credit: post.CoverImage.Credit, CreditURL: post.CoverImage.CreditURL},
			Tags:       post.Tags,
			Content:    post.Content,
		}
	}
	fmt.Printf("Got %d blog posts", len(sitePosts))
	return site.BlogListPage(sitePosts).Render(ctx, w)
}

// BlogManager handles loading and managing blog posts
type BlogManager struct {
	parser *BlogParser
	posts  map[string]*BlogPost // slug -> post
	sorted []*BlogPost          // posts sorted by date (newest first)
}

// NewBlogManager creates a new blog manager
func NewBlogManager() *BlogManager {
	return &BlogManager{
		parser: NewBlogParser(),
		posts:  make(map[string]*BlogPost),
	}
}

// LoadPosts loads all blog posts from the _posts directory
func (bm *BlogManager) LoadPosts() error {
	// Find all markdown files in _posts directory
	files, err := filepath.Glob("_posts/*.md")
	if err != nil {
		return fmt.Errorf("failed to find blog posts: %w", err)
	}

	var posts []*BlogPost

	for _, file := range files {
		fmt.Printf("Loading blog post %s\n", file)
		post, err := bm.parser.ParseFile(file)
		if err != nil {
			fmt.Printf("Warning: failed to parse %s: %v\n", file, err)
			continue
		}

		// Generate slug from filename if not provided in frontmatter
		if post.Slug == "" {
			filename := filepath.Base(file)
			post.Slug = strings.TrimSuffix(filename, ".md")
			// Remove number prefix (e.g., "04-why-ocuroot" -> "why-ocuroot")
			if idx := strings.Index(post.Slug, "-"); idx > 0 && idx < 3 {
				post.Slug = post.Slug[idx+1:]
			}
		}

		posts = append(posts, post)
		bm.posts[post.Slug] = post
	}

	// Sort posts by date (newest first)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	bm.sorted = posts
	return nil
}

// RegisterWithRenderer registers all blog routes with the ConcreteRenderer
func (bm *BlogManager) RegisterWithRenderer(r *ConcreteRenderer) {
	// Register blog list page
	r.Register("blog/index.html", &BlogListComponent{Posts: bm.sorted})

	// Register individual blog post pages
	for slug, post := range bm.posts {
		path := fmt.Sprintf("blog/%s/index.html", slug)
		r.Register(path, &BlogComponent{Post: post})
	}

	// Register blog assets (cover images, etc.)
	bm.registerBlogAssets(r)
}

// registerBlogAssets registers blog-related static assets
func (bm *BlogManager) registerBlogAssets(r *ConcreteRenderer) {
	// Walk through all files under static/assets/blog
	err := filepath.Walk("static/assets/blog", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories, only register files
		if info.IsDir() {
			return nil
		}

		// Convert static/assets/blog/why-ocuroot/cover.jpg -> assets/blog/why-ocuroot/cover.jpg
		webPath := strings.TrimPrefix(path, "static/")
		r.Register(webPath, StaticFileComponent(path))
		return nil
	})

	if err != nil {
		fmt.Printf("Warning: failed to walk blog assets directory: %v\n", err)
	}
}

// GetPost returns a blog post by slug
func (bm *BlogManager) GetPost(slug string) (*BlogPost, bool) {
	post, exists := bm.posts[slug]
	return post, exists
}

// GetAllPosts returns all posts sorted by date (newest first)
func (bm *BlogManager) GetAllPosts() []*BlogPost {
	return bm.sorted
}
