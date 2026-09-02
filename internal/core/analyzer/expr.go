package analyzer

import (
	"fmt"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

type exprType struct {
	typeOID int64
	// typeName names a type the catalog does not hold — a cast to a type no
	// dialect seeded and no schema declared, or an array of one. Analysis
	// never adds a type: it reports the name the query used and carries on,
	// so a query can be analyzed against a catalog it cannot write to.
	typeName           string
	nullable           bool
	sourceClassOID     int64
	sourceAttributeOID int64
	sourceTableAlias   string
}

func (a *analyzer) typeExpr(n ast.Node) (exprType, error) {
	switch e := n.(type) {
	case nil:
		return exprType{}, nil

	case *ast.TODO:
		return exprType{}, nil

	case *ast.A_Const:
		return a.typeConst(e)

	case *ast.ColumnRef:
		return a.typeColumnRef(e)

	case *ast.ParamRef:
		return a.typeParamRef(e)

	case *ast.A_Expr:
		return a.typeAExpr(e)

	case *ast.BoolExpr:
		return a.typeBoolExpr(e)

	case *ast.FuncCall:
		return a.typeFuncCall(e)

	case *ast.TypeCast:
		return a.typeTypeCast(e)

	case *ast.NullTest:
		if _, err := a.typeExpr(e.Arg); err != nil {
			return exprType{}, err
		}
		return a.boolType(false)

	case *ast.In:
		return a.typeIn(e)

	case *ast.BetweenExpr:
		return a.typeBetween(e)

	case *ast.CaseExpr:
		return a.typeCase(e)

	case *ast.CoalesceExpr:
		return a.typeCoalesce(e)

	case *ast.MinMaxExpr:
		return a.typeFirstOf(listItems(e.Args), false)

	case *ast.A_ArrayExpr:
		return a.typeArrayExpr(e)

	case *ast.SubLink:
		return a.typeSubLink(e)

	case *ast.CollateExpr:
		return a.typeExpr(e.Arg)

	case *ast.IntervalExpr:
		// A dialect with no interval type of its own leaves it untyped.
		oid, err := a.cat.TypeOID("interval")
		if err != nil {
			return exprType{}, nil
		}
		return exprType{typeOID: oid}, nil

	// A user variable and an ordinal in GROUP BY / ORDER BY carry no type of
	// their own, and neither is something to resolve.
	case *ast.VariableExpr, *ast.Integer, *ast.String, *ast.Float:
		return exprType{}, nil

	case *ast.List:
		for _, item := range listItems(e) {
			if _, err := a.typeExpr(item); err != nil {
				return exprType{}, err
			}
		}
		return exprType{}, nil
	}
	return exprType{}, fmt.Errorf("typeExpr: unsupported %T", n)
}

func (a *analyzer) typeConst(c *ast.A_Const) (exprType, error) {
	switch v := c.Val.(type) {
	case *ast.Integer:
		_ = v
		oid, err := a.cat.ConstTypeOID(core.ConstInteger)
		if err != nil {
			return exprType{}, err
		}
		return exprType{typeOID: oid}, nil
	case *ast.Float:
		oid, err := a.cat.ConstTypeOID(core.ConstFloat)
		if err != nil {
			return exprType{}, err
		}
		return exprType{typeOID: oid}, nil
	case *ast.String:
		oid, err := a.cat.ConstTypeOID(core.ConstString)
		if err != nil {
			return exprType{}, err
		}
		return exprType{typeOID: oid}, nil
	case *ast.Boolean:
		return a.boolType(false)
	case *ast.Null, nil:
		return exprType{nullable: true}, nil
	}
	return exprType{}, fmt.Errorf("typeConst: unsupported %T", c.Val)
}

func (a *analyzer) boolType(nullable bool) (exprType, error) {
	oid, err := a.cat.BoolTypeOID()
	if err != nil {
		return exprType{}, err
	}
	return exprType{typeOID: oid, nullable: nullable}, nil
}

func (a *analyzer) typeColumnRef(c *ast.ColumnRef) (exprType, error) {
	parts := flattenFields(c.Fields)
	if len(parts) == 0 {
		return exprType{}, fmt.Errorf("column ref: empty")
	}
	relation := ""
	column := parts[0]
	if len(parts) >= 2 {
		relation = parts[0]
		column = parts[1]
	}
	rel, col, ok, err := a.resolveColumn(relation, column)
	if err != nil {
		return exprType{}, err
	}
	if !ok {
		if relation != "" {
			return exprType{}, fmt.Errorf("unknown column %q.%q", relation, column)
		}
		// The name may be one the target list assigned rather than one a
		// relation offers.
		if t, ok, err := a.typeAlias(column); err != nil || ok {
			return t, err
		}
		return exprType{}, fmt.Errorf("unknown column %q", column)
	}
	return exprType{
		typeOID:            col.TypeOID,
		nullable:           !col.NotNull,
		sourceClassOID:     rel.classOID,
		sourceAttributeOID: col.AttOID,
		sourceTableAlias:   rel.alias,
	}, nil
}

func flattenFields(fields *ast.List) []string {
	if fields == nil {
		return nil
	}
	out := make([]string, 0, len(fields.Items))
	for _, item := range fields.Items {
		switch v := item.(type) {
		case *ast.String:
			out = append(out, v.Str)
		case *ast.A_Star:
			out = append(out, "*")
			return out
		}
	}
	return out
}

func (a *analyzer) typeParamRef(p *ast.ParamRef) (exprType, error) {
	cur, ok := a.params[p.Number]
	if !ok {
		cur = core.Parameter{Number: p.Number}
		a.params[p.Number] = cur
	}
	return exprType{typeOID: cur.TypeOID, nullable: !cur.NotNull}, nil
}

func (a *analyzer) inferParam(number int, t exprType) {
	cur, ok := a.params[number]
	if !ok {
		cur = core.Parameter{Number: number}
	}
	typed := cur.TypeOID == 0 && cur.DataType == "" && (t.typeOID != 0 || t.typeName != "")
	if typed {
		cur.TypeOID = t.typeOID
		cur.DataType, cur.IsArray = a.typeNameOf(t)
		cur.NotNull = !t.nullable
		cur.Type = a.typeExprOf(t, "")
	}
	if cur.Source == nil && t.sourceAttributeOID != 0 {
		ad, err := a.cat.LookupAttribute(t.sourceAttributeOID)
		if err == nil {
			cur.Source = &core.ColumnSource{
				Schema:     ad.Schema,
				Table:      ad.Table,
				TableAlias: t.sourceTableAlias,
				Column:     ad.Column,
			}
			if typed {
				cur.Type = a.typeExprOf(t, ad.DeclType)
			}
		}
	}
	a.params[number] = cur
}

// nameParamAfter names a placeholder compared with a function call after
// the function, the way a placeholder compared with a column is named after
// the column.
func (a *analyzer) nameParamAfter(number int, other ast.Node) {
	fc, ok := other.(*ast.FuncCall)
	if !ok {
		return
	}
	cur := a.params[number]
	if cur.Name == "" && cur.Source == nil {
		cur.Name = funcCallName(fc)
		a.params[number] = cur
	}
}

func (a *analyzer) typeAExpr(e *ast.A_Expr) (exprType, error) {
	// Not every engine classifies its operators, so a zero kind is a plain
	// operator application rather than an unset field. LIKE and its relatives
	// name a real operator too, so they resolve the same way.
	switch e.Kind {
	case ast.A_Expr_Kind_OP, ast.A_Expr_Kind_LIKE, ast.A_Expr_Kind_ILIKE, ast.A_Expr_Kind_SIMILAR, 0:
	case ast.A_Expr_Kind_OP_ANY, ast.A_Expr_Kind_OP_ALL:
		return a.typeQuantifiedExpr(e)
	case ast.A_Expr_Kind_IN, ast.A_Expr_Kind_BETWEEN, ast.A_Expr_Kind_NOT_BETWEEN,
		ast.A_Expr_Kind_BETWEEN_SYM, ast.A_Expr_Kind_NOT_BETWEEN_SYM:
		return a.typePredicateList(e)
	case ast.A_Expr_Kind_DISTINCT, ast.A_Expr_Kind_NOT_DISTINCT:
		return a.typeComparison(e)
	case ast.A_Expr_Kind_NULLIF:
		return a.typeNullIf(e)
	default:
		return exprType{}, fmt.Errorf("a_expr: unsupported kind %v", e.Kind)
	}
	opName := normalizeOperator(opNameFromList(e.Name))
	if opName == "" {
		return exprType{}, fmt.Errorf("a_expr: unnamed operator")
	}

	// Engines that do not have a dedicated boolean node report AND and OR as
	// operators. They combine predicates whatever the dialect.
	if opName == "AND" || opName == "OR" {
		if _, err := a.typeExpr(e.Lexpr); err != nil {
			return exprType{}, err
		}
		if _, err := a.typeExpr(e.Rexpr); err != nil {
			return exprType{}, err
		}
		return a.boolType(false)
	}

	leftT, err := a.typeExpr(e.Lexpr)
	if err != nil {
		return exprType{}, err
	}
	rightT, err := a.typeExpr(e.Rexpr)
	if err != nil {
		return exprType{}, err
	}

	if pr, ok := e.Lexpr.(*ast.ParamRef); ok && rightT.typeOID != 0 {
		a.inferParam(pr.Number, rightT)
		a.nameParamAfter(pr.Number, e.Rexpr)
		leftT = rightT
	}
	if pr, ok := e.Rexpr.(*ast.ParamRef); ok && leftT.typeOID != 0 {
		a.inferParam(pr.Number, leftT)
		a.nameParamAfter(pr.Number, e.Lexpr)
		rightT = leftT
	}

	overload, err := a.resolveOperator(opName, leftT.typeOID, rightT.typeOID)
	if err != nil {
		return exprType{}, err
	}
	return exprType{
		typeOID:  overload.ResultTypeOID,
		nullable: leftT.nullable || rightT.nullable,
	}, nil
}

// typeQuantifiedExpr types "x = ANY($1)" and "x > ALL(...)": the right side
// holds values of the left side's type, and the result is a predicate.
func (a *analyzer) typeQuantifiedExpr(e *ast.A_Expr) (exprType, error) {
	leftT, err := a.typeExpr(e.Lexpr)
	if err != nil {
		return exprType{}, err
	}
	if err := a.typeOperands(e.Rexpr, leftT); err != nil {
		return exprType{}, err
	}
	return a.boolType(false)
}

// typePredicateList types IN and BETWEEN, where the right side is a list of
// values compared against the left.
func (a *analyzer) typePredicateList(e *ast.A_Expr) (exprType, error) {
	leftT, err := a.typeExpr(e.Lexpr)
	if err != nil {
		return exprType{}, err
	}
	if l, ok := e.Rexpr.(*ast.List); ok {
		for _, item := range listItems(l) {
			if err := a.typeOperands(item, leftT); err != nil {
				return exprType{}, err
			}
		}
		return a.boolType(false)
	}
	if err := a.typeOperands(e.Rexpr, leftT); err != nil {
		return exprType{}, err
	}
	return a.boolType(false)
}

// typeIn types the IN node the engines that have one report, where the values
// compared against are held apart from the expression.
func (a *analyzer) typeIn(e *ast.In) (exprType, error) {
	leftT, err := a.typeExpr(e.Expr)
	if err != nil {
		return exprType{}, err
	}
	for _, item := range e.List {
		if err := a.typeOperands(item, leftT); err != nil {
			return exprType{}, err
		}
	}
	// "x IN (SELECT ...)" compares x against the subquery's column, and the
	// subquery's own placeholders are reported with the rest.
	if sel, ok := e.Sel.(*ast.SelectStmt); ok {
		cols, err := a.subqueryColumns(sel)
		if err != nil {
			return exprType{}, err
		}
		if len(cols) > 0 {
			if err := a.typeOperands(e.Expr, exprType{typeOID: cols[0].TypeOID, nullable: !cols[0].NotNull}); err != nil {
				return exprType{}, err
			}
		}
	}
	return a.boolType(false)
}

// typeBetween types the BETWEEN node the engines that have one report.
func (a *analyzer) typeBetween(e *ast.BetweenExpr) (exprType, error) {
	leftT, err := a.typeExpr(e.Expr)
	if err != nil {
		return exprType{}, err
	}
	for _, bound := range []ast.Node{e.Left, e.Right} {
		if err := a.typeOperands(bound, leftT); err != nil {
			return exprType{}, err
		}
	}
	return a.boolType(false)
}

// typeCase types CASE. Its result is the first branch's, and it is nullable
// unless every branch and the default are.
func (a *analyzer) typeCase(e *ast.CaseExpr) (exprType, error) {
	argT, err := a.typeExpr(e.Arg)
	if err != nil {
		return exprType{}, err
	}
	results := make([]ast.Node, 0, len(listItems(e.Args))+1)
	for _, item := range listItems(e.Args) {
		when, ok := item.(*ast.CaseWhen)
		if !ok {
			continue
		}
		// "CASE x WHEN y" compares y against x; "CASE WHEN y" is a predicate.
		if e.Arg != nil {
			if err := a.typeOperands(when.Expr, argT); err != nil {
				return exprType{}, err
			}
		} else if _, err := a.typeExpr(when.Expr); err != nil {
			return exprType{}, err
		}
		results = append(results, when.Result)
	}
	if e.Defresult != nil {
		results = append(results, e.Defresult)
	}
	t, err := a.typeFirstOf(results, e.Defresult == nil)
	if err != nil {
		return exprType{}, err
	}
	return t, nil
}

// typeCoalesce types COALESCE, which is its first typed argument's type and
// is null only when every argument is.
func (a *analyzer) typeCoalesce(e *ast.CoalesceExpr) (exprType, error) {
	var out exprType
	found := false
	nullable := true
	for _, n := range listItems(e.Args) {
		t, err := a.typeExpr(n)
		if err != nil {
			return exprType{}, err
		}
		if !found && t.typeOID != 0 {
			// The result is an expression's, not the column's it came from.
			out = exprType{typeOID: t.typeOID, typeName: t.typeName}
			found = true
		}
		nullable = nullable && t.nullable
	}
	out.nullable = nullable
	return out, nil
}

// typeFirstOf types a set of alternative results, taking the first one that has
// a type. The result is nullable when any alternative is, or when the caller
// says the set is not exhaustive.
func (a *analyzer) typeFirstOf(nodes []ast.Node, nullable bool) (exprType, error) {
	var out exprType
	found := false
	for _, n := range nodes {
		t, err := a.typeExpr(n)
		if err != nil {
			return exprType{}, err
		}
		if !found && t.typeOID != 0 {
			out = t
			found = true
		}
		nullable = nullable || t.nullable
	}
	out.nullable = nullable
	return out, nil
}

// typeArrayExpr types ARRAY[...], whose type is an array of its elements'.
func (a *analyzer) typeArrayExpr(e *ast.A_ArrayExpr) (exprType, error) {
	elemT, err := a.typeFirstOf(listItems(e.Elements), false)
	if err != nil {
		return exprType{}, err
	}
	element, _ := a.typeNameOf(elemT)
	if element == "" {
		return exprType{}, nil
	}
	return a.namedType(element + core.ArraySuffix), nil
}

// namedType is the type a name refers to, or the name itself when the catalog
// has no such type.
func (a *analyzer) namedType(name string) exprType {
	if oid, err := a.cat.TypeOID(name); err == nil {
		return exprType{typeOID: oid}
	}
	return exprType{typeName: name}
}

// typeNameOf reports a type's name and whether it is an array of that name,
// whether the type is one the catalog holds or one only the query named.
func (a *analyzer) typeNameOf(t exprType) (string, bool) {
	name := t.typeName
	if t.typeOID != 0 {
		var err error
		if name, err = a.cat.TypeName(t.typeOID); err != nil {
			return "", false
		}
	}
	if element, ok := strings.CutSuffix(name, core.ArraySuffix); ok {
		return element, true
	}
	return name, false
}

// typeSubLink types a subquery used as an expression: EXISTS and IN yield a
// predicate, and a scalar subquery yields its first column.
func (a *analyzer) typeSubLink(e *ast.SubLink) (exprType, error) {
	if _, err := a.typeExpr(e.Testexpr); err != nil {
		return exprType{}, err
	}
	sel, ok := e.Subselect.(*ast.SelectStmt)
	if !ok {
		if e.SubLinkType == ast.EXPR_SUBLINK || e.SubLinkType == ast.ARRAY_SUBLINK {
			return exprType{nullable: true}, nil
		}
		return a.boolType(false)
	}
	cols, err := a.subqueryColumns(sel)
	if err != nil {
		return exprType{}, err
	}
	switch e.SubLinkType {
	case ast.EXPR_SUBLINK, ast.ARRAY_SUBLINK:
		if len(cols) == 0 {
			return exprType{nullable: true}, nil
		}
		// A subquery that matches no row yields NULL.
		return exprType{typeOID: cols[0].TypeOID, nullable: true}, nil
	default:
		return a.boolType(false)
	}
}

// typeComparison types IS DISTINCT FROM and its negation, which compare any two
// values and never return NULL.
func (a *analyzer) typeComparison(e *ast.A_Expr) (exprType, error) {
	leftT, err := a.typeExpr(e.Lexpr)
	if err != nil {
		return exprType{}, err
	}
	rightT, err := a.typeExpr(e.Rexpr)
	if err != nil {
		return exprType{}, err
	}
	if err := a.typeOperands(e.Rexpr, leftT); err != nil {
		return exprType{}, err
	}
	if err := a.typeOperands(e.Lexpr, rightT); err != nil {
		return exprType{}, err
	}
	return a.boolType(false)
}

// typeNullIf types NULLIF(x, y), which is x, made nullable.
func (a *analyzer) typeNullIf(e *ast.A_Expr) (exprType, error) {
	leftT, err := a.typeExpr(e.Lexpr)
	if err != nil {
		return exprType{}, err
	}
	if err := a.typeOperands(e.Rexpr, leftT); err != nil {
		return exprType{}, err
	}
	leftT.nullable = true
	return leftT, nil
}

// typeOperands types a node standing opposite one of known type, giving a bare
// placeholder that type.
func (a *analyzer) typeOperands(n ast.Node, other exprType) error {
	if pr, ok := n.(*ast.ParamRef); ok {
		if other.typeOID != 0 || other.typeName != "" {
			a.inferParam(pr.Number, other)
		}
		return nil
	}
	_, err := a.typeExpr(n)
	return err
}

// normalizeOperator upper-cases a word operator ("like", "is not") so that the
// spelling a user wrote matches the one a dialect seeded. Symbolic operators
// are already canonical.
func normalizeOperator(name string) string {
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			return strings.ToUpper(name)
		}
	}
	return name
}

