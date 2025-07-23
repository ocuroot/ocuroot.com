package main

import (
	"context"
	"io"
)

type StaticComponent []byte

func (s StaticComponent) Render(ctx context.Context, w io.Writer) error {
	_, err := w.Write(s)
	return err
}
