package ast

type XmlSerialize struct {
	Tag NodeTag[XmlSerialize] `json:"tag"`

	Xmloption XmlOptionType `json:"xmloption"`
	Expr      Node          `json:"expr,omitempty"`
	TypeName  *TypeName     `json:"type_name,omitempty"`
	Location  int           `json:"location"`
}

func (n *XmlSerialize) Pos() int {
	return n.Location
}
