//go:build !wasmtime

package wasm

import (
	"fmt"

	"github.com/kyleconroy/sqlc/internal/plugin"
)

func Generate(req *plugin.CodeGenRequest) (*plugin.CodeGenResponse, error) {
	return nil, fmt.Errorf("sqlc built without wasmtime support")
}
