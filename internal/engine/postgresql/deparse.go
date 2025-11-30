//go:build !windows && cgo

package postgresql

import (
	"strings"

	nodes "github.com/pganalyze/pg_query_go/v6"
)

func Deparse(tree *nodes.ParseResult) (string, error) {
	output, err := nodeDeparse(tree)
	if err != nil {
		return output, err
	}
	return fixDeparse(output), nil
}

// fixDeparse corrects known bugs in pg_query_go's Deparse output
func fixDeparse(s string) string {
	// Fix missing space before SKIP LOCKED
	// pg_query_go outputs "OF tableSKIP LOCKED" instead of "OF table SKIP LOCKED"
	s = strings.ReplaceAll(s, "SKIP LOCKED", " SKIP LOCKED")
	s = strings.ReplaceAll(s, "  SKIP LOCKED", " SKIP LOCKED") // normalize double spaces
	return s
}
