package cockroach

import (
	"fmt"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type node interface {
	Format(ctx *tree.FmtCtx)
}

func convertSlice[T node](nodes []T) *ast.List {
	out := &ast.List{}
	for _, n := range nodes {
		out.Items = append(out.Items, convert(n))
	}
	return out
}

func convertAliasTableExpr(n *tree.AliasedTableExpr) ast.Node {
	if n == nil {
		return nil
	}
	return convert(n.Expr)
}

func convertCreateTable(n *tree.CreateTable) *ast.CreateTableStmt {
	if n == nil {
		return nil
	}
	create := &ast.CreateTableStmt{
		Name:        convertTableName(&n.Table),
		IfNotExists: n.IfNotExists,
	}
	for _, def := range n.Defs {
		switch d := def.(type) {
		case *tree.ColumnTableDef:
			create.Cols = append(create.Cols, &ast.ColumnDef{
				Colname: d.Name.String(),
				TypeName: &ast.TypeName{
					Name: "text",
				},
				IsNotNull: d.Nullable.Nullability == tree.NotNull,
				// IsArray:   isArray(item.ColumnDef.TypeName),
			})
		default:
			fmt.Printf("%#T\n", d)
			continue
		}
	}
	return create
}

func convertDelete(n *tree.Delete) *ast.DeleteStmt {
	if n == nil {
		return nil
	}
	return &ast.DeleteStmt{
		// Relation:      convertRangeVar(n.Relation),
		// UsingClause:   convertSlice(n.UsingClause),
		// WhereClause:   convertNode(n.WhereClause),
		// ReturningList: convertSlice(n.ReturningList),
		// WithClause:    convertWithClause(n.WithClause),
	}
}

func convertInsert(n *tree.Insert) *ast.InsertStmt {
	if n == nil {
		return nil
	}
	return &ast.InsertStmt{
		// Relation:         convertRangeVar(n.Relation),
		// Cols:             convertSlice(n.Cols),
		// SelectStmt:       convertNode(n.SelectStmt),
		// OnConflictClause: convertOnConflictClause(n.OnConflictClause),
		// ReturningList:    convertSlice(n.ReturningList),
		// WithClause:       convertWithClause(n.WithClause),
		// Override:         ast.OverridingKind(n.Override),
	}
}

func convertSelect(n *tree.Select) *ast.SelectStmt {
	if n == nil {
		return nil
	}
	var stmt *ast.SelectStmt
	switch s := n.Select.(type) {
	case *tree.SelectClause:
		stmt = &ast.SelectStmt{
			FromClause: convertSlice(s.From.Tables),
			TargetList: convertSelectExprs(s.Exprs),
		}
	default:
		fmt.Printf("%#T\n", s)
	}
	return stmt
}

func convertSelectExpr(n *tree.SelectExpr) *ast.TODO {
	return &ast.TODO{}
}

func convertSelectExprs(nodes tree.SelectExprs) *ast.List {
	out := &ast.List{}
	for _, n := range nodes {
		out.Items = append(out.Items, convert(&n))
	}
	return out
}

func convertTableName(n *tree.TableName) *ast.TableName {
	name := n.ToUnresolvedObjectName()
	switch name.NumParts {
	case 1:
		return &ast.TableName{
			Name: name.Parts[0],
		}
	case 2:
		return &ast.TableName{
			Schema: name.Parts[0],
			Name:   name.Parts[1],
		}
	default:
		return &ast.TableName{
			Catalog: name.Parts[0],
			Schema:  name.Parts[1],
			Name:    name.Parts[2],
		}
	}
}

func convert(nn node) ast.Node {
	if nn == nil {
		return &ast.TODO{}
	}
	switch n := nn.(type) {

	case *tree.AliasedTableExpr:
		return convertAliasTableExpr(n)

	case *tree.CreateTable:
		return convertCreateTable(n)

	case *tree.Delete:
		return convertDelete(n)

	case *tree.Insert:
		return convertInsert(n)

	case *tree.Select:
		return convertSelect(n)

	case *tree.SelectExpr:
		return convertSelectExpr(n)

	case *tree.TableName:
		return convertTableName(n)

	default:
		fmt.Printf("%#T\n", n)
		return &ast.TODO{}
	}
}
