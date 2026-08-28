package ast

type CreateOpClassItem struct {
	Tag NodeTag[CreateOpClassItem] `json:"tag"`

	Itemtype    int
	Name        *ObjectWithArgs
	Number      int
	OrderFamily *List
	ClassArgs   *List
	Storedtype  *TypeName
}

func (n *CreateOpClassItem) Pos() int {
	return 0
}
