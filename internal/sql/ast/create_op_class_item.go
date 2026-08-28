package ast

type CreateOpClassItem struct {
	Tag NodeTag[CreateOpClassItem] `json:"tag"`

	Itemtype    int
	Name        *ObjectWithArgs `json:",omitempty"`
	Number      int
	OrderFamily *List     `json:",omitempty"`
	ClassArgs   *List     `json:",omitempty"`
	Storedtype  *TypeName `json:",omitempty"`
}

func (n *CreateOpClassItem) Pos() int {
	return 0
}
