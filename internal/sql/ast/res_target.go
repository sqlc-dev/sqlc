package sqlc

type ResTarget struct {
	Val Node
}

func (n *ResTarget) Pos() int {
	return 0
}


