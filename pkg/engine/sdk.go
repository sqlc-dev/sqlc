// Package engine provides types and utilities for building sqlc database engine plugins.
//
// Engine plugins allow external database backends to be used with sqlc.
// Plugins communicate with sqlc via Protocol Buffers over stdin/stdout.
//
// # Compatibility
//
// Go plugins that import this package are guaranteed to be compatible with sqlc
// at compile time. If the types change incompatibly, the plugin simply won't
// compile until it's updated to match the new interface.
//
// The Protocol Buffer schema is published at buf.build/sqlc/sqlc and ensures
// binary compatibility between sqlc and plugins.
//
// # Generating engine.pb.go
//
// Run from the repository root:
//
//	make proto-engine
//
// or:
//
//	protoc --go_out=. --go_opt=module=github.com/sqlc-dev/sqlc protos/engine/engine.proto
//
// Example plugin:
//
//	package main
//
//	import "github.com/sqlc-dev/sqlc/pkg/engine"
//
//	func main() {
//		engine.Run(engine.Handler{
//			PluginName:    "my-plugin",
//			PluginVersion: "1.0.0",
//			Parse:         handleParse,
//		})
//	}
//
//go:generate protoc -I../.. --go_out=../.. --go_opt=module=github.com/sqlc-dev/sqlc protos/engine/engine.proto
package engine

import (
	"fmt"
	"io"
	"os"

	"google.golang.org/protobuf/proto"
)

// Handler contains the functions that implement the engine plugin interface.
// All types used are Protocol Buffer messages defined in engine.proto.
type Handler struct {
	PluginName    string
	PluginVersion string

	Parse func(*ParseRequest) (*ParseResponse, error)
}

// Run runs the engine plugin with the given handler.
// It reads a protobuf request from stdin and writes a protobuf response to stdout.
func Run(h Handler) {
	if err := run(h, os.Args, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(h Handler, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s <method>", args[0])
	}

	method := args[1]
	input, err := io.ReadAll(stdin)
	if err != nil {
		return fmt.Errorf("reading stdin: %w", err)
	}

	var output proto.Message

	switch method {
	case "parse":
		var req ParseRequest
		if err := proto.Unmarshal(input, &req); err != nil {
			return fmt.Errorf("parsing request: %w", err)
		}
		if h.Parse == nil {
			return fmt.Errorf("parse not implemented")
		}
		output, err = h.Parse(&req)

	default:
		return fmt.Errorf("unknown method: %s", method)
	}

	if err != nil {
		return err
	}

	data, err := proto.Marshal(output)
	if err != nil {
		return fmt.Errorf("marshaling response: %w", err)
	}

	_, err = stdout.Write(data)
	return err
}
