// Package sqlanalyze implements SQL analysis using scope graphs for name
// resolution and bidirectional type checking for type inference.
//
// It walks the sqlc AST, building a scope graph that models the visibility
// of tables, columns, and aliases. It then uses bidirectional type checking
// to infer parameter types and validate expressions.
//
// This package works with both PostgreSQL and MySQL ASTs (via the sqlc
// unified AST representation).
package sqlanalyze

import (
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/analysis/scope"
	"github.com/sqlc-dev/sqlc/internal/analysis/typecheck"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

// Result holds the analysis results for a single query.
type Result struct {
	// The root scope of the query's scope graph.
	RootScope *scope.Scope
	// Inferred types for query parameters.
	ParamTypes map[int]*typecheck.ParamTypeInference
	// Output columns with their resolved types.
	OutputColumns []OutputColumn
	// Type errors found during analysis.
	Errors []typecheck.TypeError
}

// OutputColumn describes a column in the query's result set.
type OutputColumn struct {
	Name     string
	Type     scope.Type
	TableRef string // The table this column came from, if any
}

// Analyzer performs scope-graph-based analysis on SQL queries.
type Analyzer struct {
	catalog *catalog.Catalog
	engine  config.Engine
	checker *typecheck.Checker
}

// New creates a new analyzer for the given catalog and engine.
func New(cat *catalog.Catalog, engine config.Engine) *Analyzer {
	var rules typecheck.OperatorRules
	switch engine {
	case config.EngineMySQL:
		rules = &typecheck.MySQLOperatorRules{}
	case config.EnginePostgreSQL:
		rules = &typecheck.PostgreSQLOperatorRules{}
	default:
		rules = &typecheck.DefaultOperatorRules{}
	}

	return &Analyzer{
		catalog: cat,
		engine:  engine,
		checker: typecheck.NewChecker(rules),
	}
}

// AnalyzeQuery performs full analysis on a SQL statement: builds the scope
// graph, resolves names, infers types bidirectionally, and returns results.
func (a *Analyzer) AnalyzeQuery(raw *ast.RawStmt) (*Result, error) {
	if raw == nil || raw.Stmt == nil {
		return nil, fmt.Errorf("nil statement")
	}

	result := &Result{
		ParamTypes: make(map[int]*typecheck.ParamTypeInference),
	}

	// Build the scope graph from the AST
	rootScope, err := a.buildScopeGraph(raw.Stmt)
	if err != nil {
		return nil, fmt.Errorf("building scope graph: %w", err)
	}
	result.RootScope = rootScope

	// Walk expressions to perform bidirectional type checking
	a.typeCheckStatement(raw.Stmt, rootScope)

	// Collect results
	result.ParamTypes = a.checker.ParamTypes()
	result.Errors = a.checker.Errors()

	// Compute output columns
	outputCols, err := a.computeOutputColumns(raw.Stmt, rootScope)
	if err != nil {
		return nil, fmt.Errorf("computing output columns: %w", err)
	}
	result.OutputColumns = outputCols

	return result, nil
}

// buildScopeGraph constructs the scope graph for a SQL statement.
func (a *Analyzer) buildScopeGraph(stmt ast.Node) (*scope.Scope, error) {
	switch n := stmt.(type) {
	case *ast.SelectStmt:
		return a.buildSelectScope(n)
	case *ast.InsertStmt:
		return a.buildInsertScope(n)
	case *ast.UpdateStmt:
		return a.buildUpdateScope(n)
	case *ast.DeleteStmt:
		return a.buildDeleteScope(n)
	default:
		return scope.NewScope(scope.ScopeRoot), nil
	}
}

// buildSelectScope builds the scope graph for a SELECT statement.
//
// The scope structure for SELECT is:
//
//	[CTE scope] (if WITH clause exists)
//	    |
//	[FROM scope] ← contains table declarations + aliases
//	    |
//	[WHERE scope] → PARENT → [FROM scope]
//	    |
//	[SELECT scope] → PARENT → [FROM scope]
func (a *Analyzer) buildSelectScope(sel *ast.SelectStmt) (*scope.Scope, error) {
	if sel == nil {
		return scope.NewScope(scope.ScopeRoot), nil
	}

	// Handle UNION queries
	if sel.Larg != nil {
		return a.buildSelectScope(sel.Larg)
	}

	// Build FROM scope
	fromScope := scope.NewScope(scope.ScopeFrom)

	// Process CTEs first (WITH clause)
	if sel.WithClause != nil {
		if err := a.processCTEs(sel.WithClause, fromScope); err != nil {
			return nil, err
		}
	}

	// Process FROM clause
	if sel.FromClause != nil {
		for _, item := range sel.FromClause.Items {
			if err := a.processFromItem(item, fromScope); err != nil {
				return nil, err
			}
		}
	}

	// The SELECT scope has the FROM scope as parent
	selectScope := scope.NewScope(scope.ScopeSelect)
	selectScope.AddParent(fromScope)

	return selectScope, nil
}

// buildInsertScope builds the scope graph for an INSERT statement.
func (a *Analyzer) buildInsertScope(ins *ast.InsertStmt) (*scope.Scope, error) {
	insertScope := scope.NewScope(scope.ScopeInsert)

	if ins.WithClause != nil {
		if err := a.processCTEs(ins.WithClause, insertScope); err != nil {
			return nil, err
		}
	}

	// Add the target table
	if ins.Relation != nil {
		if err := a.processFromItem(ins.Relation, insertScope); err != nil {
			return nil, err
		}
	}

	return insertScope, nil
}

// buildUpdateScope builds the scope graph for an UPDATE statement.
func (a *Analyzer) buildUpdateScope(upd *ast.UpdateStmt) (*scope.Scope, error) {
	updateScope := scope.NewScope(scope.ScopeUpdate)

	if upd.WithClause != nil {
		if err := a.processCTEs(upd.WithClause, updateScope); err != nil {
			return nil, err
		}
	}

	// Add tables from Relations
	if upd.Relations != nil {
		for _, item := range upd.Relations.Items {
			if err := a.processFromItem(item, updateScope); err != nil {
				return nil, err
			}
		}
	}

	// Add tables from FROM clause
	if upd.FromClause != nil {
		for _, item := range upd.FromClause.Items {
			if err := a.processFromItem(item, updateScope); err != nil {
				return nil, err
			}
		}
	}

	return updateScope, nil
}

// buildDeleteScope builds the scope graph for a DELETE statement.
func (a *Analyzer) buildDeleteScope(del *ast.DeleteStmt) (*scope.Scope, error) {
	deleteScope := scope.NewScope(scope.ScopeDelete)

	if del.WithClause != nil {
		if err := a.processCTEs(del.WithClause, deleteScope); err != nil {
			return nil, err
		}
	}

	if del.Relations != nil {
		for _, item := range del.Relations.Items {
			if err := a.processFromItem(item, deleteScope); err != nil {
				return nil, err
			}
		}
	}

	return deleteScope, nil
}

// processCTEs adds CTE declarations to the given scope.
func (a *Analyzer) processCTEs(with *ast.WithClause, parentScope *scope.Scope) error {
	if with == nil || with.Ctes == nil {
		return nil
	}

	for _, item := range with.Ctes.Items {
		cte, ok := item.(*ast.CommonTableExpr)
		if !ok || cte.Ctename == nil {
			continue
		}

		// Build the CTE's scope by analyzing its query
		cteQueryScope, err := a.buildScopeGraph(cte.Ctequery)
		if err != nil {
			continue // Don't fail on CTE analysis errors
		}

		// If the CTE has explicit column names, use those
		cteScope := scope.NewScope(scope.ScopeCTE)
		if cte.Aliascolnames != nil {
			for _, nameNode := range cte.Aliascolnames.Items {
				if s, ok := nameNode.(*ast.String); ok {
					cteScope.DeclareColumn(s.Str, scope.TypeUnknown, 0)
				}
			}
		} else {
			// Copy columns from the CTE's query scope
			cols := cteQueryScope.AllColumns("")
			for _, col := range cols {
				cteScope.DeclareColumn(col.Name, col.Type, col.Location)
			}
		}

		parentScope.Declare(&scope.Declaration{
			Name:  *cte.Ctename,
			Kind:  scope.DeclCTE,
			Type:  scope.TypeUnknown,
			Scope: cteScope,
		})
	}

	return nil
}

// processFromItem adds a FROM clause item (table, join, subquery) to the scope.
func (a *Analyzer) processFromItem(item ast.Node, parentScope *scope.Scope) error {
	switch n := item.(type) {
	case *ast.RangeVar:
		return a.processRangeVar(n, parentScope)

	case *ast.JoinExpr:
		return a.processJoinExpr(n, parentScope)

	case *ast.RangeSubselect:
		return a.processRangeSubselect(n, parentScope)

	case *ast.RangeFunction:
		// Function in FROM clause — add placeholder columns
		return a.processRangeFunction(n, parentScope)

	default:
		return nil // Ignore unknown FROM item types
	}
}

// processRangeVar looks up a table in the catalog and adds it to the scope.
func (a *Analyzer) processRangeVar(rv *ast.RangeVar, parentScope *scope.Scope) error {
	if rv == nil || rv.Relname == nil {
		return nil
	}

	tableName := &ast.TableName{Name: *rv.Relname}
	if rv.Schemaname != nil {
		tableName.Schema = *rv.Schemaname
	}

	// Create a scope for this table's columns
	tableScope := scope.NewScope(scope.ScopeFrom)

	// Look up the table in the catalog
	table, err := a.catalog.GetTable(tableName)
	if err != nil {
		// Table might be a CTE — check if it's already declared in parent
		_, resolveErr := parentScope.ResolveQualified(*rv.Relname, "")
		if resolveErr != nil {
			// Not found anywhere — declare empty scope so analysis can continue
			tableScope.DeclareColumn("*", scope.TypeUnknown, 0)
		}
	} else {
		// Add all columns from the catalog
		for _, col := range table.Columns {
			typ := scope.Type{
				Name:      col.Type.Name,
				Schema:    col.Type.Schema,
				NotNull:   col.IsNotNull,
				IsArray:   col.IsArray,
				ArrayDims: col.ArrayDims,
				Unsigned:  col.IsUnsigned,
				Length:    col.Length,
			}
			tableScope.DeclareColumn(col.Name, typ, 0)
		}
	}

	// Determine the name to use (alias or table name)
	name := *rv.Relname
	if rv.Alias != nil && rv.Alias.Aliasname != nil {
		alias := *rv.Alias.Aliasname
		// Register both the alias edge and the table declaration
		parentScope.AddAlias(alias, tableScope)
		parentScope.Declare(&scope.Declaration{
			Name:  alias,
			Kind:  scope.DeclAlias,
			Type:  scope.TypeUnknown,
			Scope: tableScope,
		})
	} else {
		parentScope.Declare(&scope.Declaration{
			Name:  name,
			Kind:  scope.DeclTable,
			Type:  scope.TypeUnknown,
			Scope: tableScope,
		})
	}

	return nil
}

// processJoinExpr processes a JOIN and adds both sides to the scope.
func (a *Analyzer) processJoinExpr(join *ast.JoinExpr, parentScope *scope.Scope) error {
	if join == nil {
		return nil
	}

	// Process left side
	if join.Larg != nil {
		if err := a.processFromItem(join.Larg, parentScope); err != nil {
			return err
		}
	}

	// Process right side
	if join.Rarg != nil {
		if err := a.processFromItem(join.Rarg, parentScope); err != nil {
			return err
		}
	}

	return nil
}

// processRangeSubselect processes a subquery in the FROM clause.
func (a *Analyzer) processRangeSubselect(rs *ast.RangeSubselect, parentScope *scope.Scope) error {
	if rs == nil {
		return nil
	}

	subScope := scope.NewScope(scope.ScopeSubquery)

	// Analyze the subquery
	if rs.Subquery != nil {
		subQueryScope, err := a.buildScopeGraph(rs.Subquery)
		if err == nil {
			cols := subQueryScope.AllColumns("")
			for _, col := range cols {
				subScope.DeclareColumn(col.Name, col.Type, col.Location)
			}
		}
	}

	if rs.Alias != nil && rs.Alias.Aliasname != nil {
		alias := *rs.Alias.Aliasname
		parentScope.AddAlias(alias, subScope)
		parentScope.Declare(&scope.Declaration{
			Name:  alias,
			Kind:  scope.DeclAlias,
			Type:  scope.TypeUnknown,
			Scope: subScope,
		})
	}

	return nil
}

// processRangeFunction processes a function call in the FROM clause.
func (a *Analyzer) processRangeFunction(rf *ast.RangeFunction, parentScope *scope.Scope) error {
	if rf == nil {
		return nil
	}

	funcScope := scope.NewScope(scope.ScopeFunction)

	// If there's an alias with column definitions, use those
	if rf.Alias != nil {
		if rf.Alias.Colnames != nil {
			for _, nameNode := range rf.Alias.Colnames.Items {
				if s, ok := nameNode.(*ast.String); ok {
					funcScope.DeclareColumn(s.Str, scope.TypeUnknown, 0)
				}
			}
		}
		if rf.Alias.Aliasname != nil {
			alias := *rf.Alias.Aliasname
			parentScope.AddAlias(alias, funcScope)
			parentScope.Declare(&scope.Declaration{
				Name:  alias,
				Kind:  scope.DeclAlias,
				Type:  scope.TypeUnknown,
				Scope: funcScope,
			})
		}
	}

	return nil
}

// typeCheckStatement walks the AST and performs bidirectional type checking.
func (a *Analyzer) typeCheckStatement(stmt ast.Node, rootScope *scope.Scope) {
	switch n := stmt.(type) {
	case *ast.SelectStmt:
		a.typeCheckSelect(n, rootScope)
	case *ast.InsertStmt:
		a.typeCheckInsert(n, rootScope)
	case *ast.UpdateStmt:
		a.typeCheckUpdate(n, rootScope)
	case *ast.DeleteStmt:
		a.typeCheckDelete(n, rootScope)
	}
}

// typeCheckSelect type-checks a SELECT statement's expressions.
func (a *Analyzer) typeCheckSelect(sel *ast.SelectStmt, selectScope *scope.Scope) {
	if sel == nil {
		return
	}

	// Handle UNION
	if sel.Larg != nil {
		lScope, _ := a.buildSelectScope(sel.Larg)
		if lScope != nil {
			a.typeCheckSelect(sel.Larg, lScope)
		}
		return
	}

	// Type-check WHERE clause
	if sel.WhereClause != nil {
		a.typeCheckExpr(sel.WhereClause, selectScope)
	}

	// Type-check HAVING clause
	if sel.HavingClause != nil {
		a.typeCheckExpr(sel.HavingClause, selectScope)
	}

	// Type-check LIMIT/OFFSET (they should be integer)
	if sel.LimitCount != nil {
		expr := a.astToExpr(sel.LimitCount, selectScope)
		if expr != nil {
			a.checker.Check(expr, scope.TypeInt, 0)
		}
	}
	if sel.LimitOffset != nil {
		expr := a.astToExpr(sel.LimitOffset, selectScope)
		if expr != nil {
			a.checker.Check(expr, scope.TypeInt, 0)
		}
	}
}

// typeCheckInsert type-checks an INSERT statement.
func (a *Analyzer) typeCheckInsert(ins *ast.InsertStmt, insertScope *scope.Scope) {
	if ins == nil {
		return
	}

	// For INSERT ... VALUES, check parameter types against column types
	if ins.SelectStmt != nil {
		if valSel, ok := ins.SelectStmt.(*ast.SelectStmt); ok && valSel.ValuesLists != nil {
			a.typeCheckInsertValues(ins, valSel, insertScope)
		}
	}
}

// typeCheckInsertValues infers parameter types in INSERT VALUES from column types.
func (a *Analyzer) typeCheckInsertValues(ins *ast.InsertStmt, valSel *ast.SelectStmt, insertScope *scope.Scope) {
	if ins.Relation == nil || ins.Relation.Relname == nil {
		return
	}

	tableName := &ast.TableName{Name: *ins.Relation.Relname}
	if ins.Relation.Schemaname != nil {
		tableName.Schema = *ins.Relation.Schemaname
	}

	table, err := a.catalog.GetTable(tableName)
	if err != nil {
		return
	}

	// Get the column names from the INSERT column list
	var targetCols []string
	if ins.Cols != nil {
		for _, item := range ins.Cols.Items {
			if rt, ok := item.(*ast.ResTarget); ok && rt.Name != nil {
				targetCols = append(targetCols, *rt.Name)
			}
		}
	} else {
		// No explicit columns — use all table columns in order
		for _, col := range table.Columns {
			targetCols = append(targetCols, col.Name)
		}
	}

	// Build a map of column name -> type
	colTypes := make(map[string]scope.Type)
	for _, col := range table.Columns {
		colTypes[col.Name] = scope.Type{
			Name:      col.Type.Name,
			Schema:    col.Type.Schema,
			NotNull:   col.IsNotNull,
			IsArray:   col.IsArray,
			ArrayDims: col.ArrayDims,
			Unsigned:  col.IsUnsigned,
			Length:    col.Length,
		}
	}

	// Type-check each value against its target column
	for _, row := range valSel.ValuesLists.Items {
		rowList, ok := row.(*ast.List)
		if !ok {
			continue
		}
		for i, val := range rowList.Items {
			if i >= len(targetCols) {
				break
			}
			colType, exists := colTypes[targetCols[i]]
			if !exists {
				continue
			}
			expr := a.astToExpr(val, insertScope)
			if expr != nil {
				// Use checking mode: the parameter should have the column's type
				a.checker.Check(expr, colType, 0)
			}
		}
	}
}

// typeCheckUpdate type-checks an UPDATE statement.
func (a *Analyzer) typeCheckUpdate(upd *ast.UpdateStmt, updateScope *scope.Scope) {
	if upd == nil {
		return
	}

	// Type-check SET clause values against their target columns
	if upd.TargetList != nil {
		for _, item := range upd.TargetList.Items {
			rt, ok := item.(*ast.ResTarget)
			if !ok || rt.Name == nil || rt.Val == nil {
				continue
			}

			// Look up the column type from the scope
			resolved, err := updateScope.Resolve(*rt.Name)
			if err != nil {
				continue
			}

			// Check the value against the column's type
			expr := a.astToExpr(rt.Val, updateScope)
			if expr != nil {
				a.checker.Check(expr, resolved.Declaration.Type, rt.Location)
			}
		}
	}

	// Type-check WHERE clause
	if upd.WhereClause != nil {
		a.typeCheckExpr(upd.WhereClause, updateScope)
	}
}

// typeCheckDelete type-checks a DELETE statement.
func (a *Analyzer) typeCheckDelete(del *ast.DeleteStmt, deleteScope *scope.Scope) {
	if del == nil {
		return
	}
	if del.WhereClause != nil {
		a.typeCheckExpr(del.WhereClause, deleteScope)
	}
}

// typeCheckExpr walks an AST expression, synthesizing types and
// using checking mode where appropriate.
func (a *Analyzer) typeCheckExpr(node ast.Node, sc *scope.Scope) {
	expr := a.astToExpr(node, sc)
	if expr != nil {
		a.checker.Synth(expr)
	}
}

// astToExpr converts a sqlc AST node to a type-checkable expression.
// This is the bridge between the sqlc AST and the type checker's expression language.
func (a *Analyzer) astToExpr(node ast.Node, sc *scope.Scope) typecheck.Expr {
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	case *ast.A_Const:
		return a.constToExpr(n)

	case *ast.ColumnRef:
		return a.columnRefToExpr(n, sc)

	case *ast.ParamRef:
		return &typecheck.ParamExpr{
			Number:   n.Number,
			Location: n.Location,
		}

	case *ast.A_Expr:
		return a.aExprToExpr(n, sc)

	case *ast.BoolExpr:
		return a.boolExprToExpr(n, sc)

	case *ast.FuncCall:
		return a.funcCallToExpr(n, sc)

	case *ast.TypeCast:
		return a.typeCastToExpr(n, sc)

	case *ast.SubLink:
		return a.subLinkToExpr(n)

	case *ast.NullTest:
		arg := a.astToExpr(n.Arg, sc)
		if arg == nil {
			arg = &typecheck.LiteralExpr{Type: scope.TypeUnknown}
		}
		return &typecheck.NullTestExpr{
			Arg:      arg,
			IsNot:    n.Nulltesttype == ast.NullTestTypeIsNotNull,
			Location: n.Location,
		}

	case *ast.CaseExpr:
		resultType := scope.TypeUnknown
		if n.Defresult != nil {
			if tc, ok := n.Defresult.(*ast.TypeCast); ok && tc.TypeName != nil {
				resultType = typeNameToScopeType(tc.TypeName)
			}
		}
		return &typecheck.CaseExpr{ResultType: resultType, Location: 0}

	case *ast.CoalesceExpr:
		var args []typecheck.Expr
		if n.Args != nil {
			for _, item := range n.Args.Items {
				if e := a.astToExpr(item, sc); e != nil {
					args = append(args, e)
				}
			}
		}
		return &typecheck.CoalesceExpr{Args: args}

	case *ast.In:
		return a.inToExpr(n, sc)

	case *ast.BetweenExpr:
		return a.betweenToExpr(n, sc)

	case *ast.List:
		// For lists (e.g., value lists), just check each item
		for _, item := range n.Items {
			a.astToExpr(item, sc) // Side-effecting: records param types
		}
		return nil

	default:
		return nil
	}
}

