package clickhouse

import (
	"log"
	"strconv"
	"strings"

	chparser "github.com/AfterShip/clickhouse-sql-parser/parser"

	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

type cc struct {
	paramCount    int
	subqueryCount int
	catalog       interface{} // *catalog.Catalog - using interface{} to avoid circular imports
}

func todo(n chparser.Expr) *ast.TODO {
	if debug.Active {
		log.Printf("clickhouse.convert: Unsupported AST node type %T\n", n)
		log.Printf("clickhouse.convert: This node type may not be fully supported yet. Consider using different query syntax or filing an issue.\n")
	}
	return &ast.TODO{}
}

// identifier preserves the case of identifiers as ClickHouse is case-sensitive
// for table, column, and schema names
func identifier(id string) string {
	return id
}

// normalizeFunctionName normalizes function names to lowercase for comparison
// ClickHouse function names are case-insensitive, so we normalize them for lookups
func normalizeFunctionName(name string) string {
	return strings.ToLower(name)
}

func NewIdentifier(t string) *ast.String {
	return &ast.String{Str: identifier(t)}
}

// getCatalog safely casts the interface{} catalog to a *catalog.Catalog
func (c *cc) getCatalog() *catalog.Catalog {
	if c.catalog == nil {
		return nil
	}
	cat, ok := c.catalog.(*catalog.Catalog)
	if !ok {
		return nil
	}
	return cat
}

// registerFunctionInCatalog registers or updates a function in the catalog with the given return type
func (c *cc) registerFunctionInCatalog(funcName string, returnType *ast.TypeName) {
	cat := c.getCatalog()
	if cat == nil {
		return
	}

	// Find the default schema
	var schema *catalog.Schema
	for _, s := range cat.Schemas {
		if s.Name == cat.DefaultSchema {
			schema = s
			break
		}
	}
	if schema == nil {
		return
	}

	// Check if function already exists
	for i, f := range schema.Funcs {
		if strings.ToLower(f.Name) == strings.ToLower(funcName) {
			// Update existing function
			schema.Funcs[i].ReturnType = returnType
			return
		}
	}

	// Add new function
	schema.Funcs = append(schema.Funcs, &catalog.Function{
		Name:       strings.ToLower(funcName),
		ReturnType: returnType,
	})
}

// findColumnTypeInCatalog searches all tables in the catalog for a column and returns its type
// This is used for unqualified column references (no table prefix)
// Returns the type of the first matching column found, or empty string if not found
// If multiple tables have the same column name, this is ambiguous but we return the first match
// (relying on the query to be syntactically valid from ClickHouse's perspective)
// Column names are case-sensitive in ClickHouse
func (c *cc) findColumnTypeInCatalog(columnName string) string {
	cat := c.getCatalog()
	if cat == nil {
		return ""
	}

	// Search all schemas
	for _, schema := range cat.Schemas {
		if schema == nil || schema.Tables == nil {
			continue
		}
		// Search all tables in this schema
		for _, table := range schema.Tables {
			if table == nil || table.Columns == nil {
				continue
			}
			// Search all columns in this table
			// Column names are case-sensitive
			for _, col := range table.Columns {
				if col.Name == columnName {
					return col.Type.Name
				}
			}
		}
	}

	return ""
}

// extractTypeFromColumnRef extracts the type of a column reference from the catalog
// Returns empty string if the column cannot be resolved
// Column names are case-sensitive in ClickHouse
func (c *cc) extractTypeFromColumnRef(colRef *ast.ColumnRef) string {
	if colRef == nil || colRef.Fields == nil || len(colRef.Fields.Items) == 0 {
		return ""
	}

	cat := c.getCatalog()
	if cat == nil {
		return ""
	}

	// Extract the parts of the column reference
	var parts []string
	for _, item := range colRef.Fields.Items {
		if s, ok := item.(*ast.String); ok {
			parts = append(parts, s.Str)
		}
	}

	if len(parts) == 0 {
		return ""
	}

	// Try to resolve: table.column or just column
	var tableName, columnName string
	if len(parts) == 2 {
		tableName = parts[0]
		columnName = parts[1]
	} else if len(parts) == 1 {
		columnName = parts[0]
		// Unqualified column - search all tables in all schemas
		return c.findColumnTypeInCatalog(columnName)
	} else {
		return ""
	}

	// Qualified table.column - look up the specific table
	tableRef := &ast.TableName{Name: tableName}
	table, err := cat.GetTable(tableRef)
	if err != nil {
		return ""
	}

	// Find the column - case-sensitive comparison
	for _, col := range table.Columns {
		if col.Name == columnName {
			return col.Type.Name
		}
	}

	return ""
}

// convert converts a ClickHouse AST node to a sqlc AST node
func (c *cc) convert(node chparser.Expr) ast.Node {
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	case *chparser.SelectQuery:
		result := c.convertSelectQuery(n)
		if debug.Active {
			if stmt, ok := result.(*ast.SelectStmt); ok && stmt != nil && stmt.TargetList != nil {
				isUnion := len(stmt.TargetList.Items) == 0 && stmt.Larg != nil
				log.Printf("[DEBUG] clickhouse.convert: SelectQuery converted, isUnion=%v, targets=%d", isUnion, len(stmt.TargetList.Items))
			}
		}
		return result
	case *chparser.InsertStmt:
		return c.convertInsertStmt(n)
	case *chparser.AlterTable:
		return c.convertAlterTable(n)
	case *chparser.CreateTable:
		return c.convertCreateTable(n)
	case *chparser.CreateDatabase:
		return c.convertCreateDatabase(n)
	case *chparser.CreateView:
		return c.convertCreateView(n)
	case *chparser.CreateMaterializedView:
		return c.convertCreateMaterializedView(n)
	case *chparser.DropStmt:
		return c.convertDropStmt(n)
	case *chparser.OptimizeStmt:
		return c.convertOptimizeStmt(n)
	case *chparser.DescribeStmt:
		return c.convertDescribeStmt(n)
	case *chparser.ExplainStmt:
		return c.convertExplainStmt(n)
	case *chparser.ShowStmt:
		return c.convertShowStmt(n)
	case *chparser.TruncateTable:
		return c.convertTruncateTable(n)

	// Expression nodes
	case *chparser.Ident:
		return c.convertIdent(n)
	case *chparser.Path:
		return c.convertPath(n)
	case *chparser.ColumnExpr:
		return c.convertColumnExpr(n)
	case *chparser.FunctionExpr:
		return c.convertFunctionExpr(n)
	case *chparser.BinaryOperation:
		return c.convertBinaryOperation(n)
	case *chparser.NumberLiteral:
		return c.convertNumberLiteral(n)
	case *chparser.StringLiteral:
		return c.convertStringLiteral(n)
	case *chparser.QueryParam:
		return c.convertQueryParam(n)
	case *chparser.NestedIdentifier:
		return c.convertNestedIdentifier(n)
	case *chparser.OrderExpr:
		return c.convertOrderExpr(n)
	case *chparser.PlaceHolder:
		return c.convertPlaceHolder(n)
	case *chparser.JoinTableExpr:
		return c.convertJoinTableExpr(n)

	// Additional expression nodes
	case *chparser.CastExpr:
		return c.convertCastExpr(n)
	case *chparser.CaseExpr:
		return c.convertCaseExpr(n)
	case *chparser.WindowFunctionExpr:
		return c.convertWindowFunctionExpr(n)
	case *chparser.IsNullExpr:
		return c.convertIsNullExpr(n)
	case *chparser.IsNotNullExpr:
		return c.convertIsNotNullExpr(n)
	case *chparser.UnaryExpr:
		return c.convertUnaryExpr(n)
	case *chparser.MapLiteral:
		return c.convertMapLiteral(n)
	case *chparser.ParamExprList:
		return c.convertParamExprList(n)
	case *chparser.IndexOperation:
		return c.convertIndexOperation(n)
	case *chparser.ArrayParamList:
		return c.convertArrayParamList(n)
	case *chparser.TableFunctionExpr:
		return c.convertTableFunctionExpr(n)
	case *chparser.TernaryOperation:
		return c.convertTernaryOperation(n)

	case *chparser.UsingClause:
		return c.convertUsingClause(n)

	default:
		// Return TODO for unsupported node types
		return todo(n)
	}
}

func (c *cc) convertSelectQuery(stmt *chparser.SelectQuery) ast.Node {
	selectStmt := &ast.SelectStmt{
		TargetList:   c.convertSelectItems(stmt.SelectItems),
		FromClause:   c.convertFromClause(stmt.From),
		WhereClause:  c.convertWhereClause(stmt.Where),
		GroupClause:  c.convertGroupByClause(stmt.GroupBy),
		HavingClause: c.convertHavingClause(stmt.Having),
		SortClause:   c.convertOrderByClause(stmt.OrderBy),
		WithClause:   c.convertWithClause(stmt.With),
	}

	// Handle ARRAY JOIN by integrating it into the FROM clause
	if stmt.ArrayJoin != nil {
		selectStmt.FromClause = c.mergeArrayJoinIntoFrom(selectStmt.FromClause, stmt.ArrayJoin)
	}

	// Handle DISTINCT
	if stmt.HasDistinct {
		selectStmt.DistinctClause = &ast.List{Items: []ast.Node{}}
	}

	// Handle LIMIT
	if stmt.Limit != nil {
		selectStmt.LimitCount = c.convertLimitClause(stmt.Limit)
		if stmt.Limit.Offset != nil {
			selectStmt.LimitOffset = c.convert(stmt.Limit.Offset)
		}
	}

	// Handle UNION/EXCEPT
	if stmt.UnionAll != nil || stmt.UnionDistinct != nil || stmt.Except != nil {
		// For UNION/EXCEPT queries, create a wrapper SelectStmt with no targets
		// The Larg points to the left SELECT, Rarg points to the right SELECT
		wrapperStmt := &ast.SelectStmt{
			TargetList: &ast.List{}, // Empty list, not nil
			FromClause: &ast.List{}, // Empty list, not nil
		}

		// Set the left SELECT (current selectStmt with all its targets and clauses)
		wrapperStmt.Larg = selectStmt

		// Determine the operation and set the right SELECT
		if stmt.UnionAll != nil {
			wrapperStmt.Op = ast.Union
			wrapperStmt.All = true
			wrapperStmt.Rarg = c.convertSelectQuery(stmt.UnionAll).(*ast.SelectStmt)
		} else if stmt.UnionDistinct != nil {
			wrapperStmt.Op = ast.Union
			wrapperStmt.All = false
			wrapperStmt.Rarg = c.convertSelectQuery(stmt.UnionDistinct).(*ast.SelectStmt)
		} else if stmt.Except != nil {
			wrapperStmt.Op = ast.Except
			wrapperStmt.Rarg = c.convertSelectQuery(stmt.Except).(*ast.SelectStmt)
		}

		return wrapperStmt
	}

	return selectStmt
}

