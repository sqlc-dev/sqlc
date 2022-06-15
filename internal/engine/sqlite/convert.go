package sqlite

import (
	"strconv"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	"github.com/kyleconroy/sqlc/internal/engine/sqlite/parser"
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type node interface {
	GetParser() antlr.Parser
}

func convertAlter_table_stmtContext(c *parser.Alter_table_stmtContext) ast.Node {

	if c.RENAME_() != nil {
		if newTable, ok := c.New_table_name().(*parser.New_table_nameContext); ok {
			name := newTable.Any_name().GetText()
			return &ast.RenameTableStmt{
				Table:   parseTableName(c),
				NewName: &name,
			}
		}

		if newCol, ok := c.GetNew_column_name().(*parser.Column_nameContext); ok {
			name := newCol.Any_name().GetText()
			return &ast.RenameColumnStmt{
				Table: parseTableName(c),
				Col: &ast.ColumnRef{
					Name: c.GetOld_column_name().GetText(),
				},
				NewName: &name,
			}
		}
	}

	if c.ADD_() != nil {
		if def, ok := c.Column_def().(*parser.Column_defContext); ok {
			stmt := &ast.AlterTableStmt{
				Table: parseTableName(c),
				Cmds:  &ast.List{},
			}
			name := def.Column_name().GetText()
			stmt.Cmds.Items = append(stmt.Cmds.Items, &ast.AlterTableCmd{
				Name:    &name,
				Subtype: ast.AT_AddColumn,
				Def: &ast.ColumnDef{
					Colname: name,
					TypeName: &ast.TypeName{
						Name: def.Type_name().GetText(),
					},
					IsNotNull: hasNotNullConstraint(def.AllColumn_constraint()),
				},
			})
			return stmt
		}
	}

	if c.DROP_() != nil {
		stmt := &ast.AlterTableStmt{
			Table: parseTableName(c),
			Cmds:  &ast.List{},
		}
		name := c.Column_name(0).GetText()
		//fmt.Printf("column: %s", name)
		stmt.Cmds.Items = append(stmt.Cmds.Items, &ast.AlterTableCmd{
			Name:    &name,
			Subtype: ast.AT_DropColumn,
		})
		return stmt
	}

	return &ast.TODO{}
}

func convertAttach_stmtContext(c *parser.Attach_stmtContext) ast.Node {
	name := c.Schema_name().GetText()
	return &ast.CreateSchemaStmt{
		Name: &name,
	}
}

func convertCreate_table_stmtContext(c *parser.Create_table_stmtContext) ast.Node {
	stmt := &ast.CreateTableStmt{
		Name:        parseTableName(c),
		IfNotExists: c.EXISTS_() != nil,
	}
	for _, idef := range c.AllColumn_def() {
		if def, ok := idef.(*parser.Column_defContext); ok {
			stmt.Cols = append(stmt.Cols, &ast.ColumnDef{
				Colname:   identifier(def.Column_name().GetText()),
				IsNotNull: hasNotNullConstraint(def.AllColumn_constraint()),
				TypeName:  &ast.TypeName{Name: def.Type_name().GetText()},
			})
		}
	}
	return stmt
}

func convertCreate_view_stmtContext(c *parser.Create_view_stmtContext) ast.Node {
	viewName := c.View_name().GetText()
	relation := &ast.RangeVar{
		Relname: &viewName,
	}

	if c.Schema_name() != nil {
		schemaName := c.Schema_name().GetText()
		relation.Schemaname = &schemaName
	}

	return &ast.ViewStmt{
		View:            relation,
		Aliases:         &ast.List{},
		Query:           convert(c.Select_stmt()),
		Replace:         false,
		Options:         &ast.List{},
		WithCheckOption: ast.ViewCheckOption(0),
	}
}

func convertDelete_stmtContext(c *parser.Delete_stmtContext) ast.Node {
	if qualifiedName, ok := c.Qualified_table_name().(*parser.Qualified_table_nameContext); ok {

		tableName := qualifiedName.Table_name().GetText()
		relation := &ast.RangeVar{
			Relname: &tableName,
		}

		if qualifiedName.Schema_name() != nil {
			schemaName := qualifiedName.Schema_name().GetText()
			relation.Schemaname = &schemaName
		}

		if qualifiedName.Alias() != nil {
			alias := qualifiedName.Alias().GetText()
			relation.Alias = &ast.Alias{Aliasname: &alias}
		}

		delete := &ast.DeleteStmt{
			Relation:      relation,
			ReturningList: &ast.List{},
			WithClause:    nil,
		}

		if c.WHERE_() != nil {
			if c.Expr() != nil {
				delete.WhereClause = convert(c.Expr())
			}
		}

		return delete
	}

	return &ast.TODO{}
}

func convertDrop_stmtContext(c *parser.Drop_stmtContext) ast.Node {
	if c.TABLE_() != nil || c.VIEW_() != nil {
		name := ast.TableName{
			Name: c.Any_name().GetText(),
		}
		if c.Schema_name() != nil {
			name.Schema = c.Schema_name().GetText()
		}

		return &ast.DropTableStmt{
			IfExists: c.EXISTS_() != nil,
			Tables:   []*ast.TableName{&name},
		}
	} else {
		return &ast.TODO{}
	}
}

func identifier(id string) string {
	return strings.ToLower(id)
}

func NewIdentifer(t string) *ast.String {
	return &ast.String{Str: identifier(t)}
}

func convertFuncContext(c *parser.Expr_functionContext) ast.Node {
	if name, ok := c.Function_name().(*parser.Function_nameContext); ok {
		funcName := strings.ToLower(name.GetText())

		var argNodes []ast.Node
		for _, exp := range c.AllExpr() {
			argNodes = append(argNodes, convert(exp))
		}
		args := &ast.List{Items: argNodes}

		if funcName == "coalesce" {
			return &ast.CoalesceExpr{
				Args: args,
			}
		} else {
			return &ast.FuncCall{
				Func: &ast.FuncName{
					Name: funcName,
				},
				Funcname: &ast.List{
					Items: []ast.Node{
						NewIdentifer(funcName),
					},
				},
				AggStar:     c.STAR() != nil,
				Args:        args,
				AggOrder:    &ast.List{},
				AggDistinct: c.DISTINCT_() != nil,
			}
		}
	}

	return &ast.TODO{}
}

func convertExprContext(c *parser.ExprContext) ast.Node {
	return &ast.TODO{}
}

func convertColumnNameExpr(c *parser.Expr_qualified_column_nameContext) *ast.ColumnRef {
	var items []ast.Node
	if schema, ok := c.Schema_name().(*parser.Schema_nameContext); ok {
		schemaText := schema.GetText()
		if schemaText != "" {
			items = append(items, NewIdentifer(schemaText))
		}
	}
	if table, ok := c.Table_name().(*parser.Table_nameContext); ok {
		tableName := table.GetText()
		if tableName != "" {
			items = append(items, NewIdentifer(tableName))
		}
	}
	items = append(items, NewIdentifer(c.Column_name().GetText()))
	return &ast.ColumnRef{
		Fields: &ast.List{
			Items: items,
		},
	}
}

func convertComparison(c *parser.Expr_comparisonContext) ast.Node {
	aExpr := &ast.A_Expr{
		Name: &ast.List{
			Items: []ast.Node{
				&ast.String{Str: "="}, // TODO: add actual comparison
			},
		},
		Lexpr: convert(c.Expr(0)),
		Rexpr: convert(c.Expr(1)),
	}

	return aExpr
}

func convertMultiSelect_stmtContext(c *parser.Select_stmtContext) ast.Node {
	var tables []ast.Node
	var cols []ast.Node
	var where ast.Node
	var groupBy = &ast.List{Items: []ast.Node{}}
	var having ast.Node

	for _, icore := range c.AllSelect_core() {
		core, ok := icore.(*parser.Select_coreContext)
		if !ok {
			continue
		}
		cols = append(cols, getCols(core)...)
		tables = append(tables, getTables(core)...)

		i := 0
		if core.WHERE_() != nil {
			where = convert(core.Expr(i))
			i++
		}

		if core.GROUP_() != nil {
			n := len(core.AllExpr()) - i
			if core.HAVING_() != nil {
				n--
			}

			for i < n {
				groupBy.Items = append(groupBy.Items, convert(core.Expr(i)))
				i++
			}

			if core.HAVING_() != nil {
				having = convert(core.Expr(i))
				i++
			}
		}
	}

	window := &ast.List{Items: []ast.Node{}}
	if c.Order_by_stmt() != nil {
		window.Items = append(window.Items, convert(c.Order_by_stmt()))
	}

	return &ast.SelectStmt{
		FromClause:   &ast.List{Items: tables},
		TargetList:   &ast.List{Items: cols},
		WhereClause:  where,
		GroupClause:  groupBy,
		HavingClause: having,
		WindowClause: window,
		ValuesLists:  &ast.List{},
	}
}

func getTables(core *parser.Select_coreContext) []ast.Node {
	var tables []ast.Node
	tables = append(tables, convertTablesOrSubquery(core.AllTable_or_subquery())...)

	if core.Join_clause() != nil {
		join, ok := core.Join_clause().(*parser.Join_clauseContext)
		if ok {
			tables = append(tables, convertTablesOrSubquery(join.AllTable_or_subquery())...)
		}
	}

	return tables
}

func getCols(core *parser.Select_coreContext) []ast.Node {
	var cols []ast.Node
	for _, icol := range core.AllResult_column() {
		col, ok := icol.(*parser.Result_columnContext)
		if !ok {
			continue
		}
		target := &ast.ResTarget{
			Location: col.GetStart().GetStart(),
		}
		var val ast.Node
		iexpr := col.Expr()
		switch {
		case col.STAR() != nil:
			val = convertWildCardField(col)
		case iexpr != nil:
			val = convert(iexpr)
		}

		if val == nil {
			continue
		}

		if col.AS_() != nil {
			name := col.Column_alias().GetText()
			target.Name = &name
		}

		target.Val = val
		cols = append(cols, target)
	}
	return cols
}

func convertWildCardField(c *parser.Result_columnContext) *ast.ColumnRef {
	items := []ast.Node{}
	if c.Table_name() != nil {
		items = append(items, NewIdentifer(c.Table_name().GetText()))
	}
	items = append(items, &ast.A_Star{})

	return &ast.ColumnRef{
		Fields: &ast.List{
			Items: items,
		},
		Location: c.GetStart().GetStart(),
	}
}

func convertOrderby_stmtContext(c parser.IOrder_by_stmtContext) ast.Node {
	if orderBy, ok := c.(*parser.Order_by_stmtContext); ok {
		list := &ast.List{Items: []ast.Node{}}
		for _, o := range orderBy.AllOrdering_term() {
			term, ok := o.(*parser.Ordering_termContext)
			if !ok {
				continue
			}
			list.Items = append(list.Items, &ast.CaseExpr{
				Xpr:      convert(term.Expr()),
				Location: term.Expr().GetStart().GetStart(),
			})
		}
		return list
	}
	return &ast.TODO{}
}

func convertSql_stmtContext(n *parser.Sql_stmtContext) ast.Node {
	if stmt := n.Alter_table_stmt(); stmt != nil {
		return convert(stmt)
	}
	if stmt := n.Analyze_stmt(); stmt != nil {
		return convert(stmt)
	}
	if stmt := n.Attach_stmt(); stmt != nil {
		return convert(stmt)
	}
	if stmt := n.Begin_stmt(); stmt != nil {
		return convert(stmt)
	}
	if stmt := n.Commit_stmt(); stmt != nil {
		return convert(stmt)
	}
	if stmt := n.Create_index_stmt(); stmt != nil {
		return convert(stmt)
	}
	if stmt := n.Create_table_stmt(); stmt != nil {
		return convert(stmt)
	}
	if stmt := n.Create_trigger_stmt(); stmt != nil {
		return convert(stmt)
	}
	if stmt := n.Create_view_stmt(); stmt != nil {
		return convert(stmt)
	}
	if stmt := n.Create_virtual_table_stmt(); stmt != nil {
		return convert(stmt)
	}
	if stmt := n.Delete_stmt(); stmt != nil {
		return convert(stmt)
	}
	if stmt := n.Delete_stmt_limited(); stmt != nil {
		return convert(stmt)
	}
	if stmt := n.Detach_stmt(); stmt != nil {
		return convert(stmt)
	}
	if stmt := n.Drop_stmt(); stmt != nil {
		return convert(stmt)
	}
	if stmt := n.Insert_stmt(); stmt != nil {
		return convert(stmt)
	}
	if stmt := n.Pragma_stmt(); stmt != nil {
		return convert(stmt)
	}
	if stmt := n.Reindex_stmt(); stmt != nil {
		return convert(stmt)
	}
	if stmt := n.Release_stmt(); stmt != nil {
		return convert(stmt)
	}
	if stmt := n.Rollback_stmt(); stmt != nil {
		return convert(stmt)
	}
	if stmt := n.Savepoint_stmt(); stmt != nil {
		return convert(stmt)
	}
	if stmt := n.Select_stmt(); stmt != nil {
		return convert(stmt)
	}
	if stmt := n.Update_stmt(); stmt != nil {
		return convert(stmt)
	}
	if stmt := n.Update_stmt_limited(); stmt != nil {
		return convert(stmt)
	}
	if stmt := n.Vacuum_stmt(); stmt != nil {
		return convert(stmt)
	}
	return nil
}

func convertLiteral(c *parser.Expr_literalContext) ast.Node {
	if literal, ok := c.Literal_value().(*parser.Literal_valueContext); ok {

		if literal.NUMERIC_LITERAL() != nil {
			i, _ := strconv.ParseInt(literal.GetText(), 10, 64)
			return &ast.A_Const{
				Val: &ast.Integer{Ival: i},
			}
		}

		if literal.STRING_LITERAL() != nil {
			return &ast.A_Const{
				Val: &ast.String{Str: literal.GetText()},
			}
		}

		if literal.TRUE_() != nil || literal.FALSE_() != nil {
			var i int64
			if literal.TRUE_() != nil {
				i = 1
			}

			return &ast.A_Const{
				Val: &ast.Integer{Ival: i},
			}
		}

	}
	return &ast.TODO{}
}

func convertMathOperationNode(c *parser.Expr_math_opContext) ast.Node {
	return &ast.A_Expr{
		Name: &ast.List{
			Items: []ast.Node{
				&ast.String{Str: "+"}, // todo: Convert operation types
			},
		},
		Lexpr: convert(c.Expr(0)),
		Rexpr: convert(c.Expr(1)),
	}
}

func convertBinaryNode(c *parser.Expr_binaryContext) ast.Node {
	return &ast.BoolExpr{
		// TODO: Set op
		Args: &ast.List{
			Items: []ast.Node{
				convert(c.Expr(0)),
				convert(c.Expr(1)),
			},
		},
	}
}

func convertParam(c *parser.Expr_bindContext) ast.Node {
	if c.BIND_PARAMETER() != nil {
		return &ast.ParamRef{ // TODO: Need to count these up instead of always using 0
			Location: c.GetStart().GetStart(),
		}
	}
	return &ast.TODO{}
}

func convertInsert_stmtContext(c *parser.Insert_stmtContext) ast.Node {
	tableName := c.Table_name().GetText()
	rel := &ast.RangeVar{
		Relname: &tableName,
	}
	if c.Schema_name() != nil {
		schemaName := c.Schema_name().GetText()
		rel.Schemaname = &schemaName
	}
	if c.Table_alias() != nil {
		tableAlias := c.Table_alias().GetText()
		rel.Alias = &ast.Alias{
			Aliasname: &tableAlias,
		}
	}

	insert := &ast.InsertStmt{
		Relation:      rel,
		Cols:          convertColumnNames(c.AllColumn_name()),
		ReturningList: &ast.List{},
	}

	if ss, ok := convert(c.Select_stmt()).(*ast.SelectStmt); ok {
		ss.ValuesLists = &ast.List{}
		insert.SelectStmt = ss
	} else {
		insert.SelectStmt = &ast.SelectStmt{
			FromClause:  &ast.List{},
			TargetList:  &ast.List{},
			ValuesLists: convertExprLists(c.AllExpr()),
		}
	}

	return insert
}

func convertExprLists(lists []parser.IExprContext) *ast.List {
	list := &ast.List{Items: []ast.Node{}}
	n := len(lists)
	inner := &ast.List{Items: []ast.Node{}}
	for i := 0; i < n; i++ {
		inner.Items = append(inner.Items, convert(lists[i]))
	}
	list.Items = append(list.Items, inner)
	return list
}

func convertColumnNames(cols []parser.IColumn_nameContext) *ast.List {
	list := &ast.List{Items: []ast.Node{}}
	for _, c := range cols {
		name := c.GetText()
		list.Items = append(list.Items, &ast.ResTarget{
			Name: &name,
		})
	}
	return list
}

func convertTablesOrSubquery(tableOrSubquery []parser.ITable_or_subqueryContext) []ast.Node {
	var tables []ast.Node
	for _, ifrom := range tableOrSubquery {
		from, ok := ifrom.(*parser.Table_or_subqueryContext)
		if !ok {
			continue
		}
		rel := from.Table_name().GetText()
		name := ast.RangeVar{
			Relname:  &rel,
			Location: from.GetStart().GetStart(),
		}
		if from.Schema_name() != nil {
			schema := from.Schema_name().GetText()
			name.Schemaname = &schema
		}
		if from.Table_alias() != nil {
			alias := from.Table_alias().GetText()
			name.Alias = &ast.Alias{Aliasname: &alias}
		}

		tables = append(tables, &name)
	}

	return tables
}

func convertUpdate_stmtContext(c *parser.Update_stmtContext) ast.Node {
	if c == nil {
		return nil
	}

	relations := &ast.List{}
	tableName := c.Qualified_table_name().GetText()
	rel := ast.RangeVar{
		Relname:  &tableName,
		Location: c.GetStart().GetStart(),
	}
	relations.Items = append(relations.Items, &rel)

	list := &ast.List{}
	for i, col := range c.AllColumn_name() {
		colName := col.GetText()
		target := &ast.ResTarget{
			Name: &colName,
			Val:  convert(c.Expr(i)),
		}
		list.Items = append(list.Items, target)
	}

	var where ast.Node = nil
	if c.WHERE_() != nil {
		where = convert(c.Expr(len(c.AllExpr()) - 1))
	}

	return &ast.UpdateStmt{
		Relations:     relations,
		TargetList:    list,
		WhereClause:   where,
		ReturningList: &ast.List{},
		FromClause:    &ast.List{},
		WithClause:    nil, // TODO: support with clause
	}
}

func convert(node node) ast.Node {
	switch n := node.(type) {

	case *parser.Alter_table_stmtContext:
		return convertAlter_table_stmtContext(n)

	case *parser.Attach_stmtContext:
		return convertAttach_stmtContext(n)

	case *parser.Create_table_stmtContext:
		return convertCreate_table_stmtContext(n)

	case *parser.Create_view_stmtContext:
		return convertCreate_view_stmtContext(n)

	case *parser.Drop_stmtContext:
		return convertDrop_stmtContext(n)

	case *parser.Delete_stmtContext:
		return convertDelete_stmtContext(n)

	case *parser.ExprContext:
		return convertExprContext(n)

	case *parser.Expr_functionContext:
		return convertFuncContext(n)

	case *parser.Expr_qualified_column_nameContext:
		return convertColumnNameExpr(n)

	case *parser.Expr_comparisonContext:
		return convertComparison(n)

	case *parser.Expr_bindContext:
		return convertParam(n)

	case *parser.Expr_literalContext:
		return convertLiteral(n)

	case *parser.Expr_binaryContext:
		return convertBinaryNode(n)

	case *parser.Expr_math_opContext:
		return convertMathOperationNode(n)

	case *parser.Factored_select_stmtContext:
		// TODO: need to handle this
		return &ast.TODO{}

	case *parser.Insert_stmtContext:
		return convertInsert_stmtContext(n)

	case *parser.Order_by_stmtContext:
		return convertOrderby_stmtContext(n)

	case *parser.Select_stmtContext:
		return convertMultiSelect_stmtContext(n)

	case *parser.Sql_stmtContext:
		return convertSql_stmtContext(n)

	case *parser.Update_stmtContext:
		return convertUpdate_stmtContext(n)

	default:
		return &ast.TODO{}
	}
}
