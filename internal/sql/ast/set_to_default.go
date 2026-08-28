package ast

type SetToDefault struct {
	Tag NodeTag[SetToDefault] `json:"tag"`

	Xpr       Node  `json:"xpr,omitempty"`
	TypeId    Oid   `json:"type_id"`
	TypeMod   int32 `json:"type_mod"`
	Collation Oid   `json:"collation"`
	Location  int   `json:"location"`
}

func (n *SetToDefault) Pos() int {
	return n.Location
}
