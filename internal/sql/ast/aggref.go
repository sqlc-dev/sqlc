package ast

type Aggref struct {
	Xpr           Node
	Aggfnoid      Oid
	Aggtype       Oid
	Aggcollid     Oid
	Inputcollid   Oid
	Aggargtypes   *List
	Aggdirectargs *List
	Args          *List
	Aggorder      *List
	Aggdistinct   *List
	Aggfilter     Node
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
