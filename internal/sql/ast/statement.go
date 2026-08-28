package ast

type Statement struct {
	Tag NodeTag[Statement] `json:"tag"`

	Raw *RawStmt `json:"raw,omitempty"`
}

func (n *Statement) Pos() int {
	return 0
}