func (a *Analyzer) constToExpr(n *ast.A_Const) typecheck.Expr {
	if n == nil {
		return nil
	}
	switch n.Val.(type) {
	case *ast.String:
		return &typecheck.LiteralExpr{Type: scope.TypeText}
	case *ast.Integer:
		return &typecheck.LiteralExpr{Type: scope.TypeInt}
	case *ast.Float:
		return &typecheck.LiteralExpr{Type: scope.TypeFloat}
	case *ast.Boolean:
		return &typecheck.LiteralExpr{Type: scope.TypeBool}
	case *ast.Null:
		return &typecheck.LiteralExpr{Type: scope.Type{Name: "any"}}
	default:
		return &typecheck.LiteralExpr{Type: scope.TypeUnknown}
	}
}

func (a *Analyzer) columnRefToExpr(n *ast.ColumnRef, sc *scope.Scope) typecheck.Expr {
	if n == nil {
		return nil
	}

	parts := stringSlice(n.Fields)
	resolved, err := sc.ResolveColumnRef(parts)
	if err != nil {
		return &typecheck.ColumnRefExpr{
			Parts:        parts,
			ResolvedType: scope.TypeUnknown,
			Location:     n.Location,
		}
	}

	return &typecheck.ColumnRefExpr{
		Parts:        parts,
		ResolvedType: resolved.Declaration.Type,
		Location:     n.Location,
	}
}

