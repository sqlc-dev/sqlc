package ast

type BlockIdData struct {
	Tag NodeTag[BlockIdData] `json:"tag"`

	BiHi uint16
	BiLo uint16
}

func (n *BlockIdData) Pos() int {
	return 0
}