func (c *cc) convertSelectItems(items []*chparser.SelectItem) *ast.List {
	list := &ast.List{Items: []ast.Node{}}
	for _, item := range items {
		list.Items = append(list.Items, c.convertSelectItem(item))
	}
	return list
}

func (c *cc) convertSelectItem(item *chparser.SelectItem) *ast.ResTarget {
	var name *string
	if item.Alias != nil {
		aliasName := identifier(item.Alias.Name)
		name = &aliasName
	} else {
		// If no explicit alias, try to extract a default name from the expression
		// For Path expressions like u.id, use the last part as the name
		if path, ok := item.Expr.(*chparser.Path); ok && path != nil && len(path.Fields) > 0 {
			lastName := identifier(path.Fields[len(path.Fields)-1].Name)
			name = &lastName
		}
	}

	return &ast.ResTarget{
		Name:     name,
		Val:      c.convert(item.Expr),
		Location: int(item.Pos()),
	}
}

func (c *cc) convertFromClause(from *chparser.FromClause) *ast.List {
	if from == nil {
		return &ast.List{}
	}

	list := &ast.List{Items: []ast.Node{}}

	// From.Expr can be a TableExpr, JoinExpr, or other expression
	if from.Expr != nil {
		list.Items = append(list.Items, c.convertFromExpr(from.Expr))
	}

	return list
}

func (c *cc) convertFromExpr(expr chparser.Expr) ast.Node {
	if expr == nil {
		return &ast.TODO{}
	}

	switch e := expr.(type) {
	case *chparser.TableExpr:
		return c.convertTableExpr(e)
	case *chparser.JoinTableExpr:
		// JoinTableExpr wraps a table with optional FINAL and SAMPLE clauses
		// The Table field contains the actual table reference
		if e.Table != nil {
			return c.convertTableExpr(e.Table)
		}
		return &ast.TODO{}
	case *chparser.JoinExpr:
		return c.convertJoinExpr(e)
	default:
		return c.convert(expr)
	}
}

func (c *cc) convertTableExpr(expr *chparser.TableExpr) ast.Node {
	if expr == nil {
		return &ast.TODO{}
	}

	if debug.Active {
		log.Printf("[DEBUG] convertTableExpr called with expr type: %T", expr.Expr)
	}

	// The Expr field contains the actual table reference
	var baseNode ast.Node
	var alias *string

	// Handle AliasExpr which wraps the actual table reference with an alias
	exprToProcess := expr.Expr
	if aliasExpr, ok := expr.Expr.(*chparser.AliasExpr); ok {
		// Extract the alias name
		if aliasExpr.Alias != nil {
			if aliasIdent, ok := aliasExpr.Alias.(*chparser.Ident); ok {
				aliasName := identifier(aliasIdent.Name)
				alias = &aliasName
			}
		}
		// Process the underlying expression
		exprToProcess = aliasExpr.Expr
	}

	if tableIdent, ok := exprToProcess.(*chparser.TableIdentifier); ok {
		baseNode = c.convertTableIdentifier(tableIdent)
		// Apply alias if we found one
		if alias != nil {
			if rangeVar, ok := baseNode.(*ast.RangeVar); ok {
				rangeVar.Alias = &ast.Alias{
					Aliasname: alias,
				}
			}
		}
	} else if selectQuery, ok := exprToProcess.(*chparser.SelectQuery); ok {
		// Subquery (SelectQuery)
		convertedSubquery := c.convert(selectQuery)
		if debug.Active {
			if stmt, ok := convertedSubquery.(*ast.SelectStmt); ok && stmt != nil && stmt.TargetList != nil {
				isUnion := len(stmt.TargetList.Items) == 0 && stmt.Larg != nil
				log.Printf("[DEBUG] convertTableExpr: SelectQuery converted, isUnion=%v, targets=%d", isUnion, len(stmt.TargetList.Items))
			}
		}
		rangeSubselect := &ast.RangeSubselect{
			Subquery: convertedSubquery,
		}
		if alias != nil {
			rangeSubselect.Alias = &ast.Alias{
				Aliasname: alias,
			}
		} else if expr.Alias != nil {
			if aliasIdent, ok := expr.Alias.Alias.(*chparser.Ident); ok {
				rangeSubselect.Alias = &ast.Alias{
					Aliasname: &aliasIdent.Name,
				}
			}
		} else {
			// Generate a synthetic alias for subqueries without explicit aliases
			// This is necessary for the compiler to resolve columns from the subquery
			c.subqueryCount++
			syntheticAlias := "sq_" + strconv.Itoa(c.subqueryCount)
			// IMPORTANT: Copy the string to ensure the pointer persists
			aliasCopy := syntheticAlias
			rangeSubselect.Alias = &ast.Alias{
				Aliasname: &aliasCopy,
			}
		}
		return rangeSubselect
	} else if subQuery, ok := exprToProcess.(*chparser.SubQuery); ok {
		// Subquery (SubQuery with Select field)
		if subQuery.Select == nil {
			return &ast.TODO{}
		}
		convertedSubquery := c.convert(subQuery.Select)
		if debug.Active {
			if stmt, ok := convertedSubquery.(*ast.SelectStmt); ok && stmt != nil && stmt.TargetList != nil {
				isUnion := len(stmt.TargetList.Items) == 0 && stmt.Larg != nil
				log.Printf("[DEBUG] convertTableExpr: SubQuery.Select converted, isUnion=%v, targets=%d", isUnion, len(stmt.TargetList.Items))
			}
		}
		rangeSubselect := &ast.RangeSubselect{
			Subquery: convertedSubquery,
		}
		if alias != nil {
			rangeSubselect.Alias = &ast.Alias{
				Aliasname: alias,
			}
		} else if expr.Alias != nil {
			if aliasIdent, ok := expr.Alias.Alias.(*chparser.Ident); ok {
				rangeSubselect.Alias = &ast.Alias{
					Aliasname: &aliasIdent.Name,
				}
			}
		} else {
			// Generate a synthetic alias for subqueries without explicit aliases
			// This is necessary for the compiler to resolve columns from the subquery
			c.subqueryCount++
			syntheticAlias := "sq_" + strconv.Itoa(c.subqueryCount)
			// IMPORTANT: Copy the string to ensure the pointer persists
			aliasCopy := syntheticAlias
			rangeSubselect.Alias = &ast.Alias{
				Aliasname: &aliasCopy,
			}
		}
		return rangeSubselect
	} else {
		baseNode = c.convert(exprToProcess)
	}

	return baseNode
}

func (c *cc) convertTableIdentifier(ident *chparser.TableIdentifier) *ast.RangeVar {
	var schema *string
	var table *string

	if ident.Database != nil {
		dbName := identifier(ident.Database.Name)
		schema = &dbName
	}

	if ident.Table != nil {
		tableName := identifier(ident.Table.Name)
		table = &tableName
	}

	rangeVar := &ast.RangeVar{
		Schemaname: schema,
		Relname:    table,
		Inh:        true,
		Location:   int(ident.Pos()),
	}

	return rangeVar
}

func (c *cc) convertJoinExpr(join *chparser.JoinExpr) ast.Node {
	// JoinExpr represents JOIN operations
	// Left and Right are the expressions being joined
	// Modifiers contains things like "LEFT", "RIGHT", "INNER", etc.
	// Constraints contains either an ON clause expression or a USING clause
	//
	// Note: ClickHouse's parser sometimes creates nested JoinExpr structures:
	// JoinExpr{Left: table1, Right: JoinExpr{Left: table2, Right: nil, Constraints: USING}}
	// We normalize this to match PostgreSQL's flat structure during conversion.

	// Check if Right is a nested JoinExpr with USING clause and no Right itself
	// This is a ClickHouse-specific pattern that we flatten to PostgreSQL-style
	var rarg chparser.Expr = join.Right
	var constraints chparser.Expr = join.Constraints
	
	if nestedJoin, ok := join.Right.(*chparser.JoinExpr); ok {
		// If the nested join has no Right child and has USING constraints on it,
		// we're looking at a ClickHouse nested structure that should be flattened
		if nestedJoin.Right == nil && nestedJoin.Constraints != nil {
			// Pull the table from the nested join's Left
			rarg = nestedJoin.Left
			// Pull the constraints from the nested join
			constraints = nestedJoin.Constraints
			
			// Copy modifiers from nested join if the top-level has none
			if len(join.Modifiers) == 0 && len(nestedJoin.Modifiers) > 0 {
				join.Modifiers = nestedJoin.Modifiers
			}
		}
	}

	joinNode := &ast.JoinExpr{
		Larg: c.convertFromExpr(join.Left),
		Rarg: c.convertFromExpr(rarg),
	}

	// Determine join type from modifiers
	joinType := "JOIN"
	for _, mod := range join.Modifiers {
		modUpper := strings.ToUpper(mod)
		if modUpper == "LEFT" || modUpper == "RIGHT" || modUpper == "FULL" || modUpper == "INNER" {
			joinType = modUpper + " " + joinType
		}
	}
	joinNode.Jointype = c.parseJoinType(joinType)

	// Handle constraints: either ON clause or USING clause
	if constraints != nil {
		// Check if this is a USING clause
		if usingClause, ok := constraints.(*chparser.UsingClause); ok {
			// Convert USING clause to ast.JoinExpr.UsingClause
			joinNode.UsingClause = c.convertUsingClauseToList(usingClause)
		} else {
			// Handle ON clause (regular expression)
			joinNode.Quals = c.convert(constraints)
		}
	}

	return joinNode
}

