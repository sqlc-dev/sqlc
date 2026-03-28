package typecheck

import (
	"testing"

	"github.com/sqlc-dev/sqlc/internal/analysis/scope"
)

func TestSynthLiteral(t *testing.T) {
	c := NewChecker(nil)

	tests := []struct {
		name string
		expr Expr
		want string
	}{
		{"int", &LiteralExpr{Type: scope.TypeInt}, "integer"},
		{"text", &LiteralExpr{Type: scope.TypeText}, "text"},
		{"bool", &LiteralExpr{Type: scope.TypeBool}, "boolean"},
		{"float", &LiteralExpr{Type: scope.TypeFloat}, "float"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := c.Synth(tt.expr)
			if result.Type.Name != tt.want {
				t.Errorf("Synth: got type %q, want %q", result.Type.Name, tt.want)
			}
			if result.Source != SourceLiteral {
				t.Errorf("Synth: got source %v, want SourceLiteral", result.Source)
			}
		})
	}
}

func TestSynthColumnRef(t *testing.T) {
	c := NewChecker(nil)

	expr := &ColumnRefExpr{
		Parts:        []string{"users", "name"},
		ResolvedType: scope.TypeText,
	}

	result := c.Synth(expr)
	if result.Type.Name != "text" {
		t.Errorf("got type %q, want 'text'", result.Type.Name)
	}
	if result.Source != SourceCatalog {
		t.Errorf("got source %v, want SourceCatalog", result.Source)
	}
}

func TestSynthBinaryOp_Comparison(t *testing.T) {
	c := NewChecker(nil)

	// age = 25 → boolean
	expr := &BinaryOpExpr{
		Op:    "=",
		Left:  &ColumnRefExpr{ResolvedType: scope.TypeInt},
		Right: &LiteralExpr{Type: scope.TypeInt},
	}

	result := c.Synth(expr)
	if result.Type.Name != "boolean" {
		t.Errorf("got type %q, want 'boolean'", result.Type.Name)
	}
}

func TestSynthBinaryOp_Arithmetic(t *testing.T) {
	c := NewChecker(nil)

	// age + 1 → integer
	expr := &BinaryOpExpr{
		Op:    "+",
		Left:  &ColumnRefExpr{ResolvedType: scope.TypeInt},
		Right: &LiteralExpr{Type: scope.TypeInt},
	}

	result := c.Synth(expr)
	if result.Type.Name != "integer" {
		t.Errorf("got type %q, want 'integer'", result.Type.Name)
	}
}

func TestSynthBinaryOp_Concat(t *testing.T) {
	c := NewChecker(nil)

	// name || ' suffix' → text
	expr := &BinaryOpExpr{
		Op:    "||",
		Left:  &ColumnRefExpr{ResolvedType: scope.TypeText},
		Right: &LiteralExpr{Type: scope.TypeText},
	}

	result := c.Synth(expr)
	if result.Type.Name != "text" {
		t.Errorf("got type %q, want 'text'", result.Type.Name)
	}
}

func TestCheckParamInfersType(t *testing.T) {
	c := NewChecker(nil)

	// WHERE age > $1 → $1 should be inferred as integer
	param := &ParamExpr{Number: 1, Location: 42}

	c.Check(param, scope.TypeInt, 42)

	params := c.ParamTypes()
	if len(params) != 1 {
		t.Fatalf("expected 1 param, got %d", len(params))
	}
	p, ok := params[1]
	if !ok {
		t.Fatal("param $1 not found")
	}
	if p.Type.Name != "integer" {
		t.Errorf("param $1 type: got %q, want 'integer'", p.Type.Name)
	}
	if p.Source != SourceContext {
		t.Errorf("param $1 source: got %v, want SourceContext", p.Source)
	}
}

