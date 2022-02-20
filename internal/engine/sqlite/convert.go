package sqlite

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"strconv"
	"strings"

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
	if c.TABLE_() != nil {
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

		var args []ast.Node
		for _, exp := range c.AllExpr() {
			args = append(args, convert(exp))
		}

		fn := &ast.FuncCall{
			Func: &ast.FuncName{
				Name: funcName,
			},
			Funcname: &ast.List{
				Items: []ast.Node{
					NewIdentifer(funcName),
				},
			},
			AggStar:     c.STAR() != nil,
			Args:        &ast.List{Items: args},
			AggOrder:    &ast.List{},
			AggDistinct: c.DISTINCT_() != nil,
		}

		return fn
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

func convertSimpleSelect_stmtContext(c *parser.Simple_select_stmtContext) ast.Node {
	if core, ok := c.Select_core().(*parser.Select_coreContext); ok {
		cols := getCols(core)
		tables := getTables(core)

		return &ast.SelectStmt{
			FromClause: &ast.List{Items: tables},
			TargetList: &ast.List{Items: cols},
		}
	}

	return &ast.TODO{}
}

func convertMultiSelect_stmtContext(c multiselect) ast.Node {
	var tables []ast.Node
	var cols []ast.Node
	for _, icore := range c.AllSelect_core() {
		core, ok := icore.(*parser.Select_coreContext)
		if !ok {
			continue
		}
		cols = append(cols, getCols(core)...)
		tables = append(tables, getTables(core)...)
	}
	return &ast.SelectStmt{
		FromClause: &ast.List{Items: tables},
		TargetList: &ast.List{Items: cols},
	}
}

func getTables(core *parser.Select_coreContext) []ast.Node {
	var tables []ast.Node
	for _, ifrom := range core.AllTable_or_subquery() {
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
			text := from.Schema_name().GetText()
			name.Schemaname = &text
		}
		tables = append(tables, &name)
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
			val = &ast.ColumnRef{
				Fields: &ast.List{
					Items: []ast.Node{
						&ast.A_Star{},
					},
				},
				Location: col.GetStart().GetStart(),
			}
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

func convert(node node) ast.Node {
	switch n := node.(type) {

	case *parser.Alter_table_stmtContext:
		return convertAlter_table_stmtContext(n)

	case *parser.Attach_stmtContext:
		return convertAttach_stmtContext(n)

	case *parser.Create_table_stmtContext:
		return convertCreate_table_stmtContext(n)

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

	case *parser.Select_stmtContext:
		return convertMultiSelect_stmtContext(n)

	case *parser.Sql_stmtContext:
		return convertSql_stmtContext(n)

	case *parser.Simple_select_stmtContext:
		return convertSimpleSelect_stmtContext(n)

	case *parser.Compound_select_stmtContext:
		return convertMultiSelect_stmtContext(n)

	default:
		return &ast.TODO{}
	}
}