func (c *cc) parseJoinType(joinType string) ast.JoinType {
	upperType := strings.ToUpper(joinType)
	switch {
	case strings.Contains(upperType, "LEFT"):
		return ast.JoinTypeLeft
	case strings.Contains(upperType, "RIGHT"):
		return ast.JoinTypeRight
	case strings.Contains(upperType, "FULL"):
		return ast.JoinTypeFull
	case strings.Contains(upperType, "INNER"):
		return ast.JoinTypeInner
	default:
		return ast.JoinTypeInner
	}
}

// convertUsingClause converts a ClickHouse UsingClause to an ast.List of String nodes
// This creates a representation compatible with PostgreSQL-style USING clauses
func (c *cc) convertUsingClause(using *chparser.UsingClause) ast.Node {
	if using == nil || using.Using == nil {
		return nil
	}
	return c.convertUsingClauseToList(using)
}

// convertUsingClauseToList converts a ClickHouse UsingClause to an ast.List of String nodes
// representing the column names in the USING clause
func (c *cc) convertUsingClauseToList(using *chparser.UsingClause) *ast.List {
	if using == nil || using.Using == nil || len(using.Using.Items) == 0 {
		return nil
	}

	list := &ast.List{Items: []ast.Node{}}
	for _, item := range using.Using.Items {
		// Each item should be a ColumnExpr wrapping an Ident
		colExpr, ok := item.(*chparser.ColumnExpr)
		if !ok {
			continue
		}

		// Get the column name from the ColumnExpr
		if ident, ok := colExpr.Expr.(*chparser.Ident); ok {
			colName := identifier(ident.Name)
			list.Items = append(list.Items, &ast.String{Str: colName})
		}
	}

	return list
}

func (c *cc) convertWhereClause(where *chparser.WhereClause) ast.Node {
	if where == nil {
		return nil
	}
	return c.convert(where.Expr)
}

func (c *cc) convertGroupByClause(groupBy *chparser.GroupByClause) *ast.List {
	if groupBy == nil {
		return &ast.List{}
	}

	list := &ast.List{Items: []ast.Node{}}
	// GroupBy.Expr is a single expression which might be a comma-separated list
	if groupBy.Expr != nil {
		// Just convert the expression as-is
		// The parser should handle comma-separated lists internally
		list.Items = append(list.Items, c.convert(groupBy.Expr))
	}
	return list
}

func (c *cc) convertHavingClause(having *chparser.HavingClause) ast.Node {
	if having == nil {
		return nil
	}
	return c.convert(having.Expr)
}

func (c *cc) convertOrderByClause(orderBy *chparser.OrderByClause) *ast.List {
	if orderBy == nil {
		return &ast.List{}
	}

	list := &ast.List{Items: []ast.Node{}}

	// OrderBy.Items is a slice of Expr
	// For now, just convert each item directly
	for _, item := range orderBy.Items {
		list.Items = append(list.Items, c.convert(item))
	}

	return list
}

func (c *cc) convertLimitClause(limit *chparser.LimitClause) ast.Node {
	if limit == nil || limit.Limit == nil {
		return nil
	}
	return c.convert(limit.Limit)
}

func (c *cc) convertWithClause(with *chparser.WithClause) *ast.WithClause {
	if with == nil {
		return nil
	}

	list := &ast.List{Items: []ast.Node{}}
	for _, cte := range with.CTEs {
		list.Items = append(list.Items, c.convertCTE(cte))
	}

	return &ast.WithClause{
		Ctes:     list,
		Location: int(with.Pos()),
	}
}

func (c *cc) convertCTE(cte *chparser.CTEStmt) *ast.CommonTableExpr {
	if cte == nil {
		return nil
	}

	// Extract CTE name from Expr (should be an Ident)
	var cteName *string
	if ident, ok := cte.Expr.(*chparser.Ident); ok {
		name := identifier(ident.Name)
		cteName = &name
	}

	return &ast.CommonTableExpr{
		Ctename:  cteName,
		Ctequery: c.convert(cte.Alias),
		Location: int(cte.Pos()),
	}
}

func (c *cc) convertInsertStmt(stmt *chparser.InsertStmt) ast.Node {
	insert := &ast.InsertStmt{
		Relation:      c.convertTableExprToRangeVar(stmt.Table),
		Cols:          c.convertColumnNames(stmt.ColumnNames),
		ReturningList: &ast.List{},
	}

	// Handle VALUES
	if len(stmt.Values) > 0 {
		insert.SelectStmt = &ast.SelectStmt{
			FromClause:  &ast.List{},
			TargetList:  &ast.List{},
			ValuesLists: c.convertValues(stmt.Values),
		}
	}

	// Handle INSERT INTO ... SELECT
	if stmt.SelectExpr != nil {
		insert.SelectStmt = c.convert(stmt.SelectExpr)
	}

	return insert
}

func (c *cc) convertTableExprToRangeVar(expr chparser.Expr) *ast.RangeVar {
	if tableIdent, ok := expr.(*chparser.TableIdentifier); ok {
		return c.convertTableIdentifier(tableIdent)
	}
	if ident, ok := expr.(*chparser.Ident); ok {
		name := identifier(ident.Name)
		return &ast.RangeVar{
			Relname:  &name,
			Location: int(ident.Pos()),
		}
	}
	return &ast.RangeVar{}
}

func (c *cc) convertColumnNames(colNames *chparser.ColumnNamesExpr) *ast.List {
	if colNames == nil {
		return &ast.List{}
	}

	list := &ast.List{Items: []ast.Node{}}
	for _, col := range colNames.ColumnNames {
		// ColumnNames contains NestedIdentifier which has pointers
		// Convert to ResTarget with ColumnRef so the compiler can resolve types properly
		var colName string
		if col.Ident != nil {
			colName = identifier(col.Ident.Name)
		} else if col.DotIdent != nil {
			colName = identifier(col.DotIdent.Name)
		}

		if colName != "" {
			// Create a ResTarget with a ColumnRef that the compiler can resolve
			// This allows type inference to work properly for INSERT parameters
			resTarget := &ast.ResTarget{
				Name: &colName,
				Val: &ast.ColumnRef{
					Fields: &ast.List{
						Items: []ast.Node{
							&ast.String{Str: colName},
						},
					},
				},
			}
			list.Items = append(list.Items, resTarget)
		}
	}
	return list
}

func (c *cc) convertValues(values []*chparser.AssignmentValues) *ast.List {
	list := &ast.List{Items: []ast.Node{}}
	for _, valueSet := range values {
		inner := &ast.List{Items: []ast.Node{}}
		for _, val := range valueSet.Values {
			inner.Items = append(inner.Items, c.convert(val))
		}
		list.Items = append(list.Items, inner)
	}
	return list
}

func (c *cc) convertCreateTable(stmt *chparser.CreateTable) ast.Node {
	if stmt == nil {
		return &ast.TODO{}
	}

	// Extract table name
	var schema *string
	var table *string
	if stmt.Name != nil {
		if stmt.Name.Database != nil {
			dbName := identifier(stmt.Name.Database.Name)
			schema = &dbName
		}
		if stmt.Name.Table != nil {
			tableName := identifier(stmt.Name.Table.Name)
			table = &tableName
		}
	}

	// If no schema/database specified, the table name might be in Name.Table or Name.Database
	// In ClickHouse parser, a simple "users" goes into Database field, not Table
	if table == nil && stmt.Name != nil && stmt.Name.Database != nil {
		tableName := identifier(stmt.Name.Database.Name)
		table = &tableName
		schema = nil // No schema specified, will use default
	}

	// Build TableName for CreateTableStmt
	tableName := &ast.TableName{}
	if schema != nil {
		tableName.Schema = *schema
	}
	if table != nil {
		tableName.Name = *table
	}

	createStmt := &ast.CreateTableStmt{
		Name:        tableName,
		IfNotExists: stmt.IfNotExists,
	}

	// Convert columns from TableSchema
	if stmt.TableSchema != nil && len(stmt.TableSchema.Columns) > 0 {
		cols := []*ast.ColumnDef{}
		for _, col := range stmt.TableSchema.Columns {
			if colDef, ok := col.(*chparser.ColumnDef); ok {
				if converted, ok := c.convertColumnDef(colDef).(*ast.ColumnDef); ok {
					cols = append(cols, converted)
				}
			}
		}
		createStmt.Cols = cols
	}

	// Note: ClickHouse-specific features like ENGINE, ORDER BY, PARTITION BY, and SETTINGS
	// are not stored in sqlc's CreateTableStmt as it's designed for PostgreSQL compatibility.
	// These features are parsed but not preserved in the AST for now.
	// In a full ClickHouse implementation, we might extend CreateTableStmt or create
	// ClickHouse-specific statement types.

	return createStmt
}

func (c *cc) convertCreateDatabase(stmt *chparser.CreateDatabase) ast.Node {
	if stmt == nil {
		return &ast.TODO{}
	}

	var schemaName string
	if stmt.Name != nil {
		// Name is usually an Ident
		if ident, ok := stmt.Name.(*chparser.Ident); ok {
			schemaName = identifier(ident.Name)
		}
	}

	return &ast.CreateSchemaStmt{
		Name:        &schemaName,
		IfNotExists: stmt.IfNotExists,
	}
}

