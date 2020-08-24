package ast

type RangeVar struct {
	Catalogname    *string
	Schemaname     *string
	Relname        *string
	Inh            bool
	Relpersistence byte
	Alias          *Alias
	Location       int
}

func (n *RangeVar) Pos() int {
	return n.Location
}
