package ast

type AccessPriv struct {
	Tag NodeTag[AccessPriv] `json:"tag"`

	PrivName *string `json:",omitempty"`
	Cols     *List   `json:",omitempty"`
}

func (n *AccessPriv) Pos() int {
	return 0
}