func (c *cc) convertDropStmt(stmt *chparser.DropStmt) ast.Node {
	if stmt == nil {
		return &ast.TODO{}
	}

	// ClickHouse DROP statements are mostly structural (DROP TABLE, DROP DATABASE)
	// sqlc doesn't have a dedicated DropStmt, so return TODO
	// This is expected - DROP is a DDL statement not typically used in application queries
	return &ast.TODO{}
}

func (c *cc) convertAlterTable(stmt *chparser.AlterTable) ast.Node {
	if stmt == nil {
		return &ast.TODO{}
	}

	// ClickHouse uses ALTER TABLE for modifications that would be UPDATE/DELETE in other DBs
	// sqlc doesn't have dedicated support for ALTER TABLE modifications
	// This is expected - ALTER TABLE is DDL, not typically used in application queries
	return &ast.TODO{}
}

func (c *cc) convertOptimizeStmt(stmt *chparser.OptimizeStmt) ast.Node {
	if stmt == nil {
		return &ast.TODO{}
	}

	// OPTIMIZE is a ClickHouse-specific statement for maintenance
	// Not a query statement that generates application code
	return &ast.TODO{}
}

func (c *cc) convertDescribeStmt(stmt *chparser.DescribeStmt) ast.Node {
	if stmt == nil {
		return &ast.TODO{}
	}

	// DESCRIBE/DESC is a metadata query - useful for introspection but not
	// typically used in application code generation workflows
	return &ast.TODO{}
}

func (c *cc) convertExplainStmt(stmt *chparser.ExplainStmt) ast.Node {
	if stmt == nil {
		return &ast.TODO{}
	}

	// EXPLAIN is for query analysis, not application code
	return &ast.TODO{}
}

func (c *cc) convertShowStmt(stmt *chparser.ShowStmt) ast.Node {
	if stmt == nil {
		return &ast.TODO{}
	}

	// SHOW is an introspection statement for metadata queries
	// While it returns result sets, it's not typically code-generated
	// Treating as TODO for now as it's not a primary use case
	return &ast.TODO{}
}

func (c *cc) convertTruncateTable(stmt *chparser.TruncateTable) ast.Node {
	if stmt == nil {
		return &ast.TODO{}
	}

	// TRUNCATE is a DDL statement for deleting all rows from a table
	// While executable, it's not typically generated as application code
	// Treating as TODO for now as it's a maintenance operation
	return &ast.TODO{}
}

func (c *cc) convertIdent(id *chparser.Ident) ast.Node {
	// Convert identifier to a ColumnRef (represents a column reference)
	// An identifier in a SELECT or WHERE clause refers to a column, not a string literal
	identName := identifier(id.Name)

	// Special case: * is represented as A_Star in sqlc AST
	if identName == "*" {
		return &ast.ColumnRef{
			Fields: &ast.List{
				Items: []ast.Node{
					&ast.A_Star{},
				},
			},
			Location: int(id.Pos()),
		}
	}

	return &ast.ColumnRef{
		Fields: &ast.List{
			Items: []ast.Node{
				&ast.String{Str: identName},
			},
		},
		Location: int(id.Pos()),
	}
}

func (c *cc) convertPath(path *chparser.Path) ast.Node {
	// Path represents a qualified identifier like "table.column" or "schema.table"
	// Convert it to a ColumnRef with multiple fields
	if path == nil || len(path.Fields) == 0 {
		return &ast.TODO{}
	}

	fields := &ast.List{Items: []ast.Node{}}
	for _, field := range path.Fields {
		if field != nil {
			fieldName := identifier(field.Name)
			if fieldName == "*" {
				fields.Items = append(fields.Items, &ast.A_Star{})
			} else {
				fields.Items = append(fields.Items, &ast.String{Str: fieldName})
			}
		}
	}

	return &ast.ColumnRef{
		Fields:   fields,
		Location: int(path.Pos()),
	}
}

func (c *cc) convertColumnExpr(col *chparser.ColumnExpr) ast.Node {
	// ColumnExpr wraps an expression (could be Ident, NestedIdentifier, etc.)
	// Just convert the underlying expression
	return c.convert(col.Expr)
}

func (c *cc) convertFunctionExpr(fn *chparser.FunctionExpr) ast.Node {
	// Convert function calls like COUNT(*), SUM(column), etc.
	// Normalize function names to lowercase since ClickHouse function names are case-insensitive
	originalFuncName := identifier(fn.Name.Name)
	funcNameLower := normalizeFunctionName(originalFuncName)

	// Handle sqlc_* functions (converted from sqlc.* during preprocessing)
	// Normalize back to sqlc.* schema.function format for proper AST representation
	var schema string
	var baseFuncName string

	if strings.HasPrefix(funcNameLower, "sqlc_") {
		schema = "sqlc"
		baseFuncName = strings.TrimPrefix(funcNameLower, "sqlc_")
	} else {
		baseFuncName = funcNameLower
	}

	args := &ast.List{Items: []ast.Node{}}
	var chArgs []*chparser.Expr // Keep original ClickHouse args for type analysis
	if fn.Params != nil {
		if fn.Params.Items != nil {
			for _, item := range fn.Params.Items.Items {
				chArgs = append(chArgs, &item)
				args.Items = append(args.Items, c.convert(item))
			}
		}
	}

	// Handle special context-dependent ClickHouse functions
	// For these functions, try to register them with the correct return type
	c.handleSpecialFunctionTypes(funcNameLower, fn, chArgs)

	return &ast.FuncCall{
		Func: &ast.FuncName{
			Schema: schema,
			Name:   baseFuncName,
		},
		Funcname: &ast.List{
			Items: []ast.Node{
				&ast.String{Str: funcNameLower},
			},
		},
		Args:     args,
		Location: int(fn.Pos()),
	}
}

// handleSpecialFunctionTypes handles ClickHouse functions with context-dependent return types
// funcName should already be normalized to lowercase via normalizeFunctionName
func (c *cc) handleSpecialFunctionTypes(funcName string, fn *chparser.FunctionExpr, chArgs []*chparser.Expr) {
	switch funcName {
	case "arrayjoin":
		// arrayJoin(Array(T)) returns T (the element type)
		if len(chArgs) > 0 {
			// Try to extract element type from the array argument
			elemType := c.extractArrayElementType(*chArgs[0])
			if elemType != "" {
				c.registerFunctionInCatalog("arrayjoin", &ast.TypeName{Name: elemType})
			}
		}

	case "argmin", "argmax":
		// argMin/argMax return the type of their first argument
		if len(chArgs) > 0 {
			// Try to extract type from the first argument (the value being tracked)
			argType := c.extractTypeFromChExpr(*chArgs[0])
			if argType != "" {
				c.registerFunctionInCatalog(funcName, &ast.TypeName{Name: argType})
			}
		}
	}
}

// extractTypeFromChExpr extracts a type from a ClickHouse expression
func (c *cc) extractTypeFromChExpr(expr chparser.Expr) string {
	if expr == nil {
		return ""
	}

	switch e := expr.(type) {
	case *chparser.ColumnExpr:
		// ColumnExpr wraps another expression - extract from the inner expression
		if e.Expr != nil {
			return c.extractTypeFromChExpr(e.Expr)
		}

	case *chparser.Path:
		// Path like u.id
		if len(e.Fields) >= 2 {
			colRef := &ast.ColumnRef{
				Fields: &ast.List{
					Items: []ast.Node{
						&ast.String{Str: identifier(e.Fields[0].Name)},
						&ast.String{Str: identifier(e.Fields[1].Name)},
					},
				},
			}
			return c.extractTypeFromColumnRef(colRef)
		} else if len(e.Fields) == 1 {
			// Single field - just the column name
			colRef := &ast.ColumnRef{
				Fields: &ast.List{
					Items: []ast.Node{
						&ast.String{Str: identifier(e.Fields[0].Name)},
					},
				},
			}
			return c.extractTypeFromColumnRef(colRef)
		}

	case *chparser.Ident:
		// Just an identifier
		colRef := &ast.ColumnRef{
			Fields: &ast.List{
				Items: []ast.Node{
					&ast.String{Str: identifier(e.Name)},
				},
			},
		}
		return c.extractTypeFromColumnRef(colRef)

	case *chparser.FunctionExpr:
		// Handle function calls like Array(String), CAST(x AS String), etc.
		return c.extractTypeFromFunctionCall(e)

	case *chparser.CastExpr:
		// Handle CAST(expr AS Type)
		if e.AsType != nil {
			// AsType can be a StringLiteral or ColumnType
			if stringLit, ok := e.AsType.(*chparser.StringLiteral); ok {
				return mapClickHouseType(strings.ToLower(stringLit.Literal))
			} else if colType, ok := e.AsType.(chparser.ColumnType); ok {
				return mapClickHouseType(strings.ToLower(colType.Type()))
			}
		}

	case *chparser.BinaryOperation:
		// For :: operator (PostgreSQL-style cast), extract from right side
		if string(e.Operation) == "::" {
			if ident, ok := e.RightExpr.(*chparser.Ident); ok {
				return mapClickHouseType(identifier(ident.Name))
			}
		}
	}

	return ""
}

// extractTypeFromFunctionCall extracts the return type from a function call expression
// Handles Array(ElementType), CAST(expr AS Type), and other function patterns
func (c *cc) extractTypeFromFunctionCall(fn *chparser.FunctionExpr) string {
	if fn == nil || fn.Name == nil {
		return ""
	}

	funcName := strings.ToLower(identifier(fn.Name.Name))

	// Handle Array(ElementType) - returns Array of that element type
	if funcName == "array" {
		if fn.Params != nil && fn.Params.Items != nil && len(fn.Params.Items.Items) > 0 {
			// Get the first parameter which should be the element type
			elemType := c.extractTypeFromChExpr(fn.Params.Items.Items[0])
			if elemType != "" {
				// Return as array type
				return elemType + "[]"
			}
		}
	}

	// Handle CAST(expr AS Type) or similar casting functions
	if strings.Contains(funcName, "cast") {
		if fn.Params != nil && fn.Params.Items != nil && len(fn.Params.Items.Items) > 0 {
			// Last parameter might be the type, but this is complex
			// Return empty to be safe
		}
	}

	return ""
}