func (a *Analyzer) aExprToExpr(n *ast.A_Expr, sc *scope.Scope) typecheck.Expr {
	if n == nil {
		return nil
	}

	op := astutils.Join(n.Name, ".")

	left := a.astToExpr(n.Lexpr, sc)
	right := a.astToExpr(n.Rexpr, sc)

	if left == nil {
		left = &typecheck.LiteralExpr{Type: scope.TypeUnknown}
	}
	if right == nil {
		right = &typecheck.LiteralExpr{Type: scope.TypeUnknown}
	}

	return &typecheck.BinaryOpExpr{
		Op:       op,
		Left:     left,
		Right:    right,
		Location: n.Location,
	}
}

func (a *Analyzer) boolExprToExpr(n *ast.BoolExpr, sc *scope.Scope) typecheck.Expr {
	if n == nil {
		return nil
	}

	var args []typecheck.Expr
	if n.Args != nil {
		for _, item := range n.Args.Items {
			if e := a.astToExpr(item, sc); e != nil {
				args = append(args, e)
			}
		}
	}

	var op string
	switch n.Boolop {
	case ast.BoolExprTypeAnd:
		op = "AND"
	case ast.BoolExprTypeOr:
		op = "OR"
	case ast.BoolExprTypeNot:
		op = "NOT"
	default:
		// IS NULL, IS NOT NULL checks
		op = "IS"
	}

	return &typecheck.BoolExpr{
		Op:       op,
		Args:     args,
		Location: n.Location,
	}
}

