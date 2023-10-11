package ast

type A_Expr_Kind uint

const (
	A_Expr_Kind_IN A_Expr_Kind = 7
)

func (n *A_Expr_Kind) Pos() int {
	return 0
}
