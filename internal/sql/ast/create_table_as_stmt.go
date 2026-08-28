package ast

type CreateTableAsStmt struct {
	Tag NodeTag[CreateTableAsStmt] `json:"tag"`

	Query        Node        `json:",omitempty"`
	Into         *IntoClause `json:",omitempty"`
	Relkind      ObjectType
	IsSelectInto bool
	IfNotExists  bool
}

func (n *CreateTableAsStmt) Pos() int {
	return 0
}