func (a *Analyzer) funcCallToExpr(n *ast.FuncCall, sc *scope.Scope) typecheck.Expr {
	if n == nil {
		return nil
	}

	// Type-check arguments (side effect: infers param types within args)
	var args []typecheck.Expr
	if n.Args != nil {
		for _, item := range n.Args.Items {
			if e := a.astToExpr(item, sc); e != nil {
				args = append(args, e)
			}
		}
	}

	// Try to resolve the function's return type from the catalog
	returnType := scope.TypeUnknown
	fun, err := a.catalog.ResolveFuncCall(n)
	if err == nil && fun.ReturnType != nil {
		returnType = typeNameToScopeType(fun.ReturnType)

		// Use checking mode on arguments against function parameter types
		for i, arg := range args {
			if i < len(fun.Args) && fun.Args[i].Type != nil {
				expectedType := typeNameToScopeType(fun.Args[i].Type)
				a.checker.Check(arg, expectedType, n.Location)
			}
		}
	}

	return &typecheck.FuncCallExpr{
		Name:       n.Func.Name,
		Args:       args,
		ReturnType: returnType,
		Location:   n.Location,
	}
}

func (a *Analyzer) typeCastToExpr(n *ast.TypeCast, sc *scope.Scope) typecheck.Expr {
	if n == nil {
		return nil
	}

	arg := a.astToExpr(n.Arg, sc)
	if arg == nil {
		arg = &typecheck.LiteralExpr{Type: scope.TypeUnknown}
	}

	castType := scope.TypeUnknown
	if n.TypeName != nil {
		castType = typeNameToScopeType(n.TypeName)
	}

	// If the argument is a parameter, infer its type from the cast
	if param, ok := arg.(*typecheck.ParamExpr); ok {
		a.checker.InferParamFromContext(param.Number, castType, param.Location)
	}

	return &typecheck.TypeCastExpr{
		Arg:      arg,
		CastType: castType,
		Location: 0,
	}
}

