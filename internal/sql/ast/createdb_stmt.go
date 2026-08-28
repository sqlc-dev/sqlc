package ast

type CreatedbStmt struct {
	Tag NodeTag[CreatedbStmt] `json:"tag"`

	Dbname  *string `json:",omitempty"`
	Options *List   `json:",omitempty"`
}

func (n *CreatedbStmt) Pos() int {
	return 0
}
