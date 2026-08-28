package ast

type DropFunctionStmt struct {
	Tag NodeTag[DropFunctionStmt] `json:"tag"`

	Funcs     []*FuncSpec `json:",omitempty"`
	MissingOk bool
}

func (n *DropFunctionStmt) Pos() int {
	return 0
}
