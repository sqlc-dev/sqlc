package debug

import (
	"os"

	"github.com/davecgh/go-spew/spew"
)

func Active() bool {
	return os.Getenv("SQLCDEBUG") != ""
}

func Dump(n interface{}) {
	spew.Dump(n)
}
