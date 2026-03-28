// Package typecheck implements bidirectional type checking for SQL expressions.
//
// Type information flows in two directions:
//   - Synthesis (bottom-up): "what type does this expression have?"
//   - Checking (top-down): "does this expression have the type I expect?"
//
// This is particularly useful for SQL parameter type inference. When we see
// `WHERE age > $1`, the parameter $1 is in checking mode — its expected type
// is inferred from the context (the type of `age`). When we see `SELECT age + 1`,
// the expression is in synthesis mode — we compute the result type from the operands.
//
// Reference: Dunfield & Krishnaswami, "Bidirectional Typing", ACM Computing Surveys, 2021.
package typecheck

import (
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/analysis/scope"
)

// TypeKind classifies types for the checker.
type TypeKind int

const (
	KindScalar   TypeKind = iota // A single value (int, text, bool, ...)
	KindRow                     // A row type (the result of a subquery)
	KindSet                     // A set of rows (table reference)
	KindUnknown                 // Not yet determined
)

// InferredType is the result of type inference on an expression.
type InferredType struct {
	Type   scope.Type
	Kind   TypeKind
	Source TypeSource // Where this type came from
}

// TypeSource records provenance for type inference results.
type TypeSource int

const (
	SourceCatalog    TypeSource = iota // Came from the database catalog
	SourceLiteral                     // Inferred from a literal value
	SourceOperator                    // Inferred from an operator's return type
	SourceContext                     // Inferred from surrounding context (checking mode)
	SourceFunction                    // Inferred from a function's return type
	SourceParameter                   // A parameter whose type was inferred
	SourceUnknown                     // Could not be determined
)

func (s TypeSource) String() string {
	names := [...]string{"catalog", "literal", "operator", "context", "function", "parameter", "unknown"}
	if int(s) < len(names) {
		return names[s]
	}
	return fmt.Sprintf("TypeSource(%d)", int(s))
}

// ParamTypeInference records a type inferred for a query parameter.
type ParamTypeInference struct {
	Number   int        // Parameter number ($1, $2, ...)
	Type     scope.Type // The inferred type
	Source   TypeSource // How the type was inferred
	Location int        // Source position of the parameter
}

// Checker performs bidirectional type checking on SQL expressions.
type Checker struct {
	params    map[int]*ParamTypeInference // Accumulated parameter types
	errors    []TypeError                 // Type errors found
	opRules   OperatorRules               // Engine-specific operator type rules
}

// TypeError describes a type mismatch found during checking.
type TypeError struct {
	Message  string
	Location int
	Expected scope.Type // What was expected (in check mode)
	Got      scope.Type // What was found (synthesized)
}

func (e TypeError) Error() string {
	return e.Message
}

// OperatorRules provides engine-specific type inference rules for operators.
type OperatorRules interface {
	// BinaryOp returns the result type of a binary operator given operand types.
	BinaryOp(op string, left, right scope.Type) scope.Type
	// UnaryOp returns the result type of a unary operator.
	UnaryOp(op string, operand scope.Type) scope.Type
	// CanCast returns whether a value of `from` type can be cast to `to` type.
	CanCast(from, to scope.Type) bool
}

// NewChecker creates a new type checker with the given operator rules.
// If rules is nil, default rules are used.
func NewChecker(rules OperatorRules) *Checker {
	if rules == nil {
		rules = &DefaultOperatorRules{}
	}
	return &Checker{
		params:  make(map[int]*ParamTypeInference),
		opRules: rules,
	}
}

