package sqlite

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"

	"github.com/kyleconroy/sqlc/internal/engine/sqlite/parser"
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type node interface {
	GetParser() antlr.Parser
}

func convertAlter_table_stmtContext(c *parser.Alter_table_stmtContext) ast.Node {
	if newTable, ok := c.New_table_name().(*parser.New_table_nameContext); ok {
		name := newTable.Any_name().GetText()
		return &ast.RenameTableStmt{
			Table:   parseTableName(c),
			NewName: &name,
		}
	}

	if newCol, ok := c.New_column_name().(*parser.New_column_nameContext); ok {
		name := newCol.Any_name().GetText()
		return &ast.RenameColumnStmt{
			Table: parseTableName(c),
			Col: &ast.ColumnRef{
				Name: c.Column_name().GetText(),
			},
			NewName: &name,
		}
	}

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

	return &ast.TODO{}
}

func convertAttach_stmtContext(c *parser.Attach_stmtContext) ast.Node {
	name := c.Database_name().GetText()
	return &ast.CreateSchemaStmt{
		Name: &name,
	}
}

func convertCreate_table_stmtContext(c *parser.Create_table_stmtContext) ast.Node {
	stmt := &ast.CreateTableStmt{
		Name:        parseTableName(c),
		IfNotExists: c.K_EXISTS() != nil,
	}
	for _, idef := range c.AllColumn_def() {
		if def, ok := idef.(*parser.Column_defContext); ok {
			stmt.Cols = append(stmt.Cols, &ast.ColumnDef{
				Colname:   def.Column_name().GetText(),
				IsNotNull: hasNotNullConstraint(def.AllColumn_constraint()),
				TypeName:  &ast.TypeName{Name: def.Type_name().GetText()},
			})
		}
	}
	return stmt
}

func convertDrop_table_stmtContext(c *parser.Drop_table_stmtContext) ast.Node {
	return &ast.DropTableStmt{
		IfExists: c.K_EXISTS() != nil,
		Tables:   []*ast.TableName{parseTableName(c)},
	}
}

func convertExprContext(c *parser.ExprContext) ast.Node {
	return &ast.TODO{}
}

func convertFactored_select_stmtContext(c *parser.Factored_select_stmtContext) ast.Node {
	var tables []ast.Node
	var cols []ast.Node
	for _, icore := range c.AllSelect_core() {
		core, ok := icore.(*parser.Select_coreContext)
		if !ok {
			continue
		}
		for _, icol := range core.AllResult_column() {
			col, ok := icol.(*parser.Result_columnContext)
			if !ok {
				continue
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
			cols = append(cols, &ast.ResTarget{
				Val:      val,
				Location: col.GetStart().GetStart(),
			})
		}
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
	}
	return &ast.SelectStmt{
		FromClause: &ast.List{Items: tables},
		TargetList: &ast.List{Items: cols},
	}
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
	if stmt := n.Compound_select_stmt(); stmt != nil {
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
	if stmt := n.Drop_index_stmt(); stmt != nil {
		return convert(stmt)
	}
	if stmt := n.Drop_table_stmt(); stmt != nil {
		return convert(stmt)
	}
	if stmt := n.Drop_trigger_stmt(); stmt != nil {
		return convert(stmt)
	}
	if stmt := n.Drop_view_stmt(); stmt != nil {
		return convert(stmt)
	}
	if stmt := n.Factored_select_stmt(); stmt != nil {
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
	if stmt := n.Simple_select_stmt(); stmt != nil {
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

func convert(node node) ast.Node {
	switch n := node.(type) {

	case *parser.Alter_table_stmtContext:
		return convertAlter_table_stmtContext(n)

	case *parser.Attach_stmtContext:
		return convertAttach_stmtContext(n)

	case *parser.Create_table_stmtContext:
		return convertCreate_table_stmtContext(n)

	case *parser.Drop_table_stmtContext:
		return convertDrop_table_stmtContext(n)

	case *parser.ExprContext:
		return convertExprContext(n)

	case *parser.Factored_select_stmtContext:
		return convertFactored_select_stmtContext(n)

	case *parser.Sql_stmtContext:
		return convertSql_stmtContext(n)

	default:
		return &ast.TODO{}
	}
}
