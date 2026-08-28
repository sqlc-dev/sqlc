package ast

type CopyStmt struct {
	Tag NodeTag[CopyStmt] `json:"tag"`

	Relation  *RangeVar `json:",omitempty"`
	Query     Node      `json:",omitempty"`
	Attlist   *List     `json:",omitempty"`
	IsFrom    bool
	IsProgram bool
	Filename  *string `json:",omitempty"`
	Options   *List   `json:",omitempty"`
}

func (n *CopyStmt) Pos() int {
	return 0
}
