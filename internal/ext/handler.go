package ext

import (
	"github.com/kyleconroy/sqlc/internal/plugin"
)

type Handler interface {
	Generate(*plugin.CodeGenRequest) (*plugin.CodeGenResponse, error)
}

type wrapper struct {
	fn func(*plugin.CodeGenRequest) (*plugin.CodeGenResponse, error)
}

func (w *wrapper) Generate(req *plugin.CodeGenRequest) (*plugin.CodeGenResponse, error) {
	return w.fn(req)
}

func HandleFunc(fn func(*plugin.CodeGenRequest) (*plugin.CodeGenResponse, error)) Handler {
	return &wrapper{fn}
}
