package ast

type CreateDomainStmt struct {
	Tag NodeTag[CreateDomainStmt] `json:"tag"`

	Domainname  *List          `json:"domainname,omitempty"`
	TypeName    *TypeName      `json:"type_name,omitempty"`
	CollClause  *CollateClause `json:"coll_clause,omitempty"`
	Constraints *List          `json:"constraints,omitempty"`
}

func (n *CreateDomainStmt) Pos() int {
	return 0
}
