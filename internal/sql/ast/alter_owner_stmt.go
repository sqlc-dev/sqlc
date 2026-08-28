package ast

type AlterOwnerStmt struct {
	Tag NodeTag[AlterOwnerStmt] `json:"tag"`

	ObjectType ObjectType `json:"object_type"`
	Relation   *RangeVar  `json:"relation,omitempty"`
	Object     Node       `json:"object,omitempty"`
	Newowner   *RoleSpec  `json:"newowner,omitempty"`
}

func (n *AlterOwnerStmt) Pos() int {
	return 0
}
