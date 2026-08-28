package ast

type FieldStore struct {
	Tag NodeTag[FieldStore] `json:"tag"`

	Xpr        Node  `json:"xpr,omitempty"`
	Arg        Node  `json:"arg,omitempty"`
	Newvals    *List `json:"newvals,omitempty"`
	Fieldnums  *List `json:"fieldnums,omitempty"`
	Resulttype Oid   `json:"resulttype"`
}

func (n *FieldStore) Pos() int {
	return 0
}
