package clickhouse

import (
	"strconv"
	"strings"

	chast "github.com/sqlc-dev/doubleclick/ast"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

type cc struct {
	paramCount int
}

func (c *cc) convert(node chast.Node) ast.Node {
	switch n := node.(type) {
	case *chast.SelectWithUnionQuery:
		return c.convertSelectWithUnionQuery(n)
	case *chast.SelectQuery:
		return c.convertSelectQuery(n)
	case *chast.InsertQuery:
		return c.convertInsertQuery(n)
	case *chast.CreateQuery:
		return c.convertCreateQuery(n)
	case *chast.UpdateQuery:
		return c.convertUpdateQuery(n)
	case *chast.DeleteQuery:
		return c.convertDeleteQuery(n)
	case *chast.DropQuery:
		return c.convertDropQuery(n)
	case *chast.AlterQuery:
		return c.convertAlterQuery(n)
	case *chast.TruncateQuery:
		return c.convertTruncateQuery(n)
	default:
		return todo(n)
	}
}

func (c *cc) convertSelectWithUnionQuery(n *chast.SelectWithUnionQuery) ast.Node {
	if len(n.Selects) == 0 {
		return &ast.TODO{}
	}

	// Single select without union
	if len(n.Selects) == 1 {
		return c.convert(n.Selects[0])
	}

	// Build a chain of SelectStmt with UNION operations
	var result *ast.SelectStmt
	for i, sel := range n.Selects {
		stmt, ok := c.convert(sel).(*ast.SelectStmt)
		if !ok {
			continue
		}
		if i == 0 {
			result = stmt
		} else {
			unionMode := ast.Union
			if i-1 < len(n.UnionModes) {
				switch strings.ToUpper(n.UnionModes[i-1]) {
				case "ALL":
					unionMode = ast.Union
				case "DISTINCT":
					unionMode = ast.Union
				}
			}
			result = &ast.SelectStmt{
				Op:   unionMode,
				All:  n.UnionAll || (i-1 < len(n.UnionModes) && strings.ToUpper(n.UnionModes[i-1]) == "ALL"),
				Larg: result,
				Rarg: stmt,
			}
		}
	}
	return result
}

func (c *cc) convertSelectQuery(n *chast.SelectQuery) *ast.SelectStmt {
	stmt := &ast.SelectStmt{}

	// Convert target list (SELECT columns)
	if len(n.Columns) > 0 {
		stmt.TargetList = &ast.List{}
		for _, col := range n.Columns {
			target := c.convertToResTarget(col)
			if target != nil {
				stmt.TargetList.Items = append(stmt.TargetList.Items, target)
			}
		}
	}

	// Convert FROM clause
	if n.From != nil {
		stmt.FromClause = c.convertTablesInSelectQuery(n.From)
	}

	// Convert WHERE clause
	if n.Where != nil {
		stmt.WhereClause = c.convertExpr(n.Where)
	}

	// Convert GROUP BY clause
	if len(n.GroupBy) > 0 {
		stmt.GroupClause = &ast.List{}
		for _, expr := range n.GroupBy {
			stmt.GroupClause.Items = append(stmt.GroupClause.Items, c.convertExpr(expr))
		}
	}

	// Convert HAVING clause
	if n.Having != nil {
		stmt.HavingClause = c.convertExpr(n.Having)
	}

	// Convert ORDER BY clause
	if len(n.OrderBy) > 0 {
		stmt.SortClause = &ast.List{}
		for _, orderBy := range n.OrderBy {
			stmt.SortClause.Items = append(stmt.SortClause.Items, c.convertOrderByElement(orderBy))
		}
	}

	// Convert LIMIT clause
	if n.Limit != nil {
		stmt.LimitCount = c.convertExpr(n.Limit)
	}

	// Convert OFFSET clause
	if n.Offset != nil {
		stmt.LimitOffset = c.convertExpr(n.Offset)
	}

	// Convert DISTINCT clause
	if n.Distinct {
		stmt.DistinctClause = &ast.List{}
	}

	// Convert DISTINCT ON clause
	if len(n.DistinctOn) > 0 {
		stmt.DistinctClause = &ast.List{}
		for _, expr := range n.DistinctOn {
			stmt.DistinctClause.Items = append(stmt.DistinctClause.Items, c.convertExpr(expr))
		}
	}

	// Convert WITH clause (CTEs)
	if len(n.With) > 0 {
		stmt.WithClause = &ast.WithClause{
			Ctes: &ast.List{},
		}
		for _, cte := range n.With {
			if aliased, ok := cte.(*chast.AliasedExpr); ok {
				cteNode := &ast.CommonTableExpr{
					Ctename: &aliased.Alias,
				}
				// CTE expression may be a Subquery containing the actual SELECT
				if subq, ok := aliased.Expr.(*chast.Subquery); ok {
					cteNode.Ctequery = c.convert(subq.Query)
				} else {
					// Fallback: treat the expression itself as the query
					cteNode.Ctequery = c.convertExpr(aliased.Expr)
				}
				stmt.WithClause.Ctes.Items = append(stmt.WithClause.Ctes.Items, cteNode)
			}
		}
	}

	return stmt
}

func (c *cc) convertToResTarget(expr chast.Expression) *ast.ResTarget {
	res := &ast.ResTarget{
		Location: expr.Pos().Offset,
	}

	switch e := expr.(type) {
	case *chast.Asterisk:
		if e.Table != "" {
			// table.*
			res.Val = &ast.ColumnRef{
				Fields: &ast.List{
					Items: []ast.Node{
						NewIdentifier(e.Table),
						&ast.A_Star{},
					},
				},
			}
		} else {
			// Just *
			res.Val = &ast.ColumnRef{
				Fields: &ast.List{
					Items: []ast.Node{&ast.A_Star{}},
				},
			}
		}
	case *chast.AliasedExpr:
		res.Name = &e.Alias
		res.Val = c.convertExpr(e.Expr)
	case *chast.Identifier:
		if e.Alias != "" {
			res.Name = &e.Alias
		}
		res.Val = c.convertIdentifier(e)
	case *chast.FunctionCall:
		if e.Alias != "" {
			res.Name = &e.Alias
		}
		res.Val = c.convertFunctionCall(e)
	default:
		res.Val = c.convertExpr(expr)
	}

	return res
}

func (c *cc) convertTablesInSelectQuery(n *chast.TablesInSelectQuery) *ast.List {
	if n == nil || len(n.Tables) == 0 {
		return nil
	}

	result := &ast.List{}

	for i, elem := range n.Tables {
		if elem.Table != nil {
			tableExpr := c.convertTableExpression(elem.Table)
			if i == 0 {
				result.Items = append(result.Items, tableExpr)
			} else if elem.Join != nil {
				// This element has a join
				joinExpr := c.convertTableJoin(elem.Join, result.Items[len(result.Items)-1], tableExpr)
				result.Items[len(result.Items)-1] = joinExpr
			} else {
				result.Items = append(result.Items, tableExpr)
			}
		} else if elem.Join != nil && len(result.Items) > 0 {
			// Join without table (should not happen normally)
			continue
		}
	}

	return result
}

func (c *cc) convertTableExpression(n *chast.TableExpression) ast.Node {
	var result ast.Node

	switch t := n.Table.(type) {
	case *chast.TableIdentifier:
		rv := parseTableIdentifierToRangeVar(t)
		if n.Alias != "" {
			alias := n.Alias
			rv.Alias = &ast.Alias{Aliasname: &alias}
		}
		result = rv
	case *chast.Subquery:
		subselect := &ast.RangeSubselect{
			Subquery: c.convert(t.Query),
		}
		alias := n.Alias
		if alias == "" && t.Alias != "" {
			alias = t.Alias
		}
		if alias != "" {
			subselect.Alias = &ast.Alias{Aliasname: &alias}
		}
		result = subselect
	case *chast.FunctionCall:
		// Table function like file(), url(), etc.
		rf := &ast.RangeFunction{
			Functions: &ast.List{
				Items: []ast.Node{c.convertFunctionCall(t)},
			},
		}
		if n.Alias != "" {
			alias := n.Alias
			rf.Alias = &ast.Alias{Aliasname: &alias}
		}
		result = rf
	default:
		result = &ast.TODO{}
	}

	return result
}

func (c *cc) convertTableJoin(n *chast.TableJoin, left, right ast.Node) *ast.JoinExpr {
	join := &ast.JoinExpr{
		Larg: left,
		Rarg: right,
	}

	// Convert join type
	switch n.Type {
	case chast.JoinInner:
		join.Jointype = ast.JoinTypeInner
	case chast.JoinLeft:
		join.Jointype = ast.JoinTypeLeft
	case chast.JoinRight:
		join.Jointype = ast.JoinTypeRight
	case chast.JoinFull:
		join.Jointype = ast.JoinTypeFull
	case chast.JoinCross:
		join.Jointype = ast.JoinTypeInner
		join.IsNatural = false
	default:
		join.Jointype = ast.JoinTypeInner
	}

	// Convert ON clause
	if n.On != nil {
		join.Quals = c.convertExpr(n.On)
	}

	// Convert USING clause
	if len(n.Using) > 0 {
		join.UsingClause = &ast.List{}
		for _, u := range n.Using {
			if id, ok := u.(*chast.Identifier); ok {
				join.UsingClause.Items = append(join.UsingClause.Items, NewIdentifier(id.Name()))
			}
		}
	}

	return join
}

func (c *cc) convertExpr(expr chast.Expression) ast.Node {
	if expr == nil {
		return nil
	}

	switch e := expr.(type) {
	case *chast.Identifier:
		return c.convertIdentifier(e)
	case *chast.Literal:
		return c.convertLiteral(e)
	case *chast.BinaryExpr:
		return c.convertBinaryExpr(e)
	case *chast.FunctionCall:
		return c.convertFunctionCall(e)
	case *chast.AliasedExpr:
		return c.convertExpr(e.Expr)
	case *chast.Parameter:
		return c.convertParameter(e)
	case *chast.Asterisk:
		return c.convertAsterisk(e)
	case *chast.CaseExpr:
		return c.convertCaseExpr(e)
	case *chast.CastExpr:
		return c.convertCastExpr(e)
	case *chast.BetweenExpr:
		return c.convertBetweenExpr(e)
	case *chast.InExpr:
		return c.convertInExpr(e)
	case *chast.IsNullExpr:
		return c.convertIsNullExpr(e)
	case *chast.LikeExpr:
		return c.convertLikeExpr(e)
	case *chast.Subquery:
		return c.convertSubquery(e)
	case *chast.ArrayAccess:
		return c.convertArrayAccess(e)
	case *chast.UnaryExpr:
		return c.convertUnaryExpr(e)
	case *chast.Lambda:
		// Lambda expressions are ClickHouse-specific, return as-is for now
		return &ast.TODO{}
	default:
		return &ast.TODO{}
	}
}

func (c *cc) convertIdentifier(n *chast.Identifier) *ast.ColumnRef {
	fields := &ast.List{}
	for _, part := range n.Parts {
		fields.Items = append(fields.Items, NewIdentifier(part))
	}
	return &ast.ColumnRef{
		Fields:   fields,
		Location: n.Pos().Offset,
	}
}

func (c *cc) convertLiteral(n *chast.Literal) *ast.A_Const {
	switch n.Type {
	case chast.LiteralString:
		str := n.Value.(string)
		return &ast.A_Const{
			Val:      &ast.String{Str: str},
			Location: n.Pos().Offset,
		}
	case chast.LiteralInteger:
		var ival int64
		switch v := n.Value.(type) {
		case int64:
			ival = v
		case int:
			ival = int64(v)
		case float64:
			ival = int64(v)
		case string:
			ival, _ = strconv.ParseInt(v, 10, 64)
		}
		return &ast.A_Const{
			Val:      &ast.Integer{Ival: ival},
			Location: n.Pos().Offset,
		}
	case chast.LiteralFloat:
		var fval float64
		switch v := n.Value.(type) {
		case float64:
			fval = v
		case string:
			fval, _ = strconv.ParseFloat(v, 64)
		}
		str := strconv.FormatFloat(fval, 'f', -1, 64)
		return &ast.A_Const{
			Val:      &ast.Float{Str: str},
			Location: n.Pos().Offset,
		}
	case chast.LiteralBoolean:
		// ClickHouse booleans are typically 0/1
		bval := n.Value.(bool)
		if bval {
			return &ast.A_Const{
				Val:      &ast.Integer{Ival: 1},
				Location: n.Pos().Offset,
			}
		}
		return &ast.A_Const{
			Val:      &ast.Integer{Ival: 0},
			Location: n.Pos().Offset,
		}
	case chast.LiteralNull:
		return &ast.A_Const{
			Val:      &ast.Null{},
			Location: n.Pos().Offset,
		}
	default:
		return &ast.A_Const{
			Location: n.Pos().Offset,
		}
	}
}

func (c *cc) convertBinaryExpr(n *chast.BinaryExpr) ast.Node {
	op := strings.ToUpper(n.Op)

	// Handle logical operators
	if op == "AND" || op == "OR" {
		var boolop ast.BoolExprType
		if op == "AND" {
			boolop = ast.BoolExprTypeAnd
		} else {
			boolop = ast.BoolExprTypeOr
		}
		return &ast.BoolExpr{
			Boolop: boolop,
			Args: &ast.List{
				Items: []ast.Node{
					c.convertExpr(n.Left),
					c.convertExpr(n.Right),
				},
			},
			Location: n.Pos().Offset,
		}
	}

	// Handle other operators
	return &ast.A_Expr{
		Name: &ast.List{
			Items: []ast.Node{&ast.String{Str: n.Op}},
		},
		Lexpr:    c.convertExpr(n.Left),
		Rexpr:    c.convertExpr(n.Right),
		Location: n.Pos().Offset,
	}
}

func (c *cc) convertFunctionCall(n *chast.FunctionCall) *ast.FuncCall {
	fc := &ast.FuncCall{
		Funcname: &ast.List{
			Items: []ast.Node{&ast.String{Str: n.Name}},
		},
		Location:     n.Pos().Offset,
		AggDistinct:  n.Distinct,
	}

	// Convert arguments
	if len(n.Arguments) > 0 {
		fc.Args = &ast.List{}
		for _, arg := range n.Arguments {
			fc.Args.Items = append(fc.Args.Items, c.convertExpr(arg))
		}
	}

	// Convert window function
	if n.Over != nil {
		fc.Over = &ast.WindowDef{}
		if len(n.Over.PartitionBy) > 0 {
			fc.Over.PartitionClause = &ast.List{}
			for _, p := range n.Over.PartitionBy {
				fc.Over.PartitionClause.Items = append(fc.Over.PartitionClause.Items, c.convertExpr(p))
			}
		}
		if len(n.Over.OrderBy) > 0 {
			fc.Over.OrderClause = &ast.List{}
			for _, o := range n.Over.OrderBy {
				fc.Over.OrderClause.Items = append(fc.Over.OrderClause.Items, c.convertOrderByElement(o))
			}
		}
	}

	return fc
}

func (c *cc) convertParameter(n *chast.Parameter) ast.Node {
	c.paramCount++
	// Use the parameter name if available
	name := n.Name
	if name == "" {
		name = strconv.Itoa(c.paramCount)
	}
	return &ast.ParamRef{
		Number:   c.paramCount,
		Location: n.Pos().Offset,
	}
}

func (c *cc) convertAsterisk(n *chast.Asterisk) *ast.ColumnRef {
	fields := &ast.List{}
	if n.Table != "" {
		fields.Items = append(fields.Items, NewIdentifier(n.Table))
	}
	fields.Items = append(fields.Items, &ast.A_Star{})
	return &ast.ColumnRef{
		Fields:   fields,
		Location: n.Pos().Offset,
	}
}

func (c *cc) convertCaseExpr(n *chast.CaseExpr) *ast.CaseExpr {
	ce := &ast.CaseExpr{
		Location: n.Pos().Offset,
	}

	// Convert test expression (CASE expr WHEN ...)
	if n.Operand != nil {
		ce.Arg = c.convertExpr(n.Operand)
	}

	// Convert WHEN clauses
	if len(n.Whens) > 0 {
		ce.Args = &ast.List{}
		for _, when := range n.Whens {
			caseWhen := &ast.CaseWhen{
				Expr:   c.convertExpr(when.Condition),
				Result: c.convertExpr(when.Result),
			}
			ce.Args.Items = append(ce.Args.Items, caseWhen)
		}
	}

	// Convert ELSE clause
	if n.Else != nil {
		ce.Defresult = c.convertExpr(n.Else)
	}

	return ce
}

func (c *cc) convertCastExpr(n *chast.CastExpr) *ast.TypeCast {
	tc := &ast.TypeCast{
		Arg:      c.convertExpr(n.Expr),
		Location: n.Pos().Offset,
	}

	if n.Type != nil {
		tc.TypeName = &ast.TypeName{
			Name: n.Type.Name,
		}
	}

	return tc
}

func (c *cc) convertBetweenExpr(n *chast.BetweenExpr) *ast.BetweenExpr {
	return &ast.BetweenExpr{
		Expr:     c.convertExpr(n.Expr),
		Left:     c.convertExpr(n.Low),
		Right:    c.convertExpr(n.High),
		Not:      n.Not,
		Location: n.Pos().Offset,
	}
}

func (c *cc) convertInExpr(n *chast.InExpr) *ast.In {
	in := &ast.In{
		Expr:     c.convertExpr(n.Expr),
		Not:      n.Not,
		Location: n.Pos().Offset,
	}

	// Convert the list
	if len(n.List) > 0 {
		in.List = make([]ast.Node, 0, len(n.List))
		for _, item := range n.List {
			in.List = append(in.List, c.convertExpr(item))
		}
	}

	// Handle subquery
	if n.Query != nil {
		in.Sel = c.convert(n.Query)
	}

	return in
}

func (c *cc) convertIsNullExpr(n *chast.IsNullExpr) *ast.NullTest {
	nullTest := &ast.NullTest{
		Arg:      c.convertExpr(n.Expr),
		Location: n.Pos().Offset,
	}
	if n.Not {
		nullTest.Nulltesttype = ast.NullTestTypeIsNotNull
	} else {
		nullTest.Nulltesttype = ast.NullTestTypeIsNull
	}
	return nullTest
}

func (c *cc) convertLikeExpr(n *chast.LikeExpr) *ast.A_Expr {
	kind := ast.A_Expr_Kind(0)
	opName := "~~"
	if n.CaseInsensitive {
		opName = "~~*"
	}
	if n.Not {
		opName = "!~~"
		if n.CaseInsensitive {
			opName = "!~~*"
		}
	}

	return &ast.A_Expr{
		Kind: kind,
		Name: &ast.List{
			Items: []ast.Node{&ast.String{Str: opName}},
		},
		Lexpr:    c.convertExpr(n.Expr),
		Rexpr:    c.convertExpr(n.Pattern),
		Location: n.Pos().Offset,
	}
}

func (c *cc) convertSubquery(n *chast.Subquery) *ast.SubLink {
	return &ast.SubLink{
		SubLinkType: ast.EXISTS_SUBLINK,
		Subselect:   c.convert(n.Query),
	}
}

func (c *cc) convertArrayAccess(n *chast.ArrayAccess) *ast.A_Indirection {
	return &ast.A_Indirection{
		Arg: c.convertExpr(n.Array),
		Indirection: &ast.List{
			Items: []ast.Node{
				&ast.A_Indices{
					Uidx: c.convertExpr(n.Index),
				},
			},
		},
	}
}

func (c *cc) convertUnaryExpr(n *chast.UnaryExpr) ast.Node {
	op := strings.ToUpper(n.Op)

	if op == "NOT" {
		return &ast.BoolExpr{
			Boolop: ast.BoolExprTypeNot,
			Args: &ast.List{
				Items: []ast.Node{c.convertExpr(n.Operand)},
			},
			Location: n.Pos().Offset,
		}
	}

	return &ast.A_Expr{
		Name: &ast.List{
			Items: []ast.Node{&ast.String{Str: n.Op}},
		},
		Rexpr:    c.convertExpr(n.Operand),
		Location: n.Pos().Offset,
	}
}

func (c *cc) convertOrderByElement(n *chast.OrderByElement) *ast.SortBy {
	sortBy := &ast.SortBy{
		Node:     c.convertExpr(n.Expression),
		Location: n.Expression.Pos().Offset,
	}

	if n.Descending {
		sortBy.SortbyDir = ast.SortByDirDesc
	} else {
		sortBy.SortbyDir = ast.SortByDirAsc
	}

	if n.NullsFirst != nil {
		if *n.NullsFirst {
			sortBy.SortbyNulls = ast.SortByNullsFirst
		} else {
			sortBy.SortbyNulls = ast.SortByNullsLast
		}
	}

	return sortBy
}

func (c *cc) convertInsertQuery(n *chast.InsertQuery) *ast.InsertStmt {
	stmt := &ast.InsertStmt{
		Relation: &ast.RangeVar{
			Relname: &n.Table,
		},
	}

	if n.Database != "" {
		stmt.Relation.Schemaname = &n.Database
	}

	// Convert column list
	if len(n.Columns) > 0 {
		stmt.Cols = &ast.List{}
		for _, col := range n.Columns {
			name := col.Name()
			stmt.Cols.Items = append(stmt.Cols.Items, &ast.ResTarget{
				Name: &name,
			})
		}
	}

	// Convert SELECT subquery if present
	if n.Select != nil {
		stmt.SelectStmt = c.convert(n.Select)
	}

	// Convert VALUES clause
	if len(n.Values) > 0 {
		selectStmt := &ast.SelectStmt{
			ValuesLists: &ast.List{},
		}
		for _, row := range n.Values {
			rowList := &ast.List{}
			for _, val := range row {
				rowList.Items = append(rowList.Items, c.convertExpr(val))
			}
			selectStmt.ValuesLists.Items = append(selectStmt.ValuesLists.Items, rowList)
		}
		stmt.SelectStmt = selectStmt
	}

	return stmt
}

func (c *cc) convertCreateQuery(n *chast.CreateQuery) ast.Node {
	// Handle CREATE DATABASE
	if n.CreateDatabase {
		return &ast.CreateSchemaStmt{
			Name:        &n.Database,
			IfNotExists: n.IfNotExists,
		}
	}

	// Handle CREATE TABLE
	if n.Table != "" {
		stmt := &ast.CreateTableStmt{
			Name: &ast.TableName{
				Name: identifier(n.Table),
			},
			IfNotExists: n.IfNotExists,
		}

		if n.Database != "" {
			stmt.Name.Schema = identifier(n.Database)
		}

		// Convert columns
		for _, col := range n.Columns {
			colDef := c.convertColumnDeclaration(col)
			stmt.Cols = append(stmt.Cols, colDef)
		}

		// Convert AS SELECT
		if n.AsSelect != nil {
			// This is a CREATE TABLE ... AS SELECT
			// The AsSelect field contains the SELECT statement
		}

		return stmt
	}

	// Handle CREATE VIEW
	if n.View != "" {
		return &ast.ViewStmt{
			View: &ast.RangeVar{
				Relname: &n.View,
			},
			Query:   c.convert(n.AsSelect),
			Replace: n.OrReplace,
		}
	}

	return &ast.TODO{}
}

func (c *cc) convertColumnDeclaration(n *chast.ColumnDeclaration) *ast.ColumnDef {
	colDef := &ast.ColumnDef{
		Colname:   identifier(n.Name),
		IsNotNull: isNotNull(n),
	}

	if n.Type != nil {
		colDef.TypeName = &ast.TypeName{
			Name: n.Type.Name,
		}
		// Handle type parameters (e.g., Decimal(10, 2))
		if len(n.Type.Parameters) > 0 {
			colDef.TypeName.Typmods = &ast.List{}
			for _, param := range n.Type.Parameters {
				colDef.TypeName.Typmods.Items = append(colDef.TypeName.Typmods.Items, c.convertExpr(param))
			}
		}
	}

	// Handle PRIMARY KEY constraint
	if n.PrimaryKey {
		colDef.PrimaryKey = true
	}

	// Handle DEFAULT
	if n.Default != nil {
		// colDef.RawDefault = c.convertExpr(n.Default)
	}

	// Handle comment
	if n.Comment != "" {
		colDef.Comment = n.Comment
	}

	return colDef
}

func (c *cc) convertUpdateQuery(n *chast.UpdateQuery) *ast.UpdateStmt {
	rv := &ast.RangeVar{
		Relname: &n.Table,
	}
	if n.Database != "" {
		rv.Schemaname = &n.Database
	}
	stmt := &ast.UpdateStmt{
		Relations: &ast.List{
			Items: []ast.Node{rv},
		},
	}

	// Convert assignments
	if len(n.Assignments) > 0 {
		stmt.TargetList = &ast.List{}
		for _, assign := range n.Assignments {
			name := identifier(assign.Column)
			stmt.TargetList.Items = append(stmt.TargetList.Items, &ast.ResTarget{
				Name: &name,
				Val:  c.convertExpr(assign.Value),
			})
		}
	}

	// Convert WHERE clause
	if n.Where != nil {
		stmt.WhereClause = c.convertExpr(n.Where)
	}

	return stmt
}

func (c *cc) convertDeleteQuery(n *chast.DeleteQuery) *ast.DeleteStmt {
	rv := &ast.RangeVar{
		Relname: &n.Table,
	}
	if n.Database != "" {
		rv.Schemaname = &n.Database
	}
	stmt := &ast.DeleteStmt{
		Relations: &ast.List{
			Items: []ast.Node{rv},
		},
	}

	// Convert WHERE clause
	if n.Where != nil {
		stmt.WhereClause = c.convertExpr(n.Where)
	}

	return stmt
}

func (c *cc) convertDropQuery(n *chast.DropQuery) ast.Node {
	// Handle DROP TABLE
	if n.Table != "" {
		tableName := &ast.TableName{
			Name: identifier(n.Table),
		}
		if n.Database != "" {
			tableName.Schema = identifier(n.Database)
		}
		return &ast.DropTableStmt{
			IfExists: n.IfExists,
			Tables:   []*ast.TableName{tableName},
		}
	}

	// Handle DROP TABLE with multiple tables
	if len(n.Tables) > 0 {
		tables := make([]*ast.TableName, 0, len(n.Tables))
		for _, t := range n.Tables {
			tables = append(tables, parseTableName(t))
		}
		return &ast.DropTableStmt{
			IfExists: n.IfExists,
			Tables:   tables,
		}
	}

	// Handle DROP DATABASE - return TODO for now
	// Handle DROP VIEW - return TODO for now
	return &ast.TODO{}
}

func (c *cc) convertAlterQuery(n *chast.AlterQuery) ast.Node {
	alt := &ast.AlterTableStmt{
		Table: &ast.TableName{
			Name: identifier(n.Table),
		},
		Cmds: &ast.List{},
	}

	if n.Database != "" {
		alt.Table.Schema = identifier(n.Database)
	}

	for _, cmd := range n.Commands {
		switch cmd.Type {
		case chast.AlterAddColumn:
			if cmd.Column != nil {
				name := cmd.Column.Name
				alt.Cmds.Items = append(alt.Cmds.Items, &ast.AlterTableCmd{
					Name:    &name,
					Subtype: ast.AT_AddColumn,
					Def:     c.convertColumnDeclaration(cmd.Column),
				})
			}
		case chast.AlterDropColumn:
			name := cmd.ColumnName
			alt.Cmds.Items = append(alt.Cmds.Items, &ast.AlterTableCmd{
				Name:      &name,
				Subtype:   ast.AT_DropColumn,
				MissingOk: cmd.IfExists,
			})
		case chast.AlterModifyColumn:
			if cmd.Column != nil {
				name := cmd.Column.Name
				// Drop and re-add to simulate modify
				alt.Cmds.Items = append(alt.Cmds.Items, &ast.AlterTableCmd{
					Name:    &name,
					Subtype: ast.AT_DropColumn,
				})
				alt.Cmds.Items = append(alt.Cmds.Items, &ast.AlterTableCmd{
					Name:    &name,
					Subtype: ast.AT_AddColumn,
					Def:     c.convertColumnDeclaration(cmd.Column),
				})
			}
		case chast.AlterRenameColumn:
			oldName := cmd.ColumnName
			newName := cmd.NewName
			return &ast.RenameColumnStmt{
				Table:   alt.Table,
				Col:     &ast.ColumnRef{Name: oldName},
				NewName: &newName,
			}
		}
	}

	return alt
}

func (c *cc) convertTruncateQuery(n *chast.TruncateQuery) *ast.TruncateStmt {
	stmt := &ast.TruncateStmt{
		Relations: &ast.List{},
	}

	tableName := n.Table
	schemaName := n.Database

	rv := &ast.RangeVar{
		Relname: &tableName,
	}
	if schemaName != "" {
		rv.Schemaname = &schemaName
	}

	stmt.Relations.Items = append(stmt.Relations.Items, rv)

	return stmt
}
