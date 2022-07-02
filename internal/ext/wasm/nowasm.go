//go:build !(cgo && ((linux && amd64) || (linux && arm64) || (darwin && amd64) || (darwin && arm64) || (windows && amd64)))

package wasm

import (
	"fmt"

	"github.com/kyleconroy/sqlc/internal/plugin"
)

type Runner struct {
}

func (r *Runner) Generate(req *plugin.CodeGenRequest) (*plugin.CodeGenResponse, error) {
	return nil, fmt.Errorf("sqlc built without wasmtime support")
}
