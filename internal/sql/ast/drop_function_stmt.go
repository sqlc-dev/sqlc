package ast

type DropFunctionStmt struct {
	Tag NodeTag[DropFunctionStmt] `json:"tag"`

	Funcs     []*FuncSpec `json:"funcs,omitempty"`
	MissingOk bool        `json:"missing_ok"`
}

func (n *DropFunctionStmt) Pos() int {
	return 0
}
