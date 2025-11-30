package expander

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	nodes "github.com/pganalyze/pg_query_go/v6"
)

// Expander expands SELECT * and RETURNING * queries by replacing * with explicit column names
// obtained from preparing the query against a PostgreSQL database.
type Expander struct {
	pool *pgxpool.Pool
}

// New creates a new Expander with the given connection pool.
func New(pool *pgxpool.Pool) *Expander {
	return &Expander{pool: pool}
}

// Expand takes a SQL query, and if it contains * in SELECT or RETURNING clause,
// expands it to use explicit column names. Returns the expanded query string.
func (e *Expander) Expand(ctx context.Context, query string) (string, error) {
	// Parse the query
	tree, err := parse(query)
	if err != nil {
		return "", fmt.Errorf("failed to parse query: %w", err)
	}

	if len(tree.Stmts) == 0 {
		return query, nil
	}

	stmt := tree.Stmts[0].Stmt

	// Check if there's any star in the statement (including CTEs, subqueries, etc.)
	if !hasStarAnywhere(stmt) {
		return query, nil
	}

	// Expand all stars in the statement recursively
	if err := e.expandNode(ctx, stmt); err != nil {
		return "", err
	}

	// Deparse the modified AST back to SQL
	expanded, err := deparse(tree)
	if err != nil {
		return "", fmt.Errorf("failed to deparse query: %w", err)
	}

	return expanded, nil
}

// expandNode recursively expands * in all parts of the statement
func (e *Expander) expandNode(ctx context.Context, node *nodes.Node) error {
	if node == nil {
		return nil
	}

	switch n := node.Node.(type) {
	case *nodes.Node_SelectStmt:
		return e.expandSelectStmt(ctx, n.SelectStmt)
	case *nodes.Node_InsertStmt:
		return e.expandInsertStmt(ctx, n.InsertStmt)
	case *nodes.Node_UpdateStmt:
		return e.expandUpdateStmt(ctx, n.UpdateStmt)
	case *nodes.Node_DeleteStmt:
		return e.expandDeleteStmt(ctx, n.DeleteStmt)
	case *nodes.Node_CommonTableExpr:
		return e.expandNode(ctx, n.CommonTableExpr.Ctequery)
	}
	return nil
}

// expandSelectStmt expands * in a SELECT statement including CTEs and subqueries
func (e *Expander) expandSelectStmt(ctx context.Context, stmt *nodes.SelectStmt) error {
	// First expand any CTEs - must be done in order since later CTEs may depend on earlier ones
	if stmt.WithClause != nil {
		for _, cte := range stmt.WithClause.Ctes {
			cteExpr, ok := cte.Node.(*nodes.Node_CommonTableExpr)
			if !ok {
				continue
			}
			cteSelect, ok := cteExpr.CommonTableExpr.Ctequery.Node.(*nodes.Node_SelectStmt)
			if !ok {
				continue
			}
			if hasStarInList(cteSelect.SelectStmt.TargetList) {
				// Deparse the full statement (with WITH clause context) but query just this CTE
				// We need to build a query that includes all prior CTEs for context
				columns, err := e.getCTEColumnNames(ctx, stmt, cteExpr.CommonTableExpr)
				if err != nil {
					return err
				}
				cteSelect.SelectStmt.TargetList = rewriteTargetList(cteSelect.SelectStmt.TargetList, columns)
			}
			// Recursively handle nested CTEs/subqueries in this CTE
			if err := e.expandSelectStmtInner(ctx, cteSelect.SelectStmt); err != nil {
				return err
			}
		}
	}

	// Expand subqueries in FROM clause
	for _, fromItem := range stmt.FromClause {
		if err := e.expandFromClause(ctx, fromItem); err != nil {
			return err
		}
	}

	// Expand the target list if it has stars
	if hasStarInList(stmt.TargetList) {
		// Deparse the current state to get columns
		tempTree := &nodes.ParseResult{
			Stmts: []*nodes.RawStmt{{Stmt: &nodes.Node{Node: &nodes.Node_SelectStmt{SelectStmt: stmt}}}},
		}
		tempQuery, err := deparse(tempTree)
		if err != nil {
			return fmt.Errorf("failed to deparse for column lookup: %w", err)
		}
		columns, err := e.getColumnNames(ctx, tempQuery)
		if err != nil {
			return fmt.Errorf("failed to get column names: %w", err)
		}
		stmt.TargetList = rewriteTargetList(stmt.TargetList, columns)
	}

	return nil
}

