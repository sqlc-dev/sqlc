package vitess

import (
	"fmt"

	"vitess.io/vitess/go/vt/sqlparser"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/ast/pg"
)

func convertAliasedExpr(n *sqlparser.AliasedExpr) *pg.ResTarget {
	var name string
	if !n.As.IsEmpty() {
		name = n.As.String()
	}
	return &pg.ResTarget{
		Name: &name,
		Val:  convert(n.Expr),
	}
}

func convertColName(n *sqlparser.ColName) *pg.ColumnRef {
	// TODO: Add table name if necessary
	return &pg.ColumnRef{
		Fields: &ast.List{
			Items: []ast.Node{
				&pg.String{
					Str: n.Name.String(),
				},
			},
		},
	}
}

func convertColumnType(n sqlparser.ColumnType) *ast.TypeName {
	return &ast.TypeName{
		Name: n.Type,
	}
}

func convertDDL(n *sqlparser.DDL) ast.Node {
	switch n.Action {
	case sqlparser.AddAutoIncStr:
	case sqlparser.AddColVindexStr:
	case sqlparser.AddSequenceStr:
	case sqlparser.AddVschemaTableStr:
	case sqlparser.AlterStr:
	case sqlparser.CreateStr:
		create := &ast.CreateTableStmt{
			Name:        convertTableName(n.Table),
			IfNotExists: !n.IfExists,
		}
		if n.TableSpec == nil {
			return create
		}
		for _, def := range n.TableSpec.Columns {
			create.Cols = append(create.Cols, &ast.ColumnDef{
				Colname:   def.Name.String(),
				TypeName:  convertColumnType(def.Type),
				IsNotNull: bool(def.Type.NotNull),
			})
		}
		return create
	case sqlparser.CreateVindexStr:
	case sqlparser.DropColVindexStr:
	case sqlparser.DropStr:
	case sqlparser.DropVindexStr:
	case sqlparser.DropVschemaTableStr:
	case sqlparser.FlushStr:
	case sqlparser.RenameStr:
	case sqlparser.TruncateStr:
	default:
		panic("unknown DDL action: " + n.Action)
	}
	return &ast.TODO{}
}

func convertFuncExpr(n *sqlparser.FuncExpr) *ast.FuncCall {
	// TODO: Populate additional field names
	return &ast.FuncCall{
		Func: &ast.FuncName{
			Name: n.Name.String(),
		},
		Funcname: &ast.List{
			Items: []ast.Node{
				&pg.String{Str: n.Name.String()},
			},
		},
		Args: convertSelectExprs(n.Exprs),
	}
}

func convertSelectExprs(n sqlparser.SelectExprs) *ast.List {
	exprs := make([]ast.Node, len(n))
	for i := range n {
		exprs[i] = convert(n[i])
	}
	return &ast.List{Items: exprs}
}

func convertSelectStmt(n *sqlparser.Select) *pg.SelectStmt {
	return &pg.SelectStmt{
		TargetList: convertSelectExprs(n.SelectExprs),
		FromClause: convertTableExprs(n.From),
	}
}

func convertStarExpr(n *sqlparser.StarExpr) *pg.ResTarget {
	return &pg.ResTarget{
		Val: &pg.ColumnRef{
			Fields: &ast.List{
				Items: []ast.Node{
					&pg.A_Star{},
				},
			},
		},
	}
}

func convertTableExprs(n sqlparser.TableExprs) *ast.List {
	var tables []ast.Node
	err := sqlparser.Walk(func(n sqlparser.SQLNode) (bool, error) {
		table, ok := n.(sqlparser.TableName)
		if !ok {
			return true, nil
		}
		schema := table.Qualifier.String()
		rel := table.Name.String()
		tables = append(tables, &pg.RangeVar{
			Schemaname: &schema,
			Relname:    &rel,
		})
		return false, nil
	}, n)
	if err != nil {
		panic(err)
	}
	return &ast.List{Items: tables}
}

func convertTableName(n sqlparser.TableName) *ast.TableName {
	return &ast.TableName{
		Schema: n.Qualifier.String(),
		Name:   n.Name.String(),
	}
}

func convert(node sqlparser.SQLNode) ast.Node {
	switch n := node.(type) {

	case *sqlparser.AliasedExpr:
		return convertAliasedExpr(n)

	case *sqlparser.ColName:
		return convertColName(n)

	case *sqlparser.DDL:
		return convertDDL(n)

	case *sqlparser.Select:
		return convertSelectStmt(n)

	case *sqlparser.StarExpr:
		return convertStarExpr(n)

	case *sqlparser.Insert:
		return &ast.TODO{}

	case *sqlparser.Update:
		return &ast.TODO{}

	case *sqlparser.Delete:
		return &ast.TODO{}

	default:
		fmt.Printf("%T\n", n)
		return &ast.TODO{}
	}
}
