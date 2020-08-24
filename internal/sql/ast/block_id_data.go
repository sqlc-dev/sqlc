package ast

type BlockIdData struct {
	BiHi uint16
	BiLo uint16
}

func (n *BlockIdData) Pos() int {
	return 0
}
