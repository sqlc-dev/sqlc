package ast

// In describes a 'select foo in (bar, baz)' type statement, though there are multiple important variants handled.
type In struct {
	// Expr is the value expression to be compared.
	Expr Node
	// List is the list expression in compare list.
	List []Node
	// Not is true, the expression is "not in".
	Not bool
	// Sel is the subquery, may be rewritten to other type of expression.
	Sel      Node
	Location int
}

// Pos returns the location.
func (n *In) Pos() int {
	return n.Location
}
