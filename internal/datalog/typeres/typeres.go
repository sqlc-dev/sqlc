// Package typeres expresses SQL type resolution as Datalog rules over a fixed
// schema of relations. It models how PostgreSQL resolves types for expressions,
// casts, operators, and functions.
package typeres

import (
	"fmt"
	"strconv"

	"github.com/sqlc-dev/sqlc/internal/datalog"
)

// ---------------------------------------------------------------------------
// Resolver accumulates EDB facts and, on Resolve(), evaluates the Datalog
// program to derive resolved types.
// ---------------------------------------------------------------------------

// Resolver holds the type resolution program.
type Resolver struct {
	// EDB facts stored as string tuples keyed by predicate name.
	facts map[string][][]string
}

// New creates a new Resolver.
func New() *Resolver {
	return &Resolver{
		facts: make(map[string][][]string),
	}
}

func (r *Resolver) addFact(pred string, vals ...string) {
	r.facts[pred] = append(r.facts[pred], vals)
}

// --- Methods to populate EDB (base facts) ---

// AddBaseType registers a type with its category and preferred flag.
// Category codes follow PostgreSQL: "N"=numeric, "S"=string, "B"=boolean,
// "D"=datetime, "U"=user-defined, etc.
func (r *Resolver) AddBaseType(name, category string, preferred bool) {
	r.addFact("base_type", name, category, strconv.FormatBool(preferred))
}

// AddImplicitCast registers an implicit (automatic) cast path.
func (r *Resolver) AddImplicitCast(from, to string) {
	r.addFact("implicit_cast", from, to)
}

// AddExplicitCast registers a cast that requires explicit CAST().
func (r *Resolver) AddExplicitCast(from, to string) {
	r.addFact("explicit_cast", from, to)
}

// AddAssignmentCast registers a cast allowed in assignment context.
func (r *Resolver) AddAssignmentCast(from, to string) {
	r.addFact("assignment_cast", from, to)
}

// AddOperator registers an operator signature.
func (r *Resolver) AddOperator(op, leftType, rightType, resultType string) {
	r.addFact("operator", op, leftType, rightType, resultType)
}

// AddFunction registers one parameter of a function signature.
// paramIndex is the 0-based index of the parameter.
func (r *Resolver) AddFunction(name string, paramIndex int, paramType, returnType string) {
	r.addFact("function_sig", name, strconv.Itoa(paramIndex), paramType, returnType)
}

// --- Methods to add expression structure ---

// SetExprType sets the known type of an expression (from a literal, column ref, etc.).
func (r *Resolver) SetExprType(exprID, typeName string) {
	r.addFact("expr_type", exprID, typeName)
}

// SetExprContext sets the expected type from context (e.g. INSERT target column).
func (r *Resolver) SetExprContext(exprID, contextType string) {
	r.addFact("expr_context", exprID, contextType)
}

// AddBinaryExpr registers a binary expression node.
func (r *Resolver) AddBinaryExpr(exprID, op, leftID, rightID string) {
	r.addFact("binary_expr", exprID, op, leftID, rightID)
}

// AddFuncCall registers a function call argument.
func (r *Resolver) AddFuncCall(exprID, funcName string, argIndex int, argExprID string) {
	r.addFact("func_call", exprID, funcName, strconv.Itoa(argIndex), argExprID)
}

// AddCastExpr registers an explicit CAST expression.
func (r *Resolver) AddCastExpr(exprID, innerExprID, targetType string) {
	r.addFact("cast_expr", exprID, innerExprID, targetType)
}

// AddCoalesceExpr registers a COALESCE argument.
func (r *Resolver) AddCoalesceExpr(exprID string, argIndex int, argExprID string) {
	r.addFact("coalesce_expr", exprID, strconv.Itoa(argIndex), argExprID)
}

// AddCaseExpr registers a CASE branch result expression.
func (r *Resolver) AddCaseExpr(exprID string, branchIndex int, resultExprID string) {
	r.addFact("case_expr", exprID, strconv.Itoa(branchIndex), resultExprID)
}

// ---------------------------------------------------------------------------
// Resolve
// ---------------------------------------------------------------------------

