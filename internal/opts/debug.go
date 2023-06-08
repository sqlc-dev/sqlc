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
//     trace: setting trace=<path> will output a trace

type Debug struct {
	DumpAST        bool
	DumpCatalog    bool
	Trace          string
	ProcessPlugins bool
}

func DebugFromEnv() Debug {
	return DebugFromString(os.Getenv("SQLCDEBUG"))
}

func DebugFromString(val string) Debug {
	d := Debug{
		ProcessPlugins: true,
	}
	if val == "" {
		return d
	}
	for _, pair := range strings.Split(val, ",") {
		pair = strings.TrimSpace(pair)
		switch {
		case pair == "dumpast=1":
			d.DumpAST = true
		case pair == "dumpcatalog=1":
			d.DumpCatalog = true
		case strings.HasPrefix(pair, "trace="):
			traceName := strings.TrimPrefix(pair, "trace=")
			if traceName == "1" {
				d.Trace = "trace.out"
			} else {
				d.Trace = traceName
			}
		case pair == "processplugins=0":
			d.ProcessPlugins = false
		}
	}
	return d
}
