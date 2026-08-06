package compiler

import (
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

//
// ==============================
// Internal Type System
// ==============================
//

type Kind int

const (
	KindUnknown Kind = iota // inference not supported
	KindInt
	KindFloat
	KindDecimal
	KindAny
)

type Type struct {
	Kind    Kind
	NotNull bool
	Valid   bool // explicit signal: inference succeeded
}

func unknownType() Type {
	return Type{Kind: KindUnknown, Valid: false}
}

//
// ==============================
// Entry Point
// ==============================
//

func (c *Compiler) inferExprType(node ast.Node, tables []*Table) *Column {
	if node == nil {
		return nil
	}

	switch c.conf.Engine {
	case config.EngineMySQL:
		t := c.inferMySQLExpr(node, tables)
		return c.mysqlTypeToColumn(t)

	// case config.EnginePostgreSQL:
	//     t := c.inferPostgresExpr(node, tables)
	//     return c.postgresTypeToColumn(t)

	default:
		return nil
	}
}

//
// ==============================
// MySQL Inference
// ==============================
//

func (c *Compiler) inferMySQLExpr(node ast.Node, tables []*Table) Type {
	switch n := node.(type) {
	case *ast.ColumnRef:
		return c.inferMySQLColumnRef(n, tables)

	case *ast.A_Const:
		return inferConst(n)

	case *ast.TypeCast:
		return c.inferMySQLTypeCast(n, tables)

	case *ast.A_Expr:
		return c.inferMySQLBinary(n, tables)

	default:
		return unknownType()
	}
}

//
// ------------------------------
// Leaf nodes
// ------------------------------
//

func (c *Compiler) inferMySQLColumnRef(ref *ast.ColumnRef, tables []*Table) Type {
	cols, err := outputColumnRefs(&ast.ResTarget{}, tables, ref)
	if err != nil || len(cols) == 0 {
		return unknownType()
	}

	col := cols[0]

	return Type{
		Kind:    mapMySQLKind(col.DataType),
		NotNull: col.NotNull,
		Valid:   true,
	}
}

func inferConst(node *ast.A_Const) Type {
	if node == nil || node.Val == nil {
		return unknownType()
	}

	switch node.Val.(type) {
	case *ast.Integer:
		return Type{Kind: KindInt, NotNull: true, Valid: true}

	case *ast.Float:
		return Type{Kind: KindFloat, NotNull: true, Valid: true}

	case *ast.Null:
		return Type{Kind: KindAny, NotNull: false, Valid: true}

	default:
		return unknownType()
	}
}

func (c *Compiler) inferMySQLTypeCast(node *ast.TypeCast, tables []*Table) Type {
	if node == nil || node.TypeName == nil {
		return unknownType()
	}

	// MySQL populates TypeName.Name directly; toColumn reads TypeName.Names (Postgres-style).
	kind := mapMySQLKind(node.TypeName.Name)
	if kind == KindUnknown {
		return unknownType()
	}

	arg := c.inferMySQLExpr(node.Arg, tables)

	t := Type{
		Kind:  kind,
		Valid: true,
	}

	// propagate nullability
	if arg.Valid {
		t.NotNull = arg.NotNull
	}

	// explicit NULL literal
	if constant, ok := node.Arg.(*ast.A_Const); ok {
		if _, isNull := constant.Val.(*ast.Null); isNull {
			t.NotNull = false
		}
	}

	return t
}

//
// ------------------------------
// Binary expressions
// ------------------------------
//

func (c *Compiler) inferMySQLBinary(node *ast.A_Expr, tables []*Table) Type {
	op := joinOperator(node)

	left := c.inferMySQLExpr(node.Lexpr, tables)
	right := c.inferMySQLExpr(node.Rexpr, tables)

	if !left.Valid || !right.Valid {
		return unknownType()
	}

	// NOTE: only normal division ("/") is supported for now.
	// Unsupported operators intentionally fall back to the existing behavior.
	return promoteMySQLNumeric(op, left, right)
}

//
// ==============================
// Promotion Rules (MySQL-specific for now)
// ==============================
//

// promoteMySQLNumeric applies simplified numeric promotion rules for MySQL.
// It currently only supports "/" and intentionally falls back for other operators.
func promoteMySQLNumeric(op string, a, b Type) Type {
	notNull := a.NotNull && b.NotNull

	switch op {
	case "/":
		if a.Kind == KindFloat || b.Kind == KindFloat {
			return Type{
				Kind:    KindFloat,
				NotNull: notNull,
				Valid:   true,
			}
		}

		return Type{
			Kind:    KindDecimal,
			NotNull: notNull,
			Valid:   true,
		}
	}

	return unknownType()
}

//
// ==============================
// Engine-specific Mapping
// ==============================
//

func (c *Compiler) mysqlTypeToColumn(t Type) *Column {
	if !t.Valid {
		return nil
	}

	col := &Column{
		NotNull: t.NotNull,
	}

	switch t.Kind {
	case KindInt:
		col.DataType = "int"

	case KindFloat:
		col.DataType = "float"

	case KindDecimal:
		col.DataType = "decimal"

	default:
		col.DataType = "any"
	}

	return col
}

func mapMySQLKind(dt string) Kind {
	switch dt {
	case "int", "integer", "bigint", "smallint":
		return KindInt

	case "float", "double", "real":
		return KindFloat

	case "decimal", "numeric":
		return KindDecimal

	default:
		return KindUnknown
	}
}

//
// ==============================
// AST helpers
// ==============================
//

func joinOperator(node *ast.A_Expr) string {
	if node == nil || node.Name == nil || len(node.Name.Items) == 0 {
		return ""
	}

	if s, ok := node.Name.Items[0].(*ast.String); ok {
		return s.Str
	}

	return ""
}
