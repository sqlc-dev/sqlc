package mssql

import (
	"context"
	"io"
	"strconv"
	"strings"

	"github.com/sqlc-dev/teesql/ast"
	"github.com/sqlc-dev/teesql/parser"

	"github.com/sqlc-dev/sqlc/internal/source"
	sqast "github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func NewParser() *Parser {
	return &Parser{}
}

type Parser struct{}

func (p *Parser) Parse(r io.Reader) ([]sqast.Statement, error) {
	script, err := parser.Parse(context.Background(), r)
	if err != nil {
		return nil, err
	}

	var stmts []sqast.Statement
	for _, batch := range script.Batches {
		for _, stmt := range batch.Statements {
			n := convert(stmt)
			if n == nil {
				continue
			}
			stmts = append(stmts, sqast.Statement{
				Raw: &sqast.RawStmt{
					Stmt: n,
				},
			})
		}
	}
	return stmts, nil
}

// CommentSyntax returns the comment syntax for T-SQL
// https://docs.microsoft.com/en-us/sql/t-sql/language-elements/comments-transact-sql
func (p *Parser) CommentSyntax() source.CommentSyntax {
	return source.CommentSyntax{
		Dash:      true,
		SlashStar: true,
	}
}

// IsReservedKeyword checks if the given string is a T-SQL reserved keyword
func (p *Parser) IsReservedKeyword(s string) bool {
	return reserved[strings.ToUpper(s)]
}

// Param returns the T-SQL parameter placeholder for the given position
func (p *Parser) Param(n int) string {
	return "@p" + strconv.Itoa(n)
}

// NamedParam returns the named parameter placeholder for T-SQL
func (p *Parser) NamedParam(name string) string {
	return "@" + name
}

// QuoteIdent returns a quoted identifier for T-SQL using square brackets
func (p *Parser) QuoteIdent(s string) string {
	if p.IsReservedKeyword(s) {
		return "[" + s + "]"
	}
	return s
}

// TypeName returns the SQL type name for T-SQL
func (p *Parser) TypeName(ns, name string) string {
	return name
}

// Cast formats a type cast expression for T-SQL
func (p *Parser) Cast(arg, typeName string) string {
	return "CAST(" + arg + " AS " + typeName + ")"
}

// convert converts a teesql AST statement to a sqlc AST node
func convert(stmt ast.Statement) sqast.Node {
	switch s := stmt.(type) {
	case *ast.SelectStatement:
		return convertSelectStatement(s)
	case *ast.InsertStatement:
		return convertInsertStatement(s)
	case *ast.UpdateStatement:
		return convertUpdateStatement(s)
	case *ast.DeleteStatement:
		return convertDeleteStatement(s)
	case *ast.CreateTableStatement:
		return convertCreateTableStatement(s)
	case *ast.AlterTableAddTableElementStatement:
		return convertAlterTableAddStatement(s)
	case *ast.AlterTableDropTableElementStatement:
		return convertAlterTableDropStatement(s)
	case *ast.DropTableStatement:
		return convertDropTableStatement(s)
	case *ast.DropViewStatement:
		return convertDropViewStatement(s)
	case *ast.CreateViewStatement:
		return convertCreateViewStatement(s)
	case *ast.CreateProcedureStatement:
		return convertCreateProcedureStatement(s)
	default:
		// Return a TODO node for unsupported statements
		return &sqast.TODO{}
	}
}

func convertSelectStatement(s *ast.SelectStatement) sqast.Node {
	if s == nil || s.QueryExpression == nil {
		return &sqast.TODO{}
	}
	return convertQueryExpression(s.QueryExpression)
}

func convertQueryExpression(qe ast.QueryExpression) sqast.Node {
	switch q := qe.(type) {
	case *ast.QuerySpecification:
		return convertQuerySpecification(q)
	case *ast.QueryParenthesisExpression:
		if q.QueryExpression != nil {
			return convertQueryExpression(q.QueryExpression)
		}
		return &sqast.TODO{}
	case *ast.BinaryQueryExpression:
		// Handle UNION, EXCEPT, INTERSECT
		left := convertQueryExpression(q.FirstQueryExpression)
		right := convertQueryExpression(q.SecondQueryExpression)

		// Type assert to *SelectStmt as required by SelectStmt.Larg/Rarg
		leftStmt, _ := left.(*sqast.SelectStmt)
		rightStmt, _ := right.(*sqast.SelectStmt)

		return &sqast.SelectStmt{
			Op:   convertSetOperator(q.BinaryQueryExpressionType),
			Larg: leftStmt,
			Rarg: rightStmt,
		}
	default:
		return &sqast.TODO{}
	}
}

