package cockroach

import (
	"fmt"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

func convertCreateTable(n *tree.CreateTable) *ast.CreateTableStmt {
	if n == nil {
		return nil
	}
	create := &ast.CreateTableStmt{
		Name:        convertTableName(n.Table),
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
			FromClause: convertFrom(s.From),
			TargetList: convertSelectExprs(s.Exprs),
		}
	default:
		fmt.Printf("%#T\n", s)
	}
	return stmt
}

func convertTableName(n tree.TableName) *ast.TableName {
	name := n.ToUnresolvedObjectName()
	return &ast.TableName{
		Catalog: name.Parts[0],
		Schema:  name.Parts[1],
		Name:    name.Parts[2],
	}
}

func convert(stmt tree.Statement) ast.Node {
	if stmt == nil {
		return &ast.TODO{}
	}
	switch n := stmt.(type) {

	case *tree.CreateTable:
		return convertCreateTable(n)

	case *tree.Delete:
		return convertDelete(n)

	case *tree.Insert:
		return convertInsert(n)

	case *tree.Select:
		return convertSelect(n)

	default:
		fmt.Printf("%#T\n", n)
		return &ast.TODO{}
	}
}
