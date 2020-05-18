package sqlc

type DropFunctionStmt struct {
	Funcs     []*FuncSpec
	MissingOk bool
}

func (n *DropFunctionStmt) Pos() int {
	return 0
}