func TestSynthBinaryOpWithParam(t *testing.T) {
	c := NewChecker(nil)

	// age = $1 → $1 should be inferred as integer (from column ref)
	expr := &BinaryOpExpr{
		Op:    "=",
		Left:  &ColumnRefExpr{ResolvedType: scope.TypeInt},
		Right: &ParamExpr{Number: 1, Location: 20},
	}

	c.Synth(expr)

	params := c.ParamTypes()
	p, ok := params[1]
	if !ok {
		t.Fatal("param $1 not found after synth")
	}
	if p.Type.Name != "integer" {
		t.Errorf("param $1 type: got %q, want 'integer'", p.Type.Name)
	}
}

func TestCheckMultipleParams(t *testing.T) {
	c := NewChecker(nil)

	// WHERE name = $1 AND age > $2
	nameExpr := &BinaryOpExpr{
		Op:    "=",
		Left:  &ColumnRefExpr{ResolvedType: scope.TypeText},
		Right: &ParamExpr{Number: 1, Location: 10},
	}
	ageExpr := &BinaryOpExpr{
		Op:    ">",
		Left:  &ColumnRefExpr{ResolvedType: scope.TypeInt},
		Right: &ParamExpr{Number: 2, Location: 30},
	}

	c.Synth(nameExpr)
	c.Synth(ageExpr)

	params := c.ParamTypes()
	if len(params) != 2 {
		t.Fatalf("expected 2 params, got %d", len(params))
	}
	if params[1].Type.Name != "text" {
		t.Errorf("$1: got %q, want 'text'", params[1].Type.Name)
	}
	if params[2].Type.Name != "integer" {
		t.Errorf("$2: got %q, want 'integer'", params[2].Type.Name)
	}
}

func TestSynthFuncCall(t *testing.T) {
	c := NewChecker(nil)

	expr := &FuncCallExpr{
		Name:       "count",
		ReturnType: scope.Type{Name: "bigint", NotNull: true},
	}

	result := c.Synth(expr)
	if result.Type.Name != "bigint" {
		t.Errorf("got %q, want 'bigint'", result.Type.Name)
	}
	if result.Source != SourceFunction {
		t.Errorf("got source %v, want SourceFunction", result.Source)
	}
}

func TestSynthTypeCast(t *testing.T) {
	c := NewChecker(nil)

	// $1::integer → param should get integer type, expr returns integer
	expr := &TypeCastExpr{
		Arg:      &ParamExpr{Number: 1, Location: 5},
		CastType: scope.TypeInt,
	}

	result := c.Synth(expr)
	if result.Type.Name != "integer" {
		t.Errorf("cast result: got %q, want 'integer'", result.Type.Name)
	}
}

func TestSynthSubquery(t *testing.T) {
	c := NewChecker(nil)

	// Scalar subquery with one column
	expr := &SubqueryExpr{
		Columns: []scope.Type{scope.TypeInt},
	}

	result := c.Synth(expr)
	if result.Type.Name != "integer" {
		t.Errorf("subquery: got %q, want 'integer'", result.Type.Name)
	}
	if result.Kind != KindScalar {
		t.Errorf("subquery: got kind %v, want KindScalar", result.Kind)
	}
}

func TestSynthBoolExpr(t *testing.T) {
	c := NewChecker(nil)

	expr := &BoolExpr{Op: "AND", Args: []Expr{
		&LiteralExpr{Type: scope.TypeBool},
		&LiteralExpr{Type: scope.TypeBool},
	}}

	result := c.Synth(expr)
	if result.Type.Name != "boolean" {
		t.Errorf("AND: got %q, want 'boolean'", result.Type.Name)
	}
}

func TestSynthNullTest(t *testing.T) {
	c := NewChecker(nil)

	expr := &NullTestExpr{
		Arg:   &ColumnRefExpr{ResolvedType: scope.TypeText},
		IsNot: false,
	}

	result := c.Synth(expr)
	if result.Type.Name != "boolean" {
		t.Errorf("IS NULL: got %q, want 'boolean'", result.Type.Name)
	}
}

