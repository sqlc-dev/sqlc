package dolphin

import (
	pcast "github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/types"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/ast/pg"
)

func convertAlterTableStmt(n *pcast.AlterTableStmt) ast.Node {
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

func convertCreateTableStmt(n *pcast.CreateTableStmt) ast.Node {
	create := &ast.CreateTableStmt{
		Name:        parseTableName(n.Table),
		IfNotExists: n.IfNotExists,
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

func convertDropTableStmt(n *pcast.DropTableStmt) ast.Node {
	drop := &ast.DropTableStmt{IfExists: n.IfExists}
	for _, name := range n.Tables {
		drop.Tables = append(drop.Tables, parseTableName(name))
	}
	return drop
}

func convertFieldList(n *pcast.FieldList) *ast.List {
	fields := make([]ast.Node, len(n.Fields))
	for i := range n.Fields {
		fields[i] = convertSelectField(n.Fields[i])
	}
	return &ast.List{Items: fields}
}

func convertSelectField(n *pcast.SelectField) *pg.ResTarget {
	var val ast.Node
	if n.WildCard != nil {
		val = convertWildCardField(n.WildCard)
	} else {
		val = convert(n.Expr)
	}
	var name *string
	if n.AsName.O != "" {
		name = &n.AsName.O
	}
	return &pg.ResTarget{
		// TODO: Populate Indirection field
		Name:     name,
		Val:      val,
		Location: n.Offset,
	}
}

func convertSelectStmt(n *pcast.SelectStmt) *pg.SelectStmt {
	return &pg.SelectStmt{
		TargetList: convertFieldList(n.Fields),
		FromClause: convertTableRefsClause(n.From),
	}
}

func convertTableRefsClause(n *pcast.TableRefsClause) *ast.List {
	var tables []ast.Node
	visit(n, func(n pcast.Node) {
		name, ok := n.(*pcast.TableName)
		if !ok {
			return
		}
		schema := name.Schema.String()
		rel := name.Name.String()
		tables = append(tables, &pg.RangeVar{
			Schemaname: &schema,
			Relname:    &rel,
		})
	})
	return &ast.List{Items: tables}
}

func convertWildCardField(n *pcast.WildCardField) *pg.ColumnRef {
	return &pg.ColumnRef{
		Fields: &ast.List{
			Items: []ast.Node{
				&pg.A_Star{},
			},
		},
	}
}

func convert(node pcast.Node) ast.Node {
	switch n := node.(type) {

	case *pcast.AlterTableStmt:
		return convertAlterTableStmt(n)

	case *pcast.CreateTableStmt:
		return convertCreateTableStmt(n)

	case *pcast.DropTableStmt:
		return convertDropTableStmt(n)

	case *pcast.SelectStmt:
		return convertSelectStmt(n)

	default:
		return &ast.TODO{}
	}
}