// expandSelectStmtInner expands nested structures without re-processing the target list
func (e *Expander) expandSelectStmtInner(ctx context.Context, stmt *nodes.SelectStmt) error {
	// Expand subqueries in FROM clause
	for _, fromItem := range stmt.FromClause {
		if err := e.expandFromClause(ctx, fromItem); err != nil {
			return err
		}
	}
	return nil
}

// getCTEColumnNames gets the column names for a CTE by constructing a query with proper context
func (e *Expander) getCTEColumnNames(ctx context.Context, stmt *nodes.SelectStmt, targetCTE *nodes.CommonTableExpr) ([]string, error) {
	// Build a temporary query: WITH <all CTEs up to and including target> SELECT * FROM <targetCTE>
	// This gives us the proper context for resolving column names

	var ctesToInclude []*nodes.Node
	for _, cte := range stmt.WithClause.Ctes {
		ctesToInclude = append(ctesToInclude, cte)
		cteExpr, ok := cte.Node.(*nodes.Node_CommonTableExpr)
		if ok && cteExpr.CommonTableExpr.Ctename == targetCTE.Ctename {
			break
		}
	}

	// Create a SELECT * FROM <ctename> with the relevant CTEs
	tempStmt := &nodes.SelectStmt{
		WithClause: &nodes.WithClause{
			Ctes:      ctesToInclude,
			Recursive: stmt.WithClause.Recursive,
		},
		TargetList: []*nodes.Node{
			{
				Node: &nodes.Node_ResTarget{
					ResTarget: &nodes.ResTarget{
						Val: &nodes.Node{
							Node: &nodes.Node_ColumnRef{
								ColumnRef: &nodes.ColumnRef{
									Fields: []*nodes.Node{
										{Node: &nodes.Node_AStar{AStar: &nodes.A_Star{}}},
									},
								},
							},
						},
					},
				},
			},
		},
		FromClause: []*nodes.Node{
			{
				Node: &nodes.Node_RangeVar{
					RangeVar: &nodes.RangeVar{
						Relname: targetCTE.Ctename,
						Inh:     true,
					},
				},
			},
		},
	}

	tempTree := &nodes.ParseResult{
		Stmts: []*nodes.RawStmt{{Stmt: &nodes.Node{Node: &nodes.Node_SelectStmt{SelectStmt: tempStmt}}}},
	}
	tempQuery, err := deparse(tempTree)
	if err != nil {
		return nil, fmt.Errorf("failed to deparse CTE query: %w", err)
	}

	return e.getColumnNames(ctx, tempQuery)
}

