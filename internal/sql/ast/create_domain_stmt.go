package ast

type CreateDomainStmt struct {
	Tag NodeTag[CreateDomainStmt] `json:"tag"`

	Domainname  *List
	TypeName    *TypeName
	CollClause  *CollateClause
	Constraints *List
}

func (n *CreateDomainStmt) Pos() int {
	return 0
}