func opNameFromList(l *ast.List) string {
	if l == nil {
		return ""
	}
	parts := make([]string, 0, len(l.Items))
	for _, item := range l.Items {
		if s, ok := item.(*ast.String); ok {
			parts = append(parts, s.Str)
		}
	}
	return strings.Join(parts, ".")
}

func (a *analyzer) resolveOperator(name string, leftOID, rightOID int64) (core.OperatorOverload, error) {
	candidates, err := a.cat.FindOperators(name, leftOID, rightOID)
	if err != nil {
		return core.OperatorOverload{}, err
	}
	if len(candidates) > 0 {
		return candidates[0], nil
	}

	all, err := a.cat.FindOperators(name, 0, 0)
	if err != nil {
		return core.OperatorOverload{}, err
	}
	for _, ov := range all {
		if leftOID != 0 && ov.LeftTypeOID != 0 && leftOID != ov.LeftTypeOID {
			ok, _ := a.cat.CastAllowed(leftOID, ov.LeftTypeOID, "i")
			if !ok {
				continue
			}
		}
		if rightOID != 0 && ov.RightTypeOID != 0 && rightOID != ov.RightTypeOID {
			ok, _ := a.cat.CastAllowed(rightOID, ov.RightTypeOID, "i")
			if !ok {
				continue
			}
		}
		if (leftOID == 0) != (ov.LeftTypeOID == 0) {
			continue
		}
		if (rightOID == 0) != (ov.RightTypeOID == 0) {
			continue
		}
		return ov, nil
	}

	// No dialect declares every operator over every type — extensions add
	// them, and an operator over a type the schema declared may never have
	// been registered. Rather than fail the query, assume the shape the
	// operator's name implies: a comparison yields a boolean and anything
	// else yields the type it was applied to.
	if a.cat.IsComparisonOperator(name) {
		boolOID, err := a.cat.BoolTypeOID()
		if err != nil {
			return core.OperatorOverload{}, err
		}
		return core.OperatorOverload{Name: name, ResultTypeOID: boolOID}, nil
	}
	result := leftOID
	if result == 0 {
		result = rightOID
	}
	return core.OperatorOverload{Name: name, ResultTypeOID: result}, nil
}

