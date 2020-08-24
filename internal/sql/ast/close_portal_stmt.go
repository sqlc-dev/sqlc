package pg

type ClosePortalStmt struct {
	Portalname *string
}

func (n *ClosePortalStmt) Pos() int {
	return 0
}
