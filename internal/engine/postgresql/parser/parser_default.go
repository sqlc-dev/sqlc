//go:build !windows && cgo

package parser

import "github.com/pganalyze/pg_query_go/v6/parser"

type Error = parser.Error
