package ast

type CreateDomainStmt struct {
	Tag NodeTag[CreateDomainStmt] `json:"tag"`

	Domainname  *List          `json:",omitempty"`
	TypeName    *TypeName      `json:",omitempty"`
	CollClause  *CollateClause `json:",omitempty"`
	Constraints *List          `json:",omitempty"`
}

func (n *CreateDomainStmt) Pos() int {
	return 0
}