// extractArrayElementType extracts the element type from an array type or array expression
// For Array(T), returns T. For columns of Array(T) type, returns T.
func (c *cc) extractArrayElementType(expr chparser.Expr) string {
	if expr == nil {
		return ""
	}

	// Use the general type extractor which now handles all these cases
	colType := c.extractTypeFromChExpr(expr)
	if colType != "" {
		// If it's an array type, extract the element type
		if strings.HasSuffix(colType, "[]") {
			return strings.TrimSuffix(colType, "[]")
		}
		// Otherwise return as-is (might not be an array, but the caller will handle it)
		return colType
	}

	return ""
}

func (c *cc) convertBinaryOperation(op *chparser.BinaryOperation) ast.Node {
	// Special handling for :: (type cast) operator (PostgreSQL-style ClickHouse casting)
	if string(op.Operation) == "::" {
		// Extract the type from the right side
		var typeName *ast.TypeName
		if ident, ok := op.RightExpr.(*chparser.Ident); ok {
			// The right side should be an identifier representing the type
			typeStr := identifier(ident.Name)
			mappedType := mapClickHouseType(typeStr)
			// Strip trailing [] since we'll use ArrayBounds if needed
			if strings.HasSuffix(mappedType, "[]") {
				mappedType = strings.TrimSuffix(mappedType, "[]")
			}
			typeName = &ast.TypeName{
				Name:  mappedType,
				Names: &ast.List{Items: []ast.Node{NewIdentifier(mappedType)}},
			}
		} else {
			// Fallback to text if we can't determine the type
			typeName = &ast.TypeName{
				Name:  "text",
				Names: &ast.List{Items: []ast.Node{NewIdentifier("text")}},
			}
		}

		return &ast.TypeCast{
			Arg:      c.convert(op.LeftExpr),
			TypeName: typeName,
			Location: int(op.Pos()),
		}
	}

	// Convert binary operations like =, !=, <, >, AND, OR, etc.
	return &ast.A_Expr{
		Kind: ast.A_Expr_Kind(0), // Default kind
		Name: &ast.List{
			Items: []ast.Node{
				&ast.String{Str: string(op.Operation)},
			},
		},
		Lexpr:    c.convert(op.LeftExpr),
		Rexpr:    c.convert(op.RightExpr),
		Location: int(op.Pos()),
	}
}

func (c *cc) convertNumberLiteral(num *chparser.NumberLiteral) ast.Node {
	if num == nil || num.Literal == "" {
		return &ast.A_Const{
			Val:      &ast.Integer{Ival: 0},
			Location: 0,
		}
	}

	numStr := num.Literal

	// Try to parse as integer first
	if !strings.ContainsAny(numStr, ".eE") {
		// Integer literal
		if ival, err := strconv.ParseInt(numStr, 10, 64); err == nil {
			return &ast.A_Const{
				Val:      &ast.Integer{Ival: ival},
				Location: int(num.Pos()),
			}
		}
	}

	// Try to parse as float
	if _, err := strconv.ParseFloat(numStr, 64); err == nil {
		return &ast.A_Const{
			Val:      &ast.Float{Str: numStr},
			Location: int(num.Pos()),
		}
	}

	// Fallback to integer 0 if parsing fails
	return &ast.A_Const{
		Val:      &ast.Integer{Ival: 0},
		Location: int(num.Pos()),
	}
}

func (c *cc) convertStringLiteral(str *chparser.StringLiteral) ast.Node {
	// The ClickHouse parser's StringLiteral.Pos() returns the position of the first
	// character after the opening quote. We need to adjust it to point to the opening
	// quote itself for correct location tracking in rewrite.NamedParameters, which uses
	// args[0].Pos() - 1 to find the opening paren position.
	pos := int(str.Pos())
	if pos > 0 {
		pos-- // Move from first char inside quote to the opening quote
	}
	return &ast.A_Const{
		Val: &ast.String{
			Str: str.Literal,
		},
		Location: pos,
	}
}

func (c *cc) convertQueryParam(param *chparser.QueryParam) ast.Node {
	// ClickHouse uses ? for parameters
	c.paramCount += 1
	return &ast.ParamRef{
		Number:   c.paramCount,
		Location: int(param.Pos()),
		Dollar:   false, // ClickHouse uses ? notation, not $1
	}
}

func (c *cc) convertNestedIdentifier(nested *chparser.NestedIdentifier) ast.Node {
	// NestedIdentifier represents things like "database.table" or "table.column"
	fields := &ast.List{Items: []ast.Node{}}

	if nested.Ident != nil {
		fieldName := identifier(nested.Ident.Name)
		if fieldName == "*" {
			fields.Items = append(fields.Items, &ast.A_Star{})
		} else {
			fields.Items = append(fields.Items, &ast.String{Str: fieldName})
		}
	}
	if nested.DotIdent != nil {
		fieldName := identifier(nested.DotIdent.Name)
		if fieldName == "*" {
			fields.Items = append(fields.Items, &ast.A_Star{})
		} else {
			fields.Items = append(fields.Items, &ast.String{Str: fieldName})
		}
	}

	return &ast.ColumnRef{
		Fields:   fields,
		Location: int(nested.Pos()),
	}
}

// isClickHouseTypeNullable checks if a ClickHouse column type is nullable
// In ClickHouse, columns are non-nullable by default unless wrapped in Nullable(T)
func isClickHouseTypeNullable(colType chparser.ColumnType) bool {
	if colType == nil {
		return false
	}

	// Check if it's a ComplexType (like Nullable(T), Array(T), etc.)
	if ct, ok := colType.(*chparser.ComplexType); ok {
		if ct.Name != nil {
			typeName := ct.Name.String()
			if strings.EqualFold(typeName, "Nullable") {
				return true
			}
		}
	}

	return false
}

func (c *cc) convertColumnDef(col *chparser.ColumnDef) ast.Node {
	if col == nil {
		return &ast.TODO{}
	}

	// Extract column name
	var colName string
	if col.Name != nil {
		if col.Name.Ident != nil {
			colName = identifier(col.Name.Ident.Name)
		} else if col.Name.DotIdent != nil {
			colName = identifier(col.Name.DotIdent.Name)
		}
	}

	// Convert column type
	var typeName *ast.TypeName
	if col.Type != nil {
		typeName = c.convertColumnType(col.Type)
	}

	// Extract array information from TypeName
	arrayDims := 0
	if typeName != nil && typeName.ArrayBounds != nil {
		arrayDims = len(typeName.ArrayBounds.Items)
	}

	// In ClickHouse, columns are non-nullable by default.
	// They become nullable only if explicitly wrapped in Nullable(T).
	// The Nullable wrapper is unwrapped in convertColumnType(), so we need to
	// check if the original type was Nullable to determine nullability.
	isNullable := isClickHouseTypeNullable(col.Type)

	columnDef := &ast.ColumnDef{
		Colname:   colName,
		TypeName:  typeName,
		IsNotNull: !isNullable,
		IsArray:   arrayDims > 0,
		ArrayDims: arrayDims,
	}

	return columnDef
}

func (c *cc) convertColumnType(colType chparser.ColumnType) *ast.TypeName {
	if colType == nil {
		return &ast.TypeName{
			Name:  "text",
			Names: &ast.List{Items: []ast.Node{NewIdentifier("text")}},
		}
	}

	// Extract type name - ColumnType is an interface, get the string representation
	typeName := colType.Type()

	// Handle ComplexType (e.g., LowCardinality(T), Array(T), Map(K,V), Nullable(T), etc.)
	// For LowCardinality(T), extract T directly and discard the wrapper
	// LowCardinality is a ClickHouse-specific optimization hint that doesn't affect
	// the semantic type of the data, so we unwrap it at the engine level before
	// it reaches the codegen layer. This prevents LowCardinality from leaking into
	// sqlc's type system where it would have no meaning.
	if complexType, ok := colType.(*chparser.ComplexType); ok {
		if strings.EqualFold(typeName, "LowCardinality") && len(complexType.Params) > 0 {
			innerColType := complexType.Params[0]
			return c.convertColumnType(innerColType)
		}

		// Handle Nullable(T) - unwrap and return inner type
		// Nullability is tracked via the NotNull flag in ColumnDef, not the type itself
		if strings.EqualFold(typeName, "Nullable") && len(complexType.Params) > 0 {
			innerColType := complexType.Params[0]
			return c.convertColumnType(innerColType)
		}

		// Handle Map(K, V) types
		if strings.EqualFold(typeName, "Map") && len(complexType.Params) >= 2 {
			keyColType := complexType.Params[0]
			valueColType := complexType.Params[1]

			// Get mapped type names
			keyTypeName := keyColType.Type()
			keyMappedType := mapClickHouseType(keyTypeName)

			valueTypeName := valueColType.Type()
			valueMappedType := mapClickHouseType(valueTypeName)

			// Check if the key type is valid for a Go map
			if !isValidMapKeyType(keyMappedType) {
				// If key type is not valid, fall back to map[string]interface{}
				return &ast.TypeName{
					Name:  "map[string]interface{}",
					Names: &ast.List{Items: []ast.Node{NewIdentifier("map[string]interface{}")}},
				}
			}

			// Convert database type names to valid Go type names for map syntax
			keyGoType := databaseTypeToGoType(keyMappedType)
			valueGoType := databaseTypeToGoType(valueMappedType)

			// Return map[K]V representation
			mapType := "map[" + keyGoType + "]" + valueGoType
			return &ast.TypeName{
				Name:  mapType,
				Names: &ast.List{Items: []ast.Node{NewIdentifier(mapType)}},
			}
		}
	}

	// Check if this is an array type
	lowerTypeName := strings.ToLower(typeName)
	var arrayBounds *ast.List
	if strings.HasPrefix(lowerTypeName, "array") {
		// Array types need ArrayBounds to be set for proper array dimension handling
		// Each array level adds one item to ArrayBounds
		arrayBounds = &ast.List{
			Items: []ast.Node{&ast.A_Const{}}, // One item for 1D array
		}
	}

	// Map ClickHouse types to PostgreSQL-compatible types for sqlc
	mappedType := mapClickHouseType(typeName)

	// Strip trailing [] from mappedType since we'll use ArrayBounds instead
	if strings.HasSuffix(mappedType, "[]") {
		mappedType = strings.TrimSuffix(mappedType, "[]")
	}

	return &ast.TypeName{
		Name:        mappedType,
		Names:       &ast.List{Items: []ast.Node{NewIdentifier(mappedType)}},
		ArrayBounds: arrayBounds,
	}
}

