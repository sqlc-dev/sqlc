package ast

type XmlSerialize struct {
	Tag NodeTag[XmlSerialize] `json:"tag"`

	Xmloption XmlOptionType
	Expr      Node
	TypeName  *TypeName
	Location  int
}

func (n *XmlSerialize) Pos() int {
	return n.Location
}