func convertSetOperator(op string) sqast.SetOperation {
	switch strings.ToUpper(op) {
	case "UNION":
		return sqast.Union
	case "EXCEPT":
		return sqast.Except
	case "INTERSECT":
		return sqast.Intersect
	default:
		return sqast.None
	}
}

func convertQuerySpecification(qs *ast.QuerySpecification) *sqast.SelectStmt {
	if qs == nil {
		return &sqast.SelectStmt{}
	}

	stmt := &sqast.SelectStmt{
		TargetList:  &sqast.List{},
		FromClause:  &sqast.List{},
		GroupClause: &sqast.List{},
	}

	// Convert SELECT elements (target list)
	if qs.SelectElements != nil {
		for _, elem := range qs.SelectElements {
			target := convertSelectElement(elem)
			if target != nil {
				stmt.TargetList.Items = append(stmt.TargetList.Items, target)
			}
		}
	}

	// Convert FROM clause
	if qs.FromClause != nil && qs.FromClause.TableReferences != nil {
		for _, tref := range qs.FromClause.TableReferences {
			from := convertTableReference(tref)
			if from != nil {
				stmt.FromClause.Items = append(stmt.FromClause.Items, from)
			}
		}
	}

	// Convert WHERE clause
	if qs.WhereClause != nil && qs.WhereClause.SearchCondition != nil {
		stmt.WhereClause = convertBooleanExpression(qs.WhereClause.SearchCondition)
	}

	// Convert GROUP BY clause
	if qs.GroupByClause != nil {
		for _, spec := range qs.GroupByClause.GroupingSpecifications {
			group := convertGroupingSpecification(spec)
			if group != nil {
				stmt.GroupClause.Items = append(stmt.GroupClause.Items, group)
			}
		}
	}

	// Convert HAVING clause
	if qs.HavingClause != nil && qs.HavingClause.SearchCondition != nil {
		stmt.HavingClause = convertBooleanExpression(qs.HavingClause.SearchCondition)
	}

	return stmt
}

func convertSelectElement(elem ast.SelectElement) sqast.Node {
	switch e := elem.(type) {
	case *ast.SelectStarExpression:
		return &sqast.ResTarget{
			Val: &sqast.ColumnRef{
				Fields: &sqast.List{
					Items: []sqast.Node{&sqast.A_Star{}},
				},
			},
		}
	case *ast.SelectScalarExpression:
		target := &sqast.ResTarget{}
		if e.Expression != nil {
			target.Val = convertScalarExpression(e.Expression)
		}
		if e.ColumnName != nil && e.ColumnName.Value != "" {
			name := e.ColumnName.Value
			target.Name = &name
		}
		return target
	case *ast.SelectSetVariable:
		return &sqast.TODO{}
	default:
		return &sqast.TODO{}
	}
}

func convertScalarExpression(expr ast.ScalarExpression) sqast.Node {
	switch e := expr.(type) {
	case *ast.ColumnReferenceExpression:
		return convertColumnReference(e)
	case *ast.IntegerLiteral:
		val, _ := strconv.ParseInt(e.Value, 10, 64)
		return &sqast.A_Const{
			Val: &sqast.Integer{Ival: val},
		}
	case *ast.StringLiteral:
		return &sqast.A_Const{
			Val: &sqast.String{Str: e.Value},
		}
	case *ast.NumericLiteral:
		return &sqast.A_Const{
			Val: &sqast.Float{Str: e.Value},
		}
	case *ast.NullLiteral:
		return &sqast.Null{}
	case *ast.VariableReference:
		// Convert @param to parameter reference
		return &sqast.ParamRef{
			Dollar: true,
		}
	case *ast.FunctionCall:
		return convertFunctionCall(e)
	case *ast.BinaryExpression:
		return &sqast.A_Expr{
			Name:  &sqast.List{Items: []sqast.Node{&sqast.String{Str: e.BinaryExpressionType}}},
			Lexpr: convertScalarExpression(e.FirstExpression),
			Rexpr: convertScalarExpression(e.SecondExpression),
		}
	case *ast.UnaryExpression:
		return &sqast.A_Expr{
			Name:  &sqast.List{Items: []sqast.Node{&sqast.String{Str: e.UnaryExpressionType}}},
			Rexpr: convertScalarExpression(e.Expression),
		}
	case *ast.ParenthesisExpression:
		if e.Expression != nil {
			return convertScalarExpression(e.Expression)
		}
		return &sqast.TODO{}
	case *ast.SearchedCaseExpression:
		return convertSearchedCaseExpression(e)
	case *ast.SimpleCaseExpression:
		return convertSimpleCaseExpression(e)
	case *ast.ScalarSubquery:
		if e.QueryExpression != nil {
			return &sqast.SubLink{
				SubLinkType: sqast.EXPR_SUBLINK,
				Subselect:   convertQueryExpression(e.QueryExpression),
			}
		}
		return &sqast.TODO{}
	default:
		return &sqast.TODO{}
	}
}

