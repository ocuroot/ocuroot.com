package main

import (
	"context"
	"io"
	"os"
)

type StaticComponent []byte

func (s StaticComponent) Render(ctx context.Context, w io.Writer) error {
	_, err := w.Write(s)
	return err
}

// StaticFileComponent reads and serves a file from the local filesystem
type StaticFileComponent string

func (s StaticFileComponent) Render(ctx context.Context, w io.Writer) error {
	data, err := os.ReadFile(string(s))
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

type StringComponent string

func (s StringComponent) Render(ctx context.Context, w io.Writer) error {
	_, err := w.Write([]byte(s))
	return err
}
