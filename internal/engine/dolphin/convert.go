package dolphin

import (
	"fmt"

	pcast "github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/opcode"
	driver "github.com/pingcap/parser/test_driver"
	"github.com/pingcap/parser/types"

	"github.com/kyleconroy/sqlc/internal/debug"
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type cc struct {
	paramCount int
}

func (c *cc) convertAlterTableStmt(n *pcast.AlterTableStmt) ast.Node {
	alt := &ast.AlterTableStmt{
		Table: parseTableName(n.Table),
		Cmds:  &ast.List{},
	}
	for _, spec := range n.Specs {
		switch spec.Tp {
		case pcast.AlterTableAddColumns:
			for _, def := range spec.NewColumns {
				name := def.Name.String()
				alt.Cmds.Items = append(alt.Cmds.Items, &ast.AlterTableCmd{
					Name:    &name,
					Subtype: ast.AT_AddColumn,
					Def: &ast.ColumnDef{
						Colname:   def.Name.String(),
						TypeName:  &ast.TypeName{Name: types.TypeStr(def.Tp.Tp)},
						IsNotNull: isNotNull(def),
					},
				})
			}

		case pcast.AlterTableDropColumn:
			name := spec.OldColumnName.String()
			alt.Cmds.Items = append(alt.Cmds.Items, &ast.AlterTableCmd{
				Name:    &name,
				Subtype: ast.AT_DropColumn,
				// MissingOk: spec.IfExists,
			})

		case pcast.AlterTableChangeColumn:
			// 	spew.Dump("change column", spec)

		case pcast.AlterTableModifyColumn:
			// 	spew.Dump("modify column", spec)

		case pcast.AlterTableAlterColumn:
			// 	spew.Dump("alter column", spec)

		case pcast.AlterTableAddConstraint:
			// 	spew.Dump("add const", spec)

		default:
			continue
		}
	}
	return alt
}

func (c *cc) convertAssignment(n *pcast.Assignment) *ast.ResTarget {
	name := n.Column.Name.String()
	return &ast.ResTarget{
		Name: &name,
		Val:  c.convert(n.Expr),
	}
}

func opToName(o opcode.Op) string {
	switch o {
	case opcode.EQ:
		return "="
	}
	return o.String()
}

func (c *cc) convertBinaryOperationExpr(n *pcast.BinaryOperationExpr) ast.Node {
	if n.Op == opcode.LogicAnd || n.Op == opcode.LogicOr {
		return &ast.BoolExpr{
			// TODO: Set op
			Args: &ast.List{
				Items: []ast.Node{
					c.convert(n.L),
					c.convert(n.R),
				},
			},
		}
	} else {
		return &ast.A_Expr{
			// TODO: Set kind
			Name: &ast.List{
				Items: []ast.Node{
					&ast.String{Str: opToName(n.Op)},
				},
			},
			Lexpr: c.convert(n.L),
			Rexpr: c.convert(n.R),
		}
	}
}

func (c *cc) convertCreateTableStmt(n *pcast.CreateTableStmt) ast.Node {
	create := &ast.CreateTableStmt{
		Name:        parseTableName(n.Table),
		IfNotExists: n.IfNotExists,
	}
	if n.ReferTable != nil {
		create.ReferTable = parseTableName(n.ReferTable)
	}
	for _, def := range n.Cols {
		var vals *ast.List
		if len(def.Tp.Elems) > 0 {
			vals = &ast.List{}
			for i := range def.Tp.Elems {
				vals.Items = append(vals.Items, &ast.String{
					Str: def.Tp.Elems[i],
				})
			}
		}
		create.Cols = append(create.Cols, &ast.ColumnDef{
			Colname:   def.Name.String(),
			TypeName:  &ast.TypeName{Name: types.TypeStr(def.Tp.Tp)},
			IsNotNull: isNotNull(def),
		})
	}
	return create
}

func (c *cc) convertColumnNameExpr(n *pcast.ColumnNameExpr) *ast.ColumnRef {
	return &ast.ColumnRef{
		Fields: &ast.List{
			Items: []ast.Node{
				&ast.String{Str: n.Name.Name.String()},
			},
		},
	}
}

func (c *cc) convertColumnNames(cols []*pcast.ColumnName) *ast.List {
	list := &ast.List{Items: []ast.Node{}}
	for i := range cols {
		name := cols[i].Name.String()
		list.Items = append(list.Items, &ast.ResTarget{
			Name: &name,
		})
	}
	return list
}

func (c *cc) convertDeleteStmt(n *pcast.DeleteStmt) *ast.DeleteStmt {
	rels := c.convertTableRefsClause(n.TableRefs)
	if len(rels.Items) != 1 {
		panic("expected one range var")
	}
	rel := rels.Items[0]
	rangeVar, ok := rel.(*ast.RangeVar)
	if !ok {
		panic("expected range var")
	}

	return &ast.DeleteStmt{
		Relation:      rangeVar,
		WhereClause:   c.convert(n.Where),
		ReturningList: &ast.List{},
	}
}

func (c *cc) convertDropTableStmt(n *pcast.DropTableStmt) ast.Node {
	// TODO: Remove once views are supported.
	if n.IsView {
		return &ast.TODO{}
	}
	drop := &ast.DropTableStmt{IfExists: n.IfExists}
	for _, name := range n.Tables {
		drop.Tables = append(drop.Tables, parseTableName(name))
	}
	return drop
}

func (c *cc) convertRenameTableStmt(n *pcast.RenameTableStmt) ast.Node {
	return &ast.RenameTableStmt{
		Table:   parseTableName(n.OldTable),
		NewName: &parseTableName(n.NewTable).Name,
	}
}

func (c *cc) convertExistsSubqueryExpr(n *pcast.ExistsSubqueryExpr) *ast.SubLink {
	sublink := &ast.SubLink{}
	if ss, ok := c.convert(n.Sel).(*ast.SelectStmt); ok {
		sublink.Subselect = ss
	}
	return sublink
}

func (c *cc) convertFieldList(n *pcast.FieldList) *ast.List {
	fields := make([]ast.Node, len(n.Fields))
	for i := range n.Fields {
		fields[i] = c.convertSelectField(n.Fields[i])
	}
	return &ast.List{Items: fields}
}

func (c *cc) convertFuncCallExpr(n *pcast.FuncCallExpr) *ast.FuncCall {
	schema := n.Schema.String()
	name := n.FnName.String()

	// TODO: Deprecate the usage of Funcname
	items := []ast.Node{}
	if schema != "" {
		items = append(items, &ast.String{Str: schema})
	}
	items = append(items, &ast.String{Str: name})

	fn := &ast.FuncCall{
		Args: &ast.List{},
		Func: &ast.FuncName{
			Schema: schema,
			Name:   name,
		},
		Funcname: &ast.List{
			Items: items,
		},
		Location: n.Offset,
	}
	for _, arg := range n.Args {
		fn.Args.Items = append(fn.Args.Items, c.convert(arg))
	}
	return fn
}

func (c *cc) convertInsertStmt(n *pcast.InsertStmt) *ast.InsertStmt {
	rels := c.convertTableRefsClause(n.Table)
	if len(rels.Items) != 1 {
		panic("expected one range var")
	}
	rel := rels.Items[0]
	rangeVar, ok := rel.(*ast.RangeVar)
	if !ok {
		panic("expected range var")
	}

	// debug.Dump(n)
	insert := &ast.InsertStmt{
		Relation:      rangeVar,
		Cols:          c.convertColumnNames(n.Columns),
		ReturningList: &ast.List{},
	}
	if ss, ok := c.convert(n.Select).(*ast.SelectStmt); ok {
		ss.ValuesLists = c.convertLists(n.Lists)
		insert.SelectStmt = ss
	} else {
		insert.SelectStmt = &ast.SelectStmt{
			FromClause:  &ast.List{},
			TargetList:  &ast.List{},
			ValuesLists: c.convertLists(n.Lists),
		}
	}
	return insert
}

func (c *cc) convertLists(lists [][]pcast.ExprNode) *ast.List {
	list := &ast.List{Items: []ast.Node{}}
	for _, exprs := range lists {
		inner := &ast.List{Items: []ast.Node{}}
		for _, expr := range exprs {
			inner.Items = append(inner.Items, c.convert(expr))
		}
		list.Items = append(list.Items, inner)
	}
	return list
}

func (c *cc) convertParamMarkerExpr(n *driver.ParamMarkerExpr) *ast.ParamRef {
	// Parameter numbers start at one
	c.paramCount += 1
	return &ast.ParamRef{
		Number:   c.paramCount,
		Location: n.Offset,
	}
}

func (c *cc) convertSelectField(n *pcast.SelectField) *ast.ResTarget {
	var val ast.Node
	if n.WildCard != nil {
		val = c.convertWildCardField(n.WildCard)
	} else {
		val = c.convert(n.Expr)
	}
	var name *string
	if n.AsName.O != "" {
		name = &n.AsName.O
	}
	return &ast.ResTarget{
		// TODO: Populate Indirection field
		Name:     name,
		Val:      val,
		Location: n.Offset,
	}
}

func (c *cc) convertSelectStmt(n *pcast.SelectStmt) *ast.SelectStmt {
	stmt := &ast.SelectStmt{
		TargetList:  c.convertFieldList(n.Fields),
		FromClause:  c.convertTableRefsClause(n.From),
		WhereClause: c.convert(n.Where),
	}
	if n.Limit != nil {
		stmt.LimitCount = c.convert(n.Limit.Count)
		stmt.LimitOffset = c.convert(n.Limit.Offset)
	}
	return stmt
}

func (c *cc) convertSubqueryExpr(n *pcast.SubqueryExpr) ast.Node {
	return c.convert(n.Query)
}

func (c *cc) convertTableRefsClause(n *pcast.TableRefsClause) *ast.List {
	var tables []ast.Node
	visit(n, func(n pcast.Node) {
		name, ok := n.(*pcast.TableName)
		if !ok {
			return
		}
		schema := name.Schema.String()
		rel := name.Name.String()
		tables = append(tables, &ast.RangeVar{
			Schemaname: &schema,
			Relname:    &rel,
		})
	})
	return &ast.List{Items: tables}
}

func (c *cc) convertUpdateStmt(n *pcast.UpdateStmt) *ast.UpdateStmt {
	// Relation
	rels := c.convertTableRefsClause(n.TableRefs)
	if len(rels.Items) != 1 {
		panic("expected one range var")
	}
	rel := rels.Items[0]
	rangeVar, ok := rel.(*ast.RangeVar)
	if !ok {
		panic("expected range var")
	}
	// TargetList
	list := &ast.List{}
	for _, a := range n.List {
		list.Items = append(list.Items, c.convertAssignment(a))
	}
	return &ast.UpdateStmt{
		Relation:      rangeVar,
		TargetList:    list,
		WhereClause:   c.convert(n.Where),
		FromClause:    &ast.List{},
		ReturningList: &ast.List{},
	}
}

func (c *cc) convertValueExpr(n *driver.ValueExpr) *ast.A_Const {
	return &ast.A_Const{
		Val: &ast.String{
			Str: n.Datum.GetString(),
		},
	}
}

func (c *cc) convertWildCardField(n *pcast.WildCardField) *ast.ColumnRef {
	items := []ast.Node{}
	if t := n.Table.String(); t != "" {
		items = append(items, &ast.String{Str: t})
	}
	items = append(items, &ast.A_Star{})

	return &ast.ColumnRef{
		Fields: &ast.List{
			Items: items,
		},
	}
}

func (c *cc) convert(node pcast.Node) ast.Node {
	switch n := node.(type) {

	case *driver.ParamMarkerExpr:
		return c.convertParamMarkerExpr(n)

	case *driver.ValueExpr:
		return c.convertValueExpr(n)

	case *pcast.AlterTableStmt:
		return c.convertAlterTableStmt(n)

	case *pcast.BinaryOperationExpr:
		return c.convertBinaryOperationExpr(n)

	case *pcast.ColumnNameExpr:
		return c.convertColumnNameExpr(n)

	case *pcast.CreateTableStmt:
		return c.convertCreateTableStmt(n)

	case *pcast.DeleteStmt:
		return c.convertDeleteStmt(n)

	case *pcast.DropTableStmt:
		return c.convertDropTableStmt(n)

	case *pcast.RenameTableStmt:
		return c.convertRenameTableStmt(n)

	case *pcast.ExistsSubqueryExpr:
		return c.convertExistsSubqueryExpr(n)

	case *pcast.FuncCallExpr:
		return c.convertFuncCallExpr(n)

	case *pcast.InsertStmt:
		return c.convertInsertStmt(n)

	case *pcast.SelectStmt:
		return c.convertSelectStmt(n)

	case *pcast.SubqueryExpr:
		return c.convertSubqueryExpr(n)

	case *pcast.UpdateStmt:
		return c.convertUpdateStmt(n)

	case nil:
		return nil

	default:
		if debug.Active {
			fmt.Printf("dolphin.convert: Unknown node type %T\n", n)
		}
		return &ast.TODO{}
	}
}
