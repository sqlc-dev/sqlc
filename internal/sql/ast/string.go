package ast

import "fmt"

type String struct {
	Str string
}

func (n *String) Pos() int {
	return 0
}

func (n *String) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	fmt.Fprintf(buf, "%q", n.Str)
}
