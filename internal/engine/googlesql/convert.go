package googlesql

import (
	"strconv"
	"strings"

	zjast "github.com/sqlc-dev/zetajones/ast"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

type cc struct {
	paramCount int
}

func (c *cc) convert(node zjast.Node) ast.Node {
	switch n := node.(type) {
	case *zjast.QueryStatement:
		return c.convertQuery(n.Query)
	case *zjast.InsertStatement:
		return c.convertInsertStatement(n)
	case *zjast.UpdateStatement:
		return c.convertUpdateStatement(n)
	case *zjast.DeleteStatement:
		return c.convertDeleteStatement(n)
	case *zjast.CreateTableStatement:
		return c.convertCreateTableStatement(n)
	case *zjast.DropStatement:
		return c.convertDropStatement(n)
	case *zjast.TruncateStatement:
		return c.convertTruncateStatement(n)
	case *zjast.HintedStatement:
		if n.Statement == nil {
			return todo(n)
		}
		return c.convert(n.Statement)
	default:
		return todo(n)
	}
}

func (c *cc) convertQuery(n *zjast.Query) ast.Node {
	if n == nil {
		return todo(n)
	}
	// Pipe operators (|>) have no sqlc AST equivalent.
	if len(n.PipeOperators) > 0 {
		return todo(n)
	}
	node := c.convertQueryExpr(n.QueryExpr)
	stmt, ok := node.(*ast.SelectStmt)
	if !ok {
		return node
	}
	if n.WithClause != nil {
		stmt.WithClause = c.convertWithClause(n.WithClause)
	}
	if n.OrderBy != nil {
		stmt.SortClause = c.convertOrderBy(n.OrderBy)
	}
	if n.Limit != nil {
		if n.Limit.Limit != nil {
			if _, all := n.Limit.Limit.Expr.(*zjast.LimitAll); !all {
				stmt.LimitCount = c.convertExpr(n.Limit.Limit.Expr)
			}
		}
		if n.Limit.Offset != nil {
			stmt.LimitOffset = c.convertExpr(n.Limit.Offset)
		}
	}
	return stmt
}

func (c *cc) convertQueryExpr(node zjast.Node) ast.Node {
	switch n := node.(type) {
	case *zjast.Select:
		return c.convertSelect(n)
	case *zjast.SetOperation:
		return c.convertSetOperation(n)
	case *zjast.Query:
		// Parenthesized query.
		return c.convertQuery(n)
	case *zjast.AliasedQueryExpression:
		return c.convertQuery(n.Query)
	default:
		return todo(n)
	}
}

func (c *cc) convertSetOperation(n *zjast.SetOperation) ast.Node {
	if len(n.Inputs) == 0 {
		return todo(n)
	}
	var entries []*zjast.SetOperationMetadata
	if n.Metadata != nil {
		entries = n.Metadata.Entries
	}
	result, ok := c.convertQueryExpr(n.Inputs[0]).(*ast.SelectStmt)
	if !ok {
		return todo(n)
	}
	for i, input := range n.Inputs[1:] {
		rarg, ok := c.convertQueryExpr(input).(*ast.SelectStmt)
		if !ok {
			return todo(n)
		}
		op := ast.Union
		all := false
		if i < len(entries) && entries[i] != nil {
			meta := entries[i]
			if meta.OpType != nil {
				switch strings.ToUpper(meta.OpType.Op) {
				case "UNION":
					op = ast.Union
				case "INTERSECT":
					op = ast.Intersect
				case "EXCEPT":
					op = ast.Except
				}
			}
			if meta.AllOrDistinct != nil {
				all = strings.ToUpper(meta.AllOrDistinct.Value) == "ALL"
			}
		}
		result = &ast.SelectStmt{
			Op:   op,
			All:  all,
			Larg: result,
			Rarg: rarg,
		}
	}
	return result
}

func (c *cc) convertSelect(n *zjast.Select) *ast.SelectStmt {
	stmt := &ast.SelectStmt{}

	if n.SelectList != nil {
		stmt.TargetList = &ast.List{}
		for _, col := range n.SelectList.Columns {
			if col == nil {
				continue
			}
			stmt.TargetList.Items = append(stmt.TargetList.Items, c.convertSelectColumn(col))
		}
	}

	if n.From != nil {
		stmt.FromClause = &ast.List{
			Items: []ast.Node{c.convertTableExpr(n.From.TableExpression)},
		}
	}

	if n.Where != nil {
		stmt.WhereClause = c.convertExpr(n.Where.Expr)
	}

	if n.GroupBy != nil {
		stmt.GroupClause = &ast.List{}
		for _, item := range n.GroupBy.Items {
			if item == nil {
				continue
			}
			if item.Expr != nil {
				stmt.GroupClause.Items = append(stmt.GroupClause.Items, c.convertExpr(item.Expr))
			} else {
				// ROLLUP, CUBE, GROUPING SETS, or the empty grouping item.
				stmt.GroupClause.Items = append(stmt.GroupClause.Items, todo(item))
			}
		}
	}

	if n.Having != nil {
		stmt.HavingClause = c.convertExpr(n.Having.Expr)
	}

	if n.Distinct {
		stmt.DistinctClause = &ast.List{}
	}

	return stmt
}

func (c *cc) convertSelectColumn(n *zjast.SelectColumn) *ast.ResTarget {
	res := &ast.ResTarget{
		Location: n.Pos(),
	}
	if n.Alias != nil && n.Alias.Identifier != nil {
		name := identifier(n.Alias.Identifier.Name)
		res.Name = &name
	}
	res.Val = c.convertExpr(n.Expr)
	return res
}

func (c *cc) convertTableExpr(node zjast.Node) ast.Node {
	switch n := node.(type) {
	case *zjast.TablePathExpression:
		return c.convertTablePathExpression(n)
	case *zjast.TableSubquery:
		if len(n.PostfixOperators) > 0 {
			return todo(n)
		}
		return &ast.RangeSubselect{
			Lateral:  n.IsLateral,
			Subquery: c.convertQuery(n.Query),
			Alias:    convertAlias(n.Alias),
		}
	case *zjast.Join:
		return c.convertJoin(n)
	case *zjast.ParenthesizedJoin:
		if len(n.PostfixOperators) > 0 {
			return todo(n)
		}
		return c.convertTableExpr(n.Join)
	case *zjast.TVF:
		return c.convertTVF(n)
	default:
		return todo(n)
	}
}

func (c *cc) convertTablePathExpression(n *zjast.TablePathExpression) ast.Node {
	if len(n.PostfixOperators) > 0 {
		return todo(n)
	}
	if n.UnnestExpr != nil {
		return &ast.RangeFunction{
			Functions: &ast.List{
				Items: []ast.Node{c.convertUnnestExpression(n.UnnestExpr)},
			},
			Alias: convertAlias(n.Alias),
		}
	}
	rv := parseRangeVar(n.Path)
	rv.Alias = convertAlias(n.Alias)
	return rv
}

func (c *cc) convertUnnestExpression(n *zjast.UnnestExpression) *ast.FuncCall {
	fc := &ast.FuncCall{
		Func:     &ast.FuncName{Name: "unnest"},
		Funcname: &ast.List{Items: []ast.Node{&ast.String{Str: "unnest"}}},
		Args:     &ast.List{},
		Location: n.Pos(),
	}
	for _, e := range n.Expressions {
		if e == nil {
			continue
		}
		fc.Args.Items = append(fc.Args.Items, c.convertExpr(e.Expr))
	}
	return fc
}

func (c *cc) convertTVF(n *zjast.TVF) ast.Node {
	if len(n.PostfixOperators) > 0 {
		return todo(n)
	}
	parts := pathParts(n.Name)
	if len(parts) == 0 {
		return todo(n)
	}
	fc := c.newFuncCall(parts, n.Pos())
	for _, arg := range n.Args {
		if arg == nil || arg.Expr == nil {
			continue
		}
		fc.Args.Items = append(fc.Args.Items, c.convertExpr(arg.Expr))
	}
	return &ast.RangeFunction{
		Functions: &ast.List{Items: []ast.Node{fc}},
		Alias:     convertAlias(n.Alias),
	}
}

func (c *cc) convertJoin(n *zjast.Join) ast.Node {
	join := &ast.JoinExpr{
		Larg:      c.convertTableExpr(n.Lhs),
		Rarg:      c.convertTableExpr(n.Rhs),
		IsNatural: n.Natural,
	}

	switch strings.ToUpper(n.JoinType) {
	case "", "INNER", "CROSS", "COMMA":
		join.Jointype = ast.JoinTypeInner
	case "LEFT":
		join.Jointype = ast.JoinTypeLeft
	case "RIGHT":
		join.Jointype = ast.JoinTypeRight
	case "FULL":
		join.Jointype = ast.JoinTypeFull
	default:
		join.Jointype = ast.JoinTypeInner
	}

	switch clause := n.OnOrUsingClause.(type) {
	case *zjast.OnClause:
		join.Quals = c.convertExpr(clause.Expr)
	case *zjast.UsingClause:
		join.UsingClause = &ast.List{}
		for _, key := range clause.Keys {
			if key == nil {
				continue
			}
			join.UsingClause.Items = append(join.UsingClause.Items, NewIdentifier(key.Name))
		}
	}

	return join
}

func (c *cc) convertWithClause(n *zjast.WithClause) *ast.WithClause {
	wc := &ast.WithClause{
		Recursive: n.Recursive,
		Ctes:      &ast.List{},
		Location:  n.Pos(),
	}
	for _, entry := range n.Entries {
		if entry == nil || entry.AliasedQuery == nil {
			// "name() AS GROUP ROWS" entries are not supported.
			continue
		}
		aq := entry.AliasedQuery
		cte := &ast.CommonTableExpr{
			Location: aq.Pos(),
			Ctequery: c.convertQuery(aq.Query),
		}
		if aq.Identifier != nil {
			name := identifier(aq.Identifier.Name)
			cte.Ctename = &name
		}
		wc.Ctes.Items = append(wc.Ctes.Items, cte)
	}
	return wc
}

func (c *cc) convertOrderBy(n *zjast.OrderBy) *ast.List {
	list := &ast.List{}
	for _, item := range n.Items {
		if item == nil {
			continue
		}
		sortBy := &ast.SortBy{
			Node:     c.convertExpr(item.Expr),
			Location: item.Pos(),
		}
		switch {
		case item.Descending:
			sortBy.SortbyDir = ast.SortByDirDesc
		case item.HasAsc:
			sortBy.SortbyDir = ast.SortByDirAsc
		default:
			sortBy.SortbyDir = ast.SortByDirDefault
		}
		if item.NullOrder != nil {
			if item.NullOrder.NullsFirst {
				sortBy.SortbyNulls = ast.SortByNullsFirst
			} else {
				sortBy.SortbyNulls = ast.SortByNullsLast
			}
		}
		list.Items = append(list.Items, sortBy)
	}
	return list
}

func (c *cc) convertExpr(node zjast.Node) ast.Node {
	if node == nil {
		return nil
	}
	switch n := node.(type) {
	case *zjast.Identifier:
		return &ast.ColumnRef{
			Fields:   &ast.List{Items: []ast.Node{NewIdentifier(n.Name)}},
			Location: n.Pos(),
		}
	case *zjast.PathExpression:
		return c.convertPathExpression(n)
	case *zjast.DotIdentifier:
		if fields, ok := columnRefFields(n); ok {
			return &ast.ColumnRef{
				Fields:   &ast.List{Items: fields},
				Location: n.Pos(),
			}
		}
		return todo(n)
	case *zjast.IntLiteral:
		return c.convertIntLiteral(n)
	case *zjast.FloatLiteral:
		return &ast.A_Const{
			Val:      &ast.Float{Str: n.Image},
			Location: n.Pos(),
		}
	case *zjast.StringLiteral:
		return &ast.A_Const{
			Val:      &ast.String{Str: stringLiteralValue(n)},
			Location: n.Pos(),
		}
	case *zjast.BooleanLiteral:
		return &ast.A_Const{
			Val:      &ast.Boolean{Boolval: n.Value},
			Location: n.Pos(),
		}
	case *zjast.NullLiteral:
		return &ast.A_Const{
			Val:      &ast.Null{},
			Location: n.Pos(),
		}
	case *zjast.ParameterExpr:
		return c.convertParameterExpr(n)
	case *zjast.UnaryExpression:
		return c.convertUnaryExpression(n)
	case *zjast.BinaryExpression:
		return c.convertBinaryExpression(n)
	case *zjast.BitwiseShiftExpression:
		op := ">>"
		if n.IsLeftShift {
			op = "<<"
		}
		return &ast.A_Expr{
			Name:     &ast.List{Items: []ast.Node{&ast.String{Str: op}}},
			Lexpr:    c.convertExpr(n.Lhs),
			Rexpr:    c.convertExpr(n.Rhs),
			Location: n.Pos(),
		}
	case *zjast.AndExpr:
		return c.convertBoolExpr(ast.BoolExprTypeAnd, n.Conjuncts, n.Pos())
	case *zjast.OrExpr:
		return c.convertBoolExpr(ast.BoolExprTypeOr, n.Disjuncts, n.Pos())
	case *zjast.BetweenExpression:
		return &ast.BetweenExpr{
			Expr:     c.convertExpr(n.Lhs),
			Left:     c.convertExpr(n.Low),
			Right:    c.convertExpr(n.High),
			Not:      n.IsNot,
			Location: n.Pos(),
		}
	case *zjast.InExpression:
		return c.convertInExpression(n)
	case *zjast.CaseValueExpression:
		if len(n.Arguments) == 0 {
			return todo(n)
		}
		return c.buildCaseExpr(c.convertExpr(n.Arguments[0]), n.Arguments[1:], n.Pos())
	case *zjast.CaseNoValueExpression:
		return c.buildCaseExpr(nil, n.Arguments, n.Pos())
	case *zjast.CastExpression:
		return &ast.TypeCast{
			Arg:      c.convertExpr(n.Expr),
			TypeName: &ast.TypeName{Name: typeName(n.Type)},
			Location: n.Pos(),
		}
	case *zjast.ExpressionSubquery:
		return c.convertExpressionSubquery(n)
	case *zjast.FunctionCall:
		return c.convertFunctionCall(n)
	case *zjast.AnalyticFunctionCall:
		return c.convertAnalyticFunctionCall(n)
	case *zjast.Star:
		return &ast.ColumnRef{
			Fields:   &ast.List{Items: []ast.Node{&ast.A_Star{}}},
			Location: n.Pos(),
		}
	case *zjast.DotStar:
		return c.convertDotStar(n)
	case *zjast.ArrayElement:
		return &ast.A_Indirection{
			Arg: c.convertExpr(n.Array),
			Indirection: &ast.List{
				Items: []ast.Node{
					&ast.A_Indices{Uidx: c.convertExpr(n.Position)},
				},
			},
		}
	case *zjast.NamedArgument:
		return c.convertExpr(n.Value)
	case *zjast.ExpressionWithAlias:
		return c.convertExpr(n.Expression)
	default:
		return todo(n)
	}
}

func (c *cc) convertPathExpression(n *zjast.PathExpression) ast.Node {
	fields, ok := columnRefFields(n)
	if !ok {
		return todo(n)
	}
	return &ast.ColumnRef{
		Fields:   &ast.List{Items: fields},
		Location: n.Pos(),
	}
}

func (c *cc) convertIntLiteral(n *zjast.IntLiteral) ast.Node {
	// Base 0 also accepts hexadecimal literals (0x1a).
	ival, err := strconv.ParseInt(n.Image, 0, 64)
	if err != nil {
		return &ast.A_Const{
			Val:      &ast.Float{Str: n.Image},
			Location: n.Pos(),
		}
	}
	return &ast.A_Const{
		Val:      &ast.Integer{Ival: ival},
		Location: n.Pos(),
	}
}

// convertParameterExpr converts "@name" and "?" query parameters.
//
// Named parameters are mapped to the shape the PostgreSQL and SQLite engines
// produce for "@name": an A_Expr with the operator "@" and the parameter name
// as the right operand. That is the representation sqlc's named-parameter
// rewriter (internal/sql/rewrite, named.IsParamSign) already understands.
// Positional "?" parameters become ParamRef nodes numbered by their 1-based
// position in the statement.
func (c *cc) convertParameterExpr(n *zjast.ParameterExpr) ast.Node {
	if n.Name != nil {
		return &ast.A_Expr{
			Name:     &ast.List{Items: []ast.Node{&ast.String{Str: "@"}}},
			Rexpr:    &ast.String{Str: n.Name.Name},
			Location: n.Pos(),
		}
	}
	c.paramCount++
	number := n.Position
	if number == 0 {
		number = c.paramCount
	}
	return &ast.ParamRef{
		Number:   number,
		Location: n.Pos(),
	}
}

func (c *cc) convertUnaryExpression(n *zjast.UnaryExpression) ast.Node {
	if strings.ToUpper(n.Op) == "NOT" {
		return &ast.BoolExpr{
			Boolop: ast.BoolExprTypeNot,
			Args: &ast.List{
				Items: []ast.Node{c.convertExpr(n.Operand)},
			},
			Location: n.Pos(),
		}
	}
	return &ast.A_Expr{
		Name:     &ast.List{Items: []ast.Node{&ast.String{Str: n.Op}}},
		Rexpr:    c.convertExpr(n.Operand),
		Location: n.Pos(),
	}
}

func (c *cc) convertBinaryExpression(n *zjast.BinaryExpression) ast.Node {
	// "x IS NULL" and "x IS NOT NULL" become NullTest nodes.
	if strings.ToUpper(n.Op) == "IS" {
		if _, isNull := n.Right.(*zjast.NullLiteral); isNull {
			nullTest := &ast.NullTest{
				Arg:      c.convertExpr(n.Left),
				Location: n.Pos(),
			}
			if n.IsNot {
				nullTest.Nulltesttype = ast.NullTestTypeIsNotNull
			} else {
				nullTest.Nulltesttype = ast.NullTestTypeIsNull
			}
			return nullTest
		}
	}

	var expr ast.Node = &ast.A_Expr{
		Name:     &ast.List{Items: []ast.Node{&ast.String{Str: n.Op}}},
		Lexpr:    c.convertExpr(n.Left),
		Rexpr:    c.convertExpr(n.Right),
		Location: n.Pos(),
	}
	if n.IsNot {
		// "a NOT LIKE b", "a IS NOT b", etc: wrap the comparison in a NOT.
		expr = &ast.BoolExpr{
			Boolop:   ast.BoolExprTypeNot,
			Args:     &ast.List{Items: []ast.Node{expr}},
			Location: n.Pos(),
		}
	}
	return expr
}

func (c *cc) convertBoolExpr(op ast.BoolExprType, operands []zjast.Node, loc int) *ast.BoolExpr {
	expr := &ast.BoolExpr{
		Boolop:   op,
		Args:     &ast.List{},
		Location: loc,
	}
	for _, operand := range operands {
		expr.Args.Items = append(expr.Args.Items, c.convertExpr(operand))
	}
	return expr
}

func (c *cc) convertInExpression(n *zjast.InExpression) ast.Node {
	in := &ast.In{
		Expr:     c.convertExpr(n.Lhs),
		Not:      n.IsNot,
		Location: n.Pos(),
	}
	if n.List != nil {
		for _, expr := range n.List.Exprs {
			in.List = append(in.List, c.convertExpr(expr))
		}
	}
	if n.UnnestExpr != nil {
		in.List = append(in.List, c.convertUnnestExpression(n.UnnestExpr))
	}
	if n.Query != nil {
		in.Sel = c.convertQuery(n.Query)
	}
	return in
}

// buildCaseExpr builds a CaseExpr from the alternating WHEN/THEN expression
// pairs (and optional trailing ELSE expression) of a zetajones CASE node.
func (c *cc) buildCaseExpr(arg ast.Node, rest []zjast.Node, loc int) *ast.CaseExpr {
	caseExpr := &ast.CaseExpr{
		Arg:      arg,
		Args:     &ast.List{},
		Location: loc,
	}
	i := 0
	for ; i+1 < len(rest); i += 2 {
		caseExpr.Args.Items = append(caseExpr.Args.Items, &ast.CaseWhen{
			Expr:   c.convertExpr(rest[i]),
			Result: c.convertExpr(rest[i+1]),
		})
	}
	if i < len(rest) {
		caseExpr.Defresult = c.convertExpr(rest[i])
	}
	return caseExpr
}

func (c *cc) convertExpressionSubquery(n *zjast.ExpressionSubquery) *ast.SubLink {
	subLink := &ast.SubLink{
		Subselect: c.convertQuery(n.Query),
		Location:  n.Pos(),
	}
	switch strings.ToUpper(n.Modifier) {
	case "EXISTS":
		subLink.SubLinkType = ast.EXISTS_SUBLINK
	case "ARRAY":
		subLink.SubLinkType = ast.ARRAY_SUBLINK
	default:
		subLink.SubLinkType = ast.EXPR_SUBLINK
	}
	return subLink
}

// newFuncCall builds a FuncCall from a dotted function path. GoogleSQL
// function names are case-insensitive, so they are lowercased.
func (c *cc) newFuncCall(parts []string, loc int) *ast.FuncCall {
	name := strings.ToLower(parts[len(parts)-1])
	var schema string
	if len(parts) > 1 {
		schema = strings.ToLower(strings.Join(parts[:len(parts)-1], "."))
	}
	items := []ast.Node{}
	if schema != "" {
		items = append(items, &ast.String{Str: schema})
	}
	items = append(items, &ast.String{Str: name})
	return &ast.FuncCall{
		Func:     &ast.FuncName{Schema: schema, Name: name},
		Funcname: &ast.List{Items: items},
		Args:     &ast.List{},
		Location: loc,
	}
}

func (c *cc) convertFunctionCall(n *zjast.FunctionCall) ast.Node {
	// Chained calls ("expr.method(...)") have no sqlc AST equivalent.
	if n.IsChained {
		return todo(n)
	}
	parts := pathParts(n.Function)
	if len(parts) == 0 {
		return todo(n)
	}
	fc := c.newFuncCall(parts, n.Pos())
	fc.AggDistinct = n.Distinct
	for _, arg := range n.Args {
		if _, ok := arg.(*zjast.Star); ok {
			fc.AggStar = true
			continue
		}
		fc.Args.Items = append(fc.Args.Items, c.convertExpr(arg))
	}
	if n.OrderBy != nil {
		fc.AggOrder = c.convertOrderBy(n.OrderBy)
	}
	return fc
}

func (c *cc) convertAnalyticFunctionCall(n *zjast.AnalyticFunctionCall) ast.Node {
	base := c.convertExpr(n.Expr)
	fc, ok := base.(*ast.FuncCall)
	if !ok {
		return base
	}
	over := &ast.WindowDef{Location: n.Pos()}
	if spec := n.WindowSpec; spec != nil {
		if spec.Name != nil {
			name := identifier(spec.Name.Name)
			over.Name = &name
		}
		if spec.PartitionBy != nil {
			over.PartitionClause = &ast.List{}
			for _, expr := range spec.PartitionBy.Expressions {
				over.PartitionClause.Items = append(over.PartitionClause.Items, c.convertExpr(expr))
			}
		}
		if spec.OrderBy != nil {
			over.OrderClause = c.convertOrderBy(spec.OrderBy)
		}
	}
	fc.Over = over
	return fc
}

func (c *cc) convertDotStar(n *zjast.DotStar) ast.Node {
	fields, ok := columnRefFields(n.Expr)
	if !ok {
		return todo(n)
	}
	fields = append(fields, &ast.A_Star{})
	return &ast.ColumnRef{
		Fields:   &ast.List{Items: fields},
		Location: n.Pos(),
	}
}

// targetPath extracts the table path of an INSERT/UPDATE/DELETE target.
func targetPath(node zjast.Node) (*zjast.PathExpression, bool) {
	path, ok := node.(*zjast.PathExpression)
	return path, ok && path != nil
}

func (c *cc) convertInsertStatement(n *zjast.InsertStatement) ast.Node {
	path, ok := targetPath(n.Target)
	if !ok {
		return todo(n)
	}
	stmt := &ast.InsertStmt{
		Relation: parseRangeVar(path),
	}

	if n.Columns != nil {
		stmt.Cols = &ast.List{}
		for _, col := range n.Columns.Identifiers {
			if col == nil {
				continue
			}
			name := identifier(col.Name)
			stmt.Cols.Items = append(stmt.Cols.Items, &ast.ResTarget{
				Name:     &name,
				Location: col.Pos(),
			})
		}
	}

	switch {
	case n.Rows != nil:
		selectStmt := &ast.SelectStmt{
			ValuesLists: &ast.List{},
		}
		for _, row := range n.Rows.Rows {
			if row == nil {
				continue
			}
			rowList := &ast.List{}
			for _, val := range row.Values {
				rowList.Items = append(rowList.Items, c.convertExpr(val))
			}
			selectStmt.ValuesLists.Items = append(selectStmt.ValuesLists.Items, rowList)
		}
		stmt.SelectStmt = selectStmt
	case n.Query != nil:
		stmt.SelectStmt = c.convertQuery(n.Query)
	}

	if n.Returning != nil {
		stmt.ReturningList = c.convertReturningClause(n.Returning)
	}

	return stmt
}

// convertReturningClause converts a "THEN RETURN" clause to a returning list.
func (c *cc) convertReturningClause(n *zjast.ReturningClause) *ast.List {
	list := &ast.List{}
	if n.SelectList != nil {
		for _, col := range n.SelectList.Columns {
			if col == nil {
				continue
			}
			list.Items = append(list.Items, c.convertSelectColumn(col))
		}
	}
	return list
}

func (c *cc) convertUpdateStatement(n *zjast.UpdateStatement) ast.Node {
	path, ok := targetPath(n.Target)
	if !ok {
		return todo(n)
	}
	rv := parseRangeVar(path)
	rv.Alias = convertAlias(n.Alias)
	stmt := &ast.UpdateStmt{
		Relations:  &ast.List{Items: []ast.Node{rv}},
		TargetList: &ast.List{},
	}

	if n.UpdateItemList != nil {
		for _, item := range n.UpdateItemList.Items {
			if item == nil || item.SetValue == nil {
				// Nested INSERT/UPDATE/DELETE update items are unsupported.
				continue
			}
			stmt.TargetList.Items = append(stmt.TargetList.Items, c.convertUpdateSetValue(item.SetValue))
		}
	}

	if n.From != nil {
		stmt.FromClause = &ast.List{
			Items: []ast.Node{c.convertTableExpr(n.From.TableExpression)},
		}
	}

	if n.Where != nil {
		stmt.WhereClause = c.convertExpr(whereExpr(n.Where))
	}

	if n.Returning != nil {
		stmt.ReturningList = c.convertReturningClause(n.Returning)
	}

	return stmt
}

func (c *cc) convertUpdateSetValue(n *zjast.UpdateSetValue) *ast.ResTarget {
	res := &ast.ResTarget{
		Val:      c.convertExpr(n.Value),
		Location: n.Pos(),
	}
	if path, ok := n.Path.(*zjast.PathExpression); ok {
		if parts := pathParts(path); len(parts) > 0 {
			name := identifier(parts[len(parts)-1])
			res.Name = &name
		}
	}
	return res
}

func (c *cc) convertDeleteStatement(n *zjast.DeleteStatement) ast.Node {
	path, ok := targetPath(n.Target)
	if !ok {
		return todo(n)
	}
	rv := parseRangeVar(path)
	rv.Alias = convertAlias(n.Alias)
	stmt := &ast.DeleteStmt{
		Relations: &ast.List{Items: []ast.Node{rv}},
	}

	if n.Where != nil {
		stmt.WhereClause = c.convertExpr(whereExpr(n.Where))
	}

	if n.Returning != nil {
		stmt.ReturningList = c.convertReturningClause(n.Returning)
	}

	return stmt
}

func (c *cc) convertCreateTableStatement(n *zjast.CreateTableStatement) ast.Node {
	stmt := &ast.CreateTableStmt{
		Name:        parseTableName(n.Name),
		IfNotExists: n.IsIfNotExists,
	}
	if n.TableElementList != nil {
		for _, elem := range n.TableElementList.Elements {
			col, ok := elem.(*zjast.ColumnDefinition)
			if !ok {
				// Table-level PRIMARY KEY, FOREIGN KEY, and CHECK
				// constraints are skipped.
				continue
			}
			stmt.Cols = append(stmt.Cols, c.convertColumnDefinition(col))
		}
	}
	return stmt
}

func (c *cc) convertColumnDefinition(n *zjast.ColumnDefinition) *ast.ColumnDef {
	col := &ast.ColumnDef{
		Location: n.Pos(),
		TypeName: &ast.TypeName{Name: columnSchemaTypeName(n.Schema)},
	}
	if n.Name != nil {
		col.Colname = identifier(n.Name.Name)
	}

	if attrs := columnAttributes(n.Schema); attrs != nil {
		for _, attr := range attrs.Attributes {
			switch attr.(type) {
			case *zjast.NotNullColumnAttribute:
				col.IsNotNull = true
			case *zjast.PrimaryKeyColumnAttribute:
				col.PrimaryKey = true
			}
		}
	}

	if simple, ok := n.Schema.(*zjast.SimpleColumnSchema); ok {
		// Type parameters, e.g. STRING(10) or NUMERIC(10, 2).
		if simple.TypeParameters != nil {
			col.TypeName.Typmods = &ast.List{}
			for _, param := range simple.TypeParameters.Parameters {
				col.TypeName.Typmods.Items = append(col.TypeName.Typmods.Items, c.convertExpr(param))
			}
		}
		if simple.DefaultExpression != nil {
			col.RawDefault = c.convertExpr(simple.DefaultExpression)
		}
	}

	return col
}

func (c *cc) convertDropStatement(n *zjast.DropStatement) ast.Node {
	// Only DROP TABLE maps onto the sqlc AST; other object kinds (views,
	// functions, schemas, ...) fall through to TODO.
	if strings.EqualFold(n.ObjectKind, "TABLE") && n.Path != nil {
		return &ast.DropTableStmt{
			IfExists: n.IsIfExists,
			Tables:   []*ast.TableName{parseTableName(n.Path)},
		}
	}
	return todo(n)
}

func (c *cc) convertTruncateStatement(n *zjast.TruncateStatement) ast.Node {
	if n.Target == nil {
		return todo(n)
	}
	return &ast.TruncateStmt{
		Relations: &ast.List{Items: []ast.Node{parseRangeVar(n.Target)}},
	}
}
