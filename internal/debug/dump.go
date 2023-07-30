package debug

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"

	"github.com/sqlc-dev/sqlc/internal/opts"
)

var Active bool
var Debug opts.Debug

func init() {
	Active = os.Getenv("SQLCDEBUG") != ""
	if Active {
		Debug = opts.DebugFromEnv()
	}
}

func Dump(n ...interface{}) {
	if Active {
		spew.Dump(n)
	}
}

func DumpAsJSON(a any) {
	if Active {
		out, _ := json.MarshalIndent(a, "", "  ")
		fmt.Printf("%s\n", out)
	}
}
