//go:build !windows && cgo

package expander

import (
	nodes "github.com/pganalyze/pg_query_go/v6"
)

var parse = nodes.Parse
var deparse = nodes.Deparse
