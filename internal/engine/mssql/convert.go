package mssql

import (
	"strconv"
	"strings"

	tsql "github.com/sqlc-dev/teesql/ast"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

type cc struct {
	paramCount int
	// namedParams tracks the number assigned to each "@name" so repeated
	// uses share a single parameter.
	namedParams map[string]int
}

func (c *cc) convert(node tsql.Node) ast.Node {
	switch n := node.(type) {
	case *tsql.SelectStatement:
		return c.convertSelectStatement(n)
	case *tsql.InsertStatement:
		return c.convertInsertStatement(n)
	case *tsql.UpdateStatement:
		return c.convertUpdateStatement(n)
	case *tsql.DeleteStatement:
		return c.convertDeleteStatement(n)
	case *tsql.CreateTableStatement:
		return c.convertCreateTableStatement(n)
	case *tsql.DropTableStatement:
		return c.convertDropTableStatement(n)
	case *tsql.AlterTableAddTableElementStatement:
		return c.convertAlterTableAddTableElementStatement(n)
	case *tsql.AlterTableDropTableElementStatement:
		return c.convertAlterTableDropTableElementStatement(n)
	default:
		return todo(n)
	}
}

func (c *cc) convertSelectStatement(n *tsql.SelectStatement) ast.Node {
	stmt := c.convertQueryExpression(n.QueryExpression)
	sel, ok := stmt.(*ast.SelectStmt)
	if !ok {
		return stmt
	}
	if with := c.convertWithClause(n.WithCtesAndXmlNamespaces); with != nil {
		sel.WithClause = with
	}
	return sel
}

func (c *cc) convertWithClause(n *tsql.WithCtesAndXmlNamespaces) *ast.WithClause {
	if n == nil || len(n.CommonTableExpressions) == 0 {
		return nil
	}
	with := &ast.WithClause{
		Ctes:     &ast.List{},
		Location: location(n),
	}
	for _, cte := range n.CommonTableExpressions {
		name := identifierValue(cte.ExpressionName)
		item := &ast.CommonTableExpr{
			Ctename:  &name,
			Ctequery: c.convertQueryExpression(cte.QueryExpression),
			Location: location(cte),
		}
		if len(cte.Columns) > 0 {
			item.Aliascolnames = &ast.List{}
			for _, col := range cte.Columns {
				item.Aliascolnames.Items = append(item.Aliascolnames.Items, NewIdentifier(col.Value))
			}
		}
		with.Ctes.Items = append(with.Ctes.Items, item)
	}
	return with
}

func (c *cc) convertQueryExpression(q tsql.QueryExpression) ast.Node {
	switch n := q.(type) {
	case *tsql.QuerySpecification:
		return c.convertQuerySpecification(n)
	case *tsql.BinaryQueryExpression:
		return c.convertBinaryQueryExpression(n)
	case *tsql.QueryParenthesisExpression:
		return c.convertQueryExpression(n.QueryExpression)
	default:
		return todo(q)
	}
}

func (c *cc) convertBinaryQueryExpression(n *tsql.BinaryQueryExpression) ast.Node {
	var op ast.SetOperation
	switch n.BinaryQueryExpressionType {
	case "Union":
		op = ast.Union
	case "Except":
		op = ast.Except
	case "Intersect":
		op = ast.Intersect
	default:
		return todo(n)
	}
	larg, ok := c.convertQueryExpression(n.FirstQueryExpression).(*ast.SelectStmt)
	if !ok {
		return todo(n)
	}
	rarg, ok := c.convertQueryExpression(n.SecondQueryExpression).(*ast.SelectStmt)
	if !ok {
		return todo(n)
	}
	stmt := &ast.SelectStmt{
		Op:   op,
		All:  n.All,
		Larg: larg,
		Rarg: rarg,
	}
	if n.OrderByClause != nil {
		stmt.SortClause = c.convertOrderByClause(n.OrderByClause)
	}
	return stmt
}

func (c *cc) convertQuerySpecification(n *tsql.QuerySpecification) *ast.SelectStmt {
	stmt := &ast.SelectStmt{}

	if len(n.SelectElements) > 0 {
		stmt.TargetList = &ast.List{}
		for _, elem := range n.SelectElements {
			if target := c.convertSelectElement(elem); target != nil {
				stmt.TargetList.Items = append(stmt.TargetList.Items, target)
			}
		}
	}

	if n.FromClause != nil {
		stmt.FromClause = &ast.List{}
		for _, ref := range n.FromClause.TableReferences {
			stmt.FromClause.Items = append(stmt.FromClause.Items, c.convertTableReference(ref))
		}
	}

	if n.WhereClause != nil {
		stmt.WhereClause = c.convertBooleanExpression(n.WhereClause.SearchCondition)
	}

	if n.GroupByClause != nil {
		stmt.GroupClause = &ast.List{}
		for _, spec := range n.GroupByClause.GroupingSpecifications {
			if g, ok := spec.(*tsql.ExpressionGroupingSpecification); ok {
				stmt.GroupClause.Items = append(stmt.GroupClause.Items, c.convertScalarExpression(g.Expression))
			}
		}
	}

	if n.HavingClause != nil {
		stmt.HavingClause = c.convertBooleanExpression(n.HavingClause.SearchCondition)
	}

	if n.OrderByClause != nil {
		stmt.SortClause = c.convertOrderByClause(n.OrderByClause)
	}

	// OFFSET n ROWS FETCH NEXT m ROWS ONLY
	if n.OffsetClause != nil {
		if n.OffsetClause.OffsetExpression != nil {
			stmt.LimitOffset = c.convertScalarExpression(n.OffsetClause.OffsetExpression)
		}
		if n.OffsetClause.FetchExpression != nil {
			stmt.LimitCount = c.convertScalarExpression(n.OffsetClause.FetchExpression)
		}
	}

	// TOP n; TOP n PERCENT does not translate to a row count.
	if n.TopRowFilter != nil && !n.TopRowFilter.Percent {
		stmt.LimitCount = c.convertScalarExpression(n.TopRowFilter.Expression)
	}

	if n.UniqueRowFilter == "Distinct" {
		stmt.DistinctClause = &ast.List{}
	}

	return stmt
}

func (c *cc) convertSelectElement(elem tsql.SelectElement) ast.Node {
	switch e := elem.(type) {
	case *tsql.SelectScalarExpression:
		res := &ast.ResTarget{
			Val:      c.convertScalarExpression(e.Expression),
			Location: location(e),
		}
		if name := columnAlias(e.ColumnName); name != "" {
			res.Name = &name
		}
		return res
	case *tsql.SelectStarExpression:
		fields := &ast.List{}
		if e.Qualifier != nil {
			for _, id := range e.Qualifier.Identifiers {
				fields.Items = append(fields.Items, NewIdentifier(id.Value))
			}
		}
		fields.Items = append(fields.Items, &ast.A_Star{})
		return &ast.ResTarget{
			Val: &ast.ColumnRef{
				Fields:   fields,
				Location: location(e),
			},
			Location: location(e),
		}
	default:
		return &ast.ResTarget{
			Val:      todo(elem),
			Location: location(elem),
		}
	}
}

func columnAlias(n *tsql.IdentifierOrValueExpression) string {
	if n == nil {
		return ""
	}
	if n.Identifier != nil {
		return identifierValue(n.Identifier)
	}
	if n.Value != "" {
		return identifier(n.Value)
	}
	return ""
}

func (c *cc) convertOrderByClause(n *tsql.OrderByClause) *ast.List {
	list := &ast.List{}
	for _, elem := range n.OrderByElements {
		sortBy := &ast.SortBy{
			Node:     c.convertScalarExpression(elem.Expression),
			Location: location(elem),
		}
		switch elem.SortOrder {
		case "Descending":
			sortBy.SortbyDir = ast.SortByDirDesc
		case "Ascending":
			sortBy.SortbyDir = ast.SortByDirAsc
		default:
			sortBy.SortbyDir = ast.SortByDirDefault
		}
		list.Items = append(list.Items, sortBy)
	}
	return list
}

func (c *cc) convertTableReference(ref tsql.TableReference) ast.Node {
	switch t := ref.(type) {
	case *tsql.NamedTableReference:
		rv := parseRangeVar(t.SchemaObject)
		if t.Alias != nil {
			alias := identifierValue(t.Alias)
			rv.Alias = &ast.Alias{Aliasname: &alias}
		}
		return rv
	case *tsql.QueryDerivedTable:
		sub := &ast.RangeSubselect{
			Subquery: c.convertQueryExpression(t.QueryExpression),
		}
		if t.Alias != nil {
			alias := identifierValue(t.Alias)
			sub.Alias = &ast.Alias{Aliasname: &alias}
		}
		return sub
	case *tsql.QualifiedJoin:
		join := &ast.JoinExpr{
			Larg: c.convertTableReference(t.FirstTableReference),
			Rarg: c.convertTableReference(t.SecondTableReference),
		}
		switch t.QualifiedJoinType {
		case "LeftOuter":
			join.Jointype = ast.JoinTypeLeft
		case "RightOuter":
			join.Jointype = ast.JoinTypeRight
		case "FullOuter":
			join.Jointype = ast.JoinTypeFull
		default:
			join.Jointype = ast.JoinTypeInner
		}
		if t.SearchCondition != nil {
			join.Quals = c.convertBooleanExpression(t.SearchCondition)
		}
		return join
	case *tsql.UnqualifiedJoin:
		if t.UnqualifiedJoinType != "CrossJoin" {
			return todo(ref)
		}
		return &ast.JoinExpr{
			Jointype: ast.JoinTypeInner,
			Larg:     c.convertTableReference(t.FirstTableReference),
			Rarg:     c.convertTableReference(t.SecondTableReference),
		}
	default:
		return todo(ref)
	}
}

func (c *cc) convertBooleanExpression(expr tsql.BooleanExpression) ast.Node {
	if expr == nil {
		return nil
	}
	switch e := expr.(type) {
	case *tsql.BooleanComparisonExpression:
		return &ast.A_Expr{
			Kind: ast.A_Expr_Kind_OP,
			Name: &ast.List{
				Items: []ast.Node{&ast.String{Str: comparisonOperator(e.ComparisonType)}},
			},
			Lexpr:    c.convertScalarExpression(e.FirstExpression),
			Rexpr:    c.convertScalarExpression(e.SecondExpression),
			Location: location(e),
		}
	case *tsql.BooleanBinaryExpression:
		boolop := ast.BoolExprTypeAnd
		if e.BinaryExpressionType == "Or" {
			boolop = ast.BoolExprTypeOr
		}
		return &ast.BoolExpr{
			Boolop: boolop,
			Args: &ast.List{
				Items: []ast.Node{
					c.convertBooleanExpression(e.FirstExpression),
					c.convertBooleanExpression(e.SecondExpression),
				},
			},
			Location: location(e),
		}
	case *tsql.BooleanNotExpression:
		return &ast.BoolExpr{
			Boolop: ast.BoolExprTypeNot,
			Args: &ast.List{
				Items: []ast.Node{c.convertBooleanExpression(e.Expression)},
			},
			Location: location(e),
		}
	case *tsql.BooleanParenthesisExpression:
		return c.convertBooleanExpression(e.Expression)
	case *tsql.BooleanIsNullExpression:
		test := ast.NullTestTypeIsNull
		if e.IsNot {
			test = ast.NullTestTypeIsNotNull
		}
		return &ast.NullTest{
			Arg:          c.convertScalarExpression(e.Expression),
			Nulltesttype: test,
			Location:     location(e),
		}
	case *tsql.BooleanTernaryExpression:
		return &ast.BetweenExpr{
			Expr:     c.convertScalarExpression(e.FirstExpression),
			Left:     c.convertScalarExpression(e.SecondExpression),
			Right:    c.convertScalarExpression(e.ThirdExpression),
			Not:      e.TernaryExpressionType == "NotBetween",
			Location: location(e),
		}
	case *tsql.BooleanInExpression:
		in := &ast.In{
			Expr:     c.convertScalarExpression(e.Expression),
			Not:      e.NotDefined,
			Location: location(e),
		}
		for _, item := range e.Values {
			in.List = append(in.List, c.convertScalarExpression(item))
		}
		if e.Subquery != nil {
			in.Sel = c.convertQueryExpression(e.Subquery)
		}
		return in
	case *tsql.BooleanLikeExpression:
		op := "~~"
		if e.NotDefined {
			op = "!~~"
		}
		return &ast.A_Expr{
			Kind: ast.A_Expr_Kind_OP,
			Name: &ast.List{
				Items: []ast.Node{&ast.String{Str: op}},
			},
			Lexpr:    c.convertScalarExpression(e.FirstExpression),
			Rexpr:    c.convertScalarExpression(e.SecondExpression),
			Location: location(e),
		}
	case *tsql.ExistsPredicate:
		return &ast.SubLink{
			SubLinkType: ast.EXISTS_SUBLINK,
			Subselect:   c.convertQueryExpression(e.Subquery),
			Location:    location(e),
		}
	case *tsql.BooleanScalarPlaceholder:
		return c.convertScalarExpression(e.Scalar)
	default:
		return todo(expr)
	}
}

func comparisonOperator(comparisonType string) string {
	switch comparisonType {
	case "Equals":
		return "="
	case "NotEqualToBrackets", "NotEqualToExclamation":
		return "<>"
	case "GreaterThan":
		return ">"
	case "LessThan":
		return "<"
	case "GreaterThanOrEqualTo", "NotLessThan":
		return ">="
	case "LessThanOrEqualTo", "NotGreaterThan":
		return "<="
	default:
		return "="
	}
}

func (c *cc) convertScalarExpression(expr tsql.ScalarExpression) ast.Node {
	if expr == nil {
		return nil
	}
	switch e := expr.(type) {
	case *tsql.ColumnReferenceExpression:
		return c.convertColumnReference(e)
	case *tsql.VariableReference:
		return c.convertVariableReference(e)
	case *tsql.IntegerLiteral:
		ival, _ := strconv.ParseInt(e.Value, 10, 64)
		return &ast.A_Const{
			Val:      &ast.Integer{Ival: ival},
			Location: location(e),
		}
	case *tsql.NumericLiteral:
		return &ast.A_Const{
			Val:      &ast.Float{Str: e.Value},
			Location: location(e),
		}
	case *tsql.RealLiteral:
		return &ast.A_Const{
			Val:      &ast.Float{Str: e.Value},
			Location: location(e),
		}
	case *tsql.MoneyLiteral:
		return &ast.A_Const{
			Val:      &ast.Float{Str: e.Value},
			Location: location(e),
		}
	case *tsql.StringLiteral:
		return &ast.A_Const{
			Val:      &ast.String{Str: e.Value},
			Location: location(e),
		}
	case *tsql.NullLiteral:
		return &ast.A_Const{
			Val:      &ast.Null{},
			Location: location(e),
		}
	case *tsql.BinaryExpression:
		return &ast.A_Expr{
			Kind: ast.A_Expr_Kind_OP,
			Name: &ast.List{
				Items: []ast.Node{&ast.String{Str: binaryOperator(e.BinaryExpressionType)}},
			},
			Lexpr:    c.convertScalarExpression(e.FirstExpression),
			Rexpr:    c.convertScalarExpression(e.SecondExpression),
			Location: location(e),
		}
	case *tsql.UnaryExpression:
		switch e.UnaryExpressionType {
		case "Positive":
			return c.convertScalarExpression(e.Expression)
		case "BitwiseNot":
			return &ast.A_Expr{
				Kind: ast.A_Expr_Kind_OP,
				Name: &ast.List{
					Items: []ast.Node{&ast.String{Str: "~"}},
				},
				Rexpr:    c.convertScalarExpression(e.Expression),
				Location: location(e),
			}
		default: // Negative
			return &ast.A_Expr{
				Kind: ast.A_Expr_Kind_OP,
				Name: &ast.List{
					Items: []ast.Node{&ast.String{Str: "-"}},
				},
				Rexpr:    c.convertScalarExpression(e.Expression),
				Location: location(e),
			}
		}
	case *tsql.ParenthesisExpression:
		return c.convertScalarExpression(e.Expression)
	case *tsql.FunctionCall:
		return c.convertFunctionCall(e)
	case *tsql.CastCall:
		return &ast.TypeCast{
			Arg:      c.convertScalarExpression(e.Parameter),
			TypeName: &ast.TypeName{Name: dataTypeName(e.DataType)},
			Location: location(e),
		}
	case *tsql.TryCastCall:
		return &ast.TypeCast{
			Arg:      c.convertScalarExpression(e.Parameter),
			TypeName: &ast.TypeName{Name: dataTypeName(e.DataType)},
			Location: location(e),
		}
	case *tsql.ConvertCall:
		return &ast.TypeCast{
			Arg:      c.convertScalarExpression(e.Parameter),
			TypeName: &ast.TypeName{Name: dataTypeName(e.DataType)},
			Location: location(e),
		}
	case *tsql.TryConvertCall:
		return &ast.TypeCast{
			Arg:      c.convertScalarExpression(e.Parameter),
			TypeName: &ast.TypeName{Name: dataTypeName(e.DataType)},
			Location: location(e),
		}
	case *tsql.CoalesceExpression:
		coalesce := &ast.CoalesceExpr{
			Args:     &ast.List{},
			Location: location(e),
		}
		for _, arg := range e.Expressions {
			coalesce.Args.Items = append(coalesce.Args.Items, c.convertScalarExpression(arg))
		}
		return coalesce
	case *tsql.NullIfExpression:
		return &ast.FuncCall{
			Funcname: &ast.List{
				Items: []ast.Node{&ast.String{Str: "nullif"}},
			},
			Args: &ast.List{
				Items: []ast.Node{
					c.convertScalarExpression(e.FirstExpression),
					c.convertScalarExpression(e.SecondExpression),
				},
			},
			Location: location(e),
		}
	case *tsql.IIfCall:
		return &ast.CaseExpr{
			Args: &ast.List{
				Items: []ast.Node{
					&ast.CaseWhen{
						Expr:   c.convertBooleanExpression(e.Predicate),
						Result: c.convertScalarExpression(e.ThenExpression),
					},
				},
			},
			Defresult: c.convertScalarExpression(e.ElseExpression),
			Location:  location(e),
		}
	case *tsql.SearchedCaseExpression:
		caseExpr := &ast.CaseExpr{
			Args:     &ast.List{},
			Location: location(e),
		}
		for _, when := range e.WhenClauses {
			caseExpr.Args.Items = append(caseExpr.Args.Items, &ast.CaseWhen{
				Expr:   c.convertBooleanExpression(when.WhenExpression),
				Result: c.convertScalarExpression(when.ThenExpression),
			})
		}
		if e.ElseExpression != nil {
			caseExpr.Defresult = c.convertScalarExpression(e.ElseExpression)
		}
		return caseExpr
	case *tsql.SimpleCaseExpression:
		caseExpr := &ast.CaseExpr{
			Arg:      c.convertScalarExpression(e.InputExpression),
			Args:     &ast.List{},
			Location: location(e),
		}
		for _, when := range e.WhenClauses {
			caseExpr.Args.Items = append(caseExpr.Args.Items, &ast.CaseWhen{
				Expr:   c.convertScalarExpression(when.WhenExpression),
				Result: c.convertScalarExpression(when.ThenExpression),
			})
		}
		if e.ElseExpression != nil {
			caseExpr.Defresult = c.convertScalarExpression(e.ElseExpression)
		}
		return caseExpr
	case *tsql.ScalarSubquery:
		return &ast.SubLink{
			SubLinkType: ast.EXPR_SUBLINK,
			Subselect:   c.convertQueryExpression(e.QueryExpression),
			Location:    location(e),
		}
	case *tsql.ParameterlessCall:
		return &ast.FuncCall{
			Funcname: &ast.List{
				Items: []ast.Node{&ast.String{Str: identifier(e.ParameterlessCallType)}},
			},
			Location: location(e),
		}
	default:
		return todo(expr)
	}
}

func binaryOperator(binaryExpressionType string) string {
	switch binaryExpressionType {
	case "Add":
		return "+"
	case "Subtract":
		return "-"
	case "Multiply":
		return "*"
	case "Divide":
		return "/"
	case "Modulo":
		return "%"
	case "BitwiseAnd":
		return "&"
	case "BitwiseOr":
		return "|"
	case "BitwiseXor":
		return "^"
	default:
		return "+"
	}
}

func (c *cc) convertColumnReference(n *tsql.ColumnReferenceExpression) *ast.ColumnRef {
	fields := &ast.List{}
	if n.MultiPartIdentifier != nil {
		for _, id := range n.MultiPartIdentifier.Identifiers {
			fields.Items = append(fields.Items, NewIdentifier(id.Value))
		}
	}
	if n.ColumnType == "Wildcard" {
		fields.Items = append(fields.Items, &ast.A_Star{})
	}
	return &ast.ColumnRef{
		Fields:   fields,
		Location: location(n),
	}
}

// convertVariableReference converts an "@name" query parameter. Repeated uses
// of a name share a single parameter number.
func (c *cc) convertVariableReference(n *tsql.VariableReference) ast.Node {
	name := strings.TrimPrefix(n.Name, "@")
	number, ok := c.namedParams[name]
	if !ok {
		c.paramCount++
		number = c.paramCount
		if c.namedParams == nil {
			c.namedParams = map[string]int{}
		}
		c.namedParams[name] = number
	}
	return &ast.ParamRef{
		Number:   number,
		Location: location(n),
	}
}

func (c *cc) convertFunctionCall(n *tsql.FunctionCall) *ast.FuncCall {
	fc := &ast.FuncCall{
		Funcname: &ast.List{
			Items: []ast.Node{&ast.String{Str: identifierValue(n.FunctionName)}},
		},
		Location:    location(n),
		AggDistinct: n.UniqueRowFilter == "Distinct",
	}

	for _, param := range n.Parameters {
		// COUNT(*) and friends carry the star as a wildcard column reference.
		if col, ok := param.(*tsql.ColumnReferenceExpression); ok && col.ColumnType == "Wildcard" && col.MultiPartIdentifier == nil {
			fc.AggStar = true
			continue
		}
		if fc.Args == nil {
			fc.Args = &ast.List{}
		}
		fc.Args.Items = append(fc.Args.Items, c.convertScalarExpression(param))
	}

	if n.OverClause != nil {
		fc.Over = &ast.WindowDef{Location: location(n.OverClause)}
		if len(n.OverClause.Partitions) > 0 {
			fc.Over.PartitionClause = &ast.List{}
			for _, p := range n.OverClause.Partitions {
				fc.Over.PartitionClause.Items = append(fc.Over.PartitionClause.Items, c.convertScalarExpression(p))
			}
		}
		if n.OverClause.OrderByClause != nil {
			fc.Over.OrderClause = c.convertOrderByClause(n.OverClause.OrderByClause)
		}
	}

	return fc
}

// convertOutputClause converts an OUTPUT clause to a returning list. The
// INSERTED and DELETED pseudo-table qualifiers are stripped so the columns
// resolve against the statement's target table.
func (c *cc) convertOutputClause(n *tsql.OutputClause) *ast.List {
	if n == nil || len(n.SelectColumns) == 0 {
		return nil
	}
	list := &ast.List{}
	for _, elem := range n.SelectColumns {
		target := c.convertSelectElement(elem)
		if res, ok := target.(*ast.ResTarget); ok {
			stripOutputQualifier(res.Val)
		}
		list.Items = append(list.Items, target)
	}
	return list
}

func stripOutputQualifier(node ast.Node) {
	ref, ok := node.(*ast.ColumnRef)
	if !ok || ref.Fields == nil || len(ref.Fields.Items) < 2 {
		return
	}
	if first, ok := ref.Fields.Items[0].(*ast.String); ok {
		switch first.Str {
		case "inserted", "deleted":
			ref.Fields.Items = ref.Fields.Items[1:]
		}
	}
}

func (c *cc) convertInsertStatement(n *tsql.InsertStatement) ast.Node {
	if n.InsertSpecification == nil {
		return todo(n)
	}
	spec := n.InsertSpecification

	target, ok := spec.Target.(*tsql.NamedTableReference)
	if !ok {
		return todo(n)
	}

	stmt := &ast.InsertStmt{
		Relation: parseRangeVar(target.SchemaObject),
	}

	if len(spec.Columns) > 0 {
		stmt.Cols = &ast.List{}
		for _, col := range spec.Columns {
			if col.MultiPartIdentifier == nil || len(col.MultiPartIdentifier.Identifiers) == 0 {
				continue
			}
			ids := col.MultiPartIdentifier.Identifiers
			name := identifierValue(ids[len(ids)-1])
			stmt.Cols.Items = append(stmt.Cols.Items, &ast.ResTarget{
				Name:     &name,
				Location: location(col),
			})
		}
	}

	switch src := spec.InsertSource.(type) {
	case *tsql.ValuesInsertSource:
		if !src.IsDefaultValues {
			sel := &ast.SelectStmt{
				ValuesLists: &ast.List{},
			}
			for _, row := range src.RowValues {
				rowList := &ast.List{}
				for _, val := range row.ColumnValues {
					rowList.Items = append(rowList.Items, c.convertScalarExpression(val))
				}
				sel.ValuesLists.Items = append(sel.ValuesLists.Items, rowList)
			}
			stmt.SelectStmt = sel
		}
	case *tsql.SelectInsertSource:
		stmt.SelectStmt = c.convertQueryExpression(src.Select)
	}

	if returning := c.convertOutputClause(spec.OutputClause); returning != nil {
		stmt.ReturningList = returning
	}

	return stmt
}

func (c *cc) convertUpdateStatement(n *tsql.UpdateStatement) ast.Node {
	if n.UpdateSpecification == nil {
		return todo(n)
	}
	spec := n.UpdateSpecification

	target, ok := spec.Target.(*tsql.NamedTableReference)
	if !ok {
		return todo(n)
	}

	stmt := &ast.UpdateStmt{
		Relations: &ast.List{
			Items: []ast.Node{parseRangeVar(target.SchemaObject)},
		},
		TargetList:    &ast.List{},
		WhereClause:   nil,
		FromClause:    &ast.List{},
		ReturningList: &ast.List{},
		WithClause:    c.convertWithClause(n.WithCtesAndXmlNamespaces),
	}

	for _, sc := range spec.SetClauses {
		assign, ok := sc.(*tsql.AssignmentSetClause)
		if !ok || assign.Column == nil {
			continue
		}
		ids := assign.Column.MultiPartIdentifier
		if ids == nil || len(ids.Identifiers) == 0 {
			continue
		}
		name := identifierValue(ids.Identifiers[len(ids.Identifiers)-1])
		stmt.TargetList.Items = append(stmt.TargetList.Items, &ast.ResTarget{
			Name:     &name,
			Val:      c.convertScalarExpression(assign.NewValue),
			Location: location(assign),
		})
	}

	if spec.FromClause != nil {
		for _, ref := range spec.FromClause.TableReferences {
			stmt.FromClause.Items = append(stmt.FromClause.Items, c.convertTableReference(ref))
		}
	}

	if spec.WhereClause != nil {
		stmt.WhereClause = c.convertBooleanExpression(spec.WhereClause.SearchCondition)
	}

	if returning := c.convertOutputClause(spec.OutputClause); returning != nil {
		stmt.ReturningList = returning
	}

	return stmt
}

func (c *cc) convertDeleteStatement(n *tsql.DeleteStatement) ast.Node {
	if n.DeleteSpecification == nil {
		return todo(n)
	}
	spec := n.DeleteSpecification

	target, ok := spec.Target.(*tsql.NamedTableReference)
	if !ok {
		return todo(n)
	}

	stmt := &ast.DeleteStmt{
		Relations: &ast.List{
			Items: []ast.Node{parseRangeVar(target.SchemaObject)},
		},
		ReturningList: &ast.List{},
		WithClause:    c.convertWithClause(n.WithCtesAndXmlNamespaces),
	}

	if spec.WhereClause != nil {
		stmt.WhereClause = c.convertBooleanExpression(spec.WhereClause.SearchCondition)
	}

	if returning := c.convertOutputClause(spec.OutputClause); returning != nil {
		stmt.ReturningList = returning
	}

	return stmt
}

func (c *cc) convertCreateTableStatement(n *tsql.CreateTableStatement) ast.Node {
	if n.SchemaObjectName == nil || n.Definition == nil {
		return todo(n)
	}

	stmt := &ast.CreateTableStmt{
		Name: parseTableName(n.SchemaObjectName),
	}

	// Column names in table-level PRIMARY KEY constraints are NOT NULL.
	primaryKey := map[string]bool{}
	for _, constraint := range n.Definition.TableConstraints {
		unique, ok := constraint.(*tsql.UniqueConstraintDefinition)
		if !ok || !unique.IsPrimaryKey {
			continue
		}
		for _, col := range unique.Columns {
			if col.Column == nil || col.Column.MultiPartIdentifier == nil {
				continue
			}
			ids := col.Column.MultiPartIdentifier.Identifiers
			if len(ids) == 0 {
				continue
			}
			primaryKey[identifierValue(ids[len(ids)-1])] = true
		}
	}

	for _, col := range n.Definition.ColumnDefinitions {
		stmt.Cols = append(stmt.Cols, c.convertColumnDefinition(col, primaryKey))
	}

	return stmt
}

func (c *cc) convertColumnDefinition(n *tsql.ColumnDefinition, tablePrimaryKey map[string]bool) *ast.ColumnDef {
	name := identifierValue(n.ColumnIdentifier)
	colDef := &ast.ColumnDef{
		Colname:  name,
		TypeName: &ast.TypeName{Name: dataTypeName(n.DataType)},
		Location: location(n),
	}

	// T-SQL columns are nullable unless declared otherwise.
	if n.Nullable != nil && !n.Nullable.Nullable {
		colDef.IsNotNull = true
	}
	// IDENTITY columns are always NOT NULL.
	if n.IdentityOptions != nil {
		colDef.IsNotNull = true
	}
	for _, constraint := range n.Constraints {
		switch cons := constraint.(type) {
		case *tsql.NullableConstraintDefinition:
			colDef.IsNotNull = !cons.Nullable
		case *tsql.UniqueConstraintDefinition:
			if cons.IsPrimaryKey {
				colDef.PrimaryKey = true
				colDef.IsNotNull = true
			}
		}
	}
	if tablePrimaryKey[name] {
		colDef.PrimaryKey = true
		colDef.IsNotNull = true
	}

	return colDef
}

// dataTypeName returns the lowercased base name of a column's declared type,
// e.g. "nvarchar" for NVARCHAR(100). Length and precision arguments do not
// name distinct types.
func dataTypeName(ref tsql.DataTypeReference) string {
	switch t := ref.(type) {
	case *tsql.SqlDataTypeReference:
		if t.Name != nil && t.Name.BaseIdentifier != nil {
			return identifierValue(t.Name.BaseIdentifier)
		}
		return identifier(t.SqlDataTypeOption)
	case *tsql.XmlDataTypeReference:
		return "xml"
	case *tsql.UserDataTypeReference:
		if t.Name != nil && t.Name.BaseIdentifier != nil {
			return identifierValue(t.Name.BaseIdentifier)
		}
	}
	return ""
}

func (c *cc) convertDropTableStatement(n *tsql.DropTableStatement) ast.Node {
	stmt := &ast.DropTableStmt{
		IfExists: n.IsIfExists,
	}
	for _, obj := range n.Objects {
		stmt.Tables = append(stmt.Tables, parseTableName(obj))
	}
	return stmt
}

func (c *cc) convertAlterTableAddTableElementStatement(n *tsql.AlterTableAddTableElementStatement) ast.Node {
	if n.SchemaObjectName == nil || n.Definition == nil {
		return todo(n)
	}
	stmt := &ast.AlterTableStmt{
		Table: parseTableName(n.SchemaObjectName),
		Cmds:  &ast.List{},
	}
	for _, col := range n.Definition.ColumnDefinitions {
		def := c.convertColumnDefinition(col, nil)
		stmt.Cmds.Items = append(stmt.Cmds.Items, &ast.AlterTableCmd{
			Name:    &def.Colname,
			Subtype: ast.AT_AddColumn,
			Def:     def,
		})
	}
	return stmt
}

func (c *cc) convertAlterTableDropTableElementStatement(n *tsql.AlterTableDropTableElementStatement) ast.Node {
	if n.SchemaObjectName == nil {
		return todo(n)
	}
	stmt := &ast.AlterTableStmt{
		Table: parseTableName(n.SchemaObjectName),
		Cmds:  &ast.List{},
	}
	for _, elem := range n.AlterTableDropTableElements {
		if elem.TableElementType != "Column" || elem.Name == nil {
			continue
		}
		name := identifierValue(elem.Name)
		stmt.Cmds.Items = append(stmt.Cmds.Items, &ast.AlterTableCmd{
			Name:      &name,
			Subtype:   ast.AT_DropColumn,
			MissingOk: elem.IsIfExists,
		})
	}
	if len(stmt.Cmds.Items) == 0 {
		return todo(n)
	}
	return stmt
}