func convertColumnReference(cr *ast.ColumnReferenceExpression) sqast.Node {
	if cr == nil || cr.MultiPartIdentifier == nil {
		return &sqast.TODO{}
	}

	fields := &sqast.List{}
	for _, id := range cr.MultiPartIdentifier.Identifiers {
		if id != nil {
			fields.Items = append(fields.Items, &sqast.String{Str: id.Value})
		}
	}

	return &sqast.ColumnRef{Fields: fields}
}

func convertFunctionCall(fc *ast.FunctionCall) sqast.Node {
	if fc == nil {
		return &sqast.TODO{}
	}

	fn := &sqast.FuncCall{
		Args: &sqast.List{},
	}

	// Build function name from FunctionName identifier
	if fc.FunctionName != nil && fc.FunctionName.Value != "" {
		fn.Funcname = &sqast.List{
			Items: []sqast.Node{
				&sqast.String{Str: fc.FunctionName.Value},
			},
		}
	}

	// Convert arguments
	for _, param := range fc.Parameters {
		if param != nil {
			arg := convertScalarExpression(param)
			fn.Args.Items = append(fn.Args.Items, arg)
		}
	}

	return fn
}

func convertSearchedCaseExpression(ce *ast.SearchedCaseExpression) sqast.Node {
	caseExpr := &sqast.CaseExpr{
		Args: &sqast.List{},
	}

	for _, when := range ce.WhenClauses {
		if when.WhenExpression != nil && when.ThenExpression != nil {
			caseWhen := &sqast.CaseWhen{
				Expr:   convertBooleanExpression(when.WhenExpression),
				Result: convertScalarExpression(when.ThenExpression),
			}
			caseExpr.Args.Items = append(caseExpr.Args.Items, caseWhen)
		}
	}

	if ce.ElseExpression != nil {
		caseExpr.Defresult = convertScalarExpression(ce.ElseExpression)
	}

	return caseExpr
}

func convertSimpleCaseExpression(ce *ast.SimpleCaseExpression) sqast.Node {
	caseExpr := &sqast.CaseExpr{
		Args: &sqast.List{},
	}

	if ce.InputExpression != nil {
		caseExpr.Arg = convertScalarExpression(ce.InputExpression)
	}

	for _, when := range ce.WhenClauses {
		if when.WhenExpression != nil && when.ThenExpression != nil {
			caseWhen := &sqast.CaseWhen{
				Expr:   convertScalarExpression(when.WhenExpression),
				Result: convertScalarExpression(when.ThenExpression),
			}
			caseExpr.Args.Items = append(caseExpr.Args.Items, caseWhen)
		}
	}

	if ce.ElseExpression != nil {
		caseExpr.Defresult = convertScalarExpression(ce.ElseExpression)
	}

	return caseExpr
}

