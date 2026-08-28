package ast

type CopyStmt struct {
	Tag NodeTag[CopyStmt] `json:"tag"`

	Relation  *RangeVar `json:"relation,omitempty"`
	Query     Node      `json:"query,omitempty"`
	Attlist   *List     `json:"attlist,omitempty"`
	IsFrom    bool      `json:"is_from"`
	IsProgram bool      `json:"is_program"`
	Filename  *string   `json:"filename,omitempty"`
	Options   *List     `json:"options,omitempty"`
}

func (n *CopyStmt) Pos() int {
	return 0
}
