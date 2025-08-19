package main

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ocuroot/templbuildr/site"
	"github.com/ocuroot/ui/css"

	_ "embed"
)

//go:embed site/blog.css
var blogCSS []byte

func init() {
	css.Default().Add(blogCSS)
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
	Author     Author     `yaml:"author"`
	CoverImage CoverImage `yaml:"coverImage"`
	OGImage    OGImage    `yaml:"ogImage"`
	Tags       []string   `yaml:"tags"`
}

// BlogComponent wraps a blog post for the ConcreteRenderer
type BlogComponent struct {
	Post *Content[BlogPost]
}

func (bc *BlogComponent) Render(ctx context.Context, w io.Writer) error {
	// Convert main package BlogPost to site package BlogPost
	sitePost := &site.BlogPost{
		Title:      bc.Post.FrontMatter.Title,
		Slug:       bc.Post.FrontMatter.Slug,
		Excerpt:    bc.Post.FrontMatter.Excerpt,
		Date:       bc.Post.Date,
		Author:     site.Author{Name: bc.Post.FrontMatter.Author.Name, Picture: bc.Post.FrontMatter.Author.Picture},
		CoverImage: site.CoverImage{Src: bc.Post.FrontMatter.CoverImage.Src, Alt: bc.Post.FrontMatter.CoverImage.Alt, Credit: bc.Post.FrontMatter.CoverImage.Credit, CreditURL: bc.Post.FrontMatter.CoverImage.CreditURL},
		Tags:       bc.Post.FrontMatter.Tags,
		Content:    bc.Post.Content,
	}
	return site.BlogPostPage(sitePost).Render(ctx, w)
}

// BlogListComponent wraps a list of blog posts for the ConcreteRenderer
type BlogListComponent struct {
	Posts []*Content[BlogPost]
}

func (blc *BlogListComponent) Render(ctx context.Context, w io.Writer) error {
	// Convert main package BlogPost slice to site package BlogPost slice
	sitePosts := make([]*site.BlogPost, len(blc.Posts))
	for i, post := range blc.Posts {
		sitePosts[i] = &site.BlogPost{
			Title:      post.FrontMatter.Title,
			Slug:       post.FrontMatter.Slug,
			Excerpt:    post.FrontMatter.Excerpt,
			Date:       post.Date,
			Author:     site.Author{Name: post.FrontMatter.Author.Name, Picture: post.FrontMatter.Author.Picture},
			CoverImage: site.CoverImage{Src: post.FrontMatter.CoverImage.Src, Alt: post.FrontMatter.CoverImage.Alt, Credit: post.FrontMatter.CoverImage.Credit, CreditURL: post.FrontMatter.CoverImage.CreditURL},
			Tags:       post.FrontMatter.Tags,
			Content:    post.Content,
		}
	}
	fmt.Printf("Got %d blog posts", len(sitePosts))
	return site.BlogListPage(sitePosts).Render(ctx, w)
}

// BlogManager handles loading and managing blog posts
type BlogManager struct {
	parser *Parser[BlogPost]
	posts  map[string]*Content[BlogPost] // slug -> post
	sorted []*Content[BlogPost]          // posts sorted by date (newest first)
}

// NewBlogManager creates a new blog manager
func NewBlogManager() *BlogManager {
	return &BlogManager{
		parser: NewParser[BlogPost](),
		posts:  make(map[string]*Content[BlogPost]),
	}
}

func (bm *BlogManager) LoadAndRegister(r *ConcreteRenderer) error {
	if err := bm.LoadPosts(); err != nil {
		return err
	}
	bm.RegisterWithRenderer(r)
	return nil
}

// LoadPosts loads all blog posts from the _posts directory
func (bm *BlogManager) LoadPosts() error {
	// Find all markdown files in _posts directory
	files, err := filepath.Glob("_posts/*.md")
	if err != nil {
		return fmt.Errorf("failed to find blog posts: %w", err)
	}

	var posts []*Content[BlogPost]

	for _, file := range files {
		fmt.Printf("Loading blog post %s\n", file)
		post, err := bm.parser.ParseFile(file)
		if err != nil {
			fmt.Printf("Warning: failed to parse %s: %v\n", file, err)
			continue
		}

		// Generate slug from filename if not provided in frontmatter
		if post.FrontMatter.Slug == "" {
			filename := filepath.Base(file)
			post.FrontMatter.Slug = strings.TrimSuffix(filename, ".md")
			// Remove number prefix (e.g., "04-why-ocuroot" -> "why-ocuroot")
			if idx := strings.Index(post.FrontMatter.Slug, "-"); idx > 0 && idx < 3 {
				post.FrontMatter.Slug = post.FrontMatter.Slug[idx+1:]
			}
		}

		posts = append(posts, post)
		bm.posts[post.FrontMatter.Slug] = post
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
}

// GetPost returns a blog post by slug
func (bm *BlogManager) GetPost(slug string) (*Content[BlogPost], bool) {
	post, exists := bm.posts[slug]
	return post, exists
}

// GetAllPosts returns all posts sorted by date (newest first)
func (bm *BlogManager) GetAllPosts() []*Content[BlogPost] {
	return bm.sorted
}