// Resolve runs the Datalog evaluation and returns results.
func (r *Resolver) Resolve() (*Result, error) {
	st := datalog.NewSymbolTable()
	p := datalog.NewProgram(st)

	// Populate EDB facts.
	for pred, tuples := range r.facts {
		for _, vals := range tuples {
			p.AddFact(pred, vals...)
		}
	}

	// Add the fixed set of type-resolution rules.
	addTypeRules(p, st)

	db, err := p.Evaluate()
	if err != nil {
		return nil, fmt.Errorf("typeres: evaluation failed: %w", err)
	}

	return &Result{db: db}, nil
}

// ---------------------------------------------------------------------------
// Result
// ---------------------------------------------------------------------------

// Result provides query methods over the evaluated type resolution database.
type Result struct {
	db *datalog.Database
}

// TypeOf returns the resolved type for the given expression ID.
func (res *Result) TypeOf(exprID string) (string, bool) {
	rows := res.db.Query("resolved_type")
	for _, row := range rows {
		if len(row) == 2 && row[0] == exprID {
			return row[1], true
		}
	}
	return "", false
}

// AllTypes returns a map of expression ID to resolved type name.
func (res *Result) AllTypes() map[string]string {
	m := make(map[string]string)
	rows := res.db.Query("resolved_type")
	for _, row := range rows {
		if len(row) == 2 {
			m[row[0]] = row[1]
		}
	}
	return m
}

// CanCast returns true if "from" can be implicitly cast to "to" (transitively).
func (res *Result) CanCast(from, to string) bool {
	if from == to {
		return true
	}
	return res.db.Contains("castable", from, to)
}

// CommonType returns the common supertype of two types if one exists.
func (res *Result) CommonType(a, b string) (string, bool) {
	if a == b {
		return a, true
	}
	rows := res.db.Query("common_type")
	for _, row := range rows {
		if len(row) == 3 && row[0] == a && row[1] == b {
			return row[2], true
		}
	}
	return "", false
}

// ---------------------------------------------------------------------------
// Datalog rules for type resolution
// ---------------------------------------------------------------------------

