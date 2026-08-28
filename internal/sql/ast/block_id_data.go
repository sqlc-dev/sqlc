package ast

type BlockIdData struct {
	Tag NodeTag[BlockIdData] `json:"tag"`

	BiHi uint16 `json:"bi_hi"`
	BiLo uint16 `json:"bi_lo"`
}

func (n *BlockIdData) Pos() int {
	return 0
}
