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
// Example plugin:
//
//	package main
//
//	import "github.com/sqlc-dev/sqlc/pkg/engine"
//
//	func main() {
//		engine.Run(engine.Handler{
//			PluginName:        "my-plugin",
//			PluginVersion:     "1.0.0",
//			Parse:             handleParse,
//			GetCatalog:        handleGetCatalog,
//			IsReservedKeyword: handleIsReservedKeyword,
//			GetCommentSyntax:  handleGetCommentSyntax,
//			GetDialect:        handleGetDialect,
//		})
//	}
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

	Parse             func(*ParseRequest) (*ParseResponse, error)
	GetCatalog        func(*GetCatalogRequest) (*GetCatalogResponse, error)
	IsReservedKeyword func(*IsReservedKeywordRequest) (*IsReservedKeywordResponse, error)
	GetCommentSyntax  func(*GetCommentSyntaxRequest) (*GetCommentSyntaxResponse, error)
	GetDialect        func(*GetDialectRequest) (*GetDialectResponse, error)
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

	case "get_catalog":
		var req GetCatalogRequest
		if len(input) > 0 {
			proto.Unmarshal(input, &req)
		}
		if h.GetCatalog == nil {
			return fmt.Errorf("get_catalog not implemented")
		}
		output, err = h.GetCatalog(&req)

	case "is_reserved_keyword":
		var req IsReservedKeywordRequest
		if err := proto.Unmarshal(input, &req); err != nil {
			return fmt.Errorf("parsing request: %w", err)
		}
		if h.IsReservedKeyword == nil {
			return fmt.Errorf("is_reserved_keyword not implemented")
		}
		output, err = h.IsReservedKeyword(&req)

	case "get_comment_syntax":
		var req GetCommentSyntaxRequest
		if len(input) > 0 {
			proto.Unmarshal(input, &req)
		}
		if h.GetCommentSyntax == nil {
			return fmt.Errorf("get_comment_syntax not implemented")
		}
		output, err = h.GetCommentSyntax(&req)

	case "get_dialect":
		var req GetDialectRequest
		if len(input) > 0 {
			proto.Unmarshal(input, &req)
		}
		if h.GetDialect == nil {
			return fmt.Errorf("get_dialect not implemented")
		}
		output, err = h.GetDialect(&req)

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
