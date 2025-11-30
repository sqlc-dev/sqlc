//go:build windows || !cgo

package expander

import (
	nodes "github.com/wasilibs/go-pgquery"
)

var parse = nodes.Parse
var deparse = nodes.Deparse
