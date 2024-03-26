package clickhouse

import (
	"fmt"
	"io"
	"reflect"

	chparser "github.com/AfterShip/clickhouse-sql-parser/parser"
	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func convertAlterTableCmd(node chparser.AlterTableExpr) (ast.Node, error) {
	switch n := node.(type) {
	case *chparser.AlterTableAttachPartition:
	case *chparser.AlterTableDetachPartition:
	case *chparser.AlterTableDropPartition:
	case *chparser.AlterTableFreezePartition:
	case *chparser.AlterTableAddColumn:
	case *chparser.AlterTableAddIndex:
	case *chparser.AlterTableDropColumn:
	case *chparser.AlterTableDropIndex:
	case *chparser.AlterTableRemoveTTL:
	case *chparser.AlterTableClearColumn:
	case *chparser.AlterTableClearIndex:
	case *chparser.AlterTableRenameColumn:
	case *chparser.AlterTableModifyTTL:
	case *chparser.AlterTableModifyColumn:
	case *chparser.AlterTableReplacePartition:
	default:
		_ = n
	}
	panic("not implemented")
}

func convertTableName(n *chparser.TableIdentifier) *ast.TableName {
	if n.Database != nil {
		return &ast.TableName{
			Catalog: n.Database.Name,
			Schema:  "",
			Name:    n.Table.Name,
		}
	}
	return &ast.TableName{
		Catalog: "",
		Schema:  "",
		Name:    n.Table.Name,
	}

}

type converter struct {
	output []ast.Statement

	tempColumnDef []*ast.ColumnDef
	tempTypeName  *ast.TypeName
	*chparser.DefaultASTVisitor
}

func (c *converter) VisitAlterTable(expr *chparser.AlterTable) error {
	cmdItems := make([]ast.Node, 0)
	for _, v := range expr.AlterExprs {
		n, err := convertAlterTableCmd(v)
		if err != nil {
			return err
		}
		cmdItems = append(cmdItems, n)
	}

	stmt := &ast.AlterTableStmt{
		Relation: &ast.RangeVar{},
		Table:    convertTableName(expr.TableIdentifier),
		Cmds: &ast.List{
			Items: cmdItems,
		},
		MissingOk: false,
	}
	c.output = append(c.output, ast.Statement{Raw: &ast.RawStmt{
		Stmt:         stmt,
		StmtLocation: int(expr.AlterPos),
		StmtLen:      int(expr.StatementEnd) - int(expr.AlterPos),
	}})
	return nil
}

func (c *converter) VisitCreateTable(expr *chparser.CreateTable) error {
	c.tempColumnDef = c.tempColumnDef[:]
	stmt := &ast.CreateTableStmt{
		IfNotExists: expr.IfNotExists,
		Name:        convertTableName(expr.Name),
		Cols:        c.tempColumnDef,
	}
	c.tempColumnDef = c.tempColumnDef[:]
	c.output = append(c.output, ast.Statement{Raw: &ast.RawStmt{
		Stmt:         stmt,
		StmtLocation: int(expr.CreatePos),
		StmtLen:      int(expr.StatementEnd) - int(expr.CreatePos),
	}})
	return nil
}

func (c *converter) VisitColumnTypeExpr(expr *chparser.ColumnTypeExpr) error {
	c.tempTypeName = &ast.TypeName{
		Schema: expr.Name.Name,
	}
	return nil
}
func (c *converter) VisitScalarTypeExpr(expr *chparser.ScalarTypeExpr) error {
	c.tempTypeName = &ast.TypeName{
		Name: expr.Name.Name,
	}
	return nil
}

// TODO: remove befor making pr is ready
func debugFallbackVisitor(e chparser.Expr) error {
	tn := reflect.TypeOf(e).Elem().Name()
	fmt.Println("visit: ", e.String(0), tn)
	return nil
}

func (c *converter) VisitColumn(expr *chparser.Column) error {
	c.tempTypeName = nil
	err := expr.Type.Accept(c)
	if err != nil {
		return err
	}
	c.tempColumnDef = append(c.tempColumnDef, &ast.ColumnDef{
		Colname:   expr.Name.Name,
		TypeName:  &ast.TypeName{Name: c.tempTypeName.Name},
		IsNotNull: expr.NotNull != nil,
	})
	return nil
}

func (c *converter) Convert(exprs []chparser.Expr) ([]ast.Statement, error) {
	for _, e := range exprs {
		e.Accept(c)
	}
	return c.output, nil
}

func (c *Parser) Parse(r io.Reader) ([]ast.Statement, error) {
	sql, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	chParser := chparser.NewParser(string(sql))
	statements, err := chParser.ParseStatements()
	if err != nil {
		return []ast.Statement{}, err
	}
	conv := &converter{
		DefaultASTVisitor: &chparser.DefaultASTVisitor{
			Visit: debugFallbackVisitor,
		},
	}
	return conv.Convert(statements)
}
func (c *Parser) CommentSyntax() source.CommentSyntax {
	return source.CommentSyntax{}

}
func (c *Parser) IsReservedKeyword(string) bool {
	return false
}
