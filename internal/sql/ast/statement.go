package ast

type Statement struct {
	Tag NodeTag[Statement] `json:"tag"`

	Raw *RawStmt `json:",omitempty"`
}

func (n *Statement) Pos() int {
	return 0
}