// extractArrayElementType extracts the element type from Array(ElementType)
// e.g., "array(string)" -> "string", "array(uint32)" -> "uint32"
func extractArrayElementType(chType string) string {
	chType = strings.ToLower(chType)
	// Find the content within parentheses
	start := strings.Index(chType, "(")
	end := strings.LastIndex(chType, ")")
	if start != -1 && end != -1 && end > start {
		// Extract content and handle nested types
		content := chType[start+1 : end]
		// Remove any extra whitespace and handle nested parentheses
		return strings.TrimSpace(content)
	}
	// Fallback to "string" if we can't extract the type
	return "string"
}

// databaseTypeToGoType converts database type names to valid Go syntax type names
// This is used for generating map types which need valid Go type syntax
// e.g., "double precision" -> "float64", "text" -> "string"
func databaseTypeToGoType(dbType string) string {
	// Strip any trailing [] for now and re-add later
	isArray := strings.HasSuffix(dbType, "[]")
	typeBase := dbType
	if isArray {
		typeBase = strings.TrimSuffix(dbType, "[]")
	}

	// Map database types to Go types
	goType := ""
	switch typeBase {
	// Numeric types
	case "int8":
		goType = "int8"
	case "int16":
		goType = "int16"
	case "int32":
		goType = "int32"
	case "int64":
		goType = "int64"
	case "uint8":
		goType = "uint8"
	case "uint16":
		goType = "uint16"
	case "uint32":
		goType = "uint32"
	case "uint64":
		goType = "uint64"
	case "float32", "real":
		goType = "float32"
	case "float64", "double precision":
		goType = "float64"
	case "numeric":
		goType = "string" // Decimals use string
	// String types
	case "text", "varchar", "char", "string":
		goType = "string"
	// Boolean
	case "bool", "boolean":
		goType = "bool"
	// Date/Time
	case "date", "date32", "datetime", "datetime64", "timestamp":
		goType = "time.Time"
	// UUID
	case "uuid":
		goType = "string"
	// JSON
	case "jsonb", "json":
		goType = "[]byte" // JSON types as raw bytes
	// Default
	default:
		goType = typeBase // Fall back to the original
	}

	// Re-add the array suffix if present
	if isArray {
		return "[]" + goType
	}
	return goType
}

// isValidMapKeyType checks if a Go type is valid as a map key
// In Go, valid map key types are technically arrays (comparable), but in ClickHouse
// context, we restrict map keys to practical, simple types only.
// We explicitly forbid:
// - Arrays/slices (too complex as keys)
// - Maps/nested maps (invalid in Go)
// - Unknown/complex types
func isValidMapKeyType(goType string) bool {
	// Explicitly allow common scalar types
	switch goType {
	case "int8", "int16", "int32", "int64",
		"uint8", "uint16", "uint32", "uint64",
		"float32", "float64", "real", "double precision",
		"bool", "boolean",
		"text", "string", "varchar", "char",
		"uuid",
		"numeric": // decimal types
		return true
	}

	// Date/time types are comparable
	switch goType {
	case "date", "date32", "datetime", "datetime64", "timestamp":
		return true
	}

	// Any pointer type to a scalar is valid
	if strings.HasPrefix(goType, "*") && !strings.HasPrefix(goType, "*[]") && !strings.HasPrefix(goType, "*map") {
		return true
	}

	// interface{} is technically valid (matches anything)
	if goType == "interface{}" {
		return true
	}

	// Explicitly forbid arrays, slices, and maps - too complex as map keys
	// Arrays with [] suffix
	if strings.HasSuffix(goType, "[]") {
		return false
	}

	// Maps
	if strings.HasPrefix(goType, "map[") {
		return false
	}

	// Default: assume invalid (unknown types)
	return false
}

// mapClickHouseType maps ClickHouse data types to PostgreSQL-compatible types
// that sqlc understands for Go code generation
func mapClickHouseType(chType string) string {
	chType = strings.ToLower(chType)

	switch {
	// Integer types (UInt variants - unsigned)
	case strings.HasPrefix(chType, "uint8"):
		return "uint8"
	case strings.HasPrefix(chType, "uint16"):
		return "uint16"
	case strings.HasPrefix(chType, "uint32"):
		return "uint32"
	case strings.HasPrefix(chType, "uint64"):
		return "uint64"
	// Integer types (Int variants - signed)
	case strings.HasPrefix(chType, "int8"):
		return "int8"
	case strings.HasPrefix(chType, "int16"):
		return "int16"
	case strings.HasPrefix(chType, "int32"):
		return "int32"
	case strings.HasPrefix(chType, "int64"):
		return "int64"
	case strings.HasPrefix(chType, "int128"):
		return "numeric"
	case strings.HasPrefix(chType, "int256"):
		return "numeric"

	// Float types
	case strings.HasPrefix(chType, "float32"):
		return "real"
	case strings.HasPrefix(chType, "float64"):
		return "double precision"

	// Decimal types
	case strings.HasPrefix(chType, "decimal"):
		return "numeric"

	// String types
	case chType == "string":
		return "text"
	case strings.HasPrefix(chType, "fixedstring"):
		return "varchar"

	// Date/Time types
	case chType == "date":
		return "date"
	case chType == "date32":
		return "date"
	case chType == "datetime":
		return "timestamp"
	case chType == "datetime64":
		return "timestamp"

	// Boolean
	case chType == "bool":
		return "boolean"

	// UUID
	case chType == "uuid":
		return "uuid"

	// IP address types
	case chType == "ipv4":
		return "ipv4"
	case chType == "ipv6":
		return "ipv6"

	// Array types
	case strings.HasPrefix(chType, "array"):
		// Extract element type from Array(ElementType)
		// e.g., "array(string)" -> extract "string"
		elementType := extractArrayElementType(chType)
		mappedElementType := mapClickHouseType(elementType)
		return mappedElementType + "[]"

	// JSON types
	case strings.Contains(chType, "json"):
		return "jsonb"

	// Default fallback
	default:
		return "text"
	}
}

func (c *cc) convertOrderExpr(order *chparser.OrderExpr) ast.Node {
	if order == nil {
		return &ast.TODO{}
	}

	sortBy := &ast.SortBy{
		Node:     c.convert(order.Expr),
		Location: int(order.Pos()),
	}

	// Handle sort direction
	switch order.Direction {
	case "DESC":
		sortBy.SortbyDir = ast.SortByDirDesc
	case "ASC":
		sortBy.SortbyDir = ast.SortByDirAsc
	default:
		sortBy.SortbyDir = ast.SortByDirDefault
	}

	return sortBy
}

func (c *cc) convertPlaceHolder(ph *chparser.PlaceHolder) ast.Node {
	// PlaceHolder is ClickHouse's ? parameter
	c.paramCount += 1
	return &ast.ParamRef{
		Number:   c.paramCount,
		Location: int(ph.Pos()),
		Dollar:   false, // ClickHouse uses ? notation, not $1
	}
}

func (c *cc) convertJoinTableExpr(jte *chparser.JoinTableExpr) ast.Node {
	if jte == nil || jte.Table == nil {
		return &ast.TODO{}
	}
	// JoinTableExpr is a wrapper around TableExpr with optional modifiers
	// Just extract the underlying table expression
	return c.convertTableExpr(jte.Table)
}

// convertCastExpr converts CAST expressions like CAST(column AS type)
func (c *cc) convertCastExpr(castExpr *chparser.CastExpr) ast.Node {
	if castExpr == nil {
		return &ast.TODO{}
	}

	// Convert the expression to be cast
	expr := c.convert(castExpr.Expr)

	// Convert the target type - AsType is an Expr, need to extract type information
	var typeName *ast.TypeName
	if castExpr.AsType != nil {
		// The AsType can be: ColumnType, Ident, or StringLiteral
		if stringLit, ok := castExpr.AsType.(*chparser.StringLiteral); ok {
			// CAST(x AS 'String') - extract type from string literal
			typeStr := strings.ToLower(stringLit.Literal)
			mappedType := mapClickHouseType(typeStr)
			// Strip trailing [] since we'll use ArrayBounds if needed
			if strings.HasSuffix(mappedType, "[]") {
				mappedType = strings.TrimSuffix(mappedType, "[]")
			}
			typeName = &ast.TypeName{
				Name:  mappedType,
				Names: &ast.List{Items: []ast.Node{NewIdentifier(mappedType)}},
			}
		} else if colType, ok := castExpr.AsType.(chparser.ColumnType); ok {
			// CAST(x AS ColumnType) - standard form
			typeName = c.convertColumnType(colType)
		} else if ident, ok := castExpr.AsType.(*chparser.Ident); ok {
			// Fallback: treat the identifier as a type name
			typeStr := identifier(ident.Name)
			mappedType := mapClickHouseType(typeStr)
			// Strip trailing [] since we'll use ArrayBounds if needed
			if strings.HasSuffix(mappedType, "[]") {
				mappedType = strings.TrimSuffix(mappedType, "[]")
			}
			typeName = &ast.TypeName{
				Name:  mappedType,
				Names: &ast.List{Items: []ast.Node{NewIdentifier(mappedType)}},
			}
		} else {
			// Unknown type, default to text
			typeName = &ast.TypeName{
				Name:  "text",
				Names: &ast.List{Items: []ast.Node{NewIdentifier("text")}},
			}
		}
	}

	return &ast.TypeCast{
		Arg:      expr,
		TypeName: typeName,
		Location: int(castExpr.Pos()),
	}
}

