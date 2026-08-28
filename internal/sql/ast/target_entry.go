package ast

type TargetEntry struct {
	Tag NodeTag[TargetEntry] `json:"tag"`

	Xpr             Node `json:",omitempty"`
	Expr            Node `json:",omitempty"`
	Resno           AttrNumber
	Resname         *string `json:",omitempty"`
	Ressortgroupref Index
	Resorigtbl      Oid
	Resorigcol      AttrNumber
	Resjunk         bool
}

func (n *TargetEntry) Pos() int {
	return 0
}
