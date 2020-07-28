package debug

import (
	"os"

	"github.com/davecgh/go-spew/spew"
)

var Active bool

func init() {
	Active = os.Getenv("SQLCDEBUG") != ""
}

func Dump(n ...interface{}) {
	if Active {
		spew.Dump(n)
	}
}
