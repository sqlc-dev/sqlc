package opts

import (
	"os"
	"strings"
)

// The SQLCDEBUG variable controls debugging variables within the runtime. It
// is a comma-separated list of name=val pairs setting these named variables:
//
//     dumpast: setting dumpast=1 will print the AST of every SQL statement

type Debug struct {
	DumpAST bool
}

func DebugFromEnv() (Debug, error) {
	d := Debug{}
	val := os.Getenv("SQLCDEBUG")
	if val == "" {
		return d, nil
	}
	for _, pair := range strings.Split(val, ",") {
		if pair == "dumpast=1" {
			d.DumpAST = true
		}
	}
	return d, nil
}
