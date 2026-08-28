package ast

type FieldStore struct {
	Tag NodeTag[FieldStore] `json:"tag"`

	Xpr        Node  `json:",omitempty"`
	Arg        Node  `json:",omitempty"`
	Newvals    *List `json:",omitempty"`
	Fieldnums  *List `json:",omitempty"`
	Resulttype Oid
}

func (n *FieldStore) Pos() int {
	return 0
}