func convertBooleanExpression(expr ast.BooleanExpression) sqast.Node {
	switch e := expr.(type) {
	case *ast.BooleanComparisonExpression:
		return &sqast.A_Expr{
			Name:  &sqast.List{Items: []sqast.Node{&sqast.String{Str: e.ComparisonType}}},
			Lexpr: convertScalarExpression(e.FirstExpression),
			Rexpr: convertScalarExpression(e.SecondExpression),
		}
	case *ast.BooleanBinaryExpression:
		var op string
		switch strings.ToUpper(e.BinaryExpressionType) {
		case "AND":
			op = "AND"
		case "OR":
			op = "OR"
		default:
			op = e.BinaryExpressionType
		}
		return &sqast.BoolExpr{
			Boolop: convertBoolOp(op),
			Args: &sqast.List{
				Items: []sqast.Node{
					convertBooleanExpression(e.FirstExpression),
					convertBooleanExpression(e.SecondExpression),
				},
			},
		}
	case *ast.BooleanIsNullExpression:
		nullTest := &sqast.NullTest{
			Arg: convertScalarExpression(e.Expression),
		}
		if e.IsNot {
			nullTest.Nulltesttype = sqast.NullTestTypeIsNotNull
		} else {
			nullTest.Nulltesttype = sqast.NullTestTypeIsNull
		}
		return nullTest
	case *ast.BooleanInExpression:
		return &sqast.A_Expr{
			Kind:  sqast.A_Expr_Kind(1), // A_Expr_IN
			Lexpr: convertScalarExpression(e.Expression),
		}
	case *ast.BooleanLikeExpression:
		op := "LIKE"
		if e.NotDefined {
			op = "NOT LIKE"
		}
		return &sqast.A_Expr{
			Name:  &sqast.List{Items: []sqast.Node{&sqast.String{Str: op}}},
			Lexpr: convertScalarExpression(e.FirstExpression),
			Rexpr: convertScalarExpression(e.SecondExpression),
		}
	case *ast.BooleanParenthesisExpression:
		if e.Expression != nil {
			return convertBooleanExpression(e.Expression)
		}
		return &sqast.TODO{}
	default:
		return &sqast.TODO{}
	}
}

func convertBoolOp(op string) sqast.BoolExprType {
	switch strings.ToUpper(op) {
	case "AND":
		return sqast.BoolExprTypeAnd
	case "OR":
		return sqast.BoolExprTypeOr
	case "NOT":
		return sqast.BoolExprTypeNot
	default:
		return sqast.BoolExprTypeAnd
	}
}

func convertTableReference(tref ast.TableReference) sqast.Node {
	switch t := tref.(type) {
	case *ast.NamedTableReference:
		return convertNamedTableReference(t)
	case *ast.QualifiedJoin:
		return convertQualifiedJoin(t)
	case *ast.UnqualifiedJoin:
		return convertUnqualifiedJoin(t)
	default:
		return &sqast.TODO{}
	}
}

// strPtr returns a pointer to the string, or nil if empty
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func convertNamedTableReference(ntr *ast.NamedTableReference) *sqast.RangeVar {
	if ntr == nil || ntr.SchemaObject == nil {
		return &sqast.RangeVar{}
	}

	rv := &sqast.RangeVar{}

	so := ntr.SchemaObject
	if so.DatabaseIdentifier != nil {
		rv.Catalogname = strPtr(so.DatabaseIdentifier.Value)
	}
	if so.SchemaIdentifier != nil {
		rv.Schemaname = strPtr(so.SchemaIdentifier.Value)
	}
	if so.BaseIdentifier != nil {
		rv.Relname = strPtr(so.BaseIdentifier.Value)
	}

	if ntr.Alias != nil && ntr.Alias.Value != "" {
		rv.Alias = &sqast.Alias{Aliasname: strPtr(ntr.Alias.Value)}
	}

	return rv
}

func convertQualifiedJoin(qj *ast.QualifiedJoin) sqast.Node {
	join := &sqast.JoinExpr{}

	if qj.FirstTableReference != nil {
		join.Larg = convertTableReference(qj.FirstTableReference)
	}
	if qj.SecondTableReference != nil {
		join.Rarg = convertTableReference(qj.SecondTableReference)
	}

	// Set join type
	switch strings.ToUpper(qj.QualifiedJoinType) {
	case "INNER":
		join.Jointype = sqast.JoinTypeInner
	case "LEFT":
		join.Jointype = sqast.JoinTypeLeft
	case "RIGHT":
		join.Jointype = sqast.JoinTypeRight
	case "FULL":
		join.Jointype = sqast.JoinTypeFull
	default:
		join.Jointype = sqast.JoinTypeInner
	}

	// Convert ON clause
	if qj.SearchCondition != nil {
		join.Quals = convertBooleanExpression(qj.SearchCondition)
	}

	return join
}

