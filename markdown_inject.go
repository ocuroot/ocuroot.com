package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/ocuroot/templbuildr/site"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	gast "github.com/yuin/goldmark/ast"
)

// NewTemplInjector creates a goldmark extension to add templ components
// to markdown files.
func NewTemplInjector() *TemplInjector {
	return &TemplInjector{}
}

type TemplInjector struct {
}

func NewTemplParser() parser.InlineParser {
	return &TemplParser{}
}

type TemplParser struct {
}

// Parse implements parser.InlineParser.
func (t *TemplParser) Parse(parent gast.Node, block text.Reader, pc parser.Context) gast.Node {
	line, _ := block.PeekLine()
	if strings.HasPrefix(string(line), "@") {
		block.Advance(len(line))
		return NewTemplItem(string(line)[1:])
	}
	return nil
}

// Trigger implements parser.InlineParser.
func (t *TemplParser) Trigger() []byte {
	return []byte{'@'}
}

func NewTemplHTMLRenderer() renderer.NodeRenderer {
	return &TemplHTMLRenderer{}
}

type TemplHTMLRenderer struct {
}

// RegisterFuncs implements renderer.NodeRenderer.
func (t *TemplHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindTemplItem, t.renderTemplItem)
}

func (t *TemplHTMLRenderer) renderTemplItem(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	if !entering {
		return gast.WalkContinue, nil
	}
	n := node.(*TemplItem)

	if n.Name == "ArchitectureDiagram" {
		ad := site.ArchitectureDiagram()
		ad.Render(context.Background(), w)
		return gast.WalkContinue, nil
	}

	w.WriteString("Unknown template: " + n.Name)
	return gast.WalkContinue, nil
}

func (ti *TemplInjector) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewTemplParser(), 500),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewTemplHTMLRenderer(), 500),
	))
}

type TemplItem struct {
	gast.BaseBlock
	Name string
}

// Dump implements Node.Dump.
func (n *TemplItem) Dump(source []byte, level int) {
	m := map[string]string{
		"Name": fmt.Sprintf("%v", n.Name),
	}
	gast.DumpHelper(n, source, level, m, nil)
}

// KindTaskCheckBox is a NodeKind of the TaskCheckBox node.
var KindTemplItem = gast.NewNodeKind("TemplItem")

// Kind implements Node.Kind.
func (n *TemplItem) Kind() gast.NodeKind {
	return KindTemplItem
}

// NewTemplItem returns a new TemplItem node.
func NewTemplItem(name string) *TemplItem {
	return &TemplItem{
		Name: name,
	}
}
