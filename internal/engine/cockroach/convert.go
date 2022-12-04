package cockroach

import (
	"fmt"
	"strings"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type node interface {
	Format(ctx *tree.FmtCtx)
}

func makeString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func createTableName(n *tree.TableName) *ast.TableName {
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

func convertSlice[T node](nodes []T) *ast.List {
	out := &ast.List{}
	for _, n := range nodes {
		out.Items = append(out.Items, convert(n))
	}
	return out
}

// Not sure why I have to write these functions instead of using convertSlice
func convertExprs(nodes tree.Exprs) *ast.List {
	out := &ast.List{}
	for _, n := range nodes {
		out.Items = append(out.Items, convert(n))
	}
	return out
}

func convertReturningExprs(nodes tree.ReturningExprs) *ast.List {
	out := &ast.List{}
	for _, n := range nodes {
		out.Items = append(out.Items, convert(&n))
	}
	return out
}

func convertSelectExprs(nodes tree.SelectExprs) *ast.List {
	out := &ast.List{}
	for _, n := range nodes {
		out.Items = append(out.Items, convert(&n))
	}
	return out
}

func convertNameList(nodes tree.NameList) *ast.List {
	out := &ast.List{}
	for _, n := range nodes {
		out.Items = append(out.Items, convert(&n))
	}
	return out
}

// End

func convertAliasTableExpr(n *tree.AliasedTableExpr) ast.Node {
	if n == nil {
		return nil
	}
	return convert(n.Expr)
}

func convertComparisonExpr(n *tree.ComparisonExpr) *ast.A_Expr {
	if n == nil {
		return nil
	}
	return &ast.A_Expr{
		Lexpr: convert(n.Left),
		Rexpr: convert(n.Right),
	}
}

func convertCreateTable(n *tree.CreateTable) *ast.CreateTableStmt {
	if n == nil {
		return nil
	}
	create := &ast.CreateTableStmt{
		Name:        createTableName(&n.Table),
		IfNotExists: n.IfNotExists,
	}
	for _, def := range n.Defs {
		switch d := def.(type) {
		case *tree.ColumnTableDef:
			isNotNull := d.Nullable.Nullability == tree.NotNull
			isPrimary := d.PrimaryKey.IsPrimaryKey
			create.Cols = append(create.Cols, &ast.ColumnDef{
				Colname: d.Name.String(),
				TypeName: &ast.TypeName{
					Name: strings.ToLower(d.Type.SQLString()),
				},
				IsNotNull: isNotNull || isPrimary,
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
		Relation: mustConvert[*ast.RangeVar](n.Table),
		// UsingClause:   convertSlice(n.UsingClause),
		WhereClause: convert(n.Where),
		// ReturningList: convertSlice(n.ReturningList),
		WithClause: convertWith(n.With),
	}
}

func convertInsert(n *tree.Insert) *ast.InsertStmt {
	if n == nil {
		return nil
	}
	return &ast.InsertStmt{
		Relation:   mustConvert[*ast.RangeVar](n.Table),
		Cols:       convertNameList(n.Columns),
		SelectStmt: convertSelect(n.Rows),
		// OnConflictClause: convertOnConflictClause(n.OnConflictClause),
		ReturningList: mustConvert[*ast.List](n.Returning),
		WithClause:    convertWith(n.With),
		// Override:         ast.OverridingKind(n.Override),
	}
}

func convertName(n *tree.Name) *ast.ResTarget {
	name := string(*n)
	return &ast.ResTarget{
		Name: &name,
	}
}

func convertPlaceholder(n *tree.Placeholder) *ast.ParamRef {
	if n == nil {
		return nil
	}
	var dollar bool
	// if n.Number != 0 {
	// 	dollar = true
	// }
	return &ast.ParamRef{
		Dollar: dollar,
		Number: int(n.Idx) + 1,
	}
}

func convertSelect(n *tree.Select) *ast.SelectStmt {
	if n == nil {
		return nil
	}
	stmt := mustConvert[*ast.SelectStmt](n.Select)
	return stmt
}

func convertSelectClause(n *tree.SelectClause) *ast.SelectStmt {
	if n == nil {
		return nil
	}
	return &ast.SelectStmt{
		FromClause:  convertSlice(n.From.Tables),
		TargetList:  convertSelectExprs(n.Exprs),
		WhereClause: convert(n.Where),
	}
}

func convertSelectExpr(n *tree.SelectExpr) *ast.ResTarget {
	return &ast.ResTarget{
		Val: convert(n.Expr),
	}
}

func convertTableName(n *tree.TableName) *ast.RangeVar {
	name := n.ToUnresolvedObjectName()
	switch name.NumParts {
	case 1:
		return &ast.RangeVar{
			Relname: makeString(name.Parts[0]),
		}
	case 2:
		return &ast.RangeVar{
			Schemaname: makeString(name.Parts[0]),
			Relname:    makeString(name.Parts[1]),
		}
	default:
		return &ast.RangeVar{
			Catalogname: makeString(name.Parts[0]),
			Schemaname:  makeString(name.Parts[1]),
			Relname:     makeString(name.Parts[2]),
		}
	}
}

func convertUnqualifiedStar(n *tree.UnqualifiedStar) *ast.A_Star {
	if n == nil {
		return nil
	}
	return &ast.A_Star{}
}

func convertUnresolvedName(n *tree.UnresolvedName) *ast.ColumnRef {
	items := &ast.List{}
	for _, v := range n.Parts {
		if v != "" {
			items.Items = append(items.Items, &ast.String{
				Str: v,
			})
		}
	}
	return &ast.ColumnRef{
		Fields: items,
	}
}

func convertValuesClause(n *tree.ValuesClause) *ast.SelectStmt {
	if n == nil {
		return nil
	}
	values := &ast.List{}
	for _, v := range n.Rows {
		values.Items = append(values.Items, convertExprs(v))
	}
	return &ast.SelectStmt{
		TargetList:  &ast.List{},
		ValuesLists: values,
	}
}

func convertWhere(n *tree.Where) ast.Node {
	if n == nil {
		return nil
	}
	return convert(n.Expr)
}

func convertWith(n *tree.With) *ast.WithClause {
	if n == nil {
		return nil
	}
	return &ast.WithClause{
		Recursive: n.Recursive,
		Ctes:      &ast.List{},
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

	case *tree.ComparisonExpr:
		return convertComparisonExpr(n)

	case *tree.Delete:
		return convertDelete(n)

	case *tree.Insert:
		return convertInsert(n)

	case *tree.Name:
		return convertName(n)

	case *tree.Placeholder:
		return convertPlaceholder(n)

	case *tree.ReturningExprs:
		return convertReturningExprs(*n)

	case *tree.Select:
		return convertSelect(n)

	case *tree.SelectClause:
		return convertSelectClause(n)

	case *tree.SelectExpr:
		return convertSelectExpr(n)

	case *tree.TableName:
		return convertTableName(n)

	case *tree.UnresolvedName:
		return convertUnresolvedName(n)

	case tree.UnqualifiedStar:
		return convertUnqualifiedStar(&n)

	case *tree.ValuesClause:
		return convertValuesClause(n)

	case *tree.Where:
		return convertWhere(n)

	default:
		fmt.Printf("convert unknown type %#T\n", n)
		return &ast.TODO{}
	}
}

func mustConvert[T ast.Node](in node) T {
	var out T
	n := convert(in)
	out, ok := n.(T)
	if !ok {
		panic(fmt.Sprintf("could not convert %#T to %#T", in, out))
	}
	return out
}
