package ast

type ParamRef struct {
	Number        int
	Location      int
	Dollar        bool
	IsSqlcDynamic bool
}

func (n *ParamRef) Pos() int {
	return n.Location
}
