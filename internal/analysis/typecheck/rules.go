package typecheck

import "github.com/sqlc-dev/sqlc/internal/analysis/scope"

// DefaultOperatorRules provides baseline type inference rules shared
// across PostgreSQL and MySQL. Engine-specific rules can override this.
type DefaultOperatorRules struct{}

func (r *DefaultOperatorRules) BinaryOp(op string, left, right scope.Type) scope.Type {
	switch op {
	// Comparison operators always produce boolean
	case "=", "<>", "!=", "<", ">", "<=", ">=",
		"IS", "IS NOT", "LIKE", "ILIKE", "NOT LIKE",
		"SIMILAR TO", "~", "~*", "!~", "!~*",
		"REGEXP", "NOT REGEXP":
		return scope.Type{Name: "boolean", NotNull: true}

	// Logical operators produce boolean
	case "AND", "OR":
		return scope.Type{Name: "boolean", NotNull: left.NotNull && right.NotNull}

	// Concatenation produces text
	case "||":
		return scope.Type{Name: "text", NotNull: left.NotNull && right.NotNull}

	// Arithmetic operators
	case "+", "-", "*":
		return r.arithmeticResult(left, right)
	case "/":
		result := r.arithmeticResult(left, right)
		result.NotNull = left.NotNull && right.NotNull
		return result
	case "%":
		return scope.Type{Name: "integer", NotNull: left.NotNull && right.NotNull}

	// Array operators
	case "@>", "<@", "&&":
		return scope.Type{Name: "boolean", NotNull: true}

	// JSON operators
	case "->":
		return scope.Type{Name: "jsonb", NotNull: false}
	case "->>":
		return scope.Type{Name: "text", NotNull: false}

	default:
		// For unknown operators, try to propagate the left type
		if !left.IsUnknown() {
			return left
		}
		return right
	}
}

func (r *DefaultOperatorRules) UnaryOp(op string, operand scope.Type) scope.Type {
	switch op {
	case "NOT":
		return scope.Type{Name: "boolean", NotNull: operand.NotNull}
	case "-", "+":
		return operand
	case "~": // bitwise not
		return operand
	default:
		return operand
	}
}

func (r *DefaultOperatorRules) CanCast(from, to scope.Type) bool {
	if from.IsUnknown() || to.IsUnknown() {
		return true
	}
	if from.Equals(to) {
		return true
	}

	// Numeric types are generally castable to each other
	numericTypes := map[string]bool{
		"integer": true, "int": true, "int4": true,
		"bigint": true, "int8": true,
		"smallint": true, "int2": true,
		"numeric": true, "decimal": true,
		"real": true, "float4": true,
		"float": true, "double precision": true, "float8": true,
	}
	if numericTypes[from.Name] && numericTypes[to.Name] {
		return true
	}

	// Text types are generally castable to each other
	textTypes := map[string]bool{
		"text": true, "varchar": true, "char": true,
		"character varying": true, "character": true,
		"name": true, "citext": true,
	}
	if textTypes[from.Name] && textTypes[to.Name] {
		return true
	}

	// Most things can be cast to/from text
	if from.Name == "text" || to.Name == "text" {
		return true
	}

	return true // Be permissive by default
}

func (r *DefaultOperatorRules) arithmeticResult(left, right scope.Type) scope.Type {
	// If either is unknown, use the other
	if left.IsUnknown() {
		return right
	}
	if right.IsUnknown() {
		return left
	}

	// Type promotion hierarchy
	hierarchy := map[string]int{
		"smallint": 1, "int2": 1,
		"integer": 2, "int": 2, "int4": 2,
		"bigint": 3, "int8": 3,
		"numeric": 4, "decimal": 4,
		"real": 5, "float4": 5,
		"float": 6, "double precision": 6, "float8": 6,
	}

	lp := hierarchy[left.Name]
	rp := hierarchy[right.Name]

	if lp == 0 && rp == 0 {
		return scope.Type{Name: "numeric", NotNull: left.NotNull && right.NotNull}
	}
	if lp >= rp {
		return scope.Type{Name: left.Name, NotNull: left.NotNull && right.NotNull}
	}
	return scope.Type{Name: right.Name, NotNull: left.NotNull && right.NotNull}
}

// PostgreSQLOperatorRules extends default rules with PostgreSQL-specific behavior.
type PostgreSQLOperatorRules struct {
	DefaultOperatorRules
}

func (r *PostgreSQLOperatorRules) BinaryOp(op string, left, right scope.Type) scope.Type {
	switch op {
	case "||":
		// In PostgreSQL, || is string concatenation
		return scope.Type{Name: "text", NotNull: left.NotNull && right.NotNull}
	default:
		return r.DefaultOperatorRules.BinaryOp(op, left, right)
	}
}

// MySQLOperatorRules extends default rules with MySQL-specific behavior.
type MySQLOperatorRules struct {
	DefaultOperatorRules
}

func (r *MySQLOperatorRules) BinaryOp(op string, left, right scope.Type) scope.Type {
	switch op {
	case "||":
		// In MySQL (with default settings), || is logical OR, not concat
		return scope.Type{Name: "boolean", NotNull: left.NotNull && right.NotNull}
	default:
		return r.DefaultOperatorRules.BinaryOp(op, left, right)
	}
}
