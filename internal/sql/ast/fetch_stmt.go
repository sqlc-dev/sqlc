package ast

type FetchStmt struct {
	Direction  FetchDirection
	HowMany    int64
	Portalname *string
	Ismove     bool
}

func (n *FetchStmt) Pos() int {
	return 0
}
