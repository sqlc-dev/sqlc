package ast

type A_Expr_Kind uint

const (
	A_Expr_Kind_IN   A_Expr_Kind = 7
	A_Expr_Kind_LIKE A_Expr_Kind = 8
)

func (n *A_Expr_Kind) Pos() int {
	return 0
}
