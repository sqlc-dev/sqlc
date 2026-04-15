package expander

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"
	"github.com/sqlc-dev/sqlc/internal/sql/format"
)

// Parser is an interface for SQL parsers that can parse SQL into AST statements.
type Parser interface {
	Parse(r io.Reader) ([]ast.Statement, error)
}

// ColumnGetter retrieves column names for a query by preparing it against a database.
type ColumnGetter interface {
	GetColumnNames(ctx context.Context, query string) ([]string, error)
}

// Expander expands SELECT * and RETURNING * queries by replacing * with explicit column names
// obtained from preparing the query against a database.
type Expander struct {
	colGetter ColumnGetter
	parser    Parser
	dialect   format.Dialect
}

// New creates a new Expander with the given column getter, parser, and dialect.
func New(colGetter ColumnGetter, parser Parser, dialect format.Dialect) *Expander {
	return &Expander{
		colGetter: colGetter,
		parser:    parser,
		dialect:   dialect,
	}
}

// Expand takes a SQL query, and if it contains * in SELECT or RETURNING clause,
// expands it to use explicit column names. Returns the expanded query string.
func (e *Expander) Expand(ctx context.Context, query string) (string, error) {
	// Parse the query
	stmts, err := e.parser.Parse(strings.NewReader(query))
	if err != nil {
		return "", fmt.Errorf("failed to parse query: %w", err)
	}

	if len(stmts) == 0 {
		return query, nil
	}

	stmt := stmts[0].Raw.Stmt

	// Check if there's any star in the statement (including CTEs, subqueries, etc.)
	if !hasStarAnywhere(stmt) {
		return query, nil
	}

	// Expand all stars in the statement recursively
	if err := e.expandNode(ctx, stmt); err != nil {
		return "", err
	}

	// Format the modified AST back to SQL
	expanded := ast.Format(stmts[0].Raw, e.dialect)

	return expanded, nil
}

// expandNode recursively expands * in all parts of the statement
func (e *Expander) expandNode(ctx context.Context, node ast.Node) error {
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	case *ast.SelectStmt:
		return e.expandSelectStmt(ctx, n)
	case *ast.InsertStmt:
		return e.expandInsertStmt(ctx, n)
	case *ast.UpdateStmt:
		return e.expandUpdateStmt(ctx, n)
	case *ast.DeleteStmt:
		return e.expandDeleteStmt(ctx, n)
	case *ast.CommonTableExpr:
		return e.expandNode(ctx, n.Ctequery)
	}
	return nil
}

// expandSelectStmt expands * in a SELECT statement including CTEs and subqueries
func (e *Expander) expandSelectStmt(ctx context.Context, stmt *ast.SelectStmt) error {
	// First expand any CTEs - must be done in order since later CTEs may depend on earlier ones
	if stmt.WithClause != nil && stmt.WithClause.Ctes != nil {
		for _, cteNode := range stmt.WithClause.Ctes.Items {
			cte, ok := cteNode.(*ast.CommonTableExpr)
			if !ok {
				continue
			}
			cteSelect, ok := cte.Ctequery.(*ast.SelectStmt)
			if !ok {
				continue
			}
			if hasStarInList(cteSelect.TargetList) {
				// Get column names for this CTE
				columns, err := e.getCTEColumnNames(ctx, stmt, cte)
				if err != nil {
					return err
				}
				cteSelect.TargetList = rewriteTargetList(cteSelect.TargetList, columns)
			}
			// Recursively handle nested CTEs/subqueries in this CTE
			if err := e.expandSelectStmtInner(ctx, cteSelect); err != nil {
				return err
			}
		}
	}

	// Expand subqueries in FROM clause
	if stmt.FromClause != nil {
		for _, fromItem := range stmt.FromClause.Items {
			if err := e.expandFromClause(ctx, fromItem); err != nil {
				return err
			}
		}
	}

	// Expand the target list if it has stars
	if hasStarInList(stmt.TargetList) {
		// Format the current state to get columns
		tempRaw := &ast.RawStmt{Stmt: stmt}
		tempQuery := ast.Format(tempRaw, e.dialect)
		columns, err := e.getColumnNames(ctx, tempQuery)
		if err != nil {
			return fmt.Errorf("failed to get column names: %w", err)
		}
		stmt.TargetList = rewriteTargetList(stmt.TargetList, columns)
	}

	return nil
}

// expandSelectStmtInner expands nested structures without re-processing the target list
func (e *Expander) expandSelectStmtInner(ctx context.Context, stmt *ast.SelectStmt) error {
	// Expand subqueries in FROM clause
	if stmt.FromClause != nil {
		for _, fromItem := range stmt.FromClause.Items {
			if err := e.expandFromClause(ctx, fromItem); err != nil {
				return err
			}
		}
	}
	return nil
}

