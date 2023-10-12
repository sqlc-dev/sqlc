package ast

import "fmt"

type Boolean struct {
	Boolval bool
}

func (n *Boolean) Pos() int {
	return 0
}

func (n *Boolean) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	if n.Boolval {
		fmt.Fprintf(buf, "true")
	} else {
		fmt.Fprintf(buf, "false")
	}
}
