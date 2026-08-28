package ast

type Aggref struct {
	Tag NodeTag[Aggref] `json:"tag"`

	Xpr           Node     `json:"xpr,omitempty"`
	Aggfnoid      Oid      `json:"aggfnoid"`
	Aggtype       Oid      `json:"aggtype"`
	Aggcollid     Oid      `json:"aggcollid"`
	Inputcollid   Oid      `json:"inputcollid"`
	Aggargtypes   *List    `json:"aggargtypes,omitempty"`
	Aggdirectargs *List    `json:"aggdirectargs,omitempty"`
	Args          *List    `json:"args,omitempty"`
	Aggorder      *List    `json:"aggorder,omitempty"`
	Aggdistinct   *List    `json:"aggdistinct,omitempty"`
	Aggfilter     Node     `json:"aggfilter,omitempty"`
	Aggstar       bool     `json:"aggstar"`
	Aggvariadic   bool     `json:"aggvariadic"`
	Aggkind       byte     `json:"aggkind"`
	Agglevelsup   Index    `json:"agglevelsup"`
	Aggsplit      AggSplit `json:"aggsplit"`
	Location      int      `json:"location"`
}

func (n *Aggref) Pos() int {
	return n.Location
}
