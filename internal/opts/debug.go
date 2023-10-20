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
//     processplugins: setting processplugins=0 will disable process-based plugins
//     databases: setting databases=managed will disable connections to databases via URI
//     dumpvetenv: setting dumpvetenv=1 will print the variables available to
//         a vet rule during evaluation
//     dumpexplain: setting dumpexplain=1 will print the JSON-formatted output
//         from executing EXPLAIN ... on a query during vet rule evaluation

type Debug struct {
	DumpAST              bool
	DumpCatalog          bool
	Trace                string
	ProcessPlugins       bool
	OnlyManagedDatabases bool
	DumpVetEnv           bool
	DumpExplain          bool
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
		case pair == "databases=managed":
			d.OnlyManagedDatabases = true
		case pair == "dumpvetenv=1":
			d.DumpVetEnv = true
		case pair == "dumpexplain=1":
			d.DumpExplain = true
		}
	}
	return d
}
