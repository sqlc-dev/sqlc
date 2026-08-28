package ast

type FetchStmt struct {
	Tag NodeTag[FetchStmt] `json:"tag"`

	Direction  FetchDirection `json:"direction"`
	HowMany    int64          `json:"how_many"`
	Portalname *string        `json:"portalname,omitempty"`
	Ismove     bool           `json:"ismove"`
}

func (n *FetchStmt) Pos() int {
	return 0
}
