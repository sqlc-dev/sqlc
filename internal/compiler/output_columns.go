package compiler

import (
	"errors"
	"fmt"

	"github.com/kyleconroy/sqlc/internal/sql/catalog"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/astutils"
	"github.com/kyleconroy/sqlc/internal/sql/lang"
	"github.com/kyleconroy/sqlc/internal/sql/sqlerr"
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

	var targets *ast.List
	switch n := node.(type) {
	case *ast.DeleteStmt:
		targets = n.ReturningList
	case *ast.InsertStmt:
		targets = n.ReturningList
	case *ast.SelectStmt:
		targets = n.TargetList

		if n.GroupClause != nil {
			for _, item := range n.GroupClause.Items {
				if err := findColumnForNode(item, tables, n); err != nil {
					return nil, err
				}
			}
		}
		validateOrderBy := true
		if c.conf.StrictOrderBy != nil {
			validateOrderBy = *c.conf.StrictOrderBy
		}
		if validateOrderBy {
			if n.SortClause != nil {
				for _, item := range n.SortClause.Items {
					sb, ok := item.(*ast.SortBy)
					if !ok {
						continue
					}
					if err := findColumnForNode(sb.Node, tables, n); err != nil {
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
						if err := findColumnForNode(caseExpr.Xpr, tables, n); err != nil {
							return nil, fmt.Errorf("%v: if you want to skip this validation, set 'strict_order_by' to false", err)
						}
					}
				}
			}
		}

		// For UNION queries, targets is empty and we need to look for the
		// columns in Largs.
		if len(targets.Items) == 0 && n.Larg != nil {
			return c.outputColumns(qc, n.Larg)
		}
	case *ast.CallStmt:
		targets = &ast.List{}
	case *ast.TruncateStmt, *ast.RefreshMatViewStmt, *ast.NotifyStmt, *ast.ListenStmt:
		targets = &ast.List{}
	case *ast.UpdateStmt:
		targets = n.ReturningList
	default:
		return nil, fmt.Errorf("outputColumns: unsupported node type: %T", n)
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
			switch {
			case lang.IsComparisonOperator(astutils.Join(n.Name, "")):
				// TODO: Generate a name for these operations
				cols = append(cols, &Column{Name: name, DataType: "bool", NotNull: true})
			case lang.IsMathematicalOperator(astutils.Join(n.Name, "")):
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
			if n.Boolop == ast.BoolExprTypeNot && len(n.Args.Items) == 1 {
				sublink, ok := n.Args.Items[0].(*ast.SubLink)
				if ok && sublink.SubLinkType == ast.EXISTS_SUBLINK {
					notNull = true
					if name == "" {
						name = "not_exists"
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
				if res := isTableRequired(f, col, tableRequired); res != tableNotFound {
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
		if n.Alias == nil && *n.Relname == col.Table.Name {
			return prior
		}
		if n.Alias != nil && *n.Alias.Aliasname == col.TableAlias && *n.Relname == col.Table.Name {
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

// Compute the output columns for a statement.
//
// Return an error if column references are ambiguous
// Return an error if column references don't exist
// Return an error if a table is referenced twice
// Return an error if an unknown column is referenced
func (c *Compiler) sourceTables(qc *QueryCatalog, node ast.Node) ([]*Table, error) {
	var list *ast.List
	switch n := node.(type) {
	case *ast.DeleteStmt:
		list = n.Relations
	case *ast.InsertStmt:
		list = &ast.List{
			Items: []ast.Node{n.Relation},
		}
	case *ast.SelectStmt:
		list = astutils.Search(n.FromClause, func(node ast.Node) bool {
			switch node.(type) {
			case *ast.RangeVar, *ast.RangeSubselect, *ast.RangeFunction:
				return true
			default:
				return false
			}
		})
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
		list = &ast.List{
			Items: append(n.FromClause.Items, n.Relations.Items...),
		}
	case *ast.CallStmt:
		list = &ast.List{}
	case *ast.NotifyStmt, *ast.ListenStmt:
		list = &ast.List{}
	default:
		return nil, fmt.Errorf("sourceTables: unsupported node type: %T", n)
	}

	var tables []*Table
	for _, item := range list.Items {
		switch n := item.(type) {

		case *ast.RangeFunction:
			// If the function or table can't be found, don't error out.  There
			// are many queries that depend on functions unknown to sqlc.
			var funcCall *ast.FuncCall
			switch f := n.Functions.Items[0].(type) {
			case *ast.List:
				funcCall = f.Items[0].(*ast.FuncCall)
			case *ast.FuncCall:
				funcCall = f
			default:
				return nil, fmt.Errorf("sourceTables: unsupported function call type %T", n.Functions.Items[0])
			}

			fn, err := qc.GetFunc(funcCall.Func)
			if err != nil {
				continue
			}
			table, err := qc.GetTable(&ast.TableName{
				Catalog: fn.ReturnType.Catalog,
				Schema:  fn.ReturnType.Schema,
				Name:    fn.ReturnType.Name,
			})
			if err != nil {
				if n.Alias == nil || len(n.Alias.Colnames.Items) == 0 {
					continue
				}

				table = &Table{}
				for _, colName := range n.Alias.Colnames.Items {
					table.Columns = append(table.Columns, &Column{
						Name:     colName.(*ast.String).Str,
						DataType: "any",
					})
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
			tables = append(tables, &Table{
				Rel: &ast.TableName{
					Name: *n.Alias.Aliasname,
				},
				Columns: cols,
			})

		case *ast.RangeVar:
			fqn, err := ParseTableName(n)
			if err != nil {
				return nil, err
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
	var name, alias string
	switch {
	case len(parts) == 1:
		name = parts[0]
	case len(parts) == 2:
		alias = parts[0]
		name = parts[1]
	default:
		return nil, fmt.Errorf("unknown number of fields: %d", len(parts))
	}
	var cols []*Column
	var found int
	for _, t := range tables {
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
			Message:  fmt.Sprintf("column \"%s\" does not exist", name),
			Location: res.Location,
		}
	}
	if found > 1 {
		return nil, &sqlerr.Error{
			Code:     "42703",
			Message:  fmt.Sprintf("column reference \"%s\" is ambiguous", name),
			Location: res.Location,
		}
	}
	return cols, nil
}

func findColumnForNode(item ast.Node, tables []*Table, n *ast.SelectStmt) error {
	ref, ok := item.(*ast.ColumnRef)
	if !ok {
		return nil
	}
	return findColumnForRef(ref, tables, n)
}

func findColumnForRef(ref *ast.ColumnRef, tables []*Table, selectStatement *ast.SelectStmt) error {
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
		var foundColumn bool
		for _, c := range t.Columns {
			if c.Name == name {
				found++
				foundColumn = true
			}
		}

		if foundColumn {
			continue
		}

		// Find matching alias
		for _, c := range selectStatement.TargetList.Items {
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
			Message:  fmt.Sprintf("column reference \"%s\" not found", name),
			Location: ref.Location,
		}
	}
	if found > 1 {
		return &sqlerr.Error{
			Code:     "42703",
			Message:  fmt.Sprintf("column reference \"%s\" is ambiguous", name),
			Location: ref.Location,
		}
	}

	return nil
}
