package ast

type AccessPriv struct {
	Tag NodeTag[AccessPriv] `json:"tag"`

	PrivName *string
	Cols     *List
}

func (n *AccessPriv) Pos() int {
	return 0
}
