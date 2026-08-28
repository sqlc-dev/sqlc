package ast

type AccessPriv struct {
	Tag NodeTag[AccessPriv] `json:"tag"`

	PrivName *string `json:"priv_name,omitempty"`
	Cols     *List   `json:"cols,omitempty"`
}

func (n *AccessPriv) Pos() int {
	return 0
}