func (a *Analyzer) subLinkToExpr(n *ast.SubLink) typecheck.Expr {
	if n == nil {
		return nil
	}

	switch n.SubLinkType {
	case ast.EXISTS_SUBLINK:
		return &typecheck.SubqueryExpr{IsExists: true}
	default:
		return &typecheck.SubqueryExpr{Columns: []scope.Type{scope.TypeUnknown}}
	}
}

func (a *Analyzer) inToExpr(n *ast.In, sc *scope.Scope) typecheck.Expr {
	if n == nil {
		return nil
	}

	exprNode := a.astToExpr(n.Expr, sc)
	if exprNode == nil {
		exprNode = &typecheck.LiteralExpr{Type: scope.TypeUnknown}
	}

	var values []typecheck.Expr
	for _, item := range n.List {
		if e := a.astToExpr(item, sc); e != nil {
			values = append(values, e)
		}
	}

	// For IN expressions, use checking mode on list items:
	// each item should match the type of the expression
	exprType := a.checker.Synth(exprNode)
	if !exprType.Type.IsUnknown() {
		for _, v := range values {
			a.checker.Check(v, exprType.Type, 0)
		}
	}

	return &typecheck.InExpr{
		Expr:   exprNode,
		Values: values,
	}
}

func (a *Analyzer) betweenToExpr(n *ast.BetweenExpr, sc *scope.Scope) typecheck.Expr {
	if n == nil {
		return nil
	}

	exprNode := a.astToExpr(n.Expr, sc)
	low := a.astToExpr(n.Left, sc)
	high := a.astToExpr(n.Right, sc)

	if exprNode == nil {
		exprNode = &typecheck.LiteralExpr{Type: scope.TypeUnknown}
	}
	if low == nil {
		low = &typecheck.LiteralExpr{Type: scope.TypeUnknown}
	}
	if high == nil {
		high = &typecheck.LiteralExpr{Type: scope.TypeUnknown}
	}

	// Between: low and high should match the type of the expression
	exprType := a.checker.Synth(exprNode)
	if !exprType.Type.IsUnknown() {
		a.checker.Check(low, exprType.Type, 0)
		a.checker.Check(high, exprType.Type, 0)
	}

	return &typecheck.BetweenExpr{
		Expr: exprNode,
		Low:  low,
		High: high,
	}
}

