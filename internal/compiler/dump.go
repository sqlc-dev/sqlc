package compiler

import "github.com/davecgh/go-spew/spew"

func dump(a ...interface{}) {
	spew.Dump(a...)
}
