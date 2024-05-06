//go:build !windows && cgo

package parser

import "github.com/pganalyze/pg_query_go/v5/parser"

type Error = parser.Error