// Synth synthesizes (infers bottom-up) the type of an expression.
// This is the "what type does this have?" direction.
func (c *Checker) Synth(expr Expr) InferredType {
	switch e := expr.(type) {
	case *LiteralExpr:
		return InferredType{Type: e.Type, Kind: KindScalar, Source: SourceLiteral}

	case *ColumnRefExpr:
		return InferredType{Type: e.ResolvedType, Kind: KindScalar, Source: SourceCatalog}

	case *ParamExpr:
		// In synth mode, a parameter has unknown type unless previously inferred
		if prev, ok := c.params[e.Number]; ok {
			return InferredType{Type: prev.Type, Kind: KindScalar, Source: SourceParameter}
		}
		return InferredType{Type: scope.TypeUnknown, Kind: KindUnknown, Source: SourceParameter}

	case *BinaryOpExpr:
		left := c.Synth(e.Left)
		right := c.Synth(e.Right)

		// If one side is a parameter, use checking mode to infer its type
		if lp, ok := e.Left.(*ParamExpr); ok && !left.Type.IsUnknown() {
			c.Check(e.Left, right.Type, e.Location)
			_ = lp // parameter type recorded by Check
		}
		if rp, ok := e.Right.(*ParamExpr); ok && !right.Type.IsUnknown() {
			c.Check(e.Right, left.Type, e.Location)
			_ = rp
		}

		// If one side is a param and the other is known, infer param from known
		if left.Type.IsUnknown() && !right.Type.IsUnknown() {
			if lp, ok := e.Left.(*ParamExpr); ok {
				c.recordParam(lp.Number, right.Type, SourceContext, lp.Location)
			}
			left.Type = right.Type
		}
		if right.Type.IsUnknown() && !left.Type.IsUnknown() {
			if rp, ok := e.Right.(*ParamExpr); ok {
				c.recordParam(rp.Number, left.Type, SourceContext, rp.Location)
			}
			right.Type = left.Type
		}

		resultType := c.opRules.BinaryOp(e.Op, left.Type, right.Type)
		return InferredType{Type: resultType, Kind: KindScalar, Source: SourceOperator}

	case *FuncCallExpr:
		return InferredType{Type: e.ReturnType, Kind: KindScalar, Source: SourceFunction}

	case *TypeCastExpr:
		return InferredType{Type: e.CastType, Kind: KindScalar, Source: SourceOperator}

	case *SubqueryExpr:
		if len(e.Columns) == 1 {
			return InferredType{Type: e.Columns[0], Kind: KindScalar, Source: SourceCatalog}
		}
		return InferredType{Type: scope.TypeUnknown, Kind: KindRow, Source: SourceUnknown}

	case *BoolExpr:
		return InferredType{Type: scope.TypeBool, Kind: KindScalar, Source: SourceOperator}

	case *NullTestExpr:
		return InferredType{Type: scope.TypeBool, Kind: KindScalar, Source: SourceOperator}

	case *CaseExpr:
		if e.ResultType.IsUnknown() {
			return InferredType{Type: scope.TypeUnknown, Kind: KindScalar, Source: SourceUnknown}
		}
		return InferredType{Type: e.ResultType, Kind: KindScalar, Source: SourceOperator}

	case *CoalesceExpr:
		// Type is the type of the first non-null argument
		for _, arg := range e.Args {
			t := c.Synth(arg)
			if !t.Type.IsUnknown() {
				return t
			}
		}
		return InferredType{Type: scope.TypeUnknown, Kind: KindScalar, Source: SourceUnknown}

	case *InExpr:
		return InferredType{Type: scope.TypeBool, Kind: KindScalar, Source: SourceOperator}

	case *BetweenExpr:
		return InferredType{Type: scope.TypeBool, Kind: KindScalar, Source: SourceOperator}

	default:
		return InferredType{Type: scope.TypeUnknown, Kind: KindUnknown, Source: SourceUnknown}
	}
}

// Check verifies that an expression has the expected type (top-down).
// For parameters, this records the expected type as the parameter's inferred type.
// This is the "does this have the type I expect?" direction.
func (c *Checker) Check(expr Expr, expected scope.Type, location int) bool {
	switch e := expr.(type) {
	case *ParamExpr:
		// This is the key insight of bidirectional checking for SQL:
		// when a parameter appears in a context with a known type,
		// we learn the parameter's type from the context.
		c.recordParam(e.Number, expected, SourceContext, e.Location)
		return true

	case *LiteralExpr:
		if expected.IsUnknown() {
			return true
		}
		if !c.opRules.CanCast(e.Type, expected) {
			c.addError(TypeError{
				Message:  fmt.Sprintf("literal of type %s is not compatible with expected type %s", e.Type.Name, expected.Name),
				Location: location,
				Expected: expected,
				Got:      e.Type,
			})
			return false
		}
		return true

	case *TypeCastExpr:
		// An explicit cast always succeeds at the type level
		return true

	default:
		// For other expressions, synthesize and compare
		synth := c.Synth(expr)
		if synth.Type.IsUnknown() || expected.IsUnknown() {
			return true // Can't check if either side is unknown
		}
		if !synth.Type.Equals(expected) && !c.opRules.CanCast(synth.Type, expected) {
			c.addError(TypeError{
				Message:  fmt.Sprintf("type mismatch: expected %s but got %s", expected.Name, synth.Type.Name),
				Location: location,
				Expected: expected,
				Got:      synth.Type,
			})
			return false
		}
		return true
	}
}

// InferParamFromContext infers a parameter's type from its usage context.
// This is called when we know the expected type from the surrounding expression.
func (c *Checker) InferParamFromContext(paramNum int, expectedType scope.Type, location int) {
	c.recordParam(paramNum, expectedType, SourceContext, location)
}

func (c *Checker) recordParam(number int, typ scope.Type, source TypeSource, location int) {
	if existing, ok := c.params[number]; ok {
		// If we already have a type for this parameter, keep the more specific one
		if existing.Type.IsUnknown() && !typ.IsUnknown() {
			existing.Type = typ
			existing.Source = source
		}
		return
	}
	c.params[number] = &ParamTypeInference{
		Number:   number,
		Type:     typ,
		Source:   source,
		Location: location,
	}
}

func (c *Checker) addError(err TypeError) {
	c.errors = append(c.errors, err)
}

// ParamTypes returns all inferred parameter types.
func (c *Checker) ParamTypes() map[int]*ParamTypeInference {
	return c.params
}

// Errors returns all type errors found during checking.
func (c *Checker) Errors() []TypeError {
	return c.errors
}
