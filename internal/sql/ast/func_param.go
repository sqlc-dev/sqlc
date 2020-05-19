package ast

type FuncParam struct {
	Name    *string
	Type    *TypeName
	DefExpr Node // Will always be &ast.TODO
}

func (n *FuncParam) Pos() int {
	return 0
}
