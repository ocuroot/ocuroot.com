package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/a-h/templ"
)

var ErrPathNotFound = errors.New("path not found")

type Renderer interface {
	Register(path string, c templ.Component)

	RenderAll(ctx context.Context, outputDir string) error
	RenderPath(ctx context.Context, path string) ([]byte, error)
	HasPath(path string) bool
}

type ConcreteRenderer struct {
	Paths map[string]templ.Component
}

func (r *ConcreteRenderer) Register(path string, c templ.Component) {
	r.Paths[path] = c
}

func (r *ConcreteRenderer) RenderAll(ctx context.Context, outputDir string) error {
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, 0755)
		if err != nil {
			return err
		}
	}

	for path, component := range r.Paths {
		f, err := os.Create(filepath.Join(outputDir, path))
		if err != nil {
			return err
		}
		defer f.Close()

		err = component.Render(ctx, f)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ConcreteRenderer) HasPath(path string) bool {
	_, found := r.Paths[path]
	return found
}

func (r *ConcreteRenderer) RenderPath(ctx context.Context, path string) ([]byte, error) {
	component, found := r.Paths[path]
	if !found {
		return nil, ErrPathNotFound
	}

	buf := &bytes.Buffer{}

	err := component.Render(ctx, buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
