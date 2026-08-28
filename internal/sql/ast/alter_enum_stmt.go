package ast

type AlterEnumStmt struct {
	Tag NodeTag[AlterEnumStmt] `json:"tag"`

	TypeName           *List   `json:"type_name,omitempty"`
	OldVal             *string `json:"old_val,omitempty"`
	NewVal             *string `json:"new_val,omitempty"`
	NewValNeighbor     *string `json:"new_val_neighbor,omitempty"`
	NewValIsAfter      bool    `json:"new_val_is_after"`
	SkipIfNewValExists bool    `json:"skip_if_new_val_exists"`
}

func (n *AlterEnumStmt) Pos() int {
	return 0
}
