//go:build !windows && cgo

package postgresql

import (
	nodes "github.com/pganalyze/pg_query_go/v6"
)

var Parse = nodes.Parse
var Fingerprint = nodes.Fingerprint
