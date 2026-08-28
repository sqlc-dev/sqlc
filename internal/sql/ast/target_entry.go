package ast

type TargetEntry struct {
	Tag NodeTag[TargetEntry] `json:"tag"`

	Xpr             Node       `json:"xpr,omitempty"`
	Expr            Node       `json:"expr,omitempty"`
	Resno           AttrNumber `json:"resno"`
	Resname         *string    `json:"resname,omitempty"`
	Ressortgroupref Index      `json:"ressortgroupref"`
	Resorigtbl      Oid        `json:"resorigtbl"`
	Resorigcol      AttrNumber `json:"resorigcol"`
	Resjunk         bool       `json:"resjunk"`
}

func (n *TargetEntry) Pos() int {
	return 0
}
