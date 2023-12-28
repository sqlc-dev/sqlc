//go:build !windows && cgo
// +build !windows,cgo

package postgresql

import (
	nodes "github.com/pganalyze/pg_query_go/v4"
)

var Parse = nodes.Parse
var Fingerprint = nodes.Fingerprint
