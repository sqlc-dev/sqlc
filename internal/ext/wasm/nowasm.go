//go:build nowasm || !(cgo && ((linux && amd64) || (linux && arm64) || (darwin && amd64) || (darwin && arm64) || (windows && amd64)))

package wasm

import (
	"context"
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/plugin"
)

func (r *Runner) Generate(ctx context.Context, req *plugin.CodeGenRequest) (*plugin.CodeGenResponse, error) {
	return nil, fmt.Errorf("sqlc built without wasmtime support")
}
