package ast

type CreateTableAsStmt struct {
	Tag NodeTag[CreateTableAsStmt] `json:"tag"`

	Query        Node
	Into         *IntoClause
	Relkind      ObjectType
	IsSelectInto bool
	IfNotExists  bool
}

func (n *CreateTableAsStmt) Pos() int {
	return 0
}