// getCTEColumnNames gets the column names for a CTE by constructing a query with proper context
func (e *Expander) getCTEColumnNames(ctx context.Context, stmt *ast.SelectStmt, targetCTE *ast.CommonTableExpr) ([]string, error) {
	// Build a temporary query: WITH <all CTEs up to and including target> SELECT * FROM <targetCTE>
	var ctesToInclude []ast.Node
	for _, cteNode := range stmt.WithClause.Ctes.Items {
		ctesToInclude = append(ctesToInclude, cteNode)
		cte, ok := cteNode.(*ast.CommonTableExpr)
		if ok && cte.Ctename != nil && targetCTE.Ctename != nil && *cte.Ctename == *targetCTE.Ctename {
			break
		}
	}

	// Create a SELECT * FROM <ctename> with the relevant CTEs
	cteName := ""
	if targetCTE.Ctename != nil {
		cteName = *targetCTE.Ctename
	}

	tempStmt := &ast.SelectStmt{
		WithClause: &ast.WithClause{
			Ctes:      &ast.List{Items: ctesToInclude},
			Recursive: stmt.WithClause.Recursive,
		},
		TargetList: &ast.List{
			Items: []ast.Node{
				&ast.ResTarget{
					Val: &ast.ColumnRef{
						Fields: &ast.List{
							Items: []ast.Node{&ast.A_Star{}},
						},
					},
				},
			},
		},
		FromClause: &ast.List{
			Items: []ast.Node{
				&ast.RangeVar{
					Relname: &cteName,
				},
			},
		},
	}

	tempRaw := &ast.RawStmt{Stmt: tempStmt}
	tempQuery := ast.Format(tempRaw, e.dialect)

	return e.getColumnNames(ctx, tempQuery)
}

// expandInsertStmt expands * in an INSERT statement's RETURNING clause
func (e *Expander) expandInsertStmt(ctx context.Context, stmt *ast.InsertStmt) error {
	// Expand CTEs first
	if stmt.WithClause != nil && stmt.WithClause.Ctes != nil {
		for _, cte := range stmt.WithClause.Ctes.Items {
			if err := e.expandNode(ctx, cte); err != nil {
				return err
			}
		}
	}

	// Expand the SELECT part if present
	if stmt.SelectStmt != nil {
		if err := e.expandNode(ctx, stmt.SelectStmt); err != nil {
			return err
		}
	}

	// Expand RETURNING clause
	if hasStarInList(stmt.ReturningList) {
		tempRaw := &ast.RawStmt{Stmt: stmt}
		tempQuery := ast.Format(tempRaw, e.dialect)
		columns, err := e.getColumnNames(ctx, tempQuery)
		if err != nil {
			return fmt.Errorf("failed to get column names: %w", err)
		}
		stmt.ReturningList = rewriteTargetList(stmt.ReturningList, columns)
	}

	return nil
}

// expandUpdateStmt expands * in an UPDATE statement's RETURNING clause
func (e *Expander) expandUpdateStmt(ctx context.Context, stmt *ast.UpdateStmt) error {
	// Expand CTEs first
	if stmt.WithClause != nil && stmt.WithClause.Ctes != nil {
		for _, cte := range stmt.WithClause.Ctes.Items {
			if err := e.expandNode(ctx, cte); err != nil {
				return err
			}
		}
	}

	// Expand RETURNING clause
	if hasStarInList(stmt.ReturningList) {
		tempRaw := &ast.RawStmt{Stmt: stmt}
		tempQuery := ast.Format(tempRaw, e.dialect)
		columns, err := e.getColumnNames(ctx, tempQuery)
		if err != nil {
			return fmt.Errorf("failed to get column names: %w", err)
		}
		stmt.ReturningList = rewriteTargetList(stmt.ReturningList, columns)
	}

	return nil
}

// expandDeleteStmt expands * in a DELETE statement's RETURNING clause
func (e *Expander) expandDeleteStmt(ctx context.Context, stmt *ast.DeleteStmt) error {
	// Expand CTEs first
	if stmt.WithClause != nil && stmt.WithClause.Ctes != nil {
		for _, cte := range stmt.WithClause.Ctes.Items {
			if err := e.expandNode(ctx, cte); err != nil {
				return err
			}
		}
	}

	// Expand RETURNING clause
	if hasStarInList(stmt.ReturningList) {
		tempRaw := &ast.RawStmt{Stmt: stmt}
		tempQuery := ast.Format(tempRaw, e.dialect)
		columns, err := e.getColumnNames(ctx, tempQuery)
		if err != nil {
			return fmt.Errorf("failed to get column names: %w", err)
		}
		stmt.ReturningList = rewriteTargetList(stmt.ReturningList, columns)
	}

	return nil
}

// expandFromClause expands * in subqueries within FROM clause
func (e *Expander) expandFromClause(ctx context.Context, node ast.Node) error {
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	case *ast.RangeSubselect:
		if n.Subquery != nil {
			return e.expandNode(ctx, n.Subquery)
		}
	case *ast.JoinExpr:
		if err := e.expandFromClause(ctx, n.Larg); err != nil {
			return err
		}
		if err := e.expandFromClause(ctx, n.Rarg); err != nil {
			return err
		}
	}
	return nil
}

