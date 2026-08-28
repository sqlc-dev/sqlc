package ast

type CreateOpClassItem struct {
	Tag NodeTag[CreateOpClassItem] `json:"tag"`

	Itemtype    int             `json:"itemtype"`
	Name        *ObjectWithArgs `json:"name,omitempty"`
	Number      int             `json:"number"`
	OrderFamily *List           `json:"order_family,omitempty"`
	ClassArgs   *List           `json:"class_args,omitempty"`
	Storedtype  *TypeName       `json:"storedtype,omitempty"`
}

func (n *CreateOpClassItem) Pos() int {
	return 0
}
