package ast

import ()

type AccessPriv struct {
	PrivName *string
	Cols     *List
}

func (n *AccessPriv) Pos() int {
	return 0
}
