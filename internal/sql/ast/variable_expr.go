package ast

// VariableExpr represents a MySQL user variable (e.g., @user_id)
// This is distinct from sqlc's @param named parameter syntax.
type VariableExpr struct {
	Name     string
	Location int
}

func (n *VariableExpr) Pos() int {
	return n.Location
}

func (n *VariableExpr) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.WriteString("@")
	buf.WriteString(n.Name)
}
