package sqlite

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/engine/sqlite/parser"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

type cc struct {
	paramCount int
}

type node interface {
	GetParser() antlr.Parser
}

func todo(funcname string, n node) *ast.TODO {
	if debug.Active {
		log.Printf("sqlite.%s: Unknown node type %T\n", funcname, n)
	}
	return &ast.TODO{}
}

func identifier(id string) string {
	if len(id) >= 2 && id[0] == '"' && id[len(id)-1] == '"' {
		unquoted, _ := strconv.Unquote(id)
		return unquoted
	}
	return strings.ToLower(id)
}

func NewIdentifier(t string) *ast.String {
	return &ast.String{Str: identifier(t)}
}

func (c *cc) convertAlter_table_stmtContext(n *parser.Alter_table_stmtContext) ast.Node {
	if n.RENAME_() != nil {
		if newTable, ok := n.New_table_name().(*parser.New_table_nameContext); ok {
			name := identifier(newTable.Any_name().GetText())
			return &ast.RenameTableStmt{
				Table:   parseTableName(n),
				NewName: &name,
			}
		}

		if newCol, ok := n.GetNew_column_name().(*parser.Column_nameContext); ok {
			name := identifier(newCol.Any_name().GetText())
			return &ast.RenameColumnStmt{
				Table: parseTableName(n),
				Col: &ast.ColumnRef{
					Name: identifier(n.GetOld_column_name().GetText()),
				},
				NewName: &name,
			}
		}
	}

	if n.ADD_() != nil {
		if def, ok := n.Column_def().(*parser.Column_defContext); ok {
			stmt := &ast.AlterTableStmt{
				Table: parseTableName(n),
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

	if n.DROP_() != nil {
		stmt := &ast.AlterTableStmt{
			Table: parseTableName(n),
			Cmds:  &ast.List{},
		}
		name := n.Column_name(0).GetText()
		stmt.Cmds.Items = append(stmt.Cmds.Items, &ast.AlterTableCmd{
			Name:    &name,
			Subtype: ast.AT_DropColumn,
		})
		return stmt
	}

	return todo("convertAlter_table_stmtContext", n)
}

func (c *cc) convertAttach_stmtContext(n *parser.Attach_stmtContext) ast.Node {
	name := n.Schema_name().GetText()
	return &ast.CreateSchemaStmt{
		Name: &name,
	}
}

func (c *cc) convertCreate_table_stmtContext(n *parser.Create_table_stmtContext) ast.Node {
	stmt := &ast.CreateTableStmt{
		Name:        parseTableName(n),
		IfNotExists: n.EXISTS_() != nil,
	}
	for _, idef := range n.AllColumn_def() {
		if def, ok := idef.(*parser.Column_defContext); ok {
			typeName := "any"
			if def.Type_name() != nil {
				typeName = def.Type_name().GetText()
			}
			stmt.Cols = append(stmt.Cols, &ast.ColumnDef{
				Colname:   identifier(def.Column_name().GetText()),
				IsNotNull: hasNotNullConstraint(def.AllColumn_constraint()),
				TypeName:  &ast.TypeName{Name: typeName},
			})
		}
	}
	return stmt
}

func (c *cc) convertCreate_virtual_table_stmtContext(n *parser.Create_virtual_table_stmtContext) ast.Node {
	switch moduleName := n.Module_name().GetText(); moduleName {
	case "fts5":
		// https://www.sqlite.org/fts5.html
		return c.convertCreate_virtual_table_fts5(n)
	default:
		return todo(
			fmt.Sprintf("create_virtual_table. unsupported module name: %q", moduleName),
			n,
		)
	}
}

func (c *cc) convertCreate_virtual_table_fts5(n *parser.Create_virtual_table_stmtContext) ast.Node {
	stmt := &ast.CreateTableStmt{
		Name:        parseTableName(n),
		IfNotExists: n.EXISTS_() != nil,
	}

	for _, arg := range n.AllModule_argument() {
		var columnName string

		// For example: CREATE VIRTUAL TABLE tbl_ft USING fts5(b, c UNINDEXED)
		//   * the 'b' column is parsed like Expr_qualified_column_nameContext
		//   * the 'c' column is parsed like Column_defContext
		if columnExpr, ok := arg.Expr().(*parser.Expr_qualified_column_nameContext); ok {
			columnName = columnExpr.Column_name().GetText()
		} else if columnDef, ok := arg.Column_def().(*parser.Column_defContext); ok {
			columnName = columnDef.Column_name().GetText()
		}

		if columnName != "" {
			stmt.Cols = append(stmt.Cols, &ast.ColumnDef{
				Colname: identifier(columnName),
				// you can not specify any column constraints in fts5, so we pass them manually
				IsNotNull: true,
				TypeName:  &ast.TypeName{Name: "text"},
			})
		}
	}

	return stmt
}

func (c *cc) convertCreate_view_stmtContext(n *parser.Create_view_stmtContext) ast.Node {
	viewName := n.View_name().GetText()
	relation := &ast.RangeVar{
		Relname: &viewName,
	}

	if n.Schema_name() != nil {
		schemaName := n.Schema_name().GetText()
		relation.Schemaname = &schemaName
	}

	return &ast.ViewStmt{
		View:            relation,
		Aliases:         &ast.List{},
		Query:           c.convert(n.Select_stmt()),
		Replace:         false,
		Options:         &ast.List{},
		WithCheckOption: ast.ViewCheckOption(0),
	}
}

type Delete_stmt interface {
	node

	Qualified_table_name() parser.IQualified_table_nameContext
	WHERE_() antlr.TerminalNode
	Expr() parser.IExprContext
}

func (c *cc) convertDelete_stmtContext(n Delete_stmt) ast.Node {
	if qualifiedName, ok := n.Qualified_table_name().(*parser.Qualified_table_nameContext); ok {

		tableName := identifier(qualifiedName.Table_name().GetText())
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

		relations := &ast.List{}

		relations.Items = append(relations.Items, relation)

		delete := &ast.DeleteStmt{
			Relations:  relations,
			WithClause: nil,
		}

		if n.WHERE_() != nil && n.Expr() != nil {
			delete.WhereClause = c.convert(n.Expr())
		}

		if n, ok := n.(interface {
			Returning_clause() parser.IReturning_clauseContext
		}); ok {
			delete.ReturningList = c.convertReturning_caluseContext(n.Returning_clause())
		} else {
			delete.ReturningList = c.convertReturning_caluseContext(nil)
		}
		if n, ok := n.(interface {
			Limit_stmt() parser.ILimit_stmtContext
		}); ok {
			limitCount, _ := c.convertLimit_stmtContext(n.Limit_stmt())
			delete.LimitCount = limitCount
		}

		return delete
	}

	return todo("convertDelete_stmtContext", n)
}

func (c *cc) convertDrop_stmtContext(n *parser.Drop_stmtContext) ast.Node {
	if n.TABLE_() != nil || n.VIEW_() != nil {
		name := ast.TableName{
			Name: identifier(n.Any_name().GetText()),
		}
		if n.Schema_name() != nil {
			name.Schema = n.Schema_name().GetText()
		}

		return &ast.DropTableStmt{
			IfExists: n.EXISTS_() != nil,
			Tables:   []*ast.TableName{&name},
		}
	}
	return todo("convertDrop_stmtContext", n)
}

func (c *cc) convertFuncContext(n *parser.Expr_functionContext) ast.Node {
	if name, ok := n.Qualified_function_name().(*parser.Qualified_function_nameContext); ok {
		funcName := strings.ToLower(name.Function_name().GetText())

		schema := ""
		if name.Schema_name() != nil {
			schema = name.Schema_name().GetText()
		}

		var argNodes []ast.Node
		for _, exp := range n.AllExpr() {
			argNodes = append(argNodes, c.convert(exp))
		}
		args := &ast.List{Items: argNodes}

		if funcName == "coalesce" {
			return &ast.CoalesceExpr{
				Args:     args,
				Location: name.GetStart().GetStart(),
			}
		} else {
			return &ast.FuncCall{
				Func: &ast.FuncName{
					Schema: schema,
					Name:   funcName,
				},
				Funcname: &ast.List{
					Items: []ast.Node{
						NewIdentifier(funcName),
					},
				},
				AggStar:     n.STAR() != nil,
				Args:        args,
				AggOrder:    &ast.List{},
				AggDistinct: n.DISTINCT_() != nil,
				Location:    name.GetStart().GetStart(),
			}
		}
	}

	return todo("convertFuncContext", n)
}

func (c *cc) convertExprContext(n *parser.ExprContext) ast.Node {
	return &ast.Expr{}
}

func (c *cc) convertColumnNameExpr(n *parser.Expr_qualified_column_nameContext) *ast.ColumnRef {
	var items []ast.Node
	if schema, ok := n.Schema_name().(*parser.Schema_nameContext); ok {
		schemaText := schema.GetText()
		if schemaText != "" {
			items = append(items, NewIdentifier(schemaText))
		}
	}
	if table, ok := n.Table_name().(*parser.Table_nameContext); ok {
		tableName := table.GetText()
		if tableName != "" {
			items = append(items, NewIdentifier(tableName))
		}
	}
	items = append(items, NewIdentifier(n.Column_name().GetText()))
	return &ast.ColumnRef{
		Fields: &ast.List{
			Items: items,
		},
		Location: n.GetStart().GetStart(),
	}
}

func (c *cc) convertComparison(n *parser.Expr_comparisonContext) ast.Node {
	lexpr := c.convert(n.Expr(0))

	if n.IN_() != nil {
		rexprs := []ast.Node{}
		for _, expr := range n.AllExpr()[1:] {
			e := c.convert(expr)
			switch t := e.(type) {
			case *ast.List:
				rexprs = append(rexprs, t.Items...)
			default:
				rexprs = append(rexprs, t)
			}
		}

		return &ast.In{
			Expr:     lexpr,
			List:     rexprs,
			Not:      false,
			Sel:      nil,
			Location: n.GetStart().GetStart(),
		}
	}

	return &ast.A_Expr{
		Name: &ast.List{
			Items: []ast.Node{
				&ast.String{Str: "="}, // TODO: add actual comparison
			},
		},
		Lexpr: lexpr,
		Rexpr: c.convert(n.Expr(1)),
	}
}

func (c *cc) convertMultiSelect_stmtContext(n *parser.Select_stmtContext) ast.Node {
	var ctes ast.List
	if ct := n.Common_table_stmt(); ct != nil {
		recursive := ct.RECURSIVE_() != nil
		for _, cte := range ct.AllCommon_table_expression() {
			tableName := identifier(cte.Table_name().GetText())
			var cteCols ast.List
			for _, col := range cte.AllColumn_name() {
				cteCols.Items = append(cteCols.Items, NewIdentifier(col.GetText()))
			}
			ctes.Items = append(ctes.Items, &ast.CommonTableExpr{
				Ctename:      &tableName,
				Ctequery:     c.convert(cte.Select_stmt()),
				Location:     cte.GetStart().GetStart(),
				Cterecursive: recursive,
				Ctecolnames:  &cteCols,
			})
		}
	}

	var selectStmt *ast.SelectStmt
	for s, icore := range n.AllSelect_core() {
		core, ok := icore.(*parser.Select_coreContext)
		if !ok {
			continue
		}
		cols := c.getCols(core)
		tables := c.getTables(core)

		var where ast.Node
		i := 0
		if core.WHERE_() != nil {
			where = c.convert(core.Expr(i))
			i++
		}

		var groups ast.List
		var having ast.Node
		if core.GROUP_() != nil {
			l := len(core.AllExpr()) - i
			if core.HAVING_() != nil {
				having = c.convert(core.Expr(l))
				l--
			}

			for i < l {
				groups.Items = append(groups.Items, c.convert(core.Expr(i)))
				i++
			}
		}
		var window ast.List
		if core.WINDOW_() != nil {
			for w, windowNameCtx := range core.AllWindow_name() {
				windowName := identifier(windowNameCtx.GetText())
				windowDef := core.Window_defn(w)

				_ = windowDef.Base_window_name()
				var partitionBy ast.List
				if windowDef.PARTITION_() != nil {
					for _, e := range windowDef.AllExpr() {
						partitionBy.Items = append(partitionBy.Items, c.convert(e))
					}
				}
				var orderBy ast.List
				if windowDef.ORDER_() != nil {
					for _, e := range windowDef.AllOrdering_term() {
						oterm := e.(*parser.Ordering_termContext)
						sortByDir := ast.SortByDirDefault
						if ad := oterm.Asc_desc(); ad != nil {
							if ad.ASC_() != nil {
								sortByDir = ast.SortByDirAsc
							} else {
								sortByDir = ast.SortByDirDesc
							}
						}
						sortByNulls := ast.SortByNullsDefault
						if oterm.NULLS_() != nil {
							if oterm.FIRST_() != nil {
								sortByNulls = ast.SortByNullsFirst
							} else {
								sortByNulls = ast.SortByNullsLast
							}
						}

						orderBy.Items = append(orderBy.Items, &ast.SortBy{
							Node:        c.convert(oterm.Expr()),
							SortbyDir:   sortByDir,
							SortbyNulls: sortByNulls,
							UseOp:       &ast.List{},
						})
					}
				}
				window.Items = append(window.Items, &ast.WindowDef{
					Name:            &windowName,
					PartitionClause: &partitionBy,
					OrderClause:     &orderBy,
					FrameOptions:    0, // todo
					StartOffset:     &ast.TODO{},
					EndOffset:       &ast.TODO{},
					Location:        windowNameCtx.GetStart().GetStart(),
				})
			}
		}
		sel := &ast.SelectStmt{
			FromClause:   &ast.List{Items: tables},
			TargetList:   &ast.List{Items: cols},
			WhereClause:  where,
			GroupClause:  &groups,
			HavingClause: having,
			WindowClause: &window,
			ValuesLists:  &ast.List{},
		}
		if selectStmt == nil {
			selectStmt = sel
		} else {
			co := n.Compound_operator(s - 1)
			so := ast.None
			all := false
			switch {
			case co.UNION_() != nil:
				so = ast.Union
				all = co.ALL_() != nil
			case co.INTERSECT_() != nil:
				so = ast.Intersect
			case co.EXCEPT_() != nil:
				so = ast.Except
			}
			selectStmt = &ast.SelectStmt{
				TargetList: &ast.List{},
				FromClause: &ast.List{},
				Op:         so,
				All:        all,
				Larg:       selectStmt,
				Rarg:       sel,
			}
		}
	}

	limitCount, limitOffset := c.convertLimit_stmtContext(n.Limit_stmt())
	selectStmt.LimitCount = limitCount
	selectStmt.LimitOffset = limitOffset
	selectStmt.WithClause = &ast.WithClause{Ctes: &ctes}
	return selectStmt
}

func (c *cc) convertExprListContext(n *parser.Expr_listContext) ast.Node {
	list := &ast.List{Items: []ast.Node{}}
	for _, e := range n.AllExpr() {
		list.Items = append(list.Items, c.convert(e))
	}
	return list
}

func (c *cc) getTables(core *parser.Select_coreContext) []ast.Node {
	if core.Join_clause() != nil {
		join := core.Join_clause().(*parser.Join_clauseContext)
		tables := c.convertTablesOrSubquery(join.AllTable_or_subquery())
		table := tables[0]
		for i, t := range tables[1:] {
			joinExpr := &ast.JoinExpr{
				Larg: table,
				Rarg: t,
			}
			jo := join.Join_operator(i)
			if jo.NATURAL_() != nil {
				joinExpr.IsNatural = true
			}
			switch {
			case jo.CROSS_() != nil || jo.INNER_() != nil:
				joinExpr.Jointype = ast.JoinTypeInner
			case jo.LEFT_() != nil:
				joinExpr.Jointype = ast.JoinTypeLeft
			case jo.RIGHT_() != nil:
				joinExpr.Jointype = ast.JoinTypeRight
			case jo.FULL_() != nil:
				joinExpr.Jointype = ast.JoinTypeFull
			}
			jc := join.Join_constraint(i)
			switch {
			case jc.ON_() != nil:
				joinExpr.Quals = c.convert(jc.Expr())
			case jc.USING_() != nil:
				var using ast.List
				for _, cn := range jc.AllColumn_name() {
					using.Items = append(using.Items, NewIdentifier(cn.GetText()))
				}
				joinExpr.UsingClause = &using
			}
			table = joinExpr
		}
		return []ast.Node{table}
	} else {
		return c.convertTablesOrSubquery(core.AllTable_or_subquery())
	}
}

func (c *cc) getCols(core *parser.Select_coreContext) []ast.Node {
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
			val = c.convertWildCardField(col)
		case iexpr != nil:
			val = c.convert(iexpr)
		}

		if val == nil {
			continue
		}

		if col.Column_alias() != nil {
			name := identifier(col.Column_alias().GetText())
			target.Name = &name
		}

		target.Val = val
		cols = append(cols, target)
	}
	return cols
}

func (c *cc) convertWildCardField(n *parser.Result_columnContext) *ast.ColumnRef {
	items := []ast.Node{}
	if n.Table_name() != nil {
		items = append(items, NewIdentifier(n.Table_name().GetText()))
	}
	items = append(items, &ast.A_Star{})

	return &ast.ColumnRef{
		Fields: &ast.List{
			Items: items,
		},
		Location: n.GetStart().GetStart(),
	}
}

func (c *cc) convertOrderby_stmtContext(n parser.IOrder_by_stmtContext) ast.Node {
	if orderBy, ok := n.(*parser.Order_by_stmtContext); ok {
		list := &ast.List{Items: []ast.Node{}}
		for _, o := range orderBy.AllOrdering_term() {
			term, ok := o.(*parser.Ordering_termContext)
			if !ok {
				continue
			}
			list.Items = append(list.Items, &ast.CaseExpr{
				Xpr:      c.convert(term.Expr()),
				Location: term.Expr().GetStart().GetStart(),
			})
		}
		return list
	}
	return todo("convertOrderby_stmtContext", n)
}

func (c *cc) convertLimit_stmtContext(n parser.ILimit_stmtContext) (ast.Node, ast.Node) {
	if n == nil {
		return nil, nil
	}

	var limitCount, limitOffset ast.Node
	if limit, ok := n.(*parser.Limit_stmtContext); ok {
		limitCount = c.convert(limit.Expr(0))
		if limit.OFFSET_() != nil {
			limitOffset = c.convert(limit.Expr(1))
		}
	}

	return limitCount, limitOffset
}

func (c *cc) convertSql_stmtContext(n *parser.Sql_stmtContext) ast.Node {
	if stmt := n.Alter_table_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Analyze_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Attach_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Begin_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Commit_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Create_index_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Create_table_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Create_trigger_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Create_view_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Create_virtual_table_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Delete_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Delete_stmt_limited(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Detach_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Drop_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Insert_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Pragma_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Reindex_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Release_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Rollback_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Savepoint_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Select_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Update_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Update_stmt_limited(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Vacuum_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	return nil
}

func (c *cc) convertLiteral(n *parser.Expr_literalContext) ast.Node {
	if literal, ok := n.Literal_value().(*parser.Literal_valueContext); ok {

		if literal.NUMERIC_LITERAL() != nil {
			i, _ := strconv.ParseInt(literal.GetText(), 10, 64)
			return &ast.A_Const{
				Val:      &ast.Integer{Ival: i},
				Location: n.GetStart().GetStart(),
			}
		}

		if literal.STRING_LITERAL() != nil {
			// remove surrounding single quote
			text := literal.GetText()
			return &ast.A_Const{
				Val:      &ast.String{Str: text[1 : len(text)-1]},
				Location: n.GetStart().GetStart(),
			}
		}

		if literal.TRUE_() != nil || literal.FALSE_() != nil {
			var i int64
			if literal.TRUE_() != nil {
				i = 1
			}

			return &ast.A_Const{
				Val:      &ast.Integer{Ival: i},
				Location: n.GetStart().GetStart(),
			}
		}
	}
	return todo("convertLiteral", n)
}

func (c *cc) convertBinaryNode(n *parser.Expr_binaryContext) ast.Node {
	return &ast.A_Expr{
		Name: &ast.List{
			Items: []ast.Node{
				&ast.String{Str: n.GetChild(1).(antlr.TerminalNode).GetText()},
			},
		},
		Lexpr: c.convert(n.Expr(0)),
		Rexpr: c.convert(n.Expr(1)),
	}
}

func (c *cc) convertBoolNode(n *parser.Expr_boolContext) ast.Node {
	return &ast.BoolExpr{
		// TODO: Set op
		Args: &ast.List{
			Items: []ast.Node{
				c.convert(n.Expr(0)),
				c.convert(n.Expr(1)),
			},
		},
	}
}

func (c *cc) convertParam(n *parser.Expr_bindContext) ast.Node {
	if n.NUMBERED_BIND_PARAMETER() != nil {
		// Parameter numbers start at one
		c.paramCount += 1

		text := n.GetText()
		number := c.paramCount
		if len(text) > 1 {
			number, _ = strconv.Atoi(text[1:])
		}
		return &ast.ParamRef{
			Number:   number,
			Location: n.GetStart().GetStart(),
			Dollar:   len(text) > 1,
		}
	}

	if n.NAMED_BIND_PARAMETER() != nil {
		return &ast.A_Expr{
			Name:     &ast.List{Items: []ast.Node{&ast.String{Str: "@"}}},
			Rexpr:    &ast.String{Str: n.GetText()[1:]},
			Location: n.GetStart().GetStart(),
		}
	}

	return todo("convertParam", n)
}

func (c *cc) convertInSelectNode(n *parser.Expr_in_selectContext) ast.Node {
	return c.convert(n.Select_stmt())
}

func (c *cc) convertReturning_caluseContext(n parser.IReturning_clauseContext) *ast.List {
	list := &ast.List{Items: []ast.Node{}}
	if n == nil {
		return list
	}

	r, ok := n.(*parser.Returning_clauseContext)
	if !ok {
		return list
	}

	for _, exp := range r.AllExpr() {
		list.Items = append(list.Items, &ast.ResTarget{
			Indirection: &ast.List{},
			Val:         c.convert(exp),
		})
	}

	for _, star := range r.AllSTAR() {
		list.Items = append(list.Items, &ast.ResTarget{
			Indirection: &ast.List{},
			Val: &ast.ColumnRef{
				Fields: &ast.List{
					Items: []ast.Node{&ast.A_Star{}},
				},
				Location: star.GetSymbol().GetStart(),
			},
			Location: star.GetSymbol().GetStart(),
		})
	}

	return list
}

func (c *cc) convertInsert_stmtContext(n *parser.Insert_stmtContext) ast.Node {
	tableName := identifier(n.Table_name().GetText())
	rel := &ast.RangeVar{
		Relname: &tableName,
	}
	if n.Schema_name() != nil {
		schemaName := n.Schema_name().GetText()
		rel.Schemaname = &schemaName
	}
	if n.Table_alias() != nil {
		tableAlias := identifier(n.Table_alias().GetText())
		rel.Alias = &ast.Alias{
			Aliasname: &tableAlias,
		}
	}

	insert := &ast.InsertStmt{
		Relation:      rel,
		Cols:          c.convertColumnNames(n.AllColumn_name()),
		ReturningList: c.convertReturning_caluseContext(n.Returning_clause()),
	}

	// Check if this is a DEFAULT VALUES insert
	hasDefaultValues := false
	for _, child := range n.GetChildren() {
		if term, ok := child.(antlr.TerminalNode); ok {
			if term.GetSymbol().GetTokenType() == parser.SQLiteParserDEFAULT_ {
				hasDefaultValues = true
				break
			}
		}
	}

	if hasDefaultValues {
		// For DEFAULT VALUES, create an empty select statement
		insert.SelectStmt = &ast.SelectStmt{
			FromClause:  &ast.List{},
			TargetList:  &ast.List{},
			ValuesLists: &ast.List{Items: []ast.Node{&ast.List{}}}, // Single empty values list
		}
	} else if n.Select_stmt() != nil {
		if ss, ok := c.convert(n.Select_stmt()).(*ast.SelectStmt); ok {
			ss.ValuesLists = &ast.List{}
			insert.SelectStmt = ss
		}
	} else {
		var valuesLists ast.List
		var values *ast.List
		for _, cn := range n.GetChildren() {
			switch cn := cn.(type) {
			case antlr.TerminalNode:
				switch cn.GetSymbol().GetTokenType() {
				case parser.SQLiteParserVALUES_:
					values = &ast.List{}
				case parser.SQLiteParserOPEN_PAR:
					if values != nil {
						values = &ast.List{}
					}
				case parser.SQLiteParserCOMMA:
				case parser.SQLiteParserCLOSE_PAR:
					if values != nil {
						valuesLists.Items = append(valuesLists.Items, values)
					}
				}
			case parser.IExprContext:
				if values != nil {
					values.Items = append(values.Items, c.convert(cn))
				}
			}
		}

		insert.SelectStmt = &ast.SelectStmt{
			FromClause:  &ast.List{},
			TargetList:  &ast.List{},
			ValuesLists: &valuesLists,
		}
	}

	return insert
}

func (c *cc) convertColumnNames(cols []parser.IColumn_nameContext) *ast.List {
	list := &ast.List{Items: []ast.Node{}}
	for _, c := range cols {
		name := identifier(c.GetText())
		list.Items = append(list.Items, &ast.ResTarget{
			Name: &name,
		})
	}
	return list
}

func (c *cc) convertTablesOrSubquery(n []parser.ITable_or_subqueryContext) []ast.Node {
	var tables []ast.Node
	for _, ifrom := range n {
		from, ok := ifrom.(*parser.Table_or_subqueryContext)
		if !ok {
			continue
		}

		if from.Table_name() != nil {
			rel := identifier(from.Table_name().GetText())
			rv := &ast.RangeVar{
				Relname:  &rel,
				Location: from.GetStart().GetStart(),
			}

			if from.Schema_name() != nil {
				schema := from.Schema_name().GetText()
				rv.Schemaname = &schema
			}
			if from.Table_alias() != nil {
				alias := identifier(from.Table_alias().GetText())
				rv.Alias = &ast.Alias{Aliasname: &alias}
			}
			if from.Table_alias_fallback() != nil {
				alias := identifier(from.Table_alias_fallback().GetText())
				rv.Alias = &ast.Alias{Aliasname: &alias}
			}

			tables = append(tables, rv)
		} else if from.Table_function_name() != nil {
			rel := from.Table_function_name().GetText()
			rf := &ast.RangeFunction{
				Functions: &ast.List{
					Items: []ast.Node{
						&ast.FuncCall{
							Func: &ast.FuncName{
								Name: rel,
							},
							Funcname: &ast.List{
								Items: []ast.Node{
									NewIdentifier(rel),
								},
							},
							Args: &ast.List{
								Items: []ast.Node{&ast.TODO{}},
							},
							Location: from.GetStart().GetStart(),
						},
					},
				},
			}

			if from.Table_alias() != nil {
				alias := identifier(from.Table_alias().GetText())
				rf.Alias = &ast.Alias{Aliasname: &alias}
			}

			tables = append(tables, rf)
		} else if from.Select_stmt() != nil {
			rs := &ast.RangeSubselect{
				Subquery: c.convert(from.Select_stmt()),
			}

			if from.Table_alias() != nil {
				alias := identifier(from.Table_alias().GetText())
				rs.Alias = &ast.Alias{Aliasname: &alias}
			}

			tables = append(tables, rs)
		}
	}

	return tables
}

type Update_stmt interface {
	Qualified_table_name() parser.IQualified_table_nameContext
	GetStart() antlr.Token
	AllColumn_name() []parser.IColumn_nameContext
	WHERE_() antlr.TerminalNode
	Expr(i int) parser.IExprContext
	AllExpr() []parser.IExprContext
}

func (c *cc) convertUpdate_stmtContext(n Update_stmt) ast.Node {
	if n == nil {
		return nil
	}

	relations := &ast.List{}
	tableName := identifier(n.Qualified_table_name().GetText())
	rel := ast.RangeVar{
		Relname:  &tableName,
		Location: n.GetStart().GetStart(),
	}
	relations.Items = append(relations.Items, &rel)

	list := &ast.List{}
	for i, col := range n.AllColumn_name() {
		colName := identifier(col.GetText())
		target := &ast.ResTarget{
			Name: &colName,
			Val:  c.convert(n.Expr(i)),
		}
		list.Items = append(list.Items, target)
	}

	var where ast.Node = nil
	if n.WHERE_() != nil {
		where = c.convert(n.Expr(len(n.AllExpr()) - 1))
	}

	stmt := &ast.UpdateStmt{
		Relations:   relations,
		TargetList:  list,
		WhereClause: where,
		FromClause:  &ast.List{},
		WithClause:  nil, // TODO: support with clause
	}
	if n, ok := n.(interface {
		Returning_clause() parser.IReturning_clauseContext
	}); ok {
		stmt.ReturningList = c.convertReturning_caluseContext(n.Returning_clause())
	} else {
		stmt.ReturningList = c.convertReturning_caluseContext(nil)
	}
	if n, ok := n.(interface {
		Limit_stmt() parser.ILimit_stmtContext
	}); ok {
		limitCount, _ := c.convertLimit_stmtContext(n.Limit_stmt())
		stmt.LimitCount = limitCount
	}
	return stmt
}

func (c *cc) convertBetweenExpr(n *parser.Expr_betweenContext) ast.Node {
	return &ast.BetweenExpr{
		Expr:     c.convert(n.Expr(0)),
		Left:     c.convert(n.Expr(1)),
		Right:    c.convert(n.Expr(2)),
		Location: n.GetStart().GetStart(),
		Not:      n.NOT_() != nil,
	}
}

func (c *cc) convertCastExpr(n *parser.Expr_castContext) ast.Node {
	name := n.Type_name().GetText()
	return &ast.TypeCast{
		Arg: c.convert(n.Expr()),
		TypeName: &ast.TypeName{
			Name: name,
			Names: &ast.List{Items: []ast.Node{
				NewIdentifier(name),
			}},
			ArrayBounds: &ast.List{},
		},
		Location: n.GetStart().GetStart(),
	}
}

func (c *cc) convertCollateExpr(n *parser.Expr_collateContext) ast.Node {
	return &ast.CollateExpr{
		Xpr:      c.convert(n.Expr()),
		Arg:      NewIdentifier(n.Collation_name().GetText()),
		Location: n.GetStart().GetStart(),
	}
}

func (c *cc) convertCase(n *parser.Expr_caseContext) ast.Node {
	e := &ast.CaseExpr{
		Args: &ast.List{},
	}
	es := n.AllExpr()
	if n.ELSE_() != nil {
		e.Defresult = c.convert(es[len(es)-1])
		es = es[:len(es)-1]
	}
	if len(es)%2 == 1 {
		e.Arg = c.convert(es[0])
		es = es[1:]
	}
	for i := 0; i < len(es); i += 2 {
		e.Args.Items = append(e.Args.Items, &ast.CaseWhen{
			Expr:   c.convert(es[i+0]),
			Result: c.convert(es[i+1]),
		})
	}
	return e
}

func (c *cc) convert(node node) ast.Node {
	switch n := node.(type) {

	case *parser.Alter_table_stmtContext:
		return c.convertAlter_table_stmtContext(n)

	case *parser.Attach_stmtContext:
		return c.convertAttach_stmtContext(n)

	case *parser.Create_table_stmtContext:
		return c.convertCreate_table_stmtContext(n)

	case *parser.Create_virtual_table_stmtContext:
		return c.convertCreate_virtual_table_stmtContext(n)

	case *parser.Create_view_stmtContext:
		return c.convertCreate_view_stmtContext(n)

	case *parser.Drop_stmtContext:
		return c.convertDrop_stmtContext(n)

	case *parser.Delete_stmtContext:
		return c.convertDelete_stmtContext(n)

	case *parser.Delete_stmt_limitedContext:
		return c.convertDelete_stmtContext(n)

	case *parser.ExprContext:
		return c.convertExprContext(n)

	case *parser.Expr_functionContext:
		return c.convertFuncContext(n)

	case *parser.Expr_qualified_column_nameContext:
		return c.convertColumnNameExpr(n)

	case *parser.Expr_comparisonContext:
		return c.convertComparison(n)

	case *parser.Expr_bindContext:
		return c.convertParam(n)

	case *parser.Expr_literalContext:
		return c.convertLiteral(n)

	case *parser.Expr_boolContext:
		return c.convertBoolNode(n)

	case *parser.Expr_listContext:
		return c.convertExprListContext(n)

	case *parser.Expr_binaryContext:
		return c.convertBinaryNode(n)

	case *parser.Expr_in_selectContext:
		return c.convertInSelectNode(n)

	case *parser.Expr_betweenContext:
		return c.convertBetweenExpr(n)

	case *parser.Expr_collateContext:
		return c.convertCollateExpr(n)

	case *parser.Factored_select_stmtContext:
		// TODO: need to handle this
		return todo("convert(case=parser.Factored_select_stmtContext)", n)

	case *parser.Insert_stmtContext:
		return c.convertInsert_stmtContext(n)

	case *parser.Order_by_stmtContext:
		return c.convertOrderby_stmtContext(n)

	case *parser.Select_stmtContext:
		return c.convertMultiSelect_stmtContext(n)

	case *parser.Sql_stmtContext:
		return c.convertSql_stmtContext(n)

	case *parser.Update_stmtContext:
		return c.convertUpdate_stmtContext(n)

	case *parser.Update_stmt_limitedContext:
		return c.convertUpdate_stmtContext(n)

	case *parser.Expr_castContext:
		return c.convertCastExpr(n)

	case *parser.Expr_caseContext:
		return c.convertCase(n)

	default:
		return todo("convert(case=default)", n)
	}
}
