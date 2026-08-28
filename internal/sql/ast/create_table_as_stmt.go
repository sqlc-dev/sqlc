package ast

type CreateTableAsStmt struct {
	Tag NodeTag[CreateTableAsStmt] `json:"tag"`

	Query        Node        `json:"query,omitempty"`
	Into         *IntoClause `json:"into,omitempty"`
	Relkind      ObjectType  `json:"relkind"`
	IsSelectInto bool        `json:"is_select_into"`
	IfNotExists  bool        `json:"if_not_exists"`
}

func (n *CreateTableAsStmt) Pos() int {
	return 0
}
