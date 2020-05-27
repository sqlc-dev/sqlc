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

func convertSelectStmt(n *pcast.SelectStmt) ast.Node {
	var tables []ast.Node
	visit(n.From, func(n pcast.Node) {
		name, ok := n.(*pcast.TableName)
		if !ok {
			return
		}
		tables = append(tables, parseTableName(name))
	})
	var cols []ast.Node
	visit(n.Fields, func(n pcast.Node) {
		col, ok := n.(*pcast.ColumnName)
		if !ok {
			return
		}
		cols = append(cols, &ast.ResTarget{
			Val: &ast.ColumnRef{
				Name: col.Name.String(),
			},
		})
	})
	return &pg.SelectStmt{
		FromClause: &ast.List{Items: tables},
		TargetList: &ast.List{Items: cols},
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
