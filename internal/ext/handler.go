package ext

import (
	"context"

	"github.com/kyleconroy/sqlc/internal/plugin"
)

type Handler interface {
	Generate(context.Context, *plugin.CodeGenRequest) (*plugin.CodeGenResponse, error)
}

type wrapper struct {
	fn func(context.Context, *plugin.CodeGenRequest) (*plugin.CodeGenResponse, error)
}

func (w *wrapper) Generate(ctx context.Context, req *plugin.CodeGenRequest) (*plugin.CodeGenResponse, error) {
	return w.fn(ctx, req)
}

func HandleFunc(fn func(context.Context, *plugin.CodeGenRequest) (*plugin.CodeGenResponse, error)) Handler {
	return &wrapper{fn}
}