func convertUnqualifiedJoin(uj *ast.UnqualifiedJoin) sqast.Node {
	// CROSS JOIN is represented as JoinTypeInner with no Quals in sqlc's AST
	join := &sqast.JoinExpr{
		Jointype: sqast.JoinTypeInner,
	}

	if uj.FirstTableReference != nil {
		join.Larg = convertTableReference(uj.FirstTableReference)
	}
	if uj.SecondTableReference != nil {
		join.Rarg = convertTableReference(uj.SecondTableReference)
	}

	return join
}

func convertGroupingSpecification(spec ast.GroupingSpecification) sqast.Node {
	switch s := spec.(type) {
	case *ast.ExpressionGroupingSpecification:
		if s.Expression != nil {
			return convertScalarExpression(s.Expression)
		}
		return &sqast.TODO{}
	default:
		return &sqast.TODO{}
	}
}

func convertInsertStatement(s *ast.InsertStatement) sqast.Node {
	if s == nil || s.InsertSpecification == nil {
		return &sqast.TODO{}
	}

	spec := s.InsertSpecification
	stmt := &sqast.InsertStmt{
		Cols: &sqast.List{},
	}

	// Convert target table
	if spec.Target != nil {
		if ntr, ok := spec.Target.(*ast.NamedTableReference); ok {
			stmt.Relation = convertNamedTableReference(ntr)
		}
	}

	// Convert column list
	for _, col := range spec.Columns {
		if col != nil && col.MultiPartIdentifier != nil && len(col.MultiPartIdentifier.Identifiers) > 0 {
			// Get the last identifier (column name)
			lastId := col.MultiPartIdentifier.Identifiers[len(col.MultiPartIdentifier.Identifiers)-1]
			if lastId != nil {
				stmt.Cols.Items = append(stmt.Cols.Items, &sqast.ResTarget{
					Name: strPtr(lastId.Value),
				})
			}
		}
	}

	// Convert values or select
	if spec.InsertSource != nil {
		switch src := spec.InsertSource.(type) {
		case *ast.ValuesInsertSource:
			// Handle VALUES clauses
			stmt.SelectStmt = &sqast.TODO{}
		case *ast.SelectInsertSource:
			if src.Select != nil {
				stmt.SelectStmt = convertQueryExpression(src.Select)
			}
		}
	}

	return stmt
}

func convertUpdateStatement(s *ast.UpdateStatement) sqast.Node {
	if s == nil || s.UpdateSpecification == nil {
		return &sqast.TODO{}
	}

	spec := s.UpdateSpecification
	stmt := &sqast.UpdateStmt{
		Relations:  &sqast.List{},
		TargetList: &sqast.List{},
		FromClause: &sqast.List{},
	}

	// Convert target table
	if spec.Target != nil {
		if ntr, ok := spec.Target.(*ast.NamedTableReference); ok {
			rv := convertNamedTableReference(ntr)
			stmt.Relations.Items = append(stmt.Relations.Items, rv)
		}
	}

	// Convert SET clauses
	for _, clause := range spec.SetClauses {
		if assign, ok := clause.(*ast.AssignmentSetClause); ok {
			if assign.Column != nil && assign.NewValue != nil {
				target := &sqast.ResTarget{
					Val: convertScalarExpression(assign.NewValue),
				}
				if assign.Column.MultiPartIdentifier != nil && len(assign.Column.MultiPartIdentifier.Identifiers) > 0 {
					lastId := assign.Column.MultiPartIdentifier.Identifiers[len(assign.Column.MultiPartIdentifier.Identifiers)-1]
					if lastId != nil {
						target.Name = strPtr(lastId.Value)
					}
				}
				stmt.TargetList.Items = append(stmt.TargetList.Items, target)
			}
		}
	}

	// Convert WHERE clause
	if spec.WhereClause != nil && spec.WhereClause.SearchCondition != nil {
		stmt.WhereClause = convertBooleanExpression(spec.WhereClause.SearchCondition)
	}

	// Convert FROM clause
	if spec.FromClause != nil && spec.FromClause.TableReferences != nil {
		for _, tref := range spec.FromClause.TableReferences {
			from := convertTableReference(tref)
			if from != nil {
				stmt.FromClause.Items = append(stmt.FromClause.Items, from)
			}
		}
	}

	return stmt
}

