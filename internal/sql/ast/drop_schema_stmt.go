package ast

type DropSchemaStmt struct {
	Tag NodeTag[DropSchemaStmt] `json:"tag"`

	Schemas   []*String `json:"schemas,omitempty"`
	MissingOk bool      `json:"missing_ok"`
}

func (n *DropSchemaStmt) Pos() int {
	return 0
}
