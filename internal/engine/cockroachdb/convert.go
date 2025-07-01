package cockroachdb

import (
	"fmt"
	"strings"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
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

func convertSlice[T node](nodes []T, originalSQL string) *ast.List {
	out := &ast.List{}
	for _, n := range nodes {
		out.Items = append(out.Items, convert(n, originalSQL))
	}
	return out
}

// Not sure why I have to write these functions instead of using convertSlice
func convertExprs(nodes tree.Exprs, originalSQL string) *ast.List {
	out := &ast.List{}
	for _, n := range nodes {
		out.Items = append(out.Items, convert(n, originalSQL))
	}
	return out
}

func convertReturningExprs(nodes tree.ReturningExprs, originalSQL string) *ast.List {
	out := &ast.List{}
	for _, n := range nodes {
		out.Items = append(out.Items, convert(&n, originalSQL))
	}
	return out
}

func convertSelectExprs(nodes tree.SelectExprs, originalSQL string) *ast.List {
	out := &ast.List{}
	for _, n := range nodes {
		// out.Items = append(out.Items, convert(&n))
		out.Items = append(out.Items, convertSelectExpr(&n, originalSQL))
	}
	return out
}

func convertNameList(nodes tree.NameList, originalSQL string) *ast.List {
	out := &ast.List{}
	for _, n := range nodes {
		out.Items = append(out.Items, convert(&n, originalSQL))
	}
	return out
}

// End

func convertAliasTableExpr(n *tree.AliasedTableExpr, originalSQL string) ast.Node {
	if n == nil {
		return nil
	}
	return convert(n.Expr, originalSQL)
}

func convertComparisonExpr(n *tree.ComparisonExpr, originalSQL string) *ast.A_Expr {
	if n == nil {
		return nil
	}
	return &ast.A_Expr{
		Lexpr: convert(n.Left, originalSQL),
		Rexpr: convert(n.Right, originalSQL),
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

func convertDelete(n *tree.Delete, originalSQL string) *ast.DeleteStmt {
	if n == nil {
		return nil
	}
	return &ast.DeleteStmt{
		// Relation: []mustConvert[*ast.RangeVar](n.Table),
		// Relations: &ast.List{[]ast.Node{mustConvert[*ast.RangeVar](n.Table)}},
		// UsingClause:   convertSlice(n.UsingClause),
		WhereClause: convert(n.Where, originalSQL),
		// ReturningList: convertSlice(n.ReturningList),
		WithClause: convertWith(n.With),
	}
}

func convertInsert(n *tree.Insert, originalSQL string) *ast.InsertStmt {
	if n == nil {
		return nil
	}
	return &ast.InsertStmt{
		Relation:   mustConvert[*ast.RangeVar](n.Table, originalSQL),
		Cols:       convertNameList(n.Columns, originalSQL),
		SelectStmt: convertSelect(n.Rows, originalSQL),
		// OnConflictClause: convertOnConflictClause(n.OnConflictClause),
		ReturningList: mustConvert[*ast.List](n.Returning, originalSQL),
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

func convertSelect(n *tree.Select, originalSQL string) *ast.SelectStmt {
	if n == nil {
		return nil
	}
	stmt := mustConvert[*ast.SelectStmt](n.Select, originalSQL)
	return stmt
}

func convertSelectClause(n *tree.SelectClause, originalSQL string) *ast.SelectStmt {
	if n == nil {
		return nil
	}
	return &ast.SelectStmt{
		FromClause:  convertSlice(n.From.Tables, originalSQL),
		TargetList:  convertSelectExprs(n.Exprs, originalSQL),
		WhereClause: convert(n.Where, originalSQL),
	}
}

func convertSelectExpr(n *tree.SelectExpr, originalSQL string) *ast.ResTarget {
	location := 0
	if _, ok := n.Expr.(tree.UnqualifiedStar); ok && originalSQL != "" {
		location = strings.Index(originalSQL, "*")
		if location < -1 || location >= len(originalSQL) {
			location = 0 // fallback
		}
	}
	// fmt.Println("s : ", originalSQL)
	// fmt.Println("l : ", location)
	return &ast.ResTarget{
		Val:      convert(n.Expr, originalSQL),
		Location: location + currParserIndexPos,
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

func convertUnqualifiedStar(n tree.UnqualifiedStar) *ast.ColumnRef {
	// return &ast.A_Star{}
	return &ast.ColumnRef{
		Fields: &ast.List{
			Items: []ast.Node{&ast.A_Star{}},
		},
	}
}

func convertAllColumnsSelector(n *tree.AllColumnsSelector) *ast.ColumnRef {
	tableName := n.TableName.String()
	// fmt.Println("hei : ", tableName)
	return &ast.ColumnRef{
		Fields: &ast.List{
			Items: []ast.Node{
				&ast.String{Str: tableName},
				&ast.A_Star{},
			},
		},
	}
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

func convertValuesClause(n *tree.ValuesClause, originalSQL string) *ast.SelectStmt {
	if n == nil {
		return nil
	}
	values := &ast.List{}
	for _, v := range n.Rows {
		values.Items = append(values.Items, convertExprs(v, originalSQL))
	}
	return &ast.SelectStmt{
		TargetList:  &ast.List{},
		ValuesLists: values,
	}
}

func convertWhere(n *tree.Where, originalSQL string) ast.Node {
	if n == nil {
		return nil
	}
	return convert(n.Expr, originalSQL)
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

func convert(nn node, originalSQL string) ast.Node {
	if nn == nil {
		return &ast.TODO{}
	}
	switch n := nn.(type) {

	case *tree.AliasedTableExpr:
		return convertAliasTableExpr(n, originalSQL)

	case *tree.CreateTable:
		return convertCreateTable(n)

	case *tree.ComparisonExpr:
		return convertComparisonExpr(n, originalSQL)

	case *tree.Delete:
		return convertDelete(n, originalSQL)

	case *tree.Insert:
		return convertInsert(n, originalSQL)

	case *tree.Name:
		return convertName(n)

	case *tree.Placeholder:
		return convertPlaceholder(n)

	case *tree.ReturningExprs:
		return convertReturningExprs(*n, originalSQL)

	case *tree.Select:
		return convertSelect(n, originalSQL)

	case *tree.SelectClause:
		return convertSelectClause(n, originalSQL)

	case *tree.SelectExpr:
		return convertSelectExpr(n, originalSQL)

	case tree.UnqualifiedStar:
		return convertUnqualifiedStar(n)

	case *tree.TableName:
		return convertTableName(n)

	case *tree.UnresolvedName:
		return convertUnresolvedName(n)

	case *tree.ValuesClause:
		return convertValuesClause(n, originalSQL)

	case *tree.Where:
		return convertWhere(n, originalSQL)

	case *tree.AllColumnsSelector:
		return convertAllColumnsSelector(n)

	default:
		fmt.Printf("convert unknown type %#T\n", n)
		return &ast.TODO{}
	}
}

func mustConvert[T ast.Node](in node, originalSQL string) T {
	var out T
	n := convert(in, originalSQL)
	out, ok := n.(T)
	if !ok {
		panic(fmt.Sprintf("could not convert %T to %T", in, out))
	}
	return out
}