func convertDeleteStatement(s *ast.DeleteStatement) sqast.Node {
	if s == nil || s.DeleteSpecification == nil {
		return &sqast.TODO{}
	}

	spec := s.DeleteSpecification
	stmt := &sqast.DeleteStmt{
		Relations: &sqast.List{},
	}

	// Convert target table
	if spec.Target != nil {
		if ntr, ok := spec.Target.(*ast.NamedTableReference); ok {
			rv := convertNamedTableReference(ntr)
			stmt.Relations.Items = append(stmt.Relations.Items, rv)
		}
	}

	// Convert WHERE clause
	if spec.WhereClause != nil && spec.WhereClause.SearchCondition != nil {
		stmt.WhereClause = convertBooleanExpression(spec.WhereClause.SearchCondition)
	}

	return stmt
}

// extractTypeName extracts the type name from a DataTypeReference
func extractTypeName(dt ast.DataTypeReference) string {
	if dt == nil {
		return ""
	}
	switch t := dt.(type) {
	case *ast.SqlDataTypeReference:
		if t.SqlDataTypeOption != "" {
			return t.SqlDataTypeOption
		}
		if t.Name != nil && t.Name.BaseIdentifier != nil {
			return t.Name.BaseIdentifier.Value
		}
	case *ast.UserDataTypeReference:
		if t.Name != nil && t.Name.BaseIdentifier != nil {
			return t.Name.BaseIdentifier.Value
		}
	case *ast.XmlDataTypeReference:
		return "xml"
	}
	return ""
}

// isNotNullConstraint checks if the constraint is a NOT NULL constraint
func isNotNullConstraint(constraint ast.ConstraintDefinition) bool {
	if nc, ok := constraint.(*ast.NullableConstraintDefinition); ok {
		return !nc.Nullable // Nullable=false means NOT NULL
	}
	return false
}

func convertCreateTableStatement(s *ast.CreateTableStatement) sqast.Node {
	if s == nil || s.SchemaObjectName == nil {
		return &sqast.TODO{}
	}

	stmt := &sqast.CreateTableStmt{
		Name: &sqast.TableName{},
	}

	so := s.SchemaObjectName
	if so.DatabaseIdentifier != nil {
		stmt.Name.Catalog = so.DatabaseIdentifier.Value
	}
	if so.SchemaIdentifier != nil {
		stmt.Name.Schema = so.SchemaIdentifier.Value
	}
	if so.BaseIdentifier != nil {
		stmt.Name.Name = so.BaseIdentifier.Value
	}

	// Convert columns
	if s.Definition != nil && s.Definition.ColumnDefinitions != nil {
		for _, colDef := range s.Definition.ColumnDefinitions {
			if colDef == nil {
				continue
			}
			col := &sqast.ColumnDef{}
			if colDef.ColumnIdentifier != nil {
				col.Colname = colDef.ColumnIdentifier.Value
			}
			if colDef.DataType != nil {
				col.TypeName = &sqast.TypeName{
					Name: extractTypeName(colDef.DataType),
				}
			}
			// Check for NOT NULL constraint
			for _, constraint := range colDef.Constraints {
				if constraint != nil && isNotNullConstraint(constraint) {
					col.IsNotNull = true
				}
			}
			stmt.Cols = append(stmt.Cols, col)
		}
	}

	return stmt
}

func convertAlterTableAddStatement(s *ast.AlterTableAddTableElementStatement) sqast.Node {
	if s == nil || s.SchemaObjectName == nil {
		return &sqast.TODO{}
	}

	stmt := &sqast.AlterTableStmt{
		Table: &sqast.TableName{},
		Cmds:  &sqast.List{},
	}

	so := s.SchemaObjectName
	if so.DatabaseIdentifier != nil {
		stmt.Table.Catalog = so.DatabaseIdentifier.Value
	}
	if so.SchemaIdentifier != nil {
		stmt.Table.Schema = so.SchemaIdentifier.Value
	}
	if so.BaseIdentifier != nil {
		stmt.Table.Name = so.BaseIdentifier.Value
	}

	// Convert column definitions
	if s.Definition != nil && s.Definition.ColumnDefinitions != nil {
		for _, colDef := range s.Definition.ColumnDefinitions {
			if colDef == nil {
				continue
			}
			col := &sqast.ColumnDef{}
			if colDef.ColumnIdentifier != nil {
				col.Colname = colDef.ColumnIdentifier.Value
			}
			if colDef.DataType != nil {
				col.TypeName = &sqast.TypeName{
					Name: extractTypeName(colDef.DataType),
				}
			}
			for _, constraint := range colDef.Constraints {
				if constraint != nil && isNotNullConstraint(constraint) {
					col.IsNotNull = true
				}
			}

			stmt.Cmds.Items = append(stmt.Cmds.Items, &sqast.AlterTableCmd{
				Subtype: sqast.AT_AddColumn,
				Def:     col,
			})
		}
	}

	return stmt
}

