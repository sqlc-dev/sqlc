package wasm

import "github.com/tetratelabs/wazero"

type Runner struct {
	URL    string
	SHA256 string
	Env    []string

	rt wazero.Runtime
}
