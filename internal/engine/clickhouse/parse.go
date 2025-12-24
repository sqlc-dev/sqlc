package clickhouse

import (
	"context"
	"fmt"
	"io"

	"github.com/sqlc-dev/doubleclick/ast"
	"github.com/sqlc-dev/doubleclick/parser"

	"github.com/sqlc-dev/sqlc/internal/source"
	sqlcast "github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// Parser implements the compiler.Parser interface for ClickHouse.
type Parser struct{}

// NewParser creates a new ClickHouse parser.
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses ClickHouse SQL statements and converts them to sqlc's AST.
func (p *Parser) Parse(r io.Reader) ([]sqlcast.Statement, error) {
	ctx := context.Background()
	stmts, err := parser.Parse(ctx, r)
	if err != nil {
		return nil, err
	}

	var result []sqlcast.Statement
	for _, stmt := range stmts {
		converted := p.convert(stmt)
		if converted == nil {
			continue
		}
		pos := stmt.Pos()
		result = append(result, sqlcast.Statement{
			Raw: &sqlcast.RawStmt{
				Stmt:         converted,
				StmtLocation: pos.Offset,
				StmtLen:      0,
			},
		})
	}

	// Calculate statement lengths
	for i := 0; i < len(result)-1; i++ {
		result[i].Raw.StmtLen = result[i+1].Raw.StmtLocation - result[i].Raw.StmtLocation
	}

	return result, nil
}

// CommentSyntax returns the comment syntax for ClickHouse.
func (p *Parser) CommentSyntax() source.CommentSyntax {
	return source.CommentSyntax{
		Dash:      true,
		Hash:      true,
		SlashStar: true,
	}
}

// IsReservedKeyword checks if a word is a reserved keyword in ClickHouse.
func (p *Parser) IsReservedKeyword(word string) bool {
	return isReserved(word)
}

func (p *Parser) convert(stmt ast.Statement) sqlcast.Node {
	switch s := stmt.(type) {
	case *ast.SelectWithUnionQuery:
		return p.convertSelectWithUnion(s)
	case *ast.SelectQuery:
		return p.convertSelect(s)
	case *ast.InsertQuery:
		return p.convertInsert(s)
	case *ast.CreateQuery:
		return p.convertCreate(s)
	case *ast.DropQuery:
		return p.convertDrop(s)
	case *ast.AlterQuery:
		return p.convertAlter(s)
	case *ast.TruncateQuery:
		return p.convertTruncate(s)
	default:
		return &sqlcast.TODO{}
	}
}

func (p *Parser) convertSelectWithUnion(s *ast.SelectWithUnionQuery) sqlcast.Node {
	if len(s.Selects) == 0 {
		return &sqlcast.TODO{}
	}

	first := p.convert(s.Selects[0])
	if first == nil {
		return &sqlcast.TODO{}
	}

	firstSelect, ok := first.(*sqlcast.SelectStmt)
	if !ok {
		return first
	}

	if len(s.Selects) == 1 {
		return firstSelect
	}

	left := firstSelect
	for i := 1; i < len(s.Selects); i++ {
		right := p.convert(s.Selects[i])
		rightSelect, ok := right.(*sqlcast.SelectStmt)
		if !ok {
			continue
		}
		all := false
		if i-1 < len(s.UnionModes) && s.UnionModes[i-1] == "UNION ALL" {
			all = true
		}
		left = &sqlcast.SelectStmt{
			Op:   sqlcast.Union,
			All:  all,
			Larg: left,
			Rarg: rightSelect,
		}
	}
	return left
}

func (p *Parser) convertSelect(s *ast.SelectQuery) sqlcast.Node {
	stmt := &sqlcast.SelectStmt{}

	var targets []sqlcast.Node
	for _, col := range s.Columns {
		target := p.convertColumnExpr(col)
		if target != nil {
			targets = append(targets, target)
		}
	}
	if len(targets) > 0 {
		stmt.TargetList = &sqlcast.List{Items: targets}
	}

	if s.From != nil {
		stmt.FromClause = p.convertFrom(s.From)
	}

	if s.Where != nil {
		stmt.WhereClause = p.convertExpr(s.Where)
	}

	if len(s.GroupBy) > 0 {
		var groupItems []sqlcast.Node
		for _, g := range s.GroupBy {
			groupItems = append(groupItems, p.convertExpr(g))
		}
		stmt.GroupClause = &sqlcast.List{Items: groupItems}
	}

	if s.Having != nil {
		stmt.HavingClause = p.convertExpr(s.Having)
	}

	if len(s.OrderBy) > 0 {
		var sortItems []sqlcast.Node
		for _, o := range s.OrderBy {
			sortItem := &sqlcast.SortBy{
				Node: p.convertExpr(o.Expression),
			}
			if o.Descending {
				sortItem.SortbyDir = sqlcast.SortByDirDesc
			} else {
				sortItem.SortbyDir = sqlcast.SortByDirAsc
			}
			sortItems = append(sortItems, sortItem)
		}
		stmt.SortClause = &sqlcast.List{Items: sortItems}
	}

	if s.Limit != nil {
		stmt.LimitCount = p.convertExpr(s.Limit)
	}
	if s.Offset != nil {
		stmt.LimitOffset = p.convertExpr(s.Offset)
	}

	return stmt
}

func (p *Parser) convertColumnExpr(expr ast.Expression) sqlcast.Node {
	switch e := expr.(type) {
	case *ast.Asterisk:
		return &sqlcast.ResTarget{
			Val: &sqlcast.ColumnRef{
				Fields: &sqlcast.List{
					Items: []sqlcast.Node{&sqlcast.A_Star{}},
				},
			},
		}
	case *ast.Identifier:
		colRef := &sqlcast.ColumnRef{
			Fields: &sqlcast.List{},
		}
		for _, part := range e.Parts {
			colRef.Fields.Items = append(colRef.Fields.Items, &sqlcast.String{Str: part})
		}
		target := &sqlcast.ResTarget{Val: colRef}
		if e.Alias != "" {
			target.Name = &e.Alias
		}
		return target
	default:
		converted := p.convertExpr(expr)
		if converted == nil {
			return nil
		}
		return &sqlcast.ResTarget{Val: converted}
	}
}

func (p *Parser) convertFrom(from *ast.TablesInSelectQuery) *sqlcast.List {
	if from == nil || len(from.Tables) == 0 {
		return nil
	}

	var items []sqlcast.Node
	for _, table := range from.Tables {
		if table.Table != nil {
			items = append(items, p.convertTableExpr(table.Table))
		}
	}

	return &sqlcast.List{Items: items}
}

func (p *Parser) convertTableExpr(expr *ast.TableExpression) sqlcast.Node {
	switch t := expr.Table.(type) {
	case *ast.TableIdentifier:
		rv := &sqlcast.RangeVar{
			Relname: &t.Table,
		}
		if t.Database != "" {
			rv.Schemaname = &t.Database
		}
		if expr.Alias != "" {
			rv.Alias = &sqlcast.Alias{Aliasname: &expr.Alias}
		}
		return rv
	case *ast.Subquery:
		rs := &sqlcast.RangeSubselect{
			Subquery: p.convert(t.Query),
		}
		if expr.Alias != "" {
			rs.Alias = &sqlcast.Alias{Aliasname: &expr.Alias}
		}
		return rs
	case *ast.FunctionCall:
		fc := &sqlcast.FuncCall{
			Funcname: &sqlcast.List{
				Items: []sqlcast.Node{&sqlcast.String{Str: t.Name}},
			},
		}
		for _, arg := range t.Arguments {
			if fc.Args == nil {
				fc.Args = &sqlcast.List{}
			}
			fc.Args.Items = append(fc.Args.Items, p.convertExpr(arg))
		}
		return &sqlcast.RangeFunction{
			Functions: &sqlcast.List{Items: []sqlcast.Node{fc}},
		}
	default:
		return &sqlcast.RangeVar{}
	}
}

func (p *Parser) convertExpr(expr ast.Expression) sqlcast.Node {
	if expr == nil {
		return nil
	}

	switch e := expr.(type) {
	case *ast.Identifier:
		colRef := &sqlcast.ColumnRef{Fields: &sqlcast.List{}}
		for _, part := range e.Parts {
			colRef.Fields.Items = append(colRef.Fields.Items, &sqlcast.String{Str: part})
		}
		return colRef

	case *ast.Literal:
		switch e.Type {
		case ast.LiteralString:
			return &sqlcast.A_Const{Val: &sqlcast.String{Str: fmt.Sprintf("%v", e.Value)}}
		case ast.LiteralInteger:
			if val, ok := e.Value.(int64); ok {
				return &sqlcast.A_Const{Val: &sqlcast.Integer{Ival: val}}
			}
			return &sqlcast.A_Const{Val: &sqlcast.String{Str: fmt.Sprintf("%v", e.Value)}}
		case ast.LiteralFloat:
			return &sqlcast.A_Const{Val: &sqlcast.Float{Str: fmt.Sprintf("%v", e.Value)}}
		case ast.LiteralNull:
			return &sqlcast.Null{}
		default:
			return &sqlcast.A_Const{Val: &sqlcast.String{Str: fmt.Sprintf("%v", e.Value)}}
		}

	case *ast.BinaryExpr:
		return &sqlcast.A_Expr{
			Name:  &sqlcast.List{Items: []sqlcast.Node{&sqlcast.String{Str: e.Op}}},
			Lexpr: p.convertExpr(e.Left),
			Rexpr: p.convertExpr(e.Right),
		}

	case *ast.UnaryExpr:
		return &sqlcast.A_Expr{
			Name:  &sqlcast.List{Items: []sqlcast.Node{&sqlcast.String{Str: e.Op}}},
			Rexpr: p.convertExpr(e.Operand),
		}

	case *ast.FunctionCall:
		fc := &sqlcast.FuncCall{
			Funcname: &sqlcast.List{Items: []sqlcast.Node{&sqlcast.String{Str: e.Name}}},
		}
		if len(e.Arguments) > 0 {
			fc.Args = &sqlcast.List{}
			for _, arg := range e.Arguments {
				fc.Args.Items = append(fc.Args.Items, p.convertExpr(arg))
			}
		}
		return fc

	case *ast.Subquery:
		return &sqlcast.SubLink{Subselect: p.convert(e.Query)}

	case *ast.Asterisk:
		return &sqlcast.ColumnRef{
			Fields: &sqlcast.List{Items: []sqlcast.Node{&sqlcast.A_Star{}}},
		}

	default:
		return &sqlcast.A_Const{Val: &sqlcast.String{Str: ""}}
	}
}

func (p *Parser) convertInsert(s *ast.InsertQuery) sqlcast.Node {
	stmt := &sqlcast.InsertStmt{
		Relation: &sqlcast.RangeVar{Relname: &s.Table},
	}

	if s.Database != "" {
		stmt.Relation.Schemaname = &s.Database
	}

	if len(s.Columns) > 0 {
		stmt.Cols = &sqlcast.List{}
		for _, col := range s.Columns {
			if len(col.Parts) > 0 {
				name := col.Parts[0]
				stmt.Cols.Items = append(stmt.Cols.Items, &sqlcast.ResTarget{Name: &name})
			}
		}
	}

	if s.Select != nil {
		stmt.SelectStmt = p.convert(s.Select)
	}

	return stmt
}

func (p *Parser) convertCreate(s *ast.CreateQuery) sqlcast.Node {
	if s.CreateDatabase {
		return &sqlcast.CreateSchemaStmt{
			Name:        &s.Database,
			IfNotExists: s.IfNotExists,
		}
	}

	stmt := &sqlcast.CreateTableStmt{
		Name:        &sqlcast.TableName{Name: s.Table},
		IfNotExists: s.IfNotExists,
	}

	if s.Database != "" {
		stmt.Name.Schema = s.Database
	}

	for _, col := range s.Columns {
		colDef := &sqlcast.ColumnDef{Colname: col.Name}
		if col.Type != nil {
			colDef.TypeName = &sqlcast.TypeName{Name: col.Type.Name}
		}
		if col.Nullable != nil && !*col.Nullable {
			colDef.IsNotNull = true
		}
		stmt.Cols = append(stmt.Cols, colDef)
	}

	return stmt
}

func (p *Parser) convertDrop(s *ast.DropQuery) sqlcast.Node {
	if s.DropDatabase {
		return &sqlcast.DropSchemaStmt{
			Schemas:   []*sqlcast.String{{Str: s.Database}},
			MissingOk: s.IfExists,
		}
	}

	tables := []*sqlcast.TableName{}
	if s.Table != "" {
		tableName := &sqlcast.TableName{Name: s.Table}
		if s.Database != "" {
			tableName.Schema = s.Database
		}
		tables = append(tables, tableName)
	}

	return &sqlcast.DropTableStmt{
		Tables:   tables,
		IfExists: s.IfExists,
	}
}

func (p *Parser) convertAlter(s *ast.AlterQuery) sqlcast.Node {
	tableName := &sqlcast.TableName{Name: s.Table}
	if s.Database != "" {
		tableName.Schema = s.Database
	}

	stmt := &sqlcast.AlterTableStmt{
		Table: tableName,
		Cmds:  &sqlcast.List{},
	}

	for _, cmd := range s.Commands {
		switch cmd.Type {
		case ast.AlterAddColumn:
			if cmd.Column != nil {
				colDef := &sqlcast.ColumnDef{Colname: cmd.Column.Name}
				if cmd.Column.Type != nil {
					colDef.TypeName = &sqlcast.TypeName{Name: cmd.Column.Type.Name}
				}
				altCmd := &sqlcast.AlterTableCmd{
					Subtype: sqlcast.AT_AddColumn,
					Def:     colDef,
				}
				stmt.Cmds.Items = append(stmt.Cmds.Items, altCmd)
			}
		case ast.AlterDropColumn:
			altCmd := &sqlcast.AlterTableCmd{
				Subtype:   sqlcast.AT_DropColumn,
				Name:      &cmd.ColumnName,
				MissingOk: cmd.IfExists,
			}
			stmt.Cmds.Items = append(stmt.Cmds.Items, altCmd)
		}
	}

	return stmt
}

func (p *Parser) convertTruncate(s *ast.TruncateQuery) sqlcast.Node {
	rv := &sqlcast.RangeVar{Relname: &s.Table}
	if s.Database != "" {
		rv.Schemaname = &s.Database
	}

	return &sqlcast.TruncateStmt{
		Relations: &sqlcast.List{Items: []sqlcast.Node{rv}},
	}
}
