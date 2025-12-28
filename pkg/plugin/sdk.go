// Package plugin provides types and utilities for building sqlc codegen plugins.
//
// Codegen plugins allow generating code in custom languages from sqlc.
// Plugins communicate with sqlc via Protocol Buffers over stdin/stdout.
//
// # Compatibility
//
// Go plugins that import this package are guaranteed to be compatible with sqlc
// at compile time. If the types change incompatibly, the plugin simply won't
// compile until it's updated to match the new interface.
//
// Example plugin:
//
//	package main
//
//	import "github.com/sqlc-dev/sqlc/pkg/plugin"
//
//	func main() {
//		plugin.Run(func(req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) {
//			// Generate code from req.Queries and req.Catalog
//			return &plugin.GenerateResponse{
//				Files: []*plugin.File{
//					{Name: "queries.txt", Contents: []byte("...")},
//				},
//			}, nil
//		})
//	}
package plugin

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"google.golang.org/protobuf/proto"
)

// GenerateFunc is the function signature for code generation.
type GenerateFunc func(*GenerateRequest) (*GenerateResponse, error)

// Run runs the codegen plugin with the given generate function.
// It reads a protobuf GenerateRequest from stdin and writes a GenerateResponse to stdout.
func Run(fn GenerateFunc) {
	if err := run(fn, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}

func run(fn GenerateFunc, stdin io.Reader, stdout, stderr io.Writer) error {
	reqBlob, err := io.ReadAll(stdin)
	if err != nil {
		return fmt.Errorf("reading stdin: %w", err)
	}

	var req GenerateRequest
	if err := proto.Unmarshal(reqBlob, &req); err != nil {
		return fmt.Errorf("unmarshaling request: %w", err)
	}

	resp, err := fn(&req)
	if err != nil {
		return fmt.Errorf("generating: %w", err)
	}

	respBlob, err := proto.Marshal(resp)
	if err != nil {
		return fmt.Errorf("marshaling response: %w", err)
	}

	w := bufio.NewWriter(stdout)
	if _, err := w.Write(respBlob); err != nil {
		return fmt.Errorf("writing response: %w", err)
	}
	return w.Flush()
}
