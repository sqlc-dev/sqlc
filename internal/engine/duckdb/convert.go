package duckdb

import (
	"math/big"
	"strconv"
	"strings"

	dw "github.com/sqlc-dev/darkwing/ast"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

type cc struct{}

// loc returns a node's start offset for the Location fields of the sqlc AST.
// darkwing marks synthesized nodes with a negative span.
func (c *cc) loc(n dw.Node) int {
	if n == nil {
		return 0
	}
	if pos := n.Pos(); pos >= 0 {
		return pos
	}
	return 0
}

func (c *cc) convert(node dw.Stmt) ast.Node {
	switch n := node.(type) {
	case *dw.SelectStatement:
		return c.convertSelectStatement(n)
	case *dw.InsertStatement:
		return c.convertInsertStatement(n)
	case *dw.UpdateStatement:
		return c.convertUpdateStatement(n)
	case *dw.DeleteStatement:
		return c.convertDeleteStatement(n)
	case *dw.TruncateStatement:
		return c.convertTruncateStatement(n)
	case *dw.CreateStatement:
		return c.convertCreateStatement(n)
	case *dw.DropStatement:
		return c.convertDropStatement(n)
	case *dw.AlterStatement:
		return c.convertAlterStatement(n)
	default:
		return todo(n)
	}
}

func (c *cc) convertSelectStatement(n *dw.SelectStatement) ast.Node {
	if n == nil {
		return &ast.TODO{}
	}
	return c.convertQueryNode(n.Node)
}

func (c *cc) convertQueryNode(q dw.QueryNode) ast.Node {
	switch n := q.(type) {
	case *dw.SelectNode:
		return c.convertSelectNode(n)
	case *dw.SetOperationNode:
		return c.convertSetOperationNode(n)
	case *dw.RecursiveCTENode:
		return c.convertRecursiveCTENode(n)
	default:
		return todo(q)
	}
}

// applyModifiers folds a query node's result modifiers — ORDER BY, LIMIT,
// DISTINCT — into the converted statement.
func (c *cc) applyModifiers(stmt *ast.SelectStmt, mods []dw.ResultModifier) {
	for _, mod := range mods {
		switch m := mod.(type) {
		case *dw.OrderModifier:
			stmt.SortClause = c.convertOrderBys(m.Orders)
		case *dw.LimitModifier:
			// LIMIT n% has no row-count equivalent.
			if m.LimitType == dw.LimitPercentage {
				continue
			}
			if m.Limit != nil {
				stmt.LimitCount = c.convertExpr(m.Limit)
			}
			if m.Offset != nil {
				stmt.LimitOffset = c.convertExpr(m.Offset)
			}
		case *dw.DistinctModifier:
			stmt.DistinctClause = &ast.List{}
			for _, target := range m.DistinctOnTargets {
				stmt.DistinctClause.Items = append(stmt.DistinctClause.Items, c.convertExpr(target))
			}
		}
	}
}

func (c *cc) convertWithClause(ctes dw.CTEMap) *ast.WithClause {
	if len(ctes.Entries) == 0 {
		return nil
	}
	with := &ast.WithClause{Ctes: &ast.List{}}
	for _, entry := range ctes.Entries {
		cte := entry.CTE
		name := identifier(entry.Name)
		item := &ast.CommonTableExpr{
			Ctename:  &name,
			Ctequery: c.convertQueryNode(cte.Query),
			Location: c.loc(cte),
		}
		// A recursive CTE parses into its own query node.
		if _, ok := cte.Query.(*dw.RecursiveCTENode); ok {
			with.Recursive = true
			item.Cterecursive = true
		}
		if len(cte.Aliases) > 0 {
			item.Aliascolnames = &ast.List{}
			for _, alias := range cte.Aliases {
				item.Aliascolnames.Items = append(item.Aliascolnames.Items, NewIdentifier(alias))
			}
		}
		with.Ctes.Items = append(with.Ctes.Items, item)
	}
	return with
}

// isBareStar reports a plain "*" projection with no modifiers, the select
// list darkwing synthesizes when desugaring VALUES into a select.
func isBareStar(exprs []dw.Expr) bool {
	if len(exprs) != 1 {
		return false
	}
	star, ok := exprs[0].(*dw.StarExpression)
	return ok && star.RelationName == "" && !star.Columns && star.Expr == nil &&
		len(star.ExcludeList) == 0 && len(star.QualifiedExcludeList) == 0 &&
		len(star.ReplaceList) == 0 && len(star.RenameList) == 0
}

func (c *cc) convertValuesLists(values [][]dw.Expr) *ast.List {
	lists := &ast.List{}
	for _, row := range values {
		rowList := &ast.List{}
		for _, val := range row {
			rowList.Items = append(rowList.Items, c.convertExpr(val))
		}
		lists.Items = append(lists.Items, rowList)
	}
	return lists
}

func (c *cc) convertSelectNode(n *dw.SelectNode) ast.Node {
	stmt := &ast.SelectStmt{}

	stmt.WithClause = c.convertWithClause(n.CTEs)

	// VALUES (...), (...) parses as SELECT * FROM an expression list.
	if list, ok := n.FromTable.(*dw.ExpressionListRef); ok && isBareStar(n.SelectList) {
		stmt.ValuesLists = c.convertValuesLists(list.Values)
		c.applyModifiers(stmt, n.Modifiers)
		return stmt
	}

	if len(n.SelectList) > 0 {
		stmt.TargetList = &ast.List{}
		for _, expr := range n.SelectList {
			stmt.TargetList.Items = append(stmt.TargetList.Items, c.convertResTarget(expr))
		}
	}

	if from := c.convertTableRef(n.FromTable); from != nil {
		stmt.FromClause = &ast.List{Items: []ast.Node{from}}
	}

	if n.Where != nil {
		stmt.WhereClause = c.convertExpr(n.Where)
	}

	if len(n.GroupExpressions) > 0 {
		stmt.GroupClause = &ast.List{}
		for _, expr := range n.GroupExpressions {
			stmt.GroupClause.Items = append(stmt.GroupClause.Items, c.convertExpr(expr))
		}
	}

	if n.Having != nil {
		stmt.HavingClause = c.convertExpr(n.Having)
	}

	c.applyModifiers(stmt, n.Modifiers)
	return stmt
}

func (c *cc) convertSetOperationNode(n *dw.SetOperationNode) ast.Node {
	var op ast.SetOperation
	switch n.SetOpType {
	case dw.SetOpUnion, dw.SetOpUnionByName:
		op = ast.Union
	case dw.SetOpExcept:
		op = ast.Except
	case dw.SetOpIntersect:
		op = ast.Intersect
	default:
		return todo(n)
	}

	// DuckDB v2 set operations are n-ary; sqlc's are binary, so a chain
	// folds left-associatively.
	var stmt *ast.SelectStmt
	for _, input := range n.Inputs {
		arg, ok := c.convertQueryNode(input).(*ast.SelectStmt)
		if !ok {
			return todo(n)
		}
		if stmt == nil {
			stmt = arg
			continue
		}
		stmt = &ast.SelectStmt{
			Op:   op,
			All:  n.SetOpAll,
			Larg: stmt,
			Rarg: arg,
		}
	}
	if stmt == nil {
		return todo(n)
	}
	stmt.WithClause = c.convertWithClause(n.CTEs)
	c.applyModifiers(stmt, n.Modifiers)
	return stmt
}

func (c *cc) convertRecursiveCTENode(n *dw.RecursiveCTENode) ast.Node {
	larg, ok := c.convertQueryNode(n.Left).(*ast.SelectStmt)
	if !ok {
		return todo(n)
	}
	rarg, ok := c.convertQueryNode(n.Right).(*ast.SelectStmt)
	if !ok {
		return todo(n)
	}
	stmt := &ast.SelectStmt{
		Op:   ast.Union,
		All:  n.UnionAll,
		Larg: larg,
		Rarg: rarg,
	}
	c.applyModifiers(stmt, n.Modifiers)
	return stmt
}

func (c *cc) convertResTarget(expr dw.Expr) *ast.ResTarget {
	res := &ast.ResTarget{
		Val:      c.convertExpr(expr),
		Location: c.loc(expr),
	}
	if alias := expr.GetAlias(); alias != "" {
		name := identifier(alias)
		res.Name = &name
	}
	return res
}

func (c *cc) convertOrderBys(orders []dw.OrderByNode) *ast.List {
	list := &ast.List{}
	for _, order := range orders {
		sortBy := &ast.SortBy{
			Node:     c.convertExpr(order.Expression),
			Location: c.loc(order.Expression),
		}
		switch order.Type {
		case dw.OrderAscending:
			sortBy.SortbyDir = ast.SortByDirAsc
		case dw.OrderDescending:
			sortBy.SortbyDir = ast.SortByDirDesc
		default:
			sortBy.SortbyDir = ast.SortByDirDefault
		}
		switch order.NullOrder {
		case dw.NullsFirst:
			sortBy.SortbyNulls = ast.SortByNullsFirst
		case dw.NullsLast:
			sortBy.SortbyNulls = ast.SortByNullsLast
		default:
			sortBy.SortbyNulls = ast.SortByNullsDefault
		}
		list.Items = append(list.Items, sortBy)
	}
	return list
}

func (c *cc) convertAlias(name string, columns []string) *ast.Alias {
	if name == "" && len(columns) == 0 {
		return nil
	}
	aliasName := identifier(name)
	alias := &ast.Alias{Aliasname: &aliasName}
	if len(columns) > 0 {
		alias.Colnames = &ast.List{}
		for _, col := range columns {
			alias.Colnames.Items = append(alias.Colnames.Items, NewIdentifier(col))
		}
	}
	return alias
}

func (c *cc) convertRangeVar(catalog, schema, table, alias string, loc int) *ast.RangeVar {
	name := parseTableName(catalog, schema, table)
	rv := &ast.RangeVar{
		Relname:  &name.Name,
		Location: loc,
	}
	if name.Schema != "" {
		rv.Schemaname = &name.Schema
	}
	if name.Catalog != "" {
		rv.Catalogname = &name.Catalog
	}
	if alias != "" {
		rv.Alias = c.convertAlias(alias, nil)
	}
	return rv
}

func (c *cc) convertTableRef(ref dw.TableRef) ast.Node {
	switch t := ref.(type) {
	case nil, *dw.EmptyTableRef:
		return nil
	case *dw.BaseTableRef:
		rv := c.convertRangeVar(t.CatalogName, t.SchemaName, t.TableName, "", c.loc(t))
		if alias := c.convertAlias(t.Alias, t.ColumnNameAlias); alias != nil {
			rv.Alias = alias
		}
		return rv
	case *dw.JoinRef:
		return c.convertJoinRef(t)
	case *dw.SubqueryRef:
		return &ast.RangeSubselect{
			Subquery: c.convertSelectStatement(t.Subquery),
			Alias:    c.convertAlias(t.Alias, t.ColumnNameAlias),
		}
	case *dw.TableFunctionRef:
		return c.convertTableFunctionRef(t)
	case *dw.ExpressionListRef:
		// A VALUES list used as a relation becomes a values subquery.
		sub := &ast.RangeSubselect{
			Subquery: &ast.SelectStmt{ValuesLists: c.convertValuesLists(t.Values)},
			Alias:    c.convertAlias(t.Alias, t.ExpectedNames),
		}
		return sub
	default:
		return todo(ref)
	}
}

func (c *cc) convertJoinRef(t *dw.JoinRef) ast.Node {
	join := &ast.JoinExpr{
		Larg: c.convertTableRef(t.Left),
		Rarg: c.convertTableRef(t.Right),
	}
	switch t.JoinType {
	case dw.JoinLeft:
		join.Jointype = ast.JoinTypeLeft
	case dw.JoinRight:
		join.Jointype = ast.JoinTypeRight
	case dw.JoinFull:
		join.Jointype = ast.JoinTypeFull
	default:
		join.Jointype = ast.JoinTypeInner
	}
	if t.RefType == dw.JoinRefNatural {
		join.IsNatural = true
	}
	if t.Condition != nil {
		join.Quals = c.convertExpr(t.Condition)
	}
	if len(t.UsingColumns) > 0 {
		join.UsingClause = &ast.List{}
		for _, col := range t.UsingColumns {
			join.UsingClause.Items = append(join.UsingClause.Items, NewIdentifier(col))
		}
	}
	return join
}

func (c *cc) convertTableFunctionRef(t *dw.TableFunctionRef) ast.Node {
	call, ok := c.convertExpr(t.Function).(*ast.FuncCall)
	if !ok {
		return todo(t)
	}
	return &ast.RangeFunction{
		Functions:  &ast.List{Items: []ast.Node{call}},
		Alias:      c.convertAlias(t.Alias, t.ColumnNameAlias),
		Ordinality: t.WithOrdinality,
	}
}

func (c *cc) convertExpr(expr dw.Expr) ast.Node {
	if expr == nil {
		return nil
	}
	switch e := expr.(type) {
	case *dw.ColumnRefExpression:
		return c.convertColumnRef(e)
	case *dw.ConstantExpression:
		return c.convertConstant(e)
	case *dw.ParameterExpression:
		return c.convertParameter(e)
	case *dw.FunctionExpression:
		return c.convertFunction(e)
	case *dw.OperatorExpression:
		return c.convertOperator(e)
	case *dw.ComparisonExpression:
		return c.convertComparison(e)
	case *dw.ConjunctionExpression:
		return c.convertConjunction(e)
	case *dw.CastExpression:
		return c.convertCast(e)
	case *dw.CaseExpression:
		return c.convertCase(e)
	case *dw.BetweenExpression:
		return &ast.BetweenExpr{
			Expr:     c.convertExpr(e.Input),
			Left:     c.convertExpr(e.Lower),
			Right:    c.convertExpr(e.Upper),
			Location: c.loc(e),
		}
	case *dw.SubqueryExpression:
		return c.convertSubquery(e)
	case *dw.StarExpression:
		return c.convertStar(e)
	case *dw.CollateExpression:
		return &ast.CollateExpr{
			Arg:      c.convertExpr(e.Child),
			Location: c.loc(e),
		}
	case *dw.WindowExpression:
		return c.convertWindow(e)
	default:
		return todo(expr)
	}
}

func (c *cc) convertColumnRef(e *dw.ColumnRefExpression) *ast.ColumnRef {
	fields := &ast.List{}
	for _, name := range e.ColumnNames {
		fields.Items = append(fields.Items, NewIdentifier(name))
	}
	return &ast.ColumnRef{
		Fields:   fields,
		Location: c.loc(e),
	}
}

// decimalString renders a decimal's unscaled digits with the point
// re-inserted, e.g. ("15", 1) => "1.5".
func decimalString(digits string, scale int) string {
	neg := strings.HasPrefix(digits, "-")
	if neg {
		digits = digits[1:]
	}
	if scale > 0 {
		for len(digits) <= scale {
			digits = "0" + digits
		}
		digits = digits[:len(digits)-scale] + "." + digits[len(digits)-scale:]
	}
	if neg {
		digits = "-" + digits
	}
	return digits
}

func hugeintString(h dw.Hugeint) string {
	v := new(big.Int).SetInt64(h.Upper)
	v.Lsh(v, 64)
	return v.Or(v, new(big.Int).SetUint64(h.Lower)).String()
}

func (c *cc) convertConstant(e *dw.ConstantExpression) ast.Node {
	v := e.Value
	konst := &ast.A_Const{Location: c.loc(e)}
	switch {
	case v.IsNull || v.Kind == dw.ValueNull:
		konst.Val = &ast.Null{}
	case v.Kind == dw.ValueBool:
		konst.Val = &ast.Boolean{Boolval: v.Bool}
	case v.Kind == dw.ValueInt64 && v.Type.ID == "DECIMAL":
		// A decimal literal carries its unscaled value.
		konst.Val = &ast.Float{Str: decimalString(strconv.FormatInt(v.Int64, 10), v.Type.Scale)}
	case v.Kind == dw.ValueInt64:
		konst.Val = &ast.Integer{Ival: v.Int64}
	case v.Kind == dw.ValueHugeint && v.Type.ID == "DECIMAL":
		konst.Val = &ast.Float{Str: decimalString(hugeintString(v.Hugeint), v.Type.Scale)}
	case v.Kind == dw.ValueHugeint:
		konst.Val = &ast.Float{Str: hugeintString(v.Hugeint)}
	case v.Kind == dw.ValueDouble:
		konst.Val = &ast.Float{Str: strconv.FormatFloat(v.Float64, 'g', -1, 64)}
	case v.Kind == dw.ValueString:
		konst.Val = &ast.String{Str: v.Str}
	default:
		return todo(e)
	}
	return konst
}

// convertParameter converts a prepared-statement parameter. darkwing has
// already numbered every spelling: explicit for $3 and ?3, in order of
// appearance for ?, and in order of first use for $name.
func (c *cc) convertParameter(e *dw.ParameterExpression) ast.Node {
	return &ast.ParamRef{
		Number:   e.Number,
		Dollar:   e.Kind == dw.ParameterNumbered || e.Kind == dw.ParameterNamed,
		Location: c.loc(e),
	}
}

func (c *cc) convertFunction(e *dw.FunctionExpression) ast.Node {
	name := identifier(e.FunctionName)

	// Binary and prefix operators parse as operator-named functions.
	if e.IsOperator {
		switch len(e.Arguments) {
		case 1:
			return &ast.A_Expr{
				Kind:     ast.A_Expr_Kind_OP,
				Name:     &ast.List{Items: []ast.Node{&ast.String{Str: e.FunctionName}}},
				Rexpr:    c.convertExpr(e.Arguments[0].Expr),
				Location: c.loc(e),
			}
		case 2:
			return &ast.A_Expr{
				Kind:     ast.A_Expr_Kind_OP,
				Name:     &ast.List{Items: []ast.Node{&ast.String{Str: e.FunctionName}}},
				Lexpr:    c.convertExpr(e.Arguments[0].Expr),
				Rexpr:    c.convertExpr(e.Arguments[1].Expr),
				Location: c.loc(e),
			}
		default:
			return todo(e)
		}
	}

	fc := &ast.FuncCall{
		Funcname:    &ast.List{},
		Location:    c.loc(e),
		AggDistinct: e.Distinct,
	}
	if schema := schemaName(e.Schema); schema != "" {
		fc.Funcname.Items = append(fc.Funcname.Items, &ast.String{Str: schema})
	}

	// COUNT(*) parses as the argument-less count_star.
	if name == "count_star" {
		fc.Funcname.Items = append(fc.Funcname.Items, &ast.String{Str: "count"})
		fc.AggStar = true
		return fc
	}
	fc.Funcname.Items = append(fc.Funcname.Items, &ast.String{Str: name})

	for _, arg := range e.Arguments {
		if fc.Args == nil {
			fc.Args = &ast.List{}
		}
		fc.Args.Items = append(fc.Args.Items, c.convertExpr(arg.Expr))
	}
	if e.Filter != nil {
		fc.AggFilter = c.convertExpr(e.Filter)
	}
	if len(e.OrderBys) > 0 {
		fc.AggOrder = c.convertOrderBys(e.OrderBys)
	}
	return fc
}

func (c *cc) convertOperator(e *dw.OperatorExpression) ast.Node {
	switch e.Type {
	case dw.OperatorNot:
		// NOT IN parses as NOT wrapping an IN.
		if len(e.Operands) == 1 {
			if in, ok := e.Operands[0].(*dw.OperatorExpression); ok && in.Type == dw.CompareIn {
				node := c.convertIn(in)
				node.Not = true
				return node
			}
		}
		return c.convertBoolExpr(ast.BoolExprTypeNot, e)
	case dw.OperatorIsNull, dw.OperatorIsNotNull:
		if len(e.Operands) != 1 {
			return todo(e)
		}
		test := ast.NullTestTypeIsNull
		if e.Type == dw.OperatorIsNotNull {
			test = ast.NullTestTypeIsNotNull
		}
		return &ast.NullTest{
			Arg:          c.convertExpr(e.Operands[0]),
			Nulltesttype: test,
			Location:     c.loc(e),
		}
	case dw.CompareIn, dw.CompareNotIn:
		node := c.convertIn(e)
		node.Not = e.Type == dw.CompareNotIn
		return node
	case dw.OperatorCoalesce:
		coalesce := &ast.CoalesceExpr{
			Args:     &ast.List{},
			Location: c.loc(e),
		}
		for _, operand := range e.Operands {
			coalesce.Args.Items = append(coalesce.Args.Items, c.convertExpr(operand))
		}
		return coalesce
	case dw.ArrayConstructor:
		arr := &ast.A_ArrayExpr{
			Elements: &ast.List{},
			Location: c.loc(e),
		}
		for _, operand := range e.Operands {
			arr.Elements.Items = append(arr.Elements.Items, c.convertExpr(operand))
		}
		return arr
	default:
		return todo(e)
	}
}

func (c *cc) convertIn(e *dw.OperatorExpression) *ast.In {
	in := &ast.In{Location: c.loc(e)}
	if len(e.Operands) > 0 {
		in.Expr = c.convertExpr(e.Operands[0])
		for _, item := range e.Operands[1:] {
			in.List = append(in.List, c.convertExpr(item))
		}
	}
	return in
}

func (c *cc) convertBoolExpr(op ast.BoolExprType, e *dw.OperatorExpression) ast.Node {
	expr := &ast.BoolExpr{
		Boolop:   op,
		Args:     &ast.List{},
		Location: c.loc(e),
	}
	for _, operand := range e.Operands {
		expr.Args.Items = append(expr.Args.Items, c.convertExpr(operand))
	}
	return expr
}

func comparisonOperator(t dw.ExpressionType) (string, bool) {
	switch t {
	case dw.CompareEqual:
		return "=", true
	case dw.CompareNotEqual:
		return "<>", true
	case dw.CompareLessThan:
		return "<", true
	case dw.CompareGreaterThan:
		return ">", true
	case dw.CompareLessThanOrEqual:
		return "<=", true
	case dw.CompareGreaterThanOrEqual:
		return ">=", true
	}
	return "", false
}

func (c *cc) convertComparison(e *dw.ComparisonExpression) ast.Node {
	kind := ast.A_Expr_Kind_OP
	op, ok := comparisonOperator(e.Type)
	if !ok {
		switch e.Type {
		case dw.CompareDistinctFrom:
			kind, op = ast.A_Expr_Kind_DISTINCT, "="
		case dw.CompareNotDistinctFrom:
			kind, op = ast.A_Expr_Kind_NOT_DISTINCT, "="
		default:
			return todo(e)
		}
	}
	return &ast.A_Expr{
		Kind:     kind,
		Name:     &ast.List{Items: []ast.Node{&ast.String{Str: op}}},
		Lexpr:    c.convertExpr(e.Left),
		Rexpr:    c.convertExpr(e.Right),
		Location: c.loc(e),
	}
}

func (c *cc) convertConjunction(e *dw.ConjunctionExpression) ast.Node {
	boolop := ast.BoolExprTypeAnd
	if e.Type == dw.ConjunctionOr {
		boolop = ast.BoolExprTypeOr
	}
	expr := &ast.BoolExpr{
		Boolop:   boolop,
		Args:     &ast.List{},
		Location: c.loc(e),
	}
	for _, operand := range e.Operands {
		expr.Args.Items = append(expr.Args.Items, c.convertExpr(operand))
	}
	return expr
}

func (c *cc) convertCast(e *dw.CastExpression) ast.Node {
	cast := &ast.TypeCast{
		Arg:      c.convertExpr(e.Child),
		Location: c.loc(e),
	}
	switch {
	case e.CastType != nil:
		typeName, arrayDims := c.convertTypeExpression(e.CastType)
		if arrayDims > 0 {
			typeName.ArrayBounds = &ast.List{}
			for i := 0; i < arrayDims; i++ {
				typeName.ArrayBounds.Items = append(typeName.ArrayBounds.Items, &ast.Integer{Ival: -1})
			}
		}
		cast.TypeName = typeName
	case e.ResolvedType != "":
		cast.TypeName = &ast.TypeName{Name: identifier(e.ResolvedType)}
	default:
		return todo(e)
	}
	return cast
}

func (c *cc) convertCase(e *dw.CaseExpression) ast.Node {
	caseExpr := &ast.CaseExpr{
		Args:     &ast.List{},
		Location: c.loc(e),
	}
	for _, check := range e.CaseChecks {
		caseExpr.Args.Items = append(caseExpr.Args.Items, &ast.CaseWhen{
			Expr:   c.convertExpr(check.When),
			Result: c.convertExpr(check.Then),
		})
	}
	if e.ElseExpr != nil {
		caseExpr.Defresult = c.convertExpr(e.ElseExpr)
	}
	return caseExpr
}

func (c *cc) convertSubquery(e *dw.SubqueryExpression) ast.Node {
	link := &ast.SubLink{
		Subselect: c.convertSelectStatement(e.Subquery),
		Location:  c.loc(e),
	}
	switch e.SubqueryType {
	case dw.SubqueryExists:
		link.SubLinkType = ast.EXISTS_SUBLINK
	case dw.SubqueryScalar:
		link.SubLinkType = ast.EXPR_SUBLINK
	case dw.SubqueryAny:
		// x IN (SELECT ...) and quantified comparisons desugar to ANY.
		link.SubLinkType = ast.ANY_SUBLINK
		link.Testexpr = c.convertExpr(e.Child)
		if op, ok := comparisonOperator(e.ComparisonType); ok {
			link.OperName = &ast.List{Items: []ast.Node{&ast.String{Str: op}}}
		}
	default:
		return todo(e)
	}
	return link
}

func (c *cc) convertStar(e *dw.StarExpression) ast.Node {
	// COLUMNS(...) and EXCLUDE/REPLACE/RENAME modifiers have no equivalent.
	if e.Columns || e.Expr != nil || len(e.ExcludeList) > 0 ||
		len(e.QualifiedExcludeList) > 0 || len(e.ReplaceList) > 0 || len(e.RenameList) > 0 {
		return todo(e)
	}
	fields := &ast.List{}
	if e.RelationName != "" {
		fields.Items = append(fields.Items, NewIdentifier(e.RelationName))
	}
	fields.Items = append(fields.Items, &ast.A_Star{})
	return &ast.ColumnRef{
		Fields:   fields,
		Location: c.loc(e),
	}
}

func (c *cc) convertWindow(e *dw.WindowExpression) ast.Node {
	fc := &ast.FuncCall{
		Funcname:    &ast.List{Items: []ast.Node{&ast.String{Str: identifier(e.FunctionName)}}},
		Location:    c.loc(e),
		AggDistinct: e.Distinct,
		Over:        &ast.WindowDef{Location: c.loc(e)},
	}
	for _, arg := range e.Arguments {
		if fc.Args == nil {
			fc.Args = &ast.List{}
		}
		fc.Args.Items = append(fc.Args.Items, c.convertExpr(arg.Expr))
	}
	if e.FilterExpr != nil {
		fc.AggFilter = c.convertExpr(e.FilterExpr)
	}
	if len(e.Partitions) > 0 {
		fc.Over.PartitionClause = &ast.List{}
		for _, p := range e.Partitions {
			fc.Over.PartitionClause.Items = append(fc.Over.PartitionClause.Items, c.convertExpr(p))
		}
	}
	if len(e.Orders) > 0 {
		fc.Over.OrderClause = c.convertOrderBys(e.Orders)
	}
	return fc
}

// convertTypeExpression maps an unbound DuckDB type to a sqlc type name and
// the number of list/array dimensions wrapped around it.
func (c *cc) convertTypeExpression(t *dw.TypeExpression) (*ast.TypeName, int) {
	name := identifier(t.TypeName)
	switch name {
	case "list", "array":
		// int[] is LIST(INTEGER); int[3] is ARRAY(INTEGER, 3).
		if len(t.Args) > 0 {
			if elem, ok := t.Args[0].(*dw.TypeExpression); ok {
				typeName, dims := c.convertTypeExpression(elem)
				return typeName, dims + 1
			}
		}
	}
	return &ast.TypeName{
		Schema: schemaName(t.Schema),
		Name:   name,
	}, 0
}

func (c *cc) convertReturning(returning []dw.Expr) *ast.List {
	if len(returning) == 0 {
		return nil
	}
	list := &ast.List{}
	for _, expr := range returning {
		list.Items = append(list.Items, c.convertResTarget(expr))
	}
	return list
}

func (c *cc) convertInsertStatement(n *dw.InsertStatement) ast.Node {
	stmt := &ast.InsertStmt{
		Relation:      c.convertRangeVar(n.Catalog, n.Schema, n.Table, n.TableAlias, c.loc(n)),
		WithClause:    c.convertWithClause(n.CTEs),
		DefaultValues: n.DefaultValues,
	}

	if len(n.Columns) > 0 {
		stmt.Cols = &ast.List{}
		for _, col := range n.Columns {
			name := identifier(col)
			stmt.Cols.Items = append(stmt.Cols.Items, &ast.ResTarget{Name: &name})
		}
	}

	if n.Query != nil {
		stmt.SelectStmt = c.convertSelectStatement(n.Query)
	}

	if n.OnConflict != nil {
		stmt.OnConflictClause = c.convertOnConflict(n.OnConflict)
	}

	stmt.ReturningList = c.convertReturning(n.Returning)
	return stmt
}

func (c *cc) convertOnConflict(n *dw.OnConflictInfo) *ast.OnConflictClause {
	clause := &ast.OnConflictClause{Location: c.loc(n)}
	switch n.Action {
	case dw.OnConflictUpdate:
		clause.Action = ast.OnConflictActionUpdate
	default:
		clause.Action = ast.OnConflictActionNothing
	}
	if n.SetInfo != nil {
		clause.TargetList = c.convertSetClause(n.SetInfo)
		if n.SetInfo.Condition != nil {
			clause.WhereClause = c.convertExpr(n.SetInfo.Condition)
		}
	}
	return clause
}

func (c *cc) convertSetClause(set *dw.UpdateSetInfo) *ast.List {
	list := &ast.List{}
	for i, col := range set.Columns {
		if len(col) == 0 || i >= len(set.Expressions) {
			continue
		}
		name := identifier(col[len(col)-1])
		list.Items = append(list.Items, &ast.ResTarget{
			Name: &name,
			Val:  c.convertExpr(set.Expressions[i]),
		})
	}
	return list
}

func (c *cc) convertUpdateStatement(n *dw.UpdateStatement) ast.Node {
	target := c.convertTableRef(n.Table)
	if target == nil {
		return todo(n)
	}
	stmt := &ast.UpdateStmt{
		Relations:  &ast.List{Items: []ast.Node{target}},
		TargetList: &ast.List{},
		WithClause: c.convertWithClause(n.CTEs),
	}
	if n.SetInfo != nil {
		stmt.TargetList = c.convertSetClause(n.SetInfo)
	}
	if from := c.convertTableRef(n.From); from != nil {
		stmt.FromClause = &ast.List{Items: []ast.Node{from}}
	}
	if n.Where != nil {
		stmt.WhereClause = c.convertExpr(n.Where)
	}
	stmt.ReturningList = c.convertReturning(n.Returning)
	return stmt
}

func (c *cc) convertDeleteStatement(n *dw.DeleteStatement) ast.Node {
	target := c.convertTableRef(n.Table)
	if target == nil {
		return todo(n)
	}
	stmt := &ast.DeleteStmt{
		Relations:  &ast.List{Items: []ast.Node{target}},
		WithClause: c.convertWithClause(n.CTEs),
	}
	if len(n.Using) > 0 {
		stmt.UsingClause = &ast.List{}
		for _, using := range n.Using {
			if item := c.convertTableRef(using); item != nil {
				stmt.UsingClause.Items = append(stmt.UsingClause.Items, item)
			}
		}
	}
	if n.Where != nil {
		stmt.WhereClause = c.convertExpr(n.Where)
	}
	stmt.ReturningList = c.convertReturning(n.Returning)
	return stmt
}

func (c *cc) convertTruncateStatement(n *dw.TruncateStatement) ast.Node {
	target := c.convertTableRef(n.Table)
	if target == nil {
		return todo(n)
	}
	return &ast.TruncateStmt{
		Relations: &ast.List{Items: []ast.Node{target}},
	}
}

func (c *cc) convertCreateStatement(n *dw.CreateStatement) ast.Node {
	switch info := n.Info.(type) {
	case *dw.CreateTableInfo:
		return c.convertCreateTableInfo(info)
	case *dw.CreateViewInfo:
		return c.convertCreateViewInfo(info)
	case *dw.CreateTypeInfo:
		return c.convertCreateTypeInfo(info)
	default:
		return todo(n)
	}
}

func (c *cc) convertCreateTableInfo(info *dw.CreateTableInfo) ast.Node {
	name := parseTableName(info.Catalog, info.Schema, info.Name)

	// CREATE TABLE ... AS records the table as a relation defined by its
	// query, like a view.
	if info.Query != nil {
		into := &ast.IntoClause{
			Rel: c.convertRangeVar(info.Catalog, info.Schema, info.Name, "", c.loc(info)),
		}
		if len(info.QueryAliases) > 0 {
			into.ColNames = &ast.List{}
			for _, alias := range info.QueryAliases {
				into.ColNames.Items = append(into.ColNames.Items, NewIdentifier(alias))
			}
		}
		return &ast.CreateTableAsStmt{
			Query:       c.convertSelectStatement(info.Query),
			Into:        into,
			IfNotExists: info.OnConflict == dw.CreateIgnore,
		}
	}

	stmt := &ast.CreateTableStmt{
		Name:        name,
		IfNotExists: info.OnConflict == dw.CreateIgnore,
	}

	// Column names in table-level PRIMARY KEY constraints are NOT NULL.
	primaryKey := map[string]bool{}
	for _, constraint := range info.Constraints {
		unique, ok := constraint.(*dw.UniqueConstraint)
		if !ok || !unique.IsPrimaryKey {
			continue
		}
		for _, col := range unique.Columns {
			primaryKey[identifier(col)] = true
		}
	}

	for i := range info.Columns {
		stmt.Cols = append(stmt.Cols, c.convertColumnDef(&info.Columns[i], primaryKey))
	}
	return stmt
}

func (c *cc) convertColumnDef(col *dw.ColumnDef, tablePrimaryKey map[string]bool) *ast.ColumnDef {
	name := ""
	if len(col.Names) > 0 {
		name = identifier(col.Names[len(col.Names)-1])
	}
	def := &ast.ColumnDef{
		Colname:  name,
		Location: c.loc(col),
	}
	if col.Type != nil {
		typeName, arrayDims := c.convertTypeExpression(col.Type)
		def.TypeName = typeName
		def.ArrayDims = arrayDims
		def.IsArray = arrayDims > 0
	} else {
		// A generated column may omit its type; DuckDB infers it at bind
		// time.
		def.TypeName = &ast.TypeName{Name: "any"}
	}
	for _, constraint := range col.Constraints {
		switch con := constraint.(type) {
		case *dw.NotNullConstraint:
			def.IsNotNull = !con.Null
		case *dw.UniqueConstraint:
			if con.IsPrimaryKey {
				def.PrimaryKey = true
				def.IsNotNull = true
			}
		}
	}
	if tablePrimaryKey[name] {
		def.PrimaryKey = true
		def.IsNotNull = true
	}
	return def
}

func (c *cc) convertCreateViewInfo(info *dw.CreateViewInfo) ast.Node {
	stmt := &ast.ViewStmt{
		View:    c.convertRangeVar(info.Catalog, info.Schema, info.Name, "", c.loc(info)),
		Query:   c.convertSelectStatement(info.Query),
		Replace: info.OnConflict == dw.CreateReplace,
	}
	if len(info.Aliases) > 0 {
		stmt.Aliases = &ast.List{}
		for _, alias := range info.Aliases {
			stmt.Aliases.Items = append(stmt.Aliases.Items, NewIdentifier(alias))
		}
	}
	return stmt
}

func (c *cc) convertCreateTypeInfo(info *dw.CreateTypeInfo) ast.Node {
	if !info.IsEnum {
		return todo(info)
	}
	stmt := &ast.CreateEnumStmt{
		TypeName: &ast.TypeName{
			Schema: schemaName(info.Schema),
			Name:   identifier(info.Name),
		},
		Vals: &ast.List{},
	}
	for _, val := range info.EnumValues {
		stmt.Vals.Items = append(stmt.Vals.Items, &ast.String{Str: val})
	}
	return stmt
}

func (c *cc) convertDropStatement(n *dw.DropStatement) ast.Node {
	switch n.DropType {
	case dw.DropTable, dw.DropView:
		stmt := &ast.DropTableStmt{IfExists: n.IfExists}
		for _, name := range n.Names {
			stmt.Tables = append(stmt.Tables, parseTableName(name.Catalog, name.Schema, name.Name))
		}
		return stmt
	default:
		return todo(n)
	}
}

func (c *cc) convertAlterStatement(n *dw.AlterStatement) ast.Node {
	if n.Entity != dw.AlterEntityTable {
		return todo(n)
	}
	table := parseTableName(n.Name.Catalog, n.Name.Schema, n.Name.Name)

	var stmts []ast.Node
	alter := &ast.AlterTableStmt{
		Table:     table,
		Cmds:      &ast.List{},
		MissingOk: n.IfExists,
	}
	for _, action := range n.Actions {
		switch a := action.(type) {
		case *dw.AddColumnInfo:
			def := c.convertColumnDef(&a.Column, nil)
			alter.Cmds.Items = append(alter.Cmds.Items, &ast.AlterTableCmd{
				Name:      &def.Colname,
				Subtype:   ast.AT_AddColumn,
				Def:       def,
				MissingOk: a.IfNotExists,
			})
		case *dw.DropColumnInfo:
			if len(a.Column) == 0 {
				continue
			}
			name := identifier(a.Column[len(a.Column)-1])
			alter.Cmds.Items = append(alter.Cmds.Items, &ast.AlterTableCmd{
				Name:      &name,
				Subtype:   ast.AT_DropColumn,
				MissingOk: a.IfExists,
			})
		case *dw.AlterColumnTypeInfo:
			if len(a.Column) == 0 || a.Type == nil {
				continue
			}
			name := identifier(a.Column[len(a.Column)-1])
			typeName, arrayDims := c.convertTypeExpression(a.Type)
			alter.Cmds.Items = append(alter.Cmds.Items, &ast.AlterTableCmd{
				Name:    &name,
				Subtype: ast.AT_AlterColumnType,
				Def: &ast.ColumnDef{
					Colname:   name,
					TypeName:  typeName,
					ArrayDims: arrayDims,
					IsArray:   arrayDims > 0,
				},
			})
		case *dw.SetNotNullInfo:
			if len(a.Column) == 0 {
				continue
			}
			name := identifier(a.Column[len(a.Column)-1])
			subtype := ast.AT_SetNotNull
			if a.Drop {
				subtype = ast.AT_DropNotNull
			}
			alter.Cmds.Items = append(alter.Cmds.Items, &ast.AlterTableCmd{
				Name:    &name,
				Subtype: subtype,
			})
		case *dw.RenameColumnInfo:
			if len(a.Column) == 0 {
				continue
			}
			newName := identifier(a.NewName)
			stmts = append(stmts, &ast.RenameColumnStmt{
				Table: table,
				Col: &ast.ColumnRef{
					Name: identifier(a.Column[len(a.Column)-1]),
				},
				NewName:   &newName,
				MissingOk: n.IfExists,
			})
		case *dw.RenameInfo:
			newName := identifier(a.NewName)
			stmts = append(stmts, &ast.RenameTableStmt{
				Table:     table,
				NewName:   &newName,
				MissingOk: n.IfExists,
			})
		}
	}

	if len(alter.Cmds.Items) > 0 {
		stmts = append(stmts, alter)
	}
	switch len(stmts) {
	case 0:
		return todo(n)
	case 1:
		return stmts[0]
	default:
		return &ast.List{Items: stmts}
	}
}
