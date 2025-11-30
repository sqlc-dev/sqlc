//go:build windows || !cgo

package postgresql

import (
	nodes "github.com/wasilibs/go-pgquery"
)

var Parse = nodes.Parse
var Fingerprint = nodes.Fingerprint

var nodeDeparse = nodes.Deparse
