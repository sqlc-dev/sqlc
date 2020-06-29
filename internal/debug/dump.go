package debug

import (
	"github.com/davecgh/go-spew/spew"
)

func Dump(n interface{}) {
	spew.Dump(n)
}
