package ast

type A_Expr_Kind uint

const (
	A_Expr_Kind_OP              A_Expr_Kind = 1
	A_Expr_Kind_OP_ANY          A_Expr_Kind = 2
	A_Expr_Kind_OP_ALL          A_Expr_Kind = 3
	A_Expr_Kind_DISTINCT        A_Expr_Kind = 4
	A_Expr_Kind_NOT_DISTINCT    A_Expr_Kind = 5
	A_Expr_Kind_NULLIF          A_Expr_Kind = 6
	A_Expr_Kind_IN              A_Expr_Kind = 7
	A_Expr_Kind_LIKE            A_Expr_Kind = 8
	A_Expr_Kind_ILIKE           A_Expr_Kind = 9
	A_Expr_Kind_SIMILAR         A_Expr_Kind = 10
	A_Expr_Kind_BETWEEN         A_Expr_Kind = 11
	A_Expr_Kind_NOT_BETWEEN     A_Expr_Kind = 12
	A_Expr_Kind_BETWEEN_SYM     A_Expr_Kind = 13
	A_Expr_Kind_NOT_BETWEEN_SYM A_Expr_Kind = 14
)

func (n *A_Expr_Kind) Pos() int {
	return 0
}