func (a *analyzer) typeBoolExpr(b *ast.BoolExpr) (exprType, error) {
	for _, item := range listItems(b.Args) {
		if _, err := a.typeExpr(item); err != nil {
			return exprType{}, err
		}
	}
	return a.boolType(false)
}

func (a *analyzer) typeFuncCall(f *ast.FuncCall) (exprType, error) {
	name := funcCallName(f)
	if name == "" {
		return exprType{}, fmt.Errorf("func call: missing name")
	}

	if f.AggStar && (name == "count" || name == "count.*") {
		if overloads, err := a.cat.FindProcs("count", nil); err == nil && len(overloads) > 0 {
			return exprType{typeOID: overloads[0].ReturnTypeOID, nullable: overloads[0].ReturnNullable}, nil
		}
		oid, err := a.cat.TypeOID("int8")
		if err != nil {
			return exprType{}, err
		}
		return exprType{typeOID: oid, nullable: false}, nil
	}

	args := listItems(f.Args)
	argTypes := make([]exprType, 0, len(args))
	argOIDs := make([]int64, 0, len(args))
	anyNullable := false
	for _, arg := range args {
		t, err := a.typeExpr(arg)
		if err != nil {
			return exprType{}, err
		}
		argTypes = append(argTypes, t)
		argOIDs = append(argOIDs, t.typeOID)
		anyNullable = anyNullable || t.nullable
	}

	overloads, err := a.cat.FindProcs(name, nil)
	if err != nil {
		return exprType{}, err
	}
	if len(overloads) == 0 {
		// A dialect's function list is never complete — extensions add to it,
		// and so does the user. An unknown function leaves the result untyped
		// rather than failing the query.
		return exprType{nullable: true}, nil
	}
	p := a.pickOverload(overloads, argOIDs)
	// An argument that is a bare placeholder takes the parameter's type,
	// unless the parameter is polymorphic and says nothing.
	for i, arg := range args {
		if i >= len(p.ArgTypes) {
			break
		}
		if a.isPolymorphicOID(p.ArgTypes[i]) {
			continue
		}
		if err := a.typeOperands(arg, exprType{typeOID: p.ArgTypes[i]}); err != nil {
			return exprType{}, err
		}
	}
	ret := a.returnType(p, argTypes)
	ret.nullable = p.ReturnNullable
	if !p.NeverNull && anyNullable && a.cat.PropagatesNullable() {
		ret.nullable = true
	}
	return ret, nil
}

