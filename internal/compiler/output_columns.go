package compiler

import (
	"errors"
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
	"github.com/sqlc-dev/sqlc/internal/sql/lang"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

// OutputColumns determines which columns a statement will output
func (c *Compiler) OutputColumns(stmt ast.Node) ([]*catalog.Column, error) {
	qc, err := c.buildQueryCatalog(c.catalog, stmt, nil)
	if err != nil {
		return nil, err
	}
	cols, err := c.outputColumns(qc, stmt)
	if err != nil {
		return nil, err
	}

	catCols := make([]*catalog.Column, 0, len(cols))
	for _, col := range cols {
		catCols = append(catCols, &catalog.Column{
			Name:       col.Name,
			Type:       ast.TypeName{Name: col.DataType},
			IsNotNull:  col.NotNull,
			IsUnsigned: col.Unsigned,
			IsArray:    col.IsArray,
			ArrayDims:  col.ArrayDims,
			Comment:    col.Comment,
			Length:     col.Length,
		})
	}
	return catCols, nil
}

func hasStarRef(cf *ast.ColumnRef) bool {
	for _, item := range cf.Fields.Items {
		if _, ok := item.(*ast.A_Star); ok {
			return true
		}
	}
	return false
}

// Compute the output columns for a statement.
//
// Return an error if column references are ambiguous
// Return an error if column references don't exist
func (c *Compiler) outputColumns(qc *QueryCatalog, node ast.Node) ([]*Column, error) {
	tables, err := c.sourceTables(qc, node)
	if err != nil {
		return nil, err
	}

	targets := &ast.List{}
	switch n := node.(type) {
	case *ast.DeleteStmt:
		targets = n.ReturningList
	case *ast.InsertStmt:
		targets = n.ReturningList
	case *ast.SelectStmt:
		targets = n.TargetList
		isUnion := len(targets.Items) == 0 && n.Larg != nil

		if n.GroupClause != nil {
			for _, item := range n.GroupClause.Items {
				if err := findColumnForNode(item, tables, targets); err != nil {
					return nil, err
				}
			}
		}
		validateOrderBy := true
		if c.conf.StrictOrderBy != nil {
			validateOrderBy = *c.conf.StrictOrderBy
		}
		if !isUnion && validateOrderBy {
			if n.SortClause != nil {
				for _, item := range n.SortClause.Items {
					sb, ok := item.(*ast.SortBy)
					if !ok {
						continue
					}
					if err := findColumnForNode(sb.Node, tables, targets); err != nil {
						return nil, fmt.Errorf("%v: if you want to skip this validation, set 'strict_order_by' to false", err)
					}
				}
			}
			if n.WindowClause != nil {
				for _, item := range n.WindowClause.Items {
					sb, ok := item.(*ast.List)
					if !ok {
						continue
					}
					for _, single := range sb.Items {
						caseExpr, ok := single.(*ast.CaseExpr)
						if !ok {
							continue
						}
						if err := findColumnForNode(caseExpr.Xpr, tables, targets); err != nil {
							return nil, fmt.Errorf("%v: if you want to skip this validation, set 'strict_order_by' to false", err)
						}
					}
				}
			}
		}

		// For UNION queries, targets is empty and we need to look for the
		// columns in Largs.
		if isUnion {
			return c.outputColumns(qc, n.Larg)
		}
	case *ast.UpdateStmt:
		targets = n.ReturningList
	}

	var cols []*Column

	for _, target := range targets.Items {
		res, ok := target.(*ast.ResTarget)
		if !ok {
			continue
		}
		switch n := res.Val.(type) {

		case *ast.A_Const:
			name := ""
			if res.Name != nil {
				name = *res.Name
			}
			switch n.Val.(type) {
			case *ast.String:
				cols = append(cols, &Column{Name: name, DataType: "text", NotNull: true})
			case *ast.Integer:
				cols = append(cols, &Column{Name: name, DataType: "int", NotNull: true})
			case *ast.Float:
				cols = append(cols, &Column{Name: name, DataType: "float", NotNull: true})
			case *ast.Boolean:
				cols = append(cols, &Column{Name: name, DataType: "bool", NotNull: true})
			default:
				cols = append(cols, &Column{Name: name, DataType: "any", NotNull: false})
			}

		case *ast.A_Expr:
			name := ""
			if res.Name != nil {
				name = *res.Name
			}
			switch op := astutils.Join(n.Name, ""); {
			case lang.IsComparisonOperator(op):
				// TODO: Generate a name for these operations
				cols = append(cols, &Column{Name: name, DataType: "bool", NotNull: true})
			case lang.IsMathematicalOperator(op):
				cols = append(cols, &Column{Name: name, DataType: "int", NotNull: true})
			default:
				cols = append(cols, &Column{Name: name, DataType: "any", NotNull: false})
			}

		case *ast.BoolExpr:
			name := ""
			if res.Name != nil {
				name = *res.Name
			}
			notNull := false
			if len(n.Args.Items) == 1 {
				switch n.Boolop {
				case ast.BoolExprTypeIsNull, ast.BoolExprTypeIsNotNull:
					notNull = true
				case ast.BoolExprTypeNot:
					sublink, ok := n.Args.Items[0].(*ast.SubLink)
					if ok && sublink.SubLinkType == ast.EXISTS_SUBLINK {
						notNull = true
						if name == "" {
							name = "not_exists"
						}
					}
				}
			}
			cols = append(cols, &Column{Name: name, DataType: "bool", NotNull: notNull})

		case *ast.CaseExpr:
			name := ""
			if res.Name != nil {
				name = *res.Name
			}
			// TODO: The TypeCase and A_Const code has been copied from below. Instead, we
			// need a recurse function to get the type of a node.
			if tc, ok := n.Defresult.(*ast.TypeCast); ok {
				if tc.TypeName == nil {
					return nil, errors.New("no type name type cast")
				}
				name := ""
				if ref, ok := tc.Arg.(*ast.ColumnRef); ok {
					name = astutils.Join(ref.Fields, "_")
				}
				if res.Name != nil {
					name = *res.Name
				}
				// TODO Validate column names
				col := toColumn(tc.TypeName)
				col.Name = name
				cols = append(cols, col)
			} else if aconst, ok := n.Defresult.(*ast.A_Const); ok {
				switch aconst.Val.(type) {
				case *ast.String:
					cols = append(cols, &Column{Name: name, DataType: "text", NotNull: true})
				case *ast.Integer:
					cols = append(cols, &Column{Name: name, DataType: "int", NotNull: true})
				case *ast.Float:
					cols = append(cols, &Column{Name: name, DataType: "float", NotNull: true})
				case *ast.Boolean:
					cols = append(cols, &Column{Name: name, DataType: "bool", NotNull: true})
				default:
					cols = append(cols, &Column{Name: name, DataType: "any", NotNull: false})
				}
			} else {
				cols = append(cols, &Column{Name: name, DataType: "any", NotNull: false})
			}

		case *ast.CoalesceExpr:
			name := "coalesce"
			if res.Name != nil {
				name = *res.Name
			}
			var firstColumn *Column
			var shouldNotBeNull bool
			for _, arg := range n.Args.Items {
				if _, ok := arg.(*ast.A_Const); ok {
					shouldNotBeNull = true
					continue
				}
				if ref, ok := arg.(*ast.ColumnRef); ok {
					columns, err := outputColumnRefs(res, tables, ref)
					if err != nil {
						return nil, err
					}
					for _, c := range columns {
						if firstColumn == nil {
							firstColumn = c
						}
						shouldNotBeNull = shouldNotBeNull || c.NotNull
					}
				}
			}
			if firstColumn != nil {
				firstColumn.NotNull = shouldNotBeNull
				firstColumn.skipTableRequiredCheck = true
				cols = append(cols, firstColumn)
			} else {
				cols = append(cols, &Column{Name: name, DataType: "any", NotNull: false})
			}

		case *ast.ColumnRef:
			if hasStarRef(n) {

				// add a column with a reference to an embedded table
				if embed, ok := qc.embeds.Find(n); ok {
					cols = append(cols, &Column{
						Name:       embed.Table.Name,
						EmbedTable: embed.Table,
					})
					continue
				}

				// TODO: This code is copied in func expand()
				for _, t := range tables {
					scope := astutils.Join(n.Fields, ".")
					if scope != "" && scope != t.Rel.Name {
						continue
					}
					for _, c := range t.Columns {
						cname := c.Name
						if res.Name != nil {
							cname = *res.Name
						}
						cols = append(cols, &Column{
							Name:         cname,
							OriginalName: c.Name,
							Type:         c.Type,
							Scope:        scope,
							Table:        c.Table,
							TableAlias:   t.Rel.Name,
							DataType:     c.DataType,
							NotNull:      c.NotNull,
							Unsigned:     c.Unsigned,
							IsArray:      c.IsArray,
							ArrayDims:    c.ArrayDims,
							Length:       c.Length,
						})
					}
				}
				continue
			}

			columns, err := outputColumnRefs(res, tables, n)
			if err != nil {
				return nil, err
			}
			cols = append(cols, columns...)

		case *ast.FuncCall:
			rel := n.Func
			name := rel.Name
			if res.Name != nil {
				name = *res.Name
			}
			fun, err := qc.catalog.ResolveFuncCall(n)
			if err == nil {
				cols = append(cols, &Column{
					Name:       name,
					DataType:   dataType(fun.ReturnType),
					NotNull:    !fun.ReturnTypeNullable,
					IsFuncCall: true,
				})
			} else {
				cols = append(cols, &Column{
					Name:       name,
					DataType:   "any",
					IsFuncCall: true,
				})
			}

		case *ast.SubLink:
			name := "exists"
			if res.Name != nil {
				name = *res.Name
			}
			switch n.SubLinkType {
			case ast.EXISTS_SUBLINK:
				cols = append(cols, &Column{Name: name, DataType: "bool", NotNull: true})
			case ast.EXPR_SUBLINK:
				subcols, err := c.outputColumns(qc, n.Subselect)
				if err != nil {
					return nil, err
				}
				first := subcols[0]
				if res.Name != nil {
					first.Name = *res.Name
				}
				cols = append(cols, first)
			default:
				cols = append(cols, &Column{Name: name, DataType: "any", NotNull: false})
			}

		case *ast.TypeCast:
			if n.TypeName == nil {
				return nil, errors.New("no type name type cast")
			}
			name := ""
			if ref, ok := n.Arg.(*ast.ColumnRef); ok {
				name = astutils.Join(ref.Fields, "_")
			}
			if res.Name != nil {
				name = *res.Name
			}
			// TODO Validate column names
			col := toColumn(n.TypeName)
			col.Name = name
			// TODO Add correct, real type inference
			if constant, ok := n.Arg.(*ast.A_Const); ok {
				if _, ok := constant.Val.(*ast.Null); ok {
					col.NotNull = false
				}
			}
			cols = append(cols, col)

		case *ast.SelectStmt:
			subcols, err := c.outputColumns(qc, n)
			if err != nil {
				return nil, err
			}
			first := subcols[0]
			if res.Name != nil {
				first.Name = *res.Name
			}
			cols = append(cols, first)

		default:
			name := ""
			if res.Name != nil {
				name = *res.Name
			}
			cols = append(cols, &Column{Name: name, DataType: "any", NotNull: false})

		}
	}

	if n, ok := node.(*ast.SelectStmt); ok {
		for _, col := range cols {
			if !col.NotNull || col.Table == nil || col.skipTableRequiredCheck {
				continue
			}
			for _, f := range n.FromClause.Items {
				res := isTableRequired(f, col, tableRequired)
				if res != tableNotFound {
					col.NotNull = res == tableRequired
					break
				}
			}
		}
	}

	return cols, nil
}

const (
	tableNotFound = iota
	tableRequired
	tableOptional
)

func isTableRequired(n ast.Node, col *Column, prior int) int {
	switch n := n.(type) {
	case *ast.RangeVar:
		tableMatch := *n.Relname == col.Table.Name
		aliasMatch := true
		if n.Alias != nil && col.TableAlias != "" {
			aliasMatch = *n.Alias.Aliasname == col.TableAlias
		}
		if aliasMatch && tableMatch {
			return prior
		}
	case *ast.JoinExpr:
		helper := func(l, r int) int {
			if res := isTableRequired(n.Larg, col, l); res != tableNotFound {
				return res
			}
			if res := isTableRequired(n.Rarg, col, r); res != tableNotFound {
				return res
			}
			return tableNotFound
		}
		switch n.Jointype {
		case ast.JoinTypeLeft:
			return helper(tableRequired, tableOptional)
		case ast.JoinTypeRight:
			return helper(tableOptional, tableRequired)
		case ast.JoinTypeFull:
			return helper(tableOptional, tableOptional)
		case ast.JoinTypeInner:
			return helper(tableRequired, tableRequired)
		}
	case *ast.List:
		for _, item := range n.Items {
			if res := isTableRequired(item, col, prior); res != tableNotFound {
				return res
			}
		}
	}

	return tableNotFound
}

type tableVisitor struct {
	list ast.List
}

func (r *tableVisitor) Visit(n ast.Node) astutils.Visitor {
	switch n.(type) {
	case *ast.RangeVar, *ast.RangeFunction:
		r.list.Items = append(r.list.Items, n)
		return r
	case *ast.RangeSubselect:
		r.list.Items = append(r.list.Items, n)
		return nil
	default:
		return r
	}
}

// Compute the output columns for a statement.
//
// Return an error if column references are ambiguous
// Return an error if column references don't exist
// Return an error if a table is referenced twice
// Return an error if an unknown column is referenced
func (c *Compiler) sourceTables(qc *QueryCatalog, node ast.Node) ([]*Table, error) {
	list := &ast.List{}
	switch n := node.(type) {
	case *ast.DeleteStmt:
		list = n.Relations
	case *ast.InsertStmt:
		list = &ast.List{
			Items: []ast.Node{n.Relation},
		}
	case *ast.SelectStmt:
		var tv tableVisitor
		astutils.Walk(&tv, n.FromClause)
		list = &tv.list
	case *ast.TruncateStmt:
		list = astutils.Search(n.Relations, func(node ast.Node) bool {
			_, ok := node.(*ast.RangeVar)
			return ok
		})
	case *ast.RefreshMatViewStmt:
		list = astutils.Search(n.Relation, func(node ast.Node) bool {
			_, ok := node.(*ast.RangeVar)
			return ok
		})
	case *ast.UpdateStmt:
		var tv tableVisitor
		astutils.Walk(&tv, n.FromClause)
		astutils.Walk(&tv, n.Relations)
		list = &tv.list
	}

	var tables []*Table
	for _, item := range list.Items {
		item := item
		switch n := item.(type) {

		case *ast.RangeFunction:
			var funcCall *ast.FuncCall
			switch f := n.Functions.Items[0].(type) {
			case *ast.List:
				switch fi := f.Items[0].(type) {
				case *ast.FuncCall:
					funcCall = fi
				case *ast.SQLValueFunction:
					continue // TODO handle this correctly
				default:
					continue
				}
			case *ast.FuncCall:
				funcCall = f
			default:
				return nil, fmt.Errorf("sourceTables: unsupported function call type %T", n.Functions.Items[0])
			}

			// If the function or table can't be found, don't error out.  There
			// are many queries that depend on functions unknown to sqlc.
			fn, err := qc.GetFunc(funcCall.Func)
			if err != nil {
				continue
			}
			var table *Table
			if fn.ReturnType != nil {
				table, err = qc.GetTable(&ast.TableName{
					Catalog: fn.ReturnType.Catalog,
					Schema:  fn.ReturnType.Schema,
					Name:    fn.ReturnType.Name,
				})
			}
			if table == nil || err != nil {
				if n.Alias != nil && len(n.Alias.Colnames.Items) > 0 {
					table = &Table{}
					for _, colName := range n.Alias.Colnames.Items {
						table.Columns = append(table.Columns, &Column{
							Name:     colName.(*ast.String).Str,
							DataType: "any",
						})
					}
				} else {
					colName := fn.Rel.Name
					if n.Alias != nil {
						colName = *n.Alias.Aliasname
					}
					table = &Table{
						Rel: &ast.TableName{
							Catalog: fn.Rel.Catalog,
							Schema:  fn.Rel.Schema,
							Name:    fn.Rel.Name,
						},
					}
					if len(fn.Outs) > 0 {
						for _, arg := range fn.Outs {
							table.Columns = append(table.Columns, &Column{
								Name:     arg.Name,
								DataType: arg.Type.Name,
							})
						}
					}
					if fn.ReturnType != nil {
						table.Columns = []*Column{
							{
								Name:     colName,
								DataType: fn.ReturnType.Name,
							},
						}
					}
				}
			}
			if n.Alias != nil {
				table.Rel = &ast.TableName{
					Name: *n.Alias.Aliasname,
				}
			}
			tables = append(tables, table)

		case *ast.RangeSubselect:
			cols, err := c.outputColumns(qc, n.Subquery)
			if err != nil {
				return nil, err
			}

			var tableName string
			if n.Alias != nil {
				tableName = *n.Alias.Aliasname
			}

			tables = append(tables, &Table{
				Rel: &ast.TableName{
					Name: tableName,
				},
				Columns: cols,
			})

		case *ast.RangeVar:
			fqn, err := ParseTableName(n)
			if err != nil {
				return nil, err
			}
			if qc == nil {
				return nil, fmt.Errorf("query catalog is empty")
			}
			table, cerr := qc.GetTable(fqn)
			if cerr != nil {
				// TODO: Update error location
				// cerr.Location = n.Location
				// return nil, *cerr
				return nil, cerr
			}
			if n.Alias != nil {
				table.Rel = &ast.TableName{
					Catalog: table.Rel.Catalog,
					Schema:  table.Rel.Schema,
					Name:    *n.Alias.Aliasname,
				}
			}
			tables = append(tables, table)

		default:
			return nil, fmt.Errorf("sourceTable: unsupported list item type: %T", n)
		}
	}
	return tables, nil
}

func outputColumnRefs(res *ast.ResTarget, tables []*Table, node *ast.ColumnRef) ([]*Column, error) {
	parts := stringSlice(node.Fields)
	var schema, name, alias string
	switch {
	case len(parts) == 1:
		name = parts[0]
	case len(parts) == 2:
		alias = parts[0]
		name = parts[1]
	case len(parts) == 3:
		schema = parts[0]
		alias = parts[1]
		name = parts[2]
	default:
		return nil, fmt.Errorf("unknown number of fields: %d", len(parts))
	}
	var cols []*Column
	var found int
	for _, t := range tables {
		if schema != "" && t.Rel.Schema != schema {
			continue
		}
		if alias != "" && t.Rel.Name != alias {
			continue
		}
		for _, c := range t.Columns {

			if c.Name == name {
				found += 1
				cname := c.Name
				if res.Name != nil {
					cname = *res.Name
				}
				cols = append(cols, &Column{
					Name:         cname,
					Type:         c.Type,
					Table:        c.Table,
					TableAlias:   alias,
					DataType:     c.DataType,
					NotNull:      c.NotNull,
					Unsigned:     c.Unsigned,
					IsArray:      c.IsArray,
					ArrayDims:    c.ArrayDims,
					Length:       c.Length,
					EmbedTable:   c.EmbedTable,
					OriginalName: c.Name,
				})
			}
		}
	}
	if found == 0 {
		return nil, &sqlerr.Error{
			Code:     "42703",
			Message:  fmt.Sprintf("column %q does not exist", name),
			Location: res.Location,
		}
	}
	if found > 1 {
		return nil, &sqlerr.Error{
			Code:     "42703",
			Message:  fmt.Sprintf("column reference %q is ambiguous", name),
			Location: res.Location,
		}
	}
	return cols, nil
}

func findColumnForNode(item ast.Node, tables []*Table, targetList *ast.List) error {
	ref, ok := item.(*ast.ColumnRef)
	if !ok {
		return nil
	}
	return findColumnForRef(ref, tables, targetList)
}

func findColumnForRef(ref *ast.ColumnRef, tables []*Table, targetList *ast.List) error {
	parts := stringSlice(ref.Fields)
	var alias, name string
	if len(parts) == 1 {
		name = parts[0]
	} else if len(parts) == 2 {
		alias = parts[0]
		name = parts[1]
	}

	var found int
	for _, t := range tables {
		if alias != "" && t.Rel.Name != alias {
			continue
		}

		// Find matching column
		for _, c := range t.Columns {
			if c.Name == name {
				found++
				break
			}
		}
	}

	// Find matching alias if necessary
	if found == 0 {
		for _, c := range targetList.Items {
			resTarget, ok := c.(*ast.ResTarget)
			if !ok {
				continue
			}
			if resTarget.Name != nil && *resTarget.Name == name {
				found++
			}
		}
	}

	if found == 0 {
		return &sqlerr.Error{
			Code:     "42703",
			Message:  fmt.Sprintf("column reference %q not found", name),
			Location: ref.Location,
		}
	}
	if found > 1 {
		return &sqlerr.Error{
			Code:     "42703",
			Message:  fmt.Sprintf("column reference %q is ambiguous", name),
			Location: ref.Location,
		}
	}

	return nil
}
