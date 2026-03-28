package typecheck

import "github.com/sqlc-dev/sqlc/internal/analysis/scope"

// Expr is the interface for type-checkable SQL expressions.
// This is a simplified representation of the SQL AST focused on type information.
type Expr interface {
	exprNode()
	ExprLocation() int
}

// LiteralExpr represents a constant value (string, int, float, bool, null).
type LiteralExpr struct {
	Type     scope.Type
	Location int
}

func (*LiteralExpr) exprNode()            {}
func (e *LiteralExpr) ExprLocation() int { return e.Location }

// ColumnRefExpr represents a resolved column reference.
type ColumnRefExpr struct {
	Parts        []string   // The original parts (e.g., ["u", "name"])
	ResolvedType scope.Type // The type from the catalog after name resolution
	Location     int
}

func (*ColumnRefExpr) exprNode()             {}
func (e *ColumnRefExpr) ExprLocation() int  { return e.Location }

// ParamExpr represents a query parameter ($1, $2, ?, etc.).
type ParamExpr struct {
	Number   int // 1-indexed parameter number
	Location int
}

func (*ParamExpr) exprNode()            {}
func (e *ParamExpr) ExprLocation() int { return e.Location }

// BinaryOpExpr represents a binary operation (a + b, a = b, a AND b, etc.).
type BinaryOpExpr struct {
	Op       string
	Left     Expr
	Right    Expr
	Location int
}

func (*BinaryOpExpr) exprNode()             {}
func (e *BinaryOpExpr) ExprLocation() int  { return e.Location }

// FuncCallExpr represents a function call.
type FuncCallExpr struct {
	Name       string
	Args       []Expr
	ReturnType scope.Type
	Location   int
}

func (*FuncCallExpr) exprNode()             {}
func (e *FuncCallExpr) ExprLocation() int  { return e.Location }

// TypeCastExpr represents an explicit type cast (e.g., $1::integer).
type TypeCastExpr struct {
	Arg      Expr
	CastType scope.Type
	Location int
}

func (*TypeCastExpr) exprNode()             {}
func (e *TypeCastExpr) ExprLocation() int  { return e.Location }

// SubqueryExpr represents a scalar subquery or EXISTS subquery.
type SubqueryExpr struct {
	Columns  []scope.Type // Column types of the subquery result
	IsExists bool
	Location int
}

func (*SubqueryExpr) exprNode()             {}
func (e *SubqueryExpr) ExprLocation() int  { return e.Location }

// BoolExpr represents AND, OR, NOT operations.
type BoolExpr struct {
	Op       string // "AND", "OR", "NOT"
	Args     []Expr
	Location int
}

func (*BoolExpr) exprNode()            {}
func (e *BoolExpr) ExprLocation() int { return e.Location }

// NullTestExpr represents IS NULL / IS NOT NULL.
type NullTestExpr struct {
	Arg    Expr
	IsNot  bool
	Location int
}

func (*NullTestExpr) exprNode()             {}
func (e *NullTestExpr) ExprLocation() int  { return e.Location }

// CaseExpr represents a CASE WHEN ... THEN ... ELSE ... END expression.
type CaseExpr struct {
	ResultType scope.Type
	Location   int
}

func (*CaseExpr) exprNode()            {}
func (e *CaseExpr) ExprLocation() int { return e.Location }

// CoalesceExpr represents COALESCE(a, b, c, ...).
type CoalesceExpr struct {
	Args     []Expr
	Location int
}

func (*CoalesceExpr) exprNode()             {}
func (e *CoalesceExpr) ExprLocation() int  { return e.Location }

// InExpr represents `expr IN (values...)` or `expr IN (subquery)`.
type InExpr struct {
	Expr     Expr
	Values   []Expr
	Location int
}

func (*InExpr) exprNode()            {}
func (e *InExpr) ExprLocation() int { return e.Location }

// BetweenExpr represents `expr BETWEEN low AND high`.
type BetweenExpr struct {
	Expr     Expr
	Low      Expr
	High     Expr
	Location int
}

func (*BetweenExpr) exprNode()             {}
func (e *BetweenExpr) ExprLocation() int  { return e.Location }