// computeOutputColumns determines the columns in the query's result set.
func (a *Analyzer) computeOutputColumns(stmt ast.Node, sc *scope.Scope) ([]OutputColumn, error) {
	var targets *ast.List

	switch n := stmt.(type) {
	case *ast.SelectStmt:
		if n.Larg != nil {
			return a.computeOutputColumns(n.Larg, sc)
		}
		targets = n.TargetList
	case *ast.InsertStmt:
		targets = n.ReturningList
	case *ast.UpdateStmt:
		targets = n.ReturningList
	case *ast.DeleteStmt:
		targets = n.ReturningList
	}

	if targets == nil {
		return nil, nil
	}

	var cols []OutputColumn
	for _, item := range targets.Items {
		res, ok := item.(*ast.ResTarget)
		if !ok {
			continue
		}

		switch val := res.Val.(type) {
		case *ast.ColumnRef:
			parts := stringSlice(val.Fields)

			// Handle SELECT *
			for _, field := range val.Fields.Items {
				if _, isStar := field.(*ast.A_Star); isStar {
					qualifier := ""
					if len(parts) > 0 && parts[0] != "*" {
						qualifier = parts[0]
					}
					allCols := sc.AllColumns(qualifier)
					for _, d := range allCols {
						cols = append(cols, OutputColumn{
							Name: d.Name,
							Type: d.Type,
						})
					}
					goto nextTarget
				}
			}

			// Regular column reference
			resolved, err := sc.ResolveColumnRef(parts)
			if err != nil {
				name := parts[len(parts)-1]
				if res.Name != nil {
					name = *res.Name
				}
				cols = append(cols, OutputColumn{
					Name: name,
					Type: scope.TypeUnknown,
				})
			} else {
				name := resolved.Declaration.Name
				if res.Name != nil {
					name = *res.Name
				}
				cols = append(cols, OutputColumn{
					Name: name,
					Type: resolved.Declaration.Type,
				})
			}

		case *ast.FuncCall:
			name := val.Func.Name
			if res.Name != nil {
				name = *res.Name
			}
			funcExpr := a.funcCallToExpr(val, sc)
			result := a.checker.Synth(funcExpr)
			cols = append(cols, OutputColumn{
				Name: name,
				Type: result.Type,
			})

		case *ast.A_Const:
			name := ""
			if res.Name != nil {
				name = *res.Name
			}
			constExpr := a.constToExpr(val)
			result := a.checker.Synth(constExpr)
			cols = append(cols, OutputColumn{
				Name: name,
				Type: result.Type,
			})

		default:
			name := ""
			if res.Name != nil {
				name = *res.Name
			}
			if val != nil {
				expr := a.astToExpr(val, sc)
				if expr != nil {
					result := a.checker.Synth(expr)
					cols = append(cols, OutputColumn{
						Name: name,
						Type: result.Type,
					})
					continue
				}
			}
			cols = append(cols, OutputColumn{
				Name: name,
				Type: scope.TypeUnknown,
			})
		}
	nextTarget:
	}

	return cols, nil
}

// Helper functions

func stringSlice(list *ast.List) []string {
	if list == nil {
		return nil
	}
	var result []string
	for _, item := range list.Items {
		switch n := item.(type) {
		case *ast.String:
			result = append(result, n.Str)
		case *ast.A_Star:
			// Don't include star in string slice
		}
	}
	return result
}

func typeNameToScopeType(tn *ast.TypeName) scope.Type {
	if tn == nil {
		return scope.TypeUnknown
	}
	return scope.Type{
		Name:   tn.Name,
		Schema: tn.Schema,
	}
}