// returnType resolves a polymorphic return type — max(anyelement), or a
// seed's "$2" for the type of the second argument — to the type the call was
// made with.
func (a *analyzer) returnType(p core.ProcOverload, argTypes []exprType) exprType {
	if p.ReturnTypeOID == 0 || len(argTypes) == 0 {
		return exprType{typeOID: p.ReturnTypeOID}
	}
	name, err := a.cat.TypeName(p.ReturnTypeOID)
	if err != nil {
		return exprType{typeOID: p.ReturnTypeOID}
	}
	if n, ok := argIndex(name); ok {
		if n < len(argTypes) {
			return exprType{typeOID: argTypes[n].typeOID, typeName: argTypes[n].typeName}
		}
		return exprType{}
	}
	if isPolymorphic(name) && argTypes[0].typeOID != 0 {
		return exprType{typeOID: argTypes[0].typeOID}
	}
	return exprType{typeOID: p.ReturnTypeOID}
}

// argIndex reads a seed's "$n" pseudo-type as the zero-based index of the
// argument whose type it stands for.
func argIndex(typeName string) (int, bool) {
	rest, ok := strings.CutPrefix(typeName, "$")
	if !ok {
		return 0, false
	}
	n := 0
	for _, r := range rest {
		if r < '0' || r > '9' {
			return 0, false
		}
		n = n*10 + int(r-'0')
	}
	if n == 0 {
		return 0, false
	}
	return n - 1, true
}

