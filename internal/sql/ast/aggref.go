package ast

type Aggref struct {
	Tag NodeTag[Aggref] `json:"tag"`

	Xpr           Node `json:",omitempty"`
	Aggfnoid      Oid
	Aggtype       Oid
	Aggcollid     Oid
	Inputcollid   Oid
	Aggargtypes   *List `json:",omitempty"`
	Aggdirectargs *List `json:",omitempty"`
	Args          *List `json:",omitempty"`
	Aggorder      *List `json:",omitempty"`
	Aggdistinct   *List `json:",omitempty"`
	Aggfilter     Node  `json:",omitempty"`
	Aggstar       bool
	Aggvariadic   bool
	Aggkind       byte
	Agglevelsup   Index
	Aggsplit      AggSplit
	Location      int
}

func (n *Aggref) Pos() int {
	return n.Location
}
