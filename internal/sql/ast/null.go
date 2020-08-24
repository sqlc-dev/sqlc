package pg

type Null struct {
}

func (n *Null) Pos() int {
	return 0
}
