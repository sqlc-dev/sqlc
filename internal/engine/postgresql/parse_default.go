//go:build !windows && cgo
// +build !windows,cgo

package postgresql

import (
	cgonodes "github.com/pganalyze/pg_query_go/v4"
	wasinodes "github.com/wasilibs/go-pgquery"
	"google.golang.org/protobuf/proto"
)

func parseNodes(input string) (*wasinodes.ParseResult, error) {
	cgoRes, err := cgonodes.Parse(input)
	if err != nil {
		return nil, err
	}

	// It would be too tedious to maintain conversion logic in a way that
	// can target the two types from different packages despite being identical.
	// We go ahead and take a small performance hit by marshaling through
	// protobuf to unify the types. We must use the wasilibs version because
	// the upstream version requires cgo even when only accessing the proto.

	resBytes, err := proto.Marshal(cgoRes)
	if err != nil {
		return nil, err
	}

	var wasiRes wasinodes.ParseResult
	if err := proto.Unmarshal(resBytes, &wasiRes); err != nil {
		return nil, err
	}

	return &wasiRes, nil
}