// hasStarAnywhere checks if there's a * anywhere in the statement using astutils.Search
func hasStarAnywhere(node ast.Node) bool {
	if node == nil {
		return false
	}
	// Use astutils.Search to find any A_Star node in the AST
	stars := astutils.Search(node, func(n ast.Node) bool {
		_, ok := n.(*ast.A_Star)
		return ok
	})
	return len(stars.Items) > 0
}

// hasStarInList checks if a target list contains a * expression using astutils.Search
func hasStarInList(targets *ast.List) bool {
	if targets == nil {
		return false
	}
	// Use astutils.Search to find any A_Star node in the target list
	stars := astutils.Search(targets, func(n ast.Node) bool {
		_, ok := n.(*ast.A_Star)
		return ok
	})
	return len(stars.Items) > 0
}

// getColumnNames prepares the query and returns the column names from the result
func (e *Expander) getColumnNames(ctx context.Context, query string) ([]string, error) {
	return e.colGetter.GetColumnNames(ctx, query)
}

// countStarsInList counts the number of * expressions in a target list
func countStarsInList(targets *ast.List) int {
	if targets == nil {
		return 0
	}
	count := 0
	for _, target := range targets.Items {
		resTarget, ok := target.(*ast.ResTarget)
		if !ok {
			continue
		}
		if resTarget.Val == nil {
			continue
		}
		colRef, ok := resTarget.Val.(*ast.ColumnRef)
		if !ok {
			continue
		}
		if colRef.Fields == nil {
			continue
		}
		for _, field := range colRef.Fields.Items {
			if _, ok := field.(*ast.A_Star); ok {
				count++
				break
			}
		}
	}
	return count
}

// countNonStarsInList counts the number of non-* expressions in a target list
func countNonStarsInList(targets *ast.List) int {
	if targets == nil {
		return 0
	}
	count := 0
	for _, target := range targets.Items {
		resTarget, ok := target.(*ast.ResTarget)
		if !ok {
			count++
			continue
		}
		if resTarget.Val == nil {
			count++
			continue
		}
		colRef, ok := resTarget.Val.(*ast.ColumnRef)
		if !ok {
			count++
			continue
		}
		if colRef.Fields == nil {
			count++
			continue
		}
		isStar := false
		for _, field := range colRef.Fields.Items {
			if _, ok := field.(*ast.A_Star); ok {
				isStar = true
				break
			}
		}
		if !isStar {
			count++
		}
	}
	return count
}

// rewriteTargetList replaces * in a target list with explicit column references
func rewriteTargetList(targets *ast.List, columns []string) *ast.List {
	if targets == nil {
		return nil
	}

	starCount := countStarsInList(targets)
	nonStarCount := countNonStarsInList(targets)

	// Calculate how many columns each * expands to
	// Total columns = (columns per star * number of stars) + non-star columns
	// So: columns per star = (total - non-star) / stars
	columnsPerStar := 0
	if starCount > 0 {
		columnsPerStar = (len(columns) - nonStarCount) / starCount
	}

	newItems := make([]ast.Node, 0, len(columns))
	colIndex := 0

	for _, target := range targets.Items {
		resTarget, ok := target.(*ast.ResTarget)
		if !ok {
			newItems = append(newItems, target)
			colIndex++
			continue
		}

		if resTarget.Val == nil {
			newItems = append(newItems, target)
			colIndex++
			continue
		}

		colRef, ok := resTarget.Val.(*ast.ColumnRef)
		if !ok {
			newItems = append(newItems, target)
			colIndex++
			continue
		}

		if colRef.Fields == nil {
			newItems = append(newItems, target)
			colIndex++
			continue
		}

		// Check if this is a * (with or without table qualifier)
		// and extract any table prefix
		isStar := false
		var tablePrefix []string
		for _, field := range colRef.Fields.Items {
			if _, ok := field.(*ast.A_Star); ok {
				isStar = true
				break
			}
			// Collect prefix parts (schema, table name)
			if str, ok := field.(*ast.String); ok {
				tablePrefix = append(tablePrefix, str.Str)
			}
		}

		if !isStar {
			newItems = append(newItems, target)
			colIndex++
			continue
		}

		// Replace * with explicit column references
		for i := 0; i < columnsPerStar && colIndex < len(columns); i++ {
			newItems = append(newItems, makeColumnTargetWithPrefix(columns[colIndex], tablePrefix))
			colIndex++
		}
	}

	return &ast.List{Items: newItems}
}

// makeColumnTargetWithPrefix creates a ResTarget node for a column reference with optional table prefix
func makeColumnTargetWithPrefix(colName string, prefix []string) ast.Node {
	fields := make([]ast.Node, 0, len(prefix)+1)

	// Add prefix parts (schema, table name)
	for _, p := range prefix {
		fields = append(fields, &ast.String{Str: p})
	}

	// Add column name
	fields = append(fields, &ast.String{Str: colName})

	return &ast.ResTarget{
		Val: &ast.ColumnRef{
			Fields: &ast.List{Items: fields},
		},
	}
}
