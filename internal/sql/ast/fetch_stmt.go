package ast

type FetchStmt struct {
	Tag NodeTag[FetchStmt] `json:"tag"`

	Direction  FetchDirection
	HowMany    int64
	Portalname *string `json:",omitempty"`
	Ismove     bool
}

func (n *FetchStmt) Pos() int {
	return 0
}
