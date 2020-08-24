package pg

type Integer struct {
	Ival int64
}

func (n *Integer) Pos() int {
	return 0
}
