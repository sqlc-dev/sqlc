package debug

import (
	"encoding/json"
	"fmt"

	"github.com/davecgh/go-spew/spew"

	"github.com/sqlc-dev/sqlc/internal/sqlcdebug"
)

// Active reports whether SQLCDEBUG had any value set at startup. It
// remains a global so unrelated debug-spew sites that don't tie to a
// specific setting can gate their output on "is debug mode on at all".
var Active = sqlcdebug.Any()

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
