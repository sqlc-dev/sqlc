package ast

type AlterTypeAddValueStmt struct {
	Tag NodeTag[AlterTypeAddValueStmt] `json:"tag"`

	Type               *TypeName `json:"type,omitempty"`
	NewValue           *string   `json:"new_value,omitempty"`
	NewValHasNeighbor  bool      `json:"new_val_has_neighbor"`
	NewValNeighbor     *string   `json:"new_val_neighbor,omitempty"`
	NewValIsAfter      bool      `json:"new_val_is_after"`
	SkipIfNewValExists bool      `json:"skip_if_new_val_exists"`
}

func (n *AlterTypeAddValueStmt) Pos() int {
	return 0
}
