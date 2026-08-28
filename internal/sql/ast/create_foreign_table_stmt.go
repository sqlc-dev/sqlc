package ast

type CreateForeignTableStmt struct {
	Tag NodeTag[CreateForeignTableStmt] `json:"tag"`

	Base       *CreateStmt `json:"base,omitempty"`
	Servername *string     `json:"servername,omitempty"`
	Options    *List       `json:"options,omitempty"`
}

func (n *CreateForeignTableStmt) Pos() int {
	return 0
}
