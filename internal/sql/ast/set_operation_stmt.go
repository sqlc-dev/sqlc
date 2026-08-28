package ast

type SetOperationStmt struct {
	Tag NodeTag[SetOperationStmt] `json:"tag"`

	Op            SetOperation `json:"op"`
	All           bool         `json:"all"`
	Larg          Node         `json:"larg,omitempty"`
	Rarg          Node         `json:"rarg,omitempty"`
	ColTypes      *List        `json:"col_types,omitempty"`
	ColTypmods    *List        `json:"col_typmods,omitempty"`
	ColCollations *List        `json:"col_collations,omitempty"`
	GroupClauses  *List        `json:"group_clauses,omitempty"`
}

func (n *SetOperationStmt) Pos() int {
	return 0
}
