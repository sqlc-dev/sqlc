package ast

type CreatedbStmt struct {
	Tag NodeTag[CreatedbStmt] `json:"tag"`

	Dbname  *string `json:"dbname,omitempty"`
	Options *List   `json:"options,omitempty"`
}

func (n *CreatedbStmt) Pos() int {
	return 0
}
