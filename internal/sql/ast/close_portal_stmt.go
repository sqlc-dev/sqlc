package ast

type ClosePortalStmt struct {
	Tag NodeTag[ClosePortalStmt] `json:"tag"`

	Portalname *string
}

func (n *ClosePortalStmt) Pos() int {
	return 0
}
