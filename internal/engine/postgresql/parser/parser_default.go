//go:build !windows && cgo
// +build !windows,cgo

package parser

import "github.com/pganalyze/pg_query_go/v4/parser"

type Error = parser.Error