// expandInsertStmt expands * in an INSERT statement's RETURNING clause
func (e *Expander) expandInsertStmt(ctx context.Context, stmt *nodes.InsertStmt) error {
	// Expand CTEs first
	if stmt.WithClause != nil {
		for _, cte := range stmt.WithClause.Ctes {
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
		tempTree := &nodes.ParseResult{
			Stmts: []*nodes.RawStmt{{Stmt: &nodes.Node{Node: &nodes.Node_InsertStmt{InsertStmt: stmt}}}},
		}
		tempQuery, err := deparse(tempTree)
		if err != nil {
			return fmt.Errorf("failed to deparse for column lookup: %w", err)
		}
		columns, err := e.getColumnNames(ctx, tempQuery)
		if err != nil {
			return fmt.Errorf("failed to get column names: %w", err)
		}
		stmt.ReturningList = rewriteTargetList(stmt.ReturningList, columns)
	}

	return nil
}

// expandUpdateStmt expands * in an UPDATE statement's RETURNING clause
func (e *Expander) expandUpdateStmt(ctx context.Context, stmt *nodes.UpdateStmt) error {
	// Expand CTEs first
	if stmt.WithClause != nil {
		for _, cte := range stmt.WithClause.Ctes {
			if err := e.expandNode(ctx, cte); err != nil {
				return err
			}
		}
	}

	// Expand RETURNING clause
	if hasStarInList(stmt.ReturningList) {
		tempTree := &nodes.ParseResult{
			Stmts: []*nodes.RawStmt{{Stmt: &nodes.Node{Node: &nodes.Node_UpdateStmt{UpdateStmt: stmt}}}},
		}
		tempQuery, err := deparse(tempTree)
		if err != nil {
			return fmt.Errorf("failed to deparse for column lookup: %w", err)
		}
		columns, err := e.getColumnNames(ctx, tempQuery)
		if err != nil {
			return fmt.Errorf("failed to get column names: %w", err)
		}
		stmt.ReturningList = rewriteTargetList(stmt.ReturningList, columns)
	}

	return nil
}

// expandDeleteStmt expands * in a DELETE statement's RETURNING clause
func (e *Expander) expandDeleteStmt(ctx context.Context, stmt *nodes.DeleteStmt) error {
	// Expand CTEs first
	if stmt.WithClause != nil {
		for _, cte := range stmt.WithClause.Ctes {
			if err := e.expandNode(ctx, cte); err != nil {
				return err
			}
		}
	}

	// Expand RETURNING clause
	if hasStarInList(stmt.ReturningList) {
		tempTree := &nodes.ParseResult{
			Stmts: []*nodes.RawStmt{{Stmt: &nodes.Node{Node: &nodes.Node_DeleteStmt{DeleteStmt: stmt}}}},
		}
		tempQuery, err := deparse(tempTree)
		if err != nil {
			return fmt.Errorf("failed to deparse for column lookup: %w", err)
		}
		columns, err := e.getColumnNames(ctx, tempQuery)
		if err != nil {
			return fmt.Errorf("failed to get column names: %w", err)
		}
		stmt.ReturningList = rewriteTargetList(stmt.ReturningList, columns)
	}

	return nil
}

// expandFromClause expands * in subqueries within FROM clause
func (e *Expander) expandFromClause(ctx context.Context, node *nodes.Node) error {
	if node == nil {
		return nil
	}

	switch n := node.Node.(type) {
	case *nodes.Node_RangeSubselect:
		if n.RangeSubselect.Subquery != nil {
			return e.expandNode(ctx, n.RangeSubselect.Subquery)
		}
	case *nodes.Node_JoinExpr:
		if err := e.expandFromClause(ctx, n.JoinExpr.Larg); err != nil {
			return err
		}
		if err := e.expandFromClause(ctx, n.JoinExpr.Rarg); err != nil {
			return err
		}
	}
	return nil
}

// hasStarAnywhere checks if there's a * anywhere in the statement
func hasStarAnywhere(node *nodes.Node) bool {
	if node == nil {
		return false
	}

	switch n := node.Node.(type) {
	case *nodes.Node_SelectStmt:
		if hasStarInList(n.SelectStmt.TargetList) {
			return true
		}
		if n.SelectStmt.WithClause != nil {
			for _, cte := range n.SelectStmt.WithClause.Ctes {
				if hasStarAnywhere(cte) {
					return true
				}
			}
		}
		for _, from := range n.SelectStmt.FromClause {
			if hasStarAnywhere(from) {
				return true
			}
		}
	case *nodes.Node_InsertStmt:
		if hasStarInList(n.InsertStmt.ReturningList) {
			return true
		}
		if n.InsertStmt.WithClause != nil {
			for _, cte := range n.InsertStmt.WithClause.Ctes {
				if hasStarAnywhere(cte) {
					return true
				}
			}
		}
		if hasStarAnywhere(n.InsertStmt.SelectStmt) {
			return true
		}
	case *nodes.Node_UpdateStmt:
		if hasStarInList(n.UpdateStmt.ReturningList) {
			return true
		}
		if n.UpdateStmt.WithClause != nil {
			for _, cte := range n.UpdateStmt.WithClause.Ctes {
				if hasStarAnywhere(cte) {
					return true
				}
			}
		}
	case *nodes.Node_DeleteStmt:
		if hasStarInList(n.DeleteStmt.ReturningList) {
			return true
		}
		if n.DeleteStmt.WithClause != nil {
			for _, cte := range n.DeleteStmt.WithClause.Ctes {
				if hasStarAnywhere(cte) {
					return true
				}
			}
		}
	case *nodes.Node_CommonTableExpr:
		return hasStarAnywhere(n.CommonTableExpr.Ctequery)
	case *nodes.Node_RangeSubselect:
		return hasStarAnywhere(n.RangeSubselect.Subquery)
	case *nodes.Node_JoinExpr:
		return hasStarAnywhere(n.JoinExpr.Larg) || hasStarAnywhere(n.JoinExpr.Rarg)
	}
	return false
}

// hasStarInList checks if a target list contains a * expression
func hasStarInList(targets []*nodes.Node) bool {
	for _, target := range targets {
		resTarget, ok := target.Node.(*nodes.Node_ResTarget)
		if !ok {
			continue
		}
		if resTarget.ResTarget.Val == nil {
			continue
		}
		colRef, ok := resTarget.ResTarget.Val.Node.(*nodes.Node_ColumnRef)
		if !ok {
			continue
		}
		for _, field := range colRef.ColumnRef.Fields {
			if _, ok := field.Node.(*nodes.Node_AStar); ok {
				return true
			}
		}
	}
	return false
}

// getColumnNames prepares the query and returns the column names from the result
func (e *Expander) getColumnNames(ctx context.Context, query string) ([]string, error) {
	conn, err := e.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	// Prepare the statement to get column metadata
	desc, err := conn.Conn().Prepare(ctx, "", query)
	if err != nil {
		return nil, err
	}

	columns := make([]string, len(desc.Fields))
	for i, field := range desc.Fields {
		columns[i] = field.Name
	}

	return columns, nil
}

// countStarsInList counts the number of * expressions in a target list
func countStarsInList(targets []*nodes.Node) int {
	count := 0
	for _, target := range targets {
		resTarget, ok := target.Node.(*nodes.Node_ResTarget)
		if !ok {
			continue
		}
		if resTarget.ResTarget.Val == nil {
			continue
		}
		colRef, ok := resTarget.ResTarget.Val.Node.(*nodes.Node_ColumnRef)
		if !ok {
			continue
		}
		for _, field := range colRef.ColumnRef.Fields {
			if _, ok := field.Node.(*nodes.Node_AStar); ok {
				count++
				break
			}
		}
	}
	return count
}

// countNonStarsInList counts the number of non-* expressions in a target list
func countNonStarsInList(targets []*nodes.Node) int {
	count := 0
	for _, target := range targets {
		resTarget, ok := target.Node.(*nodes.Node_ResTarget)
		if !ok {
			count++
			continue
		}
		if resTarget.ResTarget.Val == nil {
			count++
			continue
		}
		colRef, ok := resTarget.ResTarget.Val.Node.(*nodes.Node_ColumnRef)
		if !ok {
			count++
			continue
		}
		isStar := false
		for _, field := range colRef.ColumnRef.Fields {
			if _, ok := field.Node.(*nodes.Node_AStar); ok {
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
func rewriteTargetList(targets []*nodes.Node, columns []string) []*nodes.Node {
	starCount := countStarsInList(targets)
	nonStarCount := countNonStarsInList(targets)

	// Calculate how many columns each * expands to
	// Total columns = (columns per star * number of stars) + non-star columns
	// So: columns per star = (total - non-star) / stars
	columnsPerStar := 0
	if starCount > 0 {
		columnsPerStar = (len(columns) - nonStarCount) / starCount
	}

	newTargets := make([]*nodes.Node, 0, len(columns))
	colIndex := 0

	for _, target := range targets {
		resTarget, ok := target.Node.(*nodes.Node_ResTarget)
		if !ok {
			newTargets = append(newTargets, target)
			colIndex++
			continue
		}

		if resTarget.ResTarget.Val == nil {
			newTargets = append(newTargets, target)
			colIndex++
			continue
		}

		colRef, ok := resTarget.ResTarget.Val.Node.(*nodes.Node_ColumnRef)
		if !ok {
			newTargets = append(newTargets, target)
			colIndex++
			continue
		}

		// Check if this is a * (with or without table qualifier)
		// and extract any table prefix
		isStar := false
		var tablePrefix []string
		for _, field := range colRef.ColumnRef.Fields {
			if _, ok := field.Node.(*nodes.Node_AStar); ok {
				isStar = true
				break
			}
			// Collect prefix parts (schema, table name)
			if str, ok := field.Node.(*nodes.Node_String_); ok {
				tablePrefix = append(tablePrefix, str.String_.Sval)
			}
		}

		if !isStar {
			newTargets = append(newTargets, target)
			colIndex++
			continue
		}

		// Replace * with explicit column references
		for i := 0; i < columnsPerStar && colIndex < len(columns); i++ {
			newTargets = append(newTargets, makeColumnTargetWithPrefix(columns[colIndex], tablePrefix))
			colIndex++
		}
	}

	return newTargets
}

// makeColumnTargetWithPrefix creates a ResTarget node for a column reference with optional table prefix
func makeColumnTargetWithPrefix(colName string, prefix []string) *nodes.Node {
	fields := make([]*nodes.Node, 0, len(prefix)+1)

	// Add prefix parts (schema, table name)
	for _, p := range prefix {
		fields = append(fields, &nodes.Node{
			Node: &nodes.Node_String_{
				String_: &nodes.String{
					Sval: p,
				},
			},
		})
	}

	// Add column name
	fields = append(fields, &nodes.Node{
		Node: &nodes.Node_String_{
			String_: &nodes.String{
				Sval: colName,
			},
		},
	})

	return &nodes.Node{
		Node: &nodes.Node_ResTarget{
			ResTarget: &nodes.ResTarget{
				Val: &nodes.Node{
					Node: &nodes.Node_ColumnRef{
						ColumnRef: &nodes.ColumnRef{
							Fields: fields,
						},
					},
				},
			},
		},
	}
}
