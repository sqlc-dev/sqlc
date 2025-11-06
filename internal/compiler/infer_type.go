package compiler

import (
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// inferExprType recursively analyzes SQL expressions to determine their types.
//
// It handles:
//   - Column references (resolved from table schema)
//   - Literal constants (intrinsic types)
//   - Binary operations (applies type promotion and operator-specific rules)
//   - Type casts (respects explicit type annotations)
//
// Examples:
//
//	SELECT a1 / 1024              -- infers decimal or float based on operand types
//	SELECT COALESCE(a1 / 1024, 0) -- handles nested expressions recursively
//	SELECT CAST(a1 AS INT)        -- respects explicit casts
//
// Returns nil if type cannot be inferred, allowing fallback to default behavior.
func (c *Compiler) inferExprType(node ast.Node, tables []*Table) *Column {
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	case *ast.ColumnRef:
		// Try to resolve the column reference
		// Create a minimal ResTarget for outputColumnRefs
		emptyRes := &ast.ResTarget{}
		cols, err := outputColumnRefs(emptyRes, tables, n)
		if err != nil || len(cols) == 0 {
			return nil
		}
		return cols[0]

	case *ast.A_Const:
		// Infer type from constant value
		switch n.Val.(type) {
		case *ast.String:
			return &Column{DataType: "text", NotNull: true}
		case *ast.Integer:
			return &Column{DataType: "int", NotNull: true}
		case *ast.Float:
			return &Column{DataType: "float", NotNull: true}
		case *ast.Boolean:
			return &Column{DataType: "bool", NotNull: true}
		default:
			return nil
		}

	case *ast.A_Expr:
		// Recursively infer types of left and right operands
		leftCol := c.inferExprType(n.Lexpr, tables)
		rightCol := c.inferExprType(n.Rexpr, tables)

		if leftCol == nil && rightCol == nil {
			return nil
		}

		// Extract operator name
		op := ""
		if n.Name != nil && len(n.Name.Items) > 0 {
			if str, ok := n.Name.Items[0].(*ast.String); ok {
				op = str.Str
			}
		}

		// Apply database-specific type rules
		return c.combineTypes(leftCol, rightCol, op)

	case *ast.TypeCast:
		// If there's an explicit cast, use that type
		if n.TypeName != nil {
			col := toColumn(n.TypeName)
			// Check if the casted value is nullable
			if constant, ok := n.Arg.(*ast.A_Const); ok {
				if _, isNull := constant.Val.(*ast.Null); isNull {
					col.NotNull = false
				}
			}
			return col
		}
	}

	return nil
}

// combineTypes determines the result type of a binary operation.
// The logic is database-specific and handles operator semantics for each engine.
func (c *Compiler) combineTypes(left, right *Column, operator string) *Column {
	// If either operand is unknown, we can't infer the type
	if left == nil && right == nil {
		return nil
	}

	// If one operand is known, use it as a hint
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}

	// Result is nullable if any operand is nullable (SQL NULL propagation rule)
	notNull := left.NotNull && right.NotNull

	// Apply database-specific type rules
	switch c.conf.Engine {
	case config.EngineMySQL:
		return combineMySQLTypes(left, right, operator, notNull)
	default:
		// TODO: Implement type inference for PostgreSQL and SQLite
		// For now, use conservative fallback rules
		return combineGenericTypes(left, right, operator, notNull)
	}
}

// combineMySQLTypes implements MySQL-specific type inference rules.
//
// Division always returns decimal or float:
//
//	SELECT int_col / 1024        -- returns decimal
//	SELECT float_col / 1024      -- returns float
//	SELECT int_col DIV 1024      -- returns decimal (DIV is MySQL-specific)
//
// Nullability propagates (NULL op anything = NULL):
//
//	nullable_col / 1024  -- returns nullable result
//	NOT NULL col / 1024  -- returns NOT NULL result
func combineMySQLTypes(left, right *Column, operator string, notNull bool) *Column {
	// Handle nil operands
	if left == nil && right == nil {
		return nil
	}
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}

	switch operator {
	case "/", "div":
		// Division: int/int = decimal, float/anything = float
		// Note: "div" is MySQL-specific operator recognized by IsMathematicalOperator()
		if isFloatType(left.DataType) || isFloatType(right.DataType) {
			return &Column{DataType: "float", NotNull: notNull}
		}
		return &Column{DataType: "decimal", NotNull: notNull}

	case "%", "mod":
		// Modulo: returns integer if both are integers, otherwise decimal
		// Note: "mod" is MySQL-specific operator recognized by IsMathematicalOperator()
		if isIntegerType(left.DataType) && isIntegerType(right.DataType) {
			return &Column{DataType: "int", NotNull: notNull}
		}
		return &Column{DataType: "decimal", NotNull: notNull}

	case "&", "|", "<<", ">>", "~", "#", "^":
		// Bitwise operators: always integer
		return &Column{DataType: "int", NotNull: notNull}

	default:
		// Arithmetic: standard type promotion
		return promoteArithmeticTypes(left, right, notNull)
	}
}

// combineGenericTypes provides conservative fallback for unsupported databases.
// Returns nil to avoid incorrect type assumptions for PostgreSQL, SQLite, etc.
// This ensures fallback to the original behavior (interface{} or default types).
func combineGenericTypes(left, right *Column, operator string, notNull bool) *Column {
	// TODO: Implement type inference for PostgreSQL and SQLite
	// For now, return nil to use safe defaults
	return nil
}

// promoteArithmeticTypes applies standard type promotion rules for arithmetic.
// This follows the principle: int -> decimal -> float
func promoteArithmeticTypes(left, right *Column, notNull bool) *Column {
	dataType := "int"

	// Float takes precedence over all other numeric types
	if isFloatType(left.DataType) || isFloatType(right.DataType) {
		dataType = "float"
	} else if isDecimalType(left.DataType) || isDecimalType(right.DataType) {
		// Decimal takes precedence over integer
		dataType = "decimal"
	} else if isIntegerType(left.DataType) && isIntegerType(right.DataType) {
		// Both are integers
		dataType = "int"
	}

	return &Column{
		DataType: dataType,
		NotNull:  notNull,
	}
}

// isFloatType checks if a datatype is a floating-point type
func isFloatType(dataType string) bool {
	switch dataType {
	case "float", "double", "double precision", "real":
		return true
	}
	return false
}

// isDecimalType checks if a datatype is a decimal/numeric type
func isDecimalType(dataType string) bool {
	switch dataType {
	case "decimal", "numeric", "dec", "fixed":
		return true
	}
	return false
}

// isIntegerType checks if a datatype is an integer type
func isIntegerType(dataType string) bool {
	switch dataType {
	case "int", "integer", "smallint", "bigint", "tinyint", "mediumint":
		return true
	}
	return false
}
