package ast

type DropSchemaStmt struct {
	Tag NodeTag[DropSchemaStmt] `json:"tag"`

	Schemas   []*String `json:",omitempty"`
	MissingOk bool
}

func (n *DropSchemaStmt) Pos() int {
	return 0
}