// isPolymorphicOID reports whether a parameter type accepts any argument.
func (a *analyzer) isPolymorphicOID(oid int64) bool {
	if oid == 0 {
		return true
	}
	name, err := a.cat.TypeName(oid)
	if err != nil {
		return false
	}
	if _, ok := argIndex(name); ok {
		return true
	}
	return isPolymorphic(name)
}

func isPolymorphic(typeName string) bool {
	switch typeName {
	case "any", "anyelement", "anyarray", "anynonarray", "anyenum", "anyrange",
		"anymultirange", "anycompatible", "anycompatiblearray",
		"anycompatiblenonarray", "anycompatiblerange":
		return true
	}
	return false
}

// pickOverload chooses the overload whose parameters the call's arguments
// match best: an exact type match on a parameter beats a polymorphic one,
// which beats a mismatch, and any overload of the right arity beats one of
// the wrong arity.
func (a *analyzer) pickOverload(overloads []core.ProcOverload, argTypes []int64) core.ProcOverload {
	best := -1
	bestScore := -1
	for i := range overloads {
		ov := &overloads[i]
		if len(ov.ArgTypes) != len(argTypes) {
			continue
		}
		score := 0
		for j, oid := range argTypes {
			switch {
			case oid != 0 && oid == ov.ArgTypes[j]:
				score += 2
			case a.isPolymorphicOID(ov.ArgTypes[j]):
				score += 1
			}
		}
		if score > bestScore {
			best, bestScore = i, score
		}
	}
	if best >= 0 {
		return overloads[best]
	}
	return overloads[0]
}

func funcCallName(f *ast.FuncCall) string {
	if f.Funcname != nil {
		return opNameFromList(f.Funcname)
	}
	if f.Func != nil {
		return strings.ToLower(f.Func.Name)
	}
	return ""
}

func (a *analyzer) typeTypeCast(c *ast.TypeCast) (exprType, error) {
	if c.TypeName == nil {
		return exprType{}, fmt.Errorf("cast: missing target type")
	}
	name := core.TypeNameString(c.TypeName)
	if name == "" {
		return exprType{}, fmt.Errorf("cast: missing target type")
	}
	t := a.namedType(name)
	// A cast is how a query says what an otherwise untyped placeholder holds.
	if err := a.typeOperands(c.Arg, t); err != nil {
		return exprType{}, err
	}
	return t, nil
}