func convertAlterTableDropStatement(s *ast.AlterTableDropTableElementStatement) sqast.Node {
	if s == nil || s.SchemaObjectName == nil {
		return &sqast.TODO{}
	}

	stmt := &sqast.AlterTableStmt{
		Table: &sqast.TableName{},
		Cmds:  &sqast.List{},
	}

	so := s.SchemaObjectName
	if so.SchemaIdentifier != nil {
		stmt.Table.Schema = so.SchemaIdentifier.Value
	}
	if so.BaseIdentifier != nil {
		stmt.Table.Name = so.BaseIdentifier.Value
	}

	// Convert drop elements
	for _, elem := range s.AlterTableDropTableElements {
		if elem == nil || elem.Name == nil {
			continue
		}
		name := elem.Name.Value
		stmt.Cmds.Items = append(stmt.Cmds.Items, &sqast.AlterTableCmd{
			Subtype: sqast.AT_DropColumn,
			Name:    &name,
		})
	}

	return stmt
}

func convertDropTableStatement(s *ast.DropTableStatement) sqast.Node {
	if s == nil {
		return &sqast.TODO{}
	}

	stmt := &sqast.DropTableStmt{
		IfExists: s.IsIfExists,
	}
	for _, obj := range s.Objects {
		if obj == nil {
			continue
		}
		tbl := &sqast.TableName{}
		if obj.SchemaIdentifier != nil {
			tbl.Schema = obj.SchemaIdentifier.Value
		}
		if obj.BaseIdentifier != nil {
			tbl.Name = obj.BaseIdentifier.Value
		}
		stmt.Tables = append(stmt.Tables, tbl)
	}
	return stmt
}

func convertDropViewStatement(s *ast.DropViewStatement) sqast.Node {
	if s == nil {
		return &sqast.TODO{}
	}

	stmt := &sqast.DropTableStmt{
		IfExists: s.IsIfExists,
	}
	for _, obj := range s.Objects {
		if obj == nil {
			continue
		}
		tbl := &sqast.TableName{}
		if obj.SchemaIdentifier != nil {
			tbl.Schema = obj.SchemaIdentifier.Value
		}
		if obj.BaseIdentifier != nil {
			tbl.Name = obj.BaseIdentifier.Value
		}
		stmt.Tables = append(stmt.Tables, tbl)
	}
	return stmt
}

func convertCreateViewStatement(s *ast.CreateViewStatement) sqast.Node {
	if s == nil || s.SchemaObjectName == nil {
		return &sqast.TODO{}
	}

	rv := &sqast.RangeVar{}
	if s.SchemaObjectName.SchemaIdentifier != nil {
		rv.Schemaname = strPtr(s.SchemaObjectName.SchemaIdentifier.Value)
	}
	if s.SchemaObjectName.BaseIdentifier != nil {
		rv.Relname = strPtr(s.SchemaObjectName.BaseIdentifier.Value)
	}

	return &sqast.ViewStmt{
		View: rv,
	}
}

func convertCreateProcedureStatement(s *ast.CreateProcedureStatement) sqast.Node {
	if s == nil || s.ProcedureReference == nil {
		return &sqast.TODO{}
	}

	stmt := &sqast.CreateFunctionStmt{
		Func: &sqast.FuncName{},
	}

	if s.ProcedureReference.Name != nil {
		if s.ProcedureReference.Name.SchemaIdentifier != nil {
			stmt.Func.Schema = s.ProcedureReference.Name.SchemaIdentifier.Value
		}
		if s.ProcedureReference.Name.BaseIdentifier != nil {
			stmt.Func.Name = s.ProcedureReference.Name.BaseIdentifier.Value
		}
	}

	return stmt
}
