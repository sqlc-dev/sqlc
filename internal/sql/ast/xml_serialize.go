package ast

type XmlSerialize struct {
	Tag NodeTag[XmlSerialize] `json:"tag"`

	Xmloption XmlOptionType
	Expr      Node      `json:",omitempty"`
	TypeName  *TypeName `json:",omitempty"`
	Location  int
}

func (n *XmlSerialize) Pos() int {
	return n.Location
}
