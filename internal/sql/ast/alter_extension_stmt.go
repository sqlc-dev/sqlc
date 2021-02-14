package ast

type AlterExtensionStmt struct {
	Extname *string
	Options *List
}

func (n *AlterExtensionStmt) Pos() int {
	return 0
}