func addTypeRules(p *datalog.Program, st *datalog.SymbolTable) {
	V := datalog.Var
	C := datalog.Const
	rule := func(headPred string, headTerms ...datalog.TermExpr) *datalog.RuleBuilder {
		return datalog.NewRule(st, headPred, headTerms...)
	}

	// Rule 1: Direct type propagation.
	// resolved_type(E, T) :- expr_type(E, T).
	p.AddRule(rule("resolved_type", V("E"), V("T")).
		Where("expr_type", V("E"), V("T")).
		Build())

	// Rule 2a: Castable — base case from implicit_cast.
	// castable(A, B) :- implicit_cast(A, B).
	p.AddRule(rule("castable", V("A"), V("B")).
		Where("implicit_cast", V("A"), V("B")).
		Build())

	// Rule 2b: Castable — transitive closure.
	// castable(A, C) :- castable(A, B), implicit_cast(B, C).
	p.AddRule(rule("castable", V("A"), V("C")).
		Where("castable", V("A"), V("B")).
		Where("implicit_cast", V("B"), V("C")).
		Build())

	// Rule 3: Cast expressions — resolved type is the target type.
	// resolved_type(E, T) :- cast_expr(E, _, T).
	p.AddRule(rule("resolved_type", V("E"), V("T")).
		Where("cast_expr", V("E"), V("_inner"), V("T")).
		Build())

	// Rule 4: Binary operators — direct match.
	// operator_result(Op, LT, RT, Res) :- operator(Op, LT, RT, Res).
	p.AddRule(rule("operator_result", V("Op"), V("LT"), V("RT"), V("Res")).
		Where("operator", V("Op"), V("LT"), V("RT"), V("Res")).
		Build())

	// Rule 4a: Binary operators — promote left via implicit cast.
	// operator_result(Op, LT, RT, Res) :-
	//   castable(LT, LT2), operator(Op, LT2, RT, Res).
	p.AddRule(rule("operator_result", V("Op"), V("LT"), V("RT"), V("Res")).
		Where("castable", V("LT"), V("LT2")).
		Where("operator", V("Op"), V("LT2"), V("RT"), V("Res")).
		Build())

	// Rule 4b: Binary operators — promote right via implicit cast.
	// operator_result(Op, LT, RT, Res) :-
	//   castable(RT, RT2), operator(Op, LT, RT2, Res).
	p.AddRule(rule("operator_result", V("Op"), V("LT"), V("RT"), V("Res")).
		Where("castable", V("RT"), V("RT2")).
		Where("operator", V("Op"), V("LT"), V("RT2"), V("Res")).
		Build())

	// Rule 4c: Binary operators — promote both sides.
	// operator_result(Op, LT, RT, Res) :-
	//   castable(LT, LT2), castable(RT, RT2), operator(Op, LT2, RT2, Res).
	p.AddRule(rule("operator_result", V("Op"), V("LT"), V("RT"), V("Res")).
		Where("castable", V("LT"), V("LT2")).
		Where("castable", V("RT"), V("RT2")).
		Where("operator", V("Op"), V("LT2"), V("RT2"), V("Res")).
		Build())

	// Rule 4d: Resolve binary expression types.
	// resolved_type(E, Res) :-
	//   binary_expr(E, Op, L, R), resolved_type(L, LT), resolved_type(R, RT),
	//   operator_result(Op, LT, RT, Res).
	p.AddRule(rule("resolved_type", V("E"), V("Res")).
		Where("binary_expr", V("E"), V("Op"), V("L"), V("R")).
		Where("resolved_type", V("L"), V("LT")).
		Where("resolved_type", V("R"), V("RT")).
		Where("operator_result", V("Op"), V("LT"), V("RT"), V("Res")).
		Build())

	// Rule 5: Function calls — resolve via function signature.
	// For single-argument functions (simplest case):
	// resolved_type(E, RetT) :-
	//   func_call(E, Fn, Idx, ArgE),
	//   resolved_type(ArgE, ArgT),
	//   function_sig(Fn, Idx, ArgT, RetT).
	p.AddRule(rule("resolved_type", V("E"), V("RetT")).
		Where("func_call", V("E"), V("Fn"), V("Idx"), V("ArgE")).
		Where("resolved_type", V("ArgE"), V("ArgT")).
		Where("function_sig", V("Fn"), V("Idx"), V("ArgT"), V("RetT")).
		Build())

	// Rule 5a: Function calls — with implicit cast on argument.
	// resolved_type(E, RetT) :-
	//   func_call(E, Fn, Idx, ArgE),
	//   resolved_type(ArgE, ArgT),
	//   castable(ArgT, ParamT),
	//   function_sig(Fn, Idx, ParamT, RetT).
	p.AddRule(rule("resolved_type", V("E"), V("RetT")).
		Where("func_call", V("E"), V("Fn"), V("Idx"), V("ArgE")).
		Where("resolved_type", V("ArgE"), V("ArgT")).
		Where("castable", V("ArgT"), V("ParamT")).
		Where("function_sig", V("Fn"), V("Idx"), V("ParamT"), V("RetT")).
		Build())

	// Rule 6: Common type — same type.
	// common_type(A, A, A) :- base_type(A, _, _).
	p.AddRule(rule("common_type", V("A"), V("A"), V("A")).
		Where("base_type", V("A"), V("_cat"), V("_pref")).
		Build())

	// Rule 6a: Common type — A castable to B, B is preferred in its category.
	// common_type(A, B, B) :- castable(A, B), base_type(B, _, "true").
	p.AddRule(rule("common_type", V("A"), V("B"), V("B")).
		Where("castable", V("A"), V("B")).
		Where("base_type", V("B"), V("_cat"), C("true")).
		Build())

	// Rule 6b: Common type — B castable to A, A is preferred.
	// common_type(A, B, A) :- castable(B, A), base_type(A, _, "true").
	p.AddRule(rule("common_type", V("A"), V("B"), V("A")).
		Where("castable", V("B"), V("A")).
		Where("base_type", V("A"), V("_cat"), C("true")).
		Build())

	// Rule 6c: Common type — A castable to B (B not necessarily preferred).
	// common_type(A, B, B) :- castable(A, B), base_type(B, _, _).
	p.AddRule(rule("common_type", V("A"), V("B"), V("B")).
		Where("castable", V("A"), V("B")).
		Where("base_type", V("B"), V("_cat2"), V("_pref2")).
		Build())

	// Rule 6d: Common type — B castable to A.
	// common_type(A, B, A) :- castable(B, A), base_type(A, _, _).
	p.AddRule(rule("common_type", V("A"), V("B"), V("A")).
		Where("castable", V("B"), V("A")).
		Where("base_type", V("A"), V("_cat3"), V("_pref3")).
		Build())

	// Rule 6e: Common type — both castable to a third type (preferred).
	// common_type(A, B, C) :-
	//   castable(A, C), castable(B, C),
	//   base_type(C, _, "true").
	p.AddRule(rule("common_type", V("A"), V("B"), V("C")).
		Where("castable", V("A"), V("C")).
		Where("castable", V("B"), V("C")).
		Where("base_type", V("C"), V("_cat4"), C("true")).
		Build())

	// Rule 7: COALESCE — resolved type from common type of arguments.
	// For two arguments sharing the same coalesce expression:
	// resolved_type(E, CT) :-
	//   coalesce_expr(E, _, A1),
	//   coalesce_expr(E, _, A2),
	//   resolved_type(A1, T1),
	//   resolved_type(A2, T2),
	//   common_type(T1, T2, CT).
	p.AddRule(rule("resolved_type", V("E"), V("CT")).
		Where("coalesce_expr", V("E"), V("_i1"), V("A1")).
		Where("coalesce_expr", V("E"), V("_i2"), V("A2")).
		Where("resolved_type", V("A1"), V("T1")).
		Where("resolved_type", V("A2"), V("T2")).
		Where("common_type", V("T1"), V("T2"), V("CT")).
		Build())

	// Rule 8: CASE — resolved type from common type of branch results.
	// resolved_type(E, CT) :-
	//   case_expr(E, _, R1),
	//   case_expr(E, _, R2),
	//   resolved_type(R1, T1),
	//   resolved_type(R2, T2),
	//   common_type(T1, T2, CT).
	p.AddRule(rule("resolved_type", V("E"), V("CT")).
		Where("case_expr", V("E"), V("_b1"), V("R1")).
		Where("case_expr", V("E"), V("_b2"), V("R2")).
		Where("resolved_type", V("R1"), V("T1")).
		Where("resolved_type", V("R2"), V("T2")).
		Where("common_type", V("T1"), V("T2"), V("CT")).
		Build())

	// Rule 9: Context propagation — use context type if expression is castable.
	// resolved_type(E, CT) :-
	//   expr_context(E, CT),
	//   resolved_type(E, T),
	//   castable(T, CT).
	p.AddRule(rule("resolved_type", V("E"), V("CT")).
		Where("expr_context", V("E"), V("CT")).
		Where("resolved_type", V("E"), V("T")).
		Where("castable", V("T"), V("CT")).
		Build())
}

