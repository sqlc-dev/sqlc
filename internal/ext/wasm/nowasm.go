//go:build nowasm || !(cgo && ((linux && amd64) || (linux && arm64) || (darwin && amd64) || (darwin && arm64) || (windows && amd64)))

package wasm

import (
	"fmt"
    "context"

	"github.com/kyleconroy/sqlc/internal/plugin"
)

type Runner struct {
	URL    string
	SHA256 string
}

func (r *Runner) Generate(ctx context.Context, req *plugin.CodeGenRequest) (*plugin.CodeGenResponse, error) {
	return nil, fmt.Errorf("sqlc built without wasmtime support")
}
