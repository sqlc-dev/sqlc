package opts

import (
	"os"
	"strings"
)

// The SQLCDEBUG variable controls debugging variables within the runtime. It
// is a comma-separated list of name=val pairs setting these named variables:
//
//     dumpast: setting dumpast=1 will print the AST of every SQL statement
//     dumpcatalog: setting dumpcatalog=1 will print the parsed database schema

type Debug struct {
	DumpAST     bool
	DumpCatalog bool
}

func DebugFromEnv() (Debug, error) {
	d := Debug{}
	val := os.Getenv("SQLCDEBUG")
	if val == "" {
		return d, nil
	}
	for _, pair := range strings.Split(val, ",") {
		switch strings.TrimSpace(pair) {
		case "dumpast=1":
			d.DumpAST = true
		case "dumpcatalog=1":
			d.DumpCatalog = true
		}
	}
	return d, nil
}
