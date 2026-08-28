package ast

type CreateForeignTableStmt struct {
	Tag NodeTag[CreateForeignTableStmt] `json:"tag"`

	Base       *CreateStmt `json:",omitempty"`
	Servername *string     `json:",omitempty"`
	Options    *List       `json:",omitempty"`
}

func (n *CreateForeignTableStmt) Pos() int {
	return 0
}