func TestSynthCoalesce(t *testing.T) {
	c := NewChecker(nil)

	expr := &CoalesceExpr{
		Args: []Expr{
			&ColumnRefExpr{ResolvedType: scope.TypeText},
			&LiteralExpr{Type: scope.TypeText},
		},
	}

	result := c.Synth(expr)
	if result.Type.Name != "text" {
		t.Errorf("COALESCE: got %q, want 'text'", result.Type.Name)
	}
}

func TestSynthIn(t *testing.T) {
	c := NewChecker(nil)

	expr := &InExpr{
		Expr:   &ColumnRefExpr{ResolvedType: scope.TypeInt},
		Values: []Expr{&LiteralExpr{Type: scope.TypeInt}},
	}

	result := c.Synth(expr)
	if result.Type.Name != "boolean" {
		t.Errorf("IN: got %q, want 'boolean'", result.Type.Name)
	}
}

func TestSynthBetween(t *testing.T) {
	c := NewChecker(nil)

	expr := &BetweenExpr{
		Expr: &ColumnRefExpr{ResolvedType: scope.TypeInt},
		Low:  &LiteralExpr{Type: scope.TypeInt},
		High: &LiteralExpr{Type: scope.TypeInt},
	}

	result := c.Synth(expr)
	if result.Type.Name != "boolean" {
		t.Errorf("BETWEEN: got %q, want 'boolean'", result.Type.Name)
	}
}

func TestParamKeepsFirstKnownType(t *testing.T) {
	c := NewChecker(nil)

	// First use: $1 = integer column
	c.Check(&ParamExpr{Number: 1}, scope.TypeInt, 0)

	// Second use: $1 = text column (conflicting)
	c.Check(&ParamExpr{Number: 1}, scope.TypeText, 0)

	// Should keep the first known type
	params := c.ParamTypes()
	if params[1].Type.Name != "integer" {
		t.Errorf("param kept wrong type: got %q, want 'integer'", params[1].Type.Name)
	}
}

func TestMySQLConcatIsOR(t *testing.T) {
	c := NewChecker(&MySQLOperatorRules{})

	// In MySQL default mode, || is OR, not concat
	expr := &BinaryOpExpr{
		Op:    "||",
		Left:  &LiteralExpr{Type: scope.TypeBool},
		Right: &LiteralExpr{Type: scope.TypeBool},
	}

	result := c.Synth(expr)
	if result.Type.Name != "boolean" {
		t.Errorf("MySQL ||: got %q, want 'boolean'", result.Type.Name)
	}
}

func TestPostgreSQLConcatIsText(t *testing.T) {
	c := NewChecker(&PostgreSQLOperatorRules{})

	expr := &BinaryOpExpr{
		Op:    "||",
		Left:  &LiteralExpr{Type: scope.TypeText},
		Right: &LiteralExpr{Type: scope.TypeText},
	}

	result := c.Synth(expr)
	if result.Type.Name != "text" {
		t.Errorf("PostgreSQL ||: got %q, want 'text'", result.Type.Name)
	}
}

func TestArithmeticTypePromotion(t *testing.T) {
	rules := &DefaultOperatorRules{}

	tests := []struct {
		name     string
		left     scope.Type
		right    scope.Type
		wantName string
	}{
		{"int+int", scope.Type{Name: "integer"}, scope.Type{Name: "integer"}, "integer"},
		{"int+bigint", scope.Type{Name: "integer"}, scope.Type{Name: "bigint"}, "bigint"},
		{"smallint+int", scope.Type{Name: "smallint"}, scope.Type{Name: "integer"}, "integer"},
		{"int+numeric", scope.Type{Name: "integer"}, scope.Type{Name: "numeric"}, "numeric"},
		{"real+float", scope.Type{Name: "real"}, scope.Type{Name: "double precision"}, "double precision"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rules.BinaryOp("+", tt.left, tt.right)
			if result.Name != tt.wantName {
				t.Errorf("got %q, want %q", result.Name, tt.wantName)
			}
		})
	}
}