// convertCaseExpr converts CASE expressions
func (c *cc) convertCaseExpr(caseExpr *chparser.CaseExpr) ast.Node {
	if caseExpr == nil {
		return &ast.TODO{}
	}

	// Convert CASE input expression (if present)
	var arg ast.Node
	if caseExpr.Expr != nil {
		arg = c.convert(caseExpr.Expr)
	}

	// Convert WHEN clauses
	args := &ast.List{Items: []ast.Node{}}

	for _, when := range caseExpr.Whens {
		if when != nil {
			// Convert WHEN condition
			whenExpr := c.convert(when.When)
			args.Items = append(args.Items, whenExpr)

			// Convert THEN result
			thenExpr := c.convert(when.Then)
			args.Items = append(args.Items, thenExpr)
		}
	}

	// Convert ELSE clause (if present)
	var elseExpr ast.Node
	if caseExpr.Else != nil {
		elseExpr = c.convert(caseExpr.Else)
	}

	return &ast.CaseExpr{
		Arg:       arg,
		Args:      args,
		Defresult: elseExpr,
		Location:  int(caseExpr.Pos()),
	}
}

// convertWindowFunctionExpr converts window function expressions
func (c *cc) convertWindowFunctionExpr(winExpr *chparser.WindowFunctionExpr) ast.Node {
	if winExpr == nil {
		return &ast.TODO{}
	}

	// Convert the underlying function
	funcCall := c.convertFunctionExpr(winExpr.Function)

	// Convert OVER clause (OverExpr contains the window specification)
	var overClause *ast.WindowDef
	if winExpr.OverExpr != nil {
		// OverExpr might be a WindowExpr or other expression
		if winDef, ok := winExpr.OverExpr.(*chparser.WindowExpr); ok {
			overClause = c.convertWindowDef(winDef)
		}
	}

	// Wrap the function call in a window context
	if funcCall, ok := funcCall.(*ast.FuncCall); ok {
		funcCall.Over = overClause
		return funcCall
	}

	return funcCall
}

// convertWindowDef converts window definition
func (c *cc) convertWindowDef(winDef *chparser.WindowExpr) *ast.WindowDef {
	if winDef == nil {
		return nil
	}

	windowDef := &ast.WindowDef{
		Location: int(winDef.Pos()),
	}

	// Convert PARTITION BY
	if winDef.PartitionBy != nil && winDef.PartitionBy.Expr != nil {
		windowDef.PartitionClause = &ast.List{Items: []ast.Node{}}
		windowDef.PartitionClause.Items = append(windowDef.PartitionClause.Items, c.convert(winDef.PartitionBy.Expr))
	}

	// Convert ORDER BY
	if winDef.OrderBy != nil {
		windowDef.OrderClause = c.convertOrderByClause(winDef.OrderBy)
	}

	return windowDef
}

// convertIsNullExpr converts IS NULL expressions
func (c *cc) convertIsNullExpr(isNull *chparser.IsNullExpr) ast.Node {
	if isNull == nil {
		return &ast.TODO{}
	}

	return &ast.NullTest{
		Arg:          c.convert(isNull.Expr),
		Nulltesttype: ast.NullTestType(0), // IS_NULL = 0
		Location:     int(isNull.Pos()),
	}
}

// convertIsNotNullExpr converts IS NOT NULL expressions
func (c *cc) convertIsNotNullExpr(isNotNull *chparser.IsNotNullExpr) ast.Node {
	if isNotNull == nil {
		return &ast.TODO{}
	}

	return &ast.NullTest{
		Arg:          c.convert(isNotNull.Expr),
		Nulltesttype: ast.NullTestType(1), // IS_NOT_NULL = 1
		Location:     int(isNotNull.Pos()),
	}
}

// convertUnaryExpr converts unary expressions (like NOT, negation)
func (c *cc) convertUnaryExpr(unary *chparser.UnaryExpr) ast.Node {
	if unary == nil {
		return &ast.TODO{}
	}

	// Kind is a TokenKind (string)
	kindStr := string(unary.Kind)

	return &ast.A_Expr{
		Kind: ast.A_Expr_Kind(1), // AEXPR_OP_ANY or AEXPR_OP
		Name: &ast.List{
			Items: []ast.Node{
				&ast.String{Str: kindStr},
			},
		},
		Rexpr:    c.convert(unary.Expr),
		Location: int(unary.Pos()),
	}
}

// convertMapLiteral converts map/dictionary literals
func (c *cc) convertMapLiteral(mapLit *chparser.MapLiteral) ast.Node {
	if mapLit == nil {
		return &ast.TODO{}
	}

	// ClickHouse uses map literals like {'key': value, 'key2': value2}
	// Convert to a list of key-value pairs
	items := &ast.List{Items: []ast.Node{}}

	for _, kv := range mapLit.KeyValues {
		// Key is a StringLiteral value, need to convert it to a pointer
		keyLit := &kv.Key
		// Add key
		items.Items = append(items.Items, c.convert(keyLit))
		// Add value
		if kv.Value != nil {
			items.Items = append(items.Items, c.convert(kv.Value))
		}
	}

	// Return as a generic constant list (maps aren't directly supported in sqlc AST)
	return &ast.A_Const{
		Val:      items,
		Location: int(mapLit.Pos()),
	}
}

// convertParamExprList converts a parenthesized expression list to its content
// ParamExprList represents (expr1, expr2, ...) or (expr)
// We convert it by extracting and converting the items
func (c *cc) convertParamExprList(paramList *chparser.ParamExprList) ast.Node {
	if paramList == nil || paramList.Items == nil {
		return &ast.TODO{}
	}

	// If there's only one item, return that directly (unwrap the parens)
	if len(paramList.Items.Items) == 1 {
		return c.convert(paramList.Items.Items[0])
	}

	// If there are multiple items, convert them all and wrap in a list
	// This shouldn't normally happen in a WHERE clause, but handle it just in case
	items := &ast.List{Items: []ast.Node{}}
	for _, item := range paramList.Items.Items {
		if colExpr, ok := item.(*chparser.ColumnExpr); ok {
			items.Items = append(items.Items, c.convert(colExpr.Expr))
		} else {
			items.Items = append(items.Items, c.convert(item))
		}
	}
	return items
}

// mergeArrayJoinIntoFrom integrates ARRAY JOIN into the FROM clause as a special join
// ClickHouse's ARRAY JOIN is unique - it "unfolds" arrays into rows
// We represent it as a cross join with special handling
func (c *cc) mergeArrayJoinIntoFrom(fromClause *ast.List, arrayJoin *chparser.ArrayJoinClause) *ast.List {
	if fromClause == nil {
		fromClause = &ast.List{Items: []ast.Node{}}
	}

	// Convert the ARRAY JOIN expression to a join node
	arrayJoinNode := c.convertArrayJoinClause(arrayJoin)

	// Add the ARRAY JOIN to the FROM clause
	if arrayJoinNode != nil {
		fromClause.Items = append(fromClause.Items, arrayJoinNode)
	}

	return fromClause
}

// convertArrayJoinClause converts ClickHouse ARRAY JOIN to sqlc AST
// ARRAY JOIN unfolds arrays into rows - we represent it as a RangeSubselect (derived table)
// This creates a synthetic SELECT that the compiler can understand without special handling
func (c *cc) convertArrayJoinClause(arrayJoin *chparser.ArrayJoinClause) ast.Node {
	if arrayJoin == nil {
		return nil
	}

	// The Expr field contains the array expression(s) to unfold
	// It can be:
	// - A single column reference (e.g., "tags")
	// - A list of expressions with aliases (e.g., "ParsedParams AS pp" or "a AS x, b AS y")

	// Check if it's a ColumnExprList (multiple array expressions)
	if exprList, ok := arrayJoin.Expr.(*chparser.ColumnExprList); ok && len(exprList.Items) > 0 {
		// Multiple array expressions - create synthetic SELECT for each
		colnames := c.collectArrayJoinColnames(exprList.Items)
		if len(colnames) == 0 {
			return nil
		}
		return c.createArrayJoinSubquery(colnames)
	}

	// Single expression
	colnames := c.extractArrayJoinColname(arrayJoin.Expr)
	if colnames == nil {
		return nil
	}
	return c.createArrayJoinSubquery([]ast.Node{colnames})
}

// collectArrayJoinColnames extracts all column names from ARRAY JOIN expressions
// Only adds non-nil colnames to the returned list
func (c *cc) collectArrayJoinColnames(items []chparser.Expr) []ast.Node {
	var colnames []ast.Node
	for _, expr := range items {
		colname := c.extractArrayJoinColname(expr)
		// Only add non-nil colnames
		if colname != nil {
			colnames = append(colnames, colname)
		}
	}
	return colnames
}

// extractArrayJoinColname extracts the column name from an ARRAY JOIN item
// Returns a String AST node representing the column name
func (c *cc) extractArrayJoinColname(expr chparser.Expr) ast.Node {
	if expr == nil {
		return nil
	}

	// Handle ColumnExpr (most common case)
	if colExpr, ok := expr.(*chparser.ColumnExpr); ok {
		if colExpr.Alias != nil {
			// Use the explicit alias
			return &ast.String{Str: identifier(colExpr.Alias.Name)}
		}
		// Extract name from the expression itself
		if colExpr.Expr != nil {
			return c.extractNameFromExpr(colExpr.Expr)
		}
	}

	// Handle other expression types
	return c.extractNameFromExpr(expr)
}

