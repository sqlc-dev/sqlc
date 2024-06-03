package clickhouse

import (
	"fmt"
	"io"
	"reflect"
	"strconv"

	chparser "github.com/AfterShip/clickhouse-sql-parser/parser"
	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

type Parser struct {
}

func NewParser() *Parser {
	fmt.Printf("get clickhouse parser")
	return &Parser{}
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
	*chparser.DefaultASTVisitor
}

func (c *converter) convertToColumnDef(exp chparser.Expr) *ast.ColumnDef {
	switch chexp := exp.(type) {
	case *chparser.Column:
		return &ast.ColumnDef{
			Colname:   chexp.Name.Name,
			TypeName:  c.convertToTypeName(chexp.Type),
			IsNotNull: chexp.NotNull != nil,
		}
	default:
		panic("can not convert to column")
	}
}

func (c *converter) convertToTypeName(expr chparser.Expr) *ast.TypeName {
	switch chexpr := expr.(type) {
	case *chparser.ColumnTypeExpr:
		return &ast.TypeName{
			Name: chexpr.Name.Name,
		}
	case *chparser.ScalarTypeExpr:
		return &ast.TypeName{
			Name: chexpr.Name.Name,
		}
	default:
		panic("unknown type")
	}
}

func (c *converter) convertCreateTable(expr *chparser.CreateTable) ast.Node {
	stmt := &ast.CreateTableStmt{
		IfNotExists: expr.IfNotExists,
		Name:        convertTableName(expr.Name),
		Cols:        make([]*ast.ColumnDef, 0),
	}
	for _, v := range expr.TableSchema.Columns {
		stmt.Cols = append(stmt.Cols, c.convertToColumnDef(v))
	}
	return stmt
}

// TODO: remove befor making pr is ready
func debugFallbackVisitor(e chparser.Expr) error {
	// tn := reflect.TypeOf(e).Elem().Name()
	// fmt.Println("visit: ", e.String(0), tn)
	return nil
}

func (c *converter) convertToList(expr chparser.Expr) *ast.List {
	res := &ast.List{}

	switch chexpr := expr.(type) {
	case *chparser.ArrayParamList:
		return c.convertToList(chexpr.Items)
	case *chparser.TTLExprList:
	case *chparser.ColumnArgList:
	case *chparser.ColumnExprList:
		for _, v := range chexpr.Items {
			res.Items = append(res.Items, c.convertExpr(v))
		}
	case *chparser.TableArgListExpr:
		for _, v := range chexpr.Args {
			res.Items = append(res.Items, c.convertExpr(v))
		}
	case *chparser.OrderByListExpr:
		for _, v := range chexpr.Items {
			res.Items = append(res.Items, c.convertExpr(v))
		}
	case *chparser.SettingsExprList:
		for _, v := range chexpr.Items {
			res.Items = append(res.Items, c.convertExpr(v))
		}
	// case *chparser.ParamExprList:
	case *chparser.EnumValueExprList:
		for _, v := range chexpr.Enums {
			res.Items = append(res.Items, c.convertEnumValueExpr(v))
		}
	case *chparser.GroupByExpr:
		return c.convertToList(chexpr.Expr)

	case *chparser.FromExpr:
		return c.convertToList(chexpr.Expr)
	case *chparser.TableExpr:
		return &ast.List{
			Items: []ast.Node{

				c.convertExpr(chexpr.Expr),
			},
		}
	default:

		name := reflect.TypeOf(expr).Elem().Name()
		fmt.Println("missed type", name)
		panic("can not convert to list")
	}

	return nil
}

func (c *converter) convertEnumValueExpr(exp chparser.EnumValueExpr) ast.Node {
	return &ast.String{
		Str: exp.Name.Literal,
	}

}

func (c *converter) convertGroupBy(expr *chparser.GroupByExpr) ast.Node {
	return &ast.GroupingSet{
		Kind:     0,
		Content:  c.convertToList(expr.Expr),
		Location: int(expr.Pos()),
	}
}

func (c *converter) convertLimit(expr *chparser.LimitExpr) {

}

func (c *converter) convertSelectQuery(expr *chparser.SelectQuery) ast.Node {
	stmt := ast.SelectStmt{
		TargetList: c.convertToList(expr.SelectColumns),
		FromClause: c.convertToList(expr.From),
	}
	if len(expr.SelectColumns.Items) == 1 &&
		expr.SelectColumns.Items[0].String(0) == "*" {
		stmt.All = true

	}
	if expr.SelectColumns.HasDistinct {
		stmt.DistinctClause = &ast.List{
			Items: []ast.Node{
				c.convertExpr(expr.SelectColumns.Items[0]),
			},
		}
	}
	if expr.Where != nil {
		stmt.WhereClause = c.convertExpr(expr.Where)
	}
	if expr.GroupBy != nil {
		stmt.GroupClause = c.convertToList(expr.GroupBy)
	}
	if expr.Having != nil {
		stmt.HavingClause = c.convertExpr(expr.Having)
	}
	if expr.OrderBy != nil {
		stmt.SortClause = c.convertToList(expr.OrderBy)
	}

	if expr.Limit != nil {
		if expr.Limit.Limit != nil {
			stmt.LimitCount = c.convertExpr(expr.Limit.Limit)
		}
		if expr.Limit.Offset != nil {
			stmt.LimitOffset = c.convertExpr(expr.Limit.Offset)
		}
	}
	return &stmt
}

func (c *converter) convertAlterTable(e *chparser.AlterTable) ast.Node {
	for _, v := range e.AlterExprs {
		switch chexp := v.(type) {
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
			_ = chexp

		}
	}
	return nil
}

func (c *converter) convertWhere(expr *chparser.WhereExpr) ast.Node {
	return c.convertExpr(expr.Expr)
}

func (c *converter) convertBinaryExpr(expr *chparser.BinaryExpr) ast.Node {
	return &ast.BoolExpr{
		Args: &ast.List{
			Items: []ast.Node{
				c.convertExpr(expr.LeftExpr),
				c.convertExpr(expr.RightExpr),
			},
		},
	}
}

func (c *converter) convertExpr(expr chparser.Expr) ast.Node {
	switch chexp := expr.(type) {
	case *chparser.OperationExpr:
	case *chparser.TernaryExpr:
	case *chparser.BinaryExpr:
		return c.convertBinaryExpr(chexp)
	case *chparser.AlterTable:
		return c.convertAlterTable(chexp)
	case *chparser.RemovePropertyType:
	case *chparser.TableIndex:
	case *chparser.Ident:
		return &ast.String{Str: chexp.Name}
	case *chparser.UUID:
	case *chparser.CreateDatabase:
	case *chparser.CreateTable:
		return c.convertCreateTable(chexp)
	case *chparser.CreateMaterializedView:
	case *chparser.CreateView:
	case *chparser.CreateFunction:
	case *chparser.RoleName:
	case *chparser.SettingPair:
	case *chparser.RoleSetting:
	case *chparser.CreateRole:
	case *chparser.AlterRole:
	case *chparser.RoleRenamePair:
	case *chparser.DestinationExpr:
	case *chparser.ConstraintExpr:
	case *chparser.NullLiteral:
	case *chparser.NotNullLiteral:
	case *chparser.NestedIdentifier:
	case *chparser.ColumnIdentifier:
	case *chparser.TableIdentifier:
		return convertTableName(chexp)
	case *chparser.TableSchemaExpr:
	case *chparser.TableFunctionExpr:
	case *chparser.OnClusterExpr:
	case *chparser.DefaultExpr:
	case *chparser.PartitionExpr:
	case *chparser.PartitionByExpr:
	case *chparser.PrimaryKeyExpr:
	case *chparser.SampleByExpr:
	case *chparser.TTLExpr:
	case *chparser.OrderByExpr:
	case *chparser.SettingsExpr:
	case *chparser.ObjectParams:
	case *chparser.FunctionExpr:
	case *chparser.WindowFunctionExpr:
	case *chparser.Column:
	case *chparser.ScalarTypeExpr:
	case *chparser.PropertyTypeExpr:
	case *chparser.TypeWithParamsExpr:
	case *chparser.ComplexTypeExpr:
	case *chparser.NestedTypeExpr:
	case *chparser.CompressionCodec:
	case *chparser.NumberLiteral:
		numberLiteral, err := strconv.ParseInt(chexp.Literal, chexp.Base, 64)
		if err != nil {
			panic("wrong numberformat")
		}
		return &ast.A_Const{Val: &ast.Integer{Ival: numberLiteral}}
	case *chparser.StringLiteral:
		return &ast.String{Str: chexp.Literal}
	case *chparser.RatioExpr:
	case *chparser.EnumValueExpr:
	case *chparser.IntervalExpr:
	case *chparser.EngineExpr:
	case *chparser.ColumnTypeExpr:
	case *chparser.WhenExpr:
	case *chparser.CaseExpr:
	case *chparser.CastExpr:
	case *chparser.WithExpr:
	case *chparser.TopExpr:
	case *chparser.CreateLiveView:
	case *chparser.WithTimeoutExpr:
	case *chparser.TableExpr:
	case *chparser.OnExpr:
	case *chparser.UsingExpr:
	case *chparser.JoinExpr:
	case *chparser.JoinConstraintExpr:
	case *chparser.FromExpr:
	case *chparser.IsNullExpr:
	case *chparser.IsNotNullExpr:
	case *chparser.AliasExpr:
	case *chparser.WhereExpr:
		return c.convertWhere(chexp)
	case *chparser.PrewhereExpr:
	case *chparser.GroupByExpr:
	case *chparser.HavingExpr:
	case *chparser.LimitExpr:
	case *chparser.LimitByExpr:
	case *chparser.WindowConditionExpr:
	case *chparser.WindowExpr:
	case *chparser.WindowFrameExpr:
	case *chparser.WindowFrameExtendExpr:
	case *chparser.WindowFrameRangeExpr:
	case *chparser.WindowFrameCurrentRow:
	case *chparser.WindowFrameUnbounded:
	case *chparser.WindowFrameNumber:
	case *chparser.ArrayJoinExpr:
	case *chparser.SelectQuery:
		return c.convertSelectQuery(chexp)
	case *chparser.SubQueryExpr:
	case *chparser.NotExpr:
	case *chparser.NegateExpr:
	case *chparser.GlobalInExpr:
	case *chparser.ExtractExpr:
	case *chparser.DropDatabase:
	case *chparser.DropStmt:
	case *chparser.DropUserOrRole:
	case *chparser.UseExpr:
	case *chparser.CTEExpr:
	case *chparser.SetExpr:
	case *chparser.FormatExpr:
	case *chparser.OptimizeExpr:
	case *chparser.DeduplicateExpr:
	case *chparser.SystemExpr:
	case *chparser.SystemFlushExpr:
	case *chparser.SystemReloadExpr:
	case *chparser.SystemSyncExpr:
	case *chparser.SystemCtrlExpr:
	case *chparser.SystemDropExpr:
	case *chparser.TruncateTable:
	case *chparser.SampleRatioExpr:
	case *chparser.DeleteFromExpr:
	case *chparser.ColumnNamesExpr:
	case *chparser.ValuesExpr:
	case *chparser.InsertExpr:
	case *chparser.CheckExpr:
	case *chparser.UnaryExpr:
	case *chparser.RenameStmt:
	case *chparser.ExplainExpr:
	case *chparser.PrivilegeExpr:
	case *chparser.GrantPrivilegeExpr:

	}
	name := reflect.TypeOf(expr).Elem().Name()
	fmt.Println("missed type", name)
	panic("can not convert expr")
}

func (c *converter) Convert(exprs []chparser.Expr) ([]ast.Statement, error) {
	output := make([]ast.Statement, 0)
	for _, e := range exprs {
		r := c.convertExpr(e)
		output = append(output, ast.Statement{
			Raw: &ast.RawStmt{
				Stmt:         r,
				StmtLocation: r.Pos(),
				StmtLen:      int(e.End()) - int(e.Pos()),
			},
		})
	}
	return output, nil
}

func (c *Parser) Parse(r io.Reader) ([]ast.Statement, error) {
	sql, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	fmt.Println("parsing", string(sql))

	chParser := chparser.NewParser(string(sql))
	statements, err := chParser.ParseStatements()
	if err != nil {
		err = fmt.Errorf("clickhouse parser error: %w", err)
		return []ast.Statement{}, err
	}
	conv := &converter{
		DefaultASTVisitor: &chparser.DefaultASTVisitor{
			Visit: debugFallbackVisitor,
		},
	}
	result, err := conv.Convert(statements)
	if err != nil {
		err = fmt.Errorf("can not convert to sqlc ast: %w", err)
		return nil, err
	}
	for _, v := range result {
		node := reflect.TypeOf(v.Raw.Stmt).Elem().Name()
		fmt.Println("node", node)
	}
	return result, err
}
func (c *Parser) CommentSyntax() source.CommentSyntax {
	return source.CommentSyntax{}

}
func (c *Parser) IsReservedKeyword(string) bool {
	return false
}