// ---------------------------------------------------------------------------
// LoadPostgreSQLTypes pre-populates the resolver with PostgreSQL's built-in
// type catalog covering the most common ~25 types and their cast paths.
// ---------------------------------------------------------------------------

// LoadPostgreSQLTypes adds PostgreSQL's common built-in types, their categories,
// preferred flags, and implicit cast paths to the resolver.
func LoadPostgreSQLTypes(r *Resolver) {
	// --- Base types ---
	// Format: name, category, preferred
	// Categories: B=boolean, N=numeric, S=string, D=datetime, U=user-defined,
	//             V=bit-string, T=timespan

	type baseType struct {
		name      string
		category  string
		preferred bool
	}

	types := []baseType{
		// Boolean
		{"bool", "B", true},

		// Numeric
		{"int2", "N", false},
		{"int4", "N", false},
		{"int8", "N", false},
		{"float4", "N", false},
		{"float8", "N", true},
		{"numeric", "N", false},

		// String
		{"text", "S", true},
		{"varchar", "S", false},
		{"char", "S", false},
		{"name", "S", false},

		// Binary
		{"bytea", "U", false},

		// Date/Time
		{"date", "D", false},
		{"time", "D", false},
		{"timetz", "D", false},
		{"timestamp", "D", false},
		{"timestamptz", "D", true},
		{"interval", "T", true},

		// UUID
		{"uuid", "U", false},

		// JSON
		{"json", "U", false},
		{"jsonb", "U", false},

		// XML
		{"xml", "U", false},

		// OID / internal
		{"oid", "N", false},
	}

	for _, t := range types {
		r.AddBaseType(t.name, t.category, t.preferred)
	}

	// --- Implicit casts ---
	// These follow PostgreSQL's pg_cast entries with castcontext = 'i'.

	implicitCasts := [][2]string{
		// Integer promotions
		{"int2", "int4"},
		{"int2", "int8"},
		{"int2", "float4"},
		{"int2", "float8"},
		{"int2", "numeric"},
		{"int4", "int8"},
		{"int4", "float4"},
		{"int4", "float8"},
		{"int4", "numeric"},
		{"int8", "float4"},
		{"int8", "float8"},
		{"int8", "numeric"},

		// Float promotions
		{"float4", "float8"},

		// Numeric to float
		{"numeric", "float4"},
		{"numeric", "float8"},

		// String conversions
		{"char", "varchar"},
		{"char", "text"},
		{"varchar", "text"},
		{"name", "text"},

		// Date/time promotions
		{"date", "timestamp"},
		{"date", "timestamptz"},
		{"timestamp", "timestamptz"},
		{"time", "timetz"},
		{"time", "interval"},

		// OID
		{"int4", "oid"},
		{"int8", "oid"},

		// JSON
		{"json", "jsonb"},
	}

	for _, c := range implicitCasts {
		r.AddImplicitCast(c[0], c[1])
	}

	// --- Assignment casts ---
	assignCasts := [][2]string{
		{"int4", "int2"},
		{"int8", "int4"},
		{"int8", "int2"},
		{"float8", "float4"},
		{"float8", "numeric"},
		{"float4", "numeric"},
		{"numeric", "int2"},
		{"numeric", "int4"},
		{"numeric", "int8"},
		{"text", "varchar"},
		{"text", "char"},
		{"varchar", "char"},
		{"timestamptz", "timestamp"},
		{"timestamptz", "date"},
		{"timestamp", "date"},
		{"timetz", "time"},
	}

	for _, c := range assignCasts {
		r.AddAssignmentCast(c[0], c[1])
	}

	// --- Explicit casts ---
	explicitCasts := [][2]string{
		{"text", "int4"},
		{"text", "int8"},
		{"text", "float8"},
		{"text", "numeric"},
		{"text", "bool"},
		{"text", "date"},
		{"text", "timestamp"},
		{"text", "timestamptz"},
		{"text", "interval"},
		{"text", "uuid"},
		{"text", "json"},
		{"text", "jsonb"},
		{"int4", "bool"},
		{"bool", "int4"},
		{"float8", "int4"},
		{"float8", "int8"},
		{"float8", "int2"},
		{"jsonb", "json"},
	}

	for _, c := range explicitCasts {
		r.AddExplicitCast(c[0], c[1])
	}

	// --- Common operators ---
	// Arithmetic: +, -, *, /
	arithOps := []string{"+", "-", "*", "/"}
	numericTypes := []string{"int2", "int4", "int8", "float4", "float8", "numeric"}

	for _, op := range arithOps {
		for _, t := range numericTypes {
			r.AddOperator(op, t, t, t)
		}
	}

	// Comparison: =, <>, <, >, <=, >=
	cmpOps := []string{"=", "<>", "<", ">", "<=", ">="}
	allComparable := append(numericTypes, "text", "varchar", "char", "date",
		"time", "timetz", "timestamp", "timestamptz", "interval",
		"bool", "uuid", "bytea")

	for _, op := range cmpOps {
		for _, t := range allComparable {
			r.AddOperator(op, t, t, "bool")
		}
	}

	// String concatenation
	r.AddOperator("||", "text", "text", "text")
	r.AddOperator("||", "varchar", "varchar", "text")
	r.AddOperator("||", "text", "varchar", "text")
	r.AddOperator("||", "varchar", "text", "text")

	// Boolean operators
	r.AddOperator("AND", "bool", "bool", "bool")
	r.AddOperator("OR", "bool", "bool", "bool")

	// --- Common functions ---
	// length(text) -> int4
	r.AddFunction("length", 0, "text", "int4")
	r.AddFunction("length", 0, "bytea", "int4")

	// upper/lower(text) -> text
	r.AddFunction("upper", 0, "text", "text")
	r.AddFunction("lower", 0, "text", "text")

	// trim(text) -> text
	r.AddFunction("trim", 0, "text", "text")
	r.AddFunction("btrim", 0, "text", "text")
	r.AddFunction("ltrim", 0, "text", "text")
	r.AddFunction("rtrim", 0, "text", "text")

	// substring(text, int4, int4) -> text
	r.AddFunction("substring", 0, "text", "text")
	r.AddFunction("substring", 1, "int4", "text")
	r.AddFunction("substring", 2, "int4", "text")

	// replace(text, text, text) -> text
	r.AddFunction("replace", 0, "text", "text")
	r.AddFunction("replace", 1, "text", "text")
	r.AddFunction("replace", 2, "text", "text")

	// abs(numeric) -> numeric
	r.AddFunction("abs", 0, "int4", "int4")
	r.AddFunction("abs", 0, "int8", "int8")
	r.AddFunction("abs", 0, "float8", "float8")
	r.AddFunction("abs", 0, "numeric", "numeric")

	// round/ceil/floor
	r.AddFunction("round", 0, "numeric", "numeric")
	r.AddFunction("round", 0, "float8", "float8")
	r.AddFunction("ceil", 0, "numeric", "numeric")
	r.AddFunction("ceil", 0, "float8", "float8")
	r.AddFunction("floor", 0, "numeric", "numeric")
	r.AddFunction("floor", 0, "float8", "float8")

	// now() -> timestamptz
	r.AddFunction("now", 0, "void", "timestamptz")

	// coalesce is handled structurally, not as a function

	// count(*) -> int8 (aggregate)
	r.AddFunction("count", 0, "any", "int8")

	// sum
	r.AddFunction("sum", 0, "int4", "int8")
	r.AddFunction("sum", 0, "int8", "numeric")
	r.AddFunction("sum", 0, "float4", "float4")
	r.AddFunction("sum", 0, "float8", "float8")
	r.AddFunction("sum", 0, "numeric", "numeric")

	// avg
	r.AddFunction("avg", 0, "int4", "numeric")
	r.AddFunction("avg", 0, "int8", "numeric")
	r.AddFunction("avg", 0, "float4", "float8")
	r.AddFunction("avg", 0, "float8", "float8")
	r.AddFunction("avg", 0, "numeric", "numeric")

	// min/max — same type in, same type out
	for _, fn := range []string{"min", "max"} {
		for _, t := range allComparable {
			r.AddFunction(fn, 0, t, t)
		}
	}

	// string_agg(text, text) -> text
	r.AddFunction("string_agg", 0, "text", "text")
	r.AddFunction("string_agg", 1, "text", "text")

	// array_agg is complex; skip for now

	// to_char, to_date, to_timestamp, to_number
	r.AddFunction("to_char", 0, "timestamptz", "text")
	r.AddFunction("to_char", 0, "timestamp", "text")
	r.AddFunction("to_char", 0, "interval", "text")
	r.AddFunction("to_char", 0, "numeric", "text")
	r.AddFunction("to_char", 0, "int4", "text")
	r.AddFunction("to_char", 0, "int8", "text")
	r.AddFunction("to_char", 0, "float8", "text")
	r.AddFunction("to_char", 1, "text", "text")

	r.AddFunction("to_date", 0, "text", "date")
	r.AddFunction("to_date", 1, "text", "date")

	r.AddFunction("to_timestamp", 0, "text", "timestamptz")
	r.AddFunction("to_timestamp", 0, "float8", "timestamptz")
	r.AddFunction("to_timestamp", 1, "text", "timestamptz")

	r.AddFunction("to_number", 0, "text", "numeric")
	r.AddFunction("to_number", 1, "text", "numeric")

	// gen_random_uuid() -> uuid
	r.AddFunction("gen_random_uuid", 0, "void", "uuid")

	// jsonb_build_object, json_build_object — variadic, simplified
	r.AddFunction("jsonb_build_object", 0, "any", "jsonb")
	r.AddFunction("json_build_object", 0, "any", "json")

	// date_trunc(text, timestamptz) -> timestamptz
	r.AddFunction("date_trunc", 0, "text", "timestamptz")
	r.AddFunction("date_trunc", 1, "timestamptz", "timestamptz")
	r.AddFunction("date_trunc", 0, "text", "timestamp")
	r.AddFunction("date_trunc", 1, "timestamp", "timestamp")
}