// extractNameFromExpr extracts a name from an arbitrary expression
func (c *cc) extractNameFromExpr(expr chparser.Expr) ast.Node {
	if expr == nil {
		return nil
	}

	// Path expression (e.g., u.tags)
	if path, ok := expr.(*chparser.Path); ok && len(path.Fields) > 0 {
		lastName := path.Fields[len(path.Fields)-1].Name
		return &ast.String{Str: identifier(lastName)}
	}

	// Simple identifier
	if ident, ok := expr.(*chparser.Ident); ok {
		return &ast.String{Str: identifier(ident.Name)}
	}

	return nil
}

// createArrayJoinSubquery creates a synthetic RangeSubselect representing ARRAY JOIN output
// We create a synthetic SelectStmt with ResTargets that have the column names
// The compiler will evaluate this SelectStmt normally via outputColumns logic
func (c *cc) createArrayJoinSubquery(colnames []ast.Node) ast.Node {
	// Filter out any nil column names
	validColnames := []ast.Node{}
	for _, colname := range colnames {
		if colname != nil {
			validColnames = append(validColnames, colname)
		}
	}

	if len(validColnames) == 0 {
		return nil
	}

	// Create a synthetic SELECT statement with the column names
	// SELECT colname1, colname2, ...
	// This allows the compiler's existing outputColumns logic to extract the columns
	targetList := &ast.List{Items: []ast.Node{}}
	for _, colname := range validColnames {
		if strNode, ok := colname.(*ast.String); ok {
			// Create a ResTarget for each column name
			// Use A_Const with a String value - this doesn't require table lookups
			colName := strNode.Str
			targetList.Items = append(targetList.Items, &ast.ResTarget{
				Name: &colName,
				Val: &ast.A_Const{
					Val: &ast.String{Str: colName},
				},
			})
		}
	}

	// Create synthetic SelectStmt for this ARRAY JOIN
	// Initialize with empty Lists to avoid nil pointer dereferences
	syntheticSelect := &ast.SelectStmt{
		TargetList:   targetList,
		FromClause:   &ast.List{},
		GroupClause:  &ast.List{},
		WindowClause: &ast.List{},
		SortClause:   &ast.List{},
	}

	// Wrap in RangeSubselect (derived table)
	// The compiler will call outputColumns on this subquery to get the columns
	return &ast.RangeSubselect{
		Lateral:  false, // ARRAY JOIN is not a lateral subquery
		Subquery: syntheticSelect,
		Alias:    nil, // No need for Colnames since we have a proper Subquery
	}
}

// convertArrayJoinItemToFunc converts a single ARRAY JOIN item to a FuncCall and optional colnames
// Returns the FuncCall and a list of column name nodes (StringNodes) for the alias(es)
func (c *cc) convertArrayJoinItemToFunc(expr chparser.Expr) (*ast.FuncCall, []ast.Node) {
	if expr == nil {
		return nil, nil
	}

	var arrayExpr ast.Node
	var colnames []ast.Node

	// Handle ColumnExpr (which can have an alias) - this is what ARRAY JOIN produces
	if colExpr, ok := expr.(*chparser.ColumnExpr); ok {
		// Extract the expression and alias
		arrayExpr = c.convert(colExpr.Expr)

		if colExpr.Alias != nil {
			columnName := identifier(colExpr.Alias.Name)
			colnames = append(colnames, &ast.String{Str: columnName})
		}
	} else if selectItem, ok := expr.(*chparser.SelectItem); ok {
		// Also handle SelectItem for compatibility
		arrayExpr = c.convert(selectItem.Expr)

		if selectItem.Alias != nil {
			columnName := identifier(selectItem.Alias.Name)
			colnames = append(colnames, &ast.String{Str: columnName})
		}
	} else {
		// Direct column reference without alias
		arrayExpr = c.convert(expr)
	}

	if arrayExpr == nil {
		return nil, nil
	}

	// Create a function call representing the array unnesting
	// We use a special function name "arrayjoin" to indicate this is an ARRAY JOIN
	funcCall := &ast.FuncCall{
		Func: &ast.FuncName{
			Name: "arrayjoin",
		},
		Args: &ast.List{
			Items: []ast.Node{arrayExpr},
		},
	}

	return funcCall, colnames
}

// convertIndexOperation converts array/tuple indexing like arr[1] or tuple[2]
func (c *cc) convertIndexOperation(idxOp *chparser.IndexOperation) ast.Node {
	if idxOp == nil {
		return &ast.TODO{}
	}

	// Convert the index expression
	idx := c.convert(idxOp.Index)

	// Create an A_Indices node representing array/tuple indexing
	// IsSlice is false for single-element access like arr[1]
	// It would be true for range access like arr[1:5] (if supported)
	return &ast.A_Indices{
		IsSlice: false,
		Lidx:    idx,
		Uidx:    nil, // No upper index for single-element access
	}
}

// convertArrayParamList converts array literals like [1, 2, 3] or ['a', 'b']
func (c *cc) convertArrayParamList(arrList *chparser.ArrayParamList) ast.Node {
	if arrList == nil || arrList.Items == nil {
		return &ast.TODO{}
	}

	// Convert each item in the array
	items := &ast.List{Items: []ast.Node{}}
	for _, item := range arrList.Items.Items {
		// Each item is a ColumnExpr, extract the underlying expression
		converted := c.convert(item)
		items.Items = append(items.Items, converted)
	}

	// Return an A_ArrayExpr representing the array literal
	return &ast.A_ArrayExpr{
		Elements: items,
		Location: int(arrList.Pos()),
	}
}

// convertTableFunctionExpr converts table functions like SELECT * FROM numbers(10)
// These are ClickHouse-specific functions that return table-like results
func (c *cc) convertTableFunctionExpr(tfn *chparser.TableFunctionExpr) ast.Node {
	if tfn == nil {
		return &ast.TODO{}
	}

	// TableFunctionExpr has a Name (which is an Expr) and Args (TableArgListExpr)
	// We convert it to a RangeFunction to represent a function in FROM clause context

	// Get the function name by converting the Name expression
	// Usually it's a simple Ident, but could be more complex
	var funcName string
	if tfn.Name != nil {
		if ident, ok := tfn.Name.(*chparser.Ident); ok {
			funcName = identifier(ident.Name)
		} else {
			funcName = "table_function"
		}
	} else {
		funcName = "unknown"
	}

	// Convert arguments if present
	args := &ast.List{Items: []ast.Node{}}
	if tfn.Args != nil && tfn.Args.Args != nil {
		for _, arg := range tfn.Args.Args {
			args.Items = append(args.Items, c.convert(arg))
		}
	}

	// Create a FuncCall representing the table function
	funcCall := &ast.FuncCall{
		Func: &ast.FuncName{
			Name: funcName,
		},
		Args: args,
	}

	// Wrap in a RangeFunction to represent it in FROM clause context
	return &ast.RangeFunction{
		Functions: &ast.List{
			Items: []ast.Node{funcCall},
		},
	}
}

// convertTernaryOperation converts ternary conditional expressions
// These are similar to CASE expressions but use a different structure
func (c *cc) convertTernaryOperation(ternary *chparser.TernaryOperation) ast.Node {
	if ternary == nil {
		return &ast.TODO{}
	}

	// Convert to a CaseExpr structure for consistency with sqlc AST
	// A ternary operation is: condition ? true_expr : false_expr
	// This maps to: CASE WHEN condition THEN true_expr ELSE false_expr END

	// Convert the condition and expressions
	condition := c.convert(ternary.Condition)
	trueExpr := c.convert(ternary.TrueExpr)
	falseExpr := c.convert(ternary.FalseExpr)

	return &ast.CaseExpr{
		Arg: nil, // No CASE expr, just WHEN conditions
		Args: &ast.List{
			Items: []ast.Node{
				condition,
				trueExpr,
			},
		},
		Defresult: falseExpr, // ELSE clause
		Location:  int(ternary.Pos()),
	}
}

// convertCreateView converts CREATE VIEW statements
// ClickHouse views are similar to other SQL databases but may have specific features
func (c *cc) convertCreateView(view *chparser.CreateView) ast.Node {
	if view == nil {
		return &ast.TODO{}
	}

	// Extract view name from TableIdentifier
	var viewName string
	if view.Name != nil {
		if view.Name.Table != nil {
			viewName = identifier(view.Name.Table.Name)
		}
	}

	// Convert the SELECT query from SubQuery
	var selectStmt ast.Node
	if view.SubQuery != nil && view.SubQuery.Select != nil {
		selectStmt = c.convert(view.SubQuery.Select)
	}

	// For now, return a TODO since sqlc AST doesn't have a specific View representation
	// The SelectStmt is converted for reference
	_ = selectStmt
	_ = viewName

	return &ast.TODO{}
}

// convertCreateMaterializedView converts CREATE MATERIALIZED VIEW statements
// These are ClickHouse-specific materialized views
func (c *cc) convertCreateMaterializedView(matView *chparser.CreateMaterializedView) ast.Node {
	if matView == nil {
		return &ast.TODO{}
	}

	// Extract view name from TableIdentifier
	var viewName string
	if matView.Name != nil {
		if matView.Name.Table != nil {
			viewName = identifier(matView.Name.Table.Name)
		}
	}

	// Convert the SELECT query from SubQuery
	var selectStmt ast.Node
	if matView.SubQuery != nil && matView.SubQuery.Select != nil {
		selectStmt = c.convert(matView.SubQuery.Select)
	}

	// For now, return a TODO since sqlc AST doesn't have a specific MaterializedView representation
	// The SelectStmt is converted for reference
	_ = selectStmt
	_ = viewName

	return &ast.TODO{}
}
