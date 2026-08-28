package ast

type ClosePortalStmt struct {
	Tag NodeTag[ClosePortalStmt] `json:"tag"`

	Portalname *string `json:"portalname,omitempty"`
}

func (n *ClosePortalStmt) Pos() int {
	return 0
}
