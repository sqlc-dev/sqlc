package ast

type TargetEntry struct {
	Tag NodeTag[TargetEntry] `json:"tag"`

	Xpr             Node
	Expr            Node
	Resno           AttrNumber
	Resname         *string
	Ressortgroupref Index
	Resorigtbl      Oid
	Resorigcol      AttrNumber
	Resjunk         bool
}

func (n *TargetEntry) Pos() int {
	return 0
}
