package ast

type varatt_external struct {
	VaRawsize    int32
	VaExtsize    int32
	VaValueid    Oid
	VaToastrelid Oid
}

func (n *varatt_external) Pos() int {
	return 0
}
