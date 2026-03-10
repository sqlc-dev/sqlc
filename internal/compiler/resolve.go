package compiler

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
	"github.com/sqlc-dev/sqlc/internal/sql/named"
	"github.com/sqlc-dev/sqlc/internal/sql/rewrite"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

func dataType(n *ast.TypeName) string {
	if n.Schema != "" {
		return n.Schema + "." + n.Name
	} else {
		return n.Name
	}
}

func hasConcreteParamType(col *Column) bool {
	return col != nil && col.DataType != "" && col.DataType != "any"
}

func (comp *Compiler) paramTypeString(col *Column) string {
	if !hasConcreteParamType(col) {
		return "any"
	}

	arraySuffix := strings.Repeat("[]", col.ArrayDims)
	if col.Type != nil && col.Type.Name != "" {
		return comp.parser.TypeName(col.Type.Schema, col.Type.Name) + arraySuffix
	}

	if rel, err := ParseRelationString(col.DataType); err == nil && rel.Catalog == "" {
		return comp.parser.TypeName(rel.Schema, rel.Name) + arraySuffix
	}

	return col.DataType + arraySuffix
}

func compatibleParamTypes(a, b *Column) bool {
	if !hasConcreteParamType(a) || !hasConcreteParamType(b) {
		return true
	}
	return a.DataType == b.DataType &&
		a.Unsigned == b.Unsigned &&
		a.IsArray == b.IsArray &&
		a.ArrayDims == b.ArrayDims
}

func sameTypeName(a, b *ast.TypeName) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return a.Catalog == b.Catalog && a.Schema == b.Schema && a.Name == b.Name
}

func matchingFuncCallOverloads(c *catalog.Catalog, call *ast.FuncCall) []catalog.Function {
	funs, err := c.ListFuncsByName(call.Func)
	if err != nil {
		return nil
	}

	var positional []ast.Node
	var named []*ast.NamedArgExpr
	if call.Args != nil {
		for _, arg := range call.Args.Items {
			if narg, ok := arg.(*ast.NamedArgExpr); ok {
				named = append(named, narg)
				continue
			}
			if len(named) > 0 {
				return nil
			}
			positional = append(positional, arg)
		}
	}

	var matches []catalog.Function
	for _, fun := range funs {
		args := fun.InArgs()
		var defaults int
		var variadic bool
		known := map[string]struct{}{}
		for _, arg := range args {
			if arg.HasDefault {
				defaults += 1
			}
			if arg.Mode == ast.FuncParamVariadic {
				variadic = true
				defaults += 1
			}
			if arg.Name != "" {
				known[arg.Name] = struct{}{}
			}
		}

		argc := len(named) + len(positional)
		if variadic {
			if argc < (len(args) - defaults) {
				continue
			}
		} else {
			if argc > len(args) || argc < (len(args)-defaults) {
				continue
			}
		}

		var unknownArgName bool
		for _, expr := range named {
			if expr.Name != nil {
				if _, found := known[*expr.Name]; !found {
					unknownArgName = true
				}
			}
		}
		if unknownArgName {
			continue
		}

		matches = append(matches, fun)
	}

	return matches
}

func stableFuncCallArgType(c *catalog.Catalog, call *ast.FuncCall, argIndex int, argName string) *ast.TypeName {
	var stable *ast.TypeName
	var seen bool

	for _, fun := range matchingFuncCallOverloads(c, call) {
		args := fun.InArgs()
		var current *ast.TypeName
		if argName == "" {
			if argIndex >= len(args) {
				return nil
			}
			current = args[argIndex].Type
		} else {
			for _, arg := range args {
				if arg.Name == argName {
					current = arg.Type
					break
				}
			}
			if current == nil {
				return nil
			}
		}

		if !seen {
			stable = current
			seen = true
			continue
		}
		if !sameTypeName(stable, current) {
			return nil
		}
	}

	return stable
}

func resolvedFuncCallArgType(fun *catalog.Function, argIndex int, argName string) *ast.TypeName {
	if fun == nil {
		return nil
	}
	if argName == "" {
		if argIndex < len(fun.Args) {
			return fun.Args[argIndex].Type
		}
		return nil
	}
	for _, arg := range fun.Args {
		if arg.Name == argName {
			return arg.Type
		}
	}
	return nil
}

func mergeResolvedParam(existing, incoming Parameter) Parameter {
	if existing.Column == nil {
		return incoming
	}
	if incoming.Column == nil {
		return existing
	}

	base := existing
	other := incoming
	if hasConcreteParamType(incoming.Column) && !hasConcreteParamType(existing.Column) {
		base = incoming
		other = existing
	}

	col := *base.Column
	if col.Name == "" {
		col.Name = other.Column.Name
	}
	if col.OriginalName == "" {
		col.OriginalName = other.Column.OriginalName
	}
	if col.Table == nil {
		col.Table = other.Column.Table
	}
	if col.Type == nil {
		col.Type = other.Column.Type
	}
	if col.Length == nil {
		col.Length = other.Column.Length
	}
	col.IsNamedParam = col.IsNamedParam || other.Column.IsNamedParam
	col.IsSqlcSlice = col.IsSqlcSlice || other.Column.IsSqlcSlice

	base.Column = &col
	return base
}

func (comp *Compiler) incompatibleParamRefError(ref paramRef, existing, incoming Parameter) error {
	return &sqlerr.Error{
		Code: "42P08",
		Message: fmt.Sprintf(
			"parameter $%d has incompatible types: %s, %s",
			ref.ref.Number,
			comp.paramTypeString(existing.Column),
			comp.paramTypeString(incoming.Column),
		),
		Location: ref.ref.Location,
	}
}

func (comp *Compiler) resolveCatalogRefs(qc *QueryCatalog, rvs []*ast.RangeVar, args []paramRef, params *named.ParamSet, embeds rewrite.EmbedSet) ([]Parameter, error) {
	c := comp.catalog

	aliasMap := map[string]*ast.TableName{}
	// TODO: Deprecate defaultTable
	var defaultTable *ast.TableName
	var tables []*ast.TableName

	typeMap := map[string]map[string]map[string]*catalog.Column{}
	indexTable := func(table catalog.Table) error {
		tables = append(tables, table.Rel)
		if defaultTable == nil {
			defaultTable = table.Rel
		}
		schema := table.Rel.Schema
		if schema == "" {
			schema = c.DefaultSchema
		}
		if _, exists := typeMap[schema]; !exists {
			typeMap[schema] = map[string]map[string]*catalog.Column{}
		}
		typeMap[schema][table.Rel.Name] = map[string]*catalog.Column{}
		for _, c := range table.Columns {
			cc := c
			typeMap[schema][table.Rel.Name][c.Name] = cc
		}
		return nil
	}

	for _, rv := range rvs {
		if rv.Relname == nil {
			continue
		}
		fqn, err := ParseTableName(rv)
		if err != nil {
			return nil, err
		}
		if _, found := aliasMap[fqn.Name]; found {
			continue
		}
		table, err := c.GetTable(fqn)
		if err != nil {
			if qc == nil {
				continue
			}
			// If the table name doesn't exist, first check if it's a CTE
			if _, qcerr := qc.GetTable(fqn); qcerr != nil {
				return nil, err
			}
			continue
		}
		err = indexTable(table)
		if err != nil {
			return nil, err
		}
		if rv.Alias != nil {
			aliasMap[*rv.Alias.Aliasname] = fqn
		}
	}

	// resolve a table for an embed
	for _, embed := range embeds {
		table, err := c.GetTable(embed.Table)
		if err == nil {
			embed.Table = table.Rel
			continue
		}

		if alias, ok := aliasMap[embed.Table.Name]; ok {
			embed.Table = alias
			continue
		}

		return nil, fmt.Errorf("unable to resolve table with %q: %w", embed.Orig(), err)
	}

	var a []Parameter
	seen := map[int]int{}
	paramCounts := map[int]int{}
	for _, ref := range args {
		paramCounts[ref.ref.Number] += 1
	}

	addParam := func(ref paramRef, p Parameter) error {
		if idx, ok := seen[p.Number]; ok {
			if !compatibleParamTypes(a[idx].Column, p.Column) {
				return comp.incompatibleParamRefError(ref, a[idx], p)
			}
			a[idx] = mergeResolvedParam(a[idx], p)
			return nil
		}
		seen[p.Number] = len(a)
		a = append(a, p)
		return nil
	}

	addUnknownParam := func(ref paramRef) error {
		defaultP := named.NewInferredParam(ref.name, false)
		p, isNamed := params.FetchMerge(ref.ref.Number, defaultP)
		return addParam(ref, Parameter{
			Number: ref.ref.Number,
			Column: &Column{
				Name:         p.Name(),
				DataType:     "any",
				IsNamedParam: isNamed,
			},
		})
	}

	addColumnParam := func(ref paramRef, key string, location int) error {
		var schema, rel string
		// TODO: Deprecate defaultTable
		if defaultTable != nil {
			schema = defaultTable.Schema
			rel = defaultTable.Name
		}
		if ref.rv != nil {
			fqn, err := ParseTableName(ref.rv)
			if err != nil {
				return err
			}
			schema = fqn.Schema
			rel = fqn.Name
		}
		if schema == "" {
			schema = c.DefaultSchema
		}

		tableMap, ok := typeMap[schema][rel]
		if !ok {
			return sqlerr.RelationNotFound(rel)
		}

		if c, ok := tableMap[key]; ok {
			defaultP := named.NewInferredParam(key, c.IsNotNull)
			p, isNamed := params.FetchMerge(ref.ref.Number, defaultP)
			return addParam(ref, Parameter{
				Number: ref.ref.Number,
				Column: &Column{
					Name:         p.Name(),
					OriginalName: c.Name,
					DataType:     dataType(&c.Type),
					NotNull:      p.NotNull(),
					Unsigned:     c.IsUnsigned,
					IsArray:      c.IsArray,
					ArrayDims:    c.ArrayDims,
					Table:        &ast.TableName{Schema: schema, Name: rel},
					Length:       c.Length,
					IsNamedParam: isNamed,
					IsSqlcSlice:  p.IsSqlcSlice(),
				},
			})
		}

		return &sqlerr.Error{
			Code:     "42703",
			Message:  fmt.Sprintf("column %q does not exist", key),
			Location: location,
		}
	}

	for _, ref := range args {
		switch n := ref.parent.(type) {

		case *limitOffset:
			defaultP := named.NewInferredParam("offset", true)
			p, isNamed := params.FetchMerge(ref.ref.Number, defaultP)
			if err := addParam(ref, Parameter{
				Number: ref.ref.Number,
				Column: &Column{
					Name:         p.Name(),
					DataType:     "integer",
					NotNull:      p.NotNull(),
					IsNamedParam: isNamed,
				},
			}); err != nil {
				return nil, err
			}

		case *limitCount:
			defaultP := named.NewInferredParam("limit", true)
			p, isNamed := params.FetchMerge(ref.ref.Number, defaultP)
			if err := addParam(ref, Parameter{
				Number: ref.ref.Number,
				Column: &Column{
					Name:         p.Name(),
					DataType:     "integer",
					NotNull:      p.NotNull(),
					IsNamedParam: isNamed,
				},
			}); err != nil {
				return nil, err
			}

		case *ast.A_Expr:
			// TODO: While this works for a wide range of simple expressions,
			// more complicated expressions will cause this logic to fail.
			list := astutils.Search(n.Lexpr, func(node ast.Node) bool {
				_, ok := node.(*ast.ColumnRef)
				return ok
			})
			if len(list.Items) == 0 {
				list = astutils.Search(n.Rexpr, func(node ast.Node) bool {
					_, ok := node.(*ast.ColumnRef)
					return ok
				})
			}

			if len(list.Items) == 0 {
				// TODO: Move this to database-specific engine package
				dataType := "any"
				if astutils.Join(n.Name, ".") == "||" {
					dataType = "text"
				}

				defaultP := named.NewParam("")
				p, isNamed := params.FetchMerge(ref.ref.Number, defaultP)
				if err := addParam(ref, Parameter{
					Number: ref.ref.Number,
					Column: &Column{
						Name:         p.Name(),
						DataType:     dataType,
						IsNamedParam: isNamed,
						NotNull:      p.NotNull(),
						IsSqlcSlice:  p.IsSqlcSlice(),
					},
				}); err != nil {
					return nil, err
				}
				continue
			}

			switch node := list.Items[0].(type) {
			case *ast.ColumnRef:
				items := stringSlice(node.Fields)
				var key, alias string
				switch len(items) {
				case 1:
					key = items[0]
				case 2:
					alias = items[0]
					key = items[1]
				case 3:
					// schema := items[0]
					alias = items[1]
					key = items[2]
				default:
					panic("too many field items: " + strconv.Itoa(len(items)))
				}

				search := tables
				if alias == "" && ref.rv != nil {
					fqn, err := ParseTableName(ref.rv)
					if err != nil {
						return nil, err
					}
					search = []*ast.TableName{fqn}
				} else if alias != "" {
					if original, ok := aliasMap[alias]; ok {
						search = []*ast.TableName{original}
					} else {
						var located bool
						for _, fqn := range tables {
							if fqn.Name == alias {
								located = true
								search = []*ast.TableName{fqn}
							}
						}
						if !located {
							return nil, &sqlerr.Error{
								Code:     "42703",
								Message:  fmt.Sprintf("table alias %q does not exist", alias),
								Location: node.Location,
							}
						}
					}
				}

				var found int
				for _, table := range search {
					schema := table.Schema
					if schema == "" {
						schema = c.DefaultSchema
					}
					if c, ok := typeMap[schema][table.Name][key]; ok {
						found += 1
						if ref.name != "" {
							key = ref.name
						}

						defaultP := named.NewInferredParam(key, c.IsNotNull)
						p, isNamed := params.FetchMerge(ref.ref.Number, defaultP)
						if err := addParam(ref, Parameter{
							Number: ref.ref.Number,
							Column: &Column{
								Name:         p.Name(),
								OriginalName: c.Name,
								DataType:     dataType(&c.Type),
								NotNull:      p.NotNull(),
								Unsigned:     c.IsUnsigned,
								IsArray:      c.IsArray,
								ArrayDims:    c.ArrayDims,
								Length:       c.Length,
								Table:        table,
								IsNamedParam: isNamed,
								IsSqlcSlice:  p.IsSqlcSlice(),
							},
						}); err != nil {
							return nil, err
						}
					}
				}

				if found == 0 {
					return nil, &sqlerr.Error{
						Code:     "42703",
						Message:  fmt.Sprintf("column %q does not exist", key),
						Location: node.Location,
					}
				}
				if found > 1 {
					return nil, &sqlerr.Error{
						Code:     "42703",
						Message:  fmt.Sprintf("column reference %q is ambiguous", key),
						Location: node.Location,
					}
				}
			}

		case *ast.BetweenExpr:
			if n == nil || n.Expr == nil || n.Left == nil || n.Right == nil {
				fmt.Println("ast.BetweenExpr is nil")
				continue
			}

			var key string
			if ref, ok := n.Expr.(*ast.ColumnRef); ok {
				itemsCount := len(ref.Fields.Items)
				if str, ok := ref.Fields.Items[itemsCount-1].(*ast.String); ok {
					key = str.Str
				}
			}

			for _, table := range tables {
				schema := table.Schema
				if schema == "" {
					schema = c.DefaultSchema
				}

				if c, ok := typeMap[schema][table.Name][key]; ok {
					defaultP := named.NewInferredParam(key, c.IsNotNull)
					p, isNamed := params.FetchMerge(ref.ref.Number, defaultP)
					var namePrefix string
					if !isNamed {
						switch ref.ref {
						case n.Left:
							namePrefix = "from_"
						case n.Right:
							namePrefix = "to_"
						}
					}

					if err := addParam(ref, Parameter{
						Number: ref.ref.Number,
						Column: &Column{
							Name:         namePrefix + p.Name(),
							DataType:     dataType(&c.Type),
							NotNull:      p.NotNull(),
							Unsigned:     c.IsUnsigned,
							IsArray:      c.IsArray,
							ArrayDims:    c.ArrayDims,
							Table:        table,
							IsNamedParam: isNamed,
							IsSqlcSlice:  p.IsSqlcSlice(),
						},
					}); err != nil {
						return nil, err
					}
				}
			}

		case *ast.FuncCall:
			fun, resolveErr := c.ResolveFuncCall(n)
			if resolveErr != nil {
				// Synthesize a function on the fly to avoid returning with an error
				// for an unknown Postgres function (e.g. defined in an extension)
				var args []*catalog.Argument
				for range n.Args.Items {
					args = append(args, &catalog.Argument{
						Type: &ast.TypeName{Name: "any"},
					})
				}
				fun = &catalog.Function{
					Name:       n.Func.Name,
					Args:       args,
					ReturnType: &ast.TypeName{Name: "any"},
				}
			}

			var added bool
			for i, item := range n.Args.Items {
				funcName := fun.Name
				var argName string
				switch inode := item.(type) {
				case *ast.ParamRef:
					if inode.Number != ref.ref.Number {
						continue
					}
				case *ast.TypeCast:
					pr, ok := inode.Arg.(*ast.ParamRef)
					if !ok {
						continue
					}
					if pr.Number != ref.ref.Number {
						continue
					}
				case *ast.NamedArgExpr:
					pr, ok := inode.Arg.(*ast.ParamRef)
					if !ok {
						continue
					}
					if pr.Number != ref.ref.Number {
						continue
					}
					if inode.Name != nil {
						argName = *inode.Name
					}
				default:
					continue
				}

				if fun.Args == nil {
					defaultName := funcName
					if argName != "" {
						defaultName = argName
					}

					defaultP := named.NewInferredParam(defaultName, false)
					p, isNamed := params.FetchMerge(ref.ref.Number, defaultP)
					added = true
					if err := addParam(ref, Parameter{
						Number: ref.ref.Number,
						Column: &Column{
							Name:         p.Name(),
							DataType:     "any",
							IsNamedParam: isNamed,
							NotNull:      p.NotNull(),
							IsSqlcSlice:  p.IsSqlcSlice(),
						},
					}); err != nil {
						return nil, err
					}
					continue
				}

				var paramName string
				var paramType *ast.TypeName

				if argName == "" {
					if i < len(fun.Args) {
						paramName = fun.Args[i].Name
					}
				} else {
					paramName = argName
				}
				if paramName == "" {
					paramName = funcName
				}
				if resolveErr == nil {
					if paramCounts[ref.ref.Number] > 1 {
						paramType = stableFuncCallArgType(c, n, i, argName)
					} else {
						paramType = resolvedFuncCallArgType(fun, i, argName)
					}
				}
				if paramType == nil {
					paramType = &ast.TypeName{Name: ""}
				}

				defaultP := named.NewInferredParam(paramName, true)
				p, isNamed := params.FetchMerge(ref.ref.Number, defaultP)
				added = true
				if err := addParam(ref, Parameter{
					Number: ref.ref.Number,
					Column: &Column{
						Name:         p.Name(),
						DataType:     dataType(paramType),
						NotNull:      p.NotNull(),
						IsNamedParam: isNamed,
						IsSqlcSlice:  p.IsSqlcSlice(),
					},
				}); err != nil {
					return nil, err
				}
			}

			if fun.ReturnType == nil {
				if !added {
					if err := addUnknownParam(ref); err != nil {
						return nil, err
					}
				}
				continue
			}

			table, err := c.GetTable(&ast.TableName{
				Catalog: fun.ReturnType.Catalog,
				Schema:  fun.ReturnType.Schema,
				Name:    fun.ReturnType.Name,
			})
			if err != nil {
				if !added {
					if err := addUnknownParam(ref); err != nil {
						return nil, err
					}
				}
				continue
			}
			err = indexTable(table)
			if err != nil {
				return nil, err
			}

		case *ast.ResTarget:
			if n.Name == nil {
				return nil, fmt.Errorf("*ast.ResTarget has nil name")
			}
			if err := addColumnParam(ref, *n.Name, n.Location); err != nil {
				return nil, err
			}

		case *ast.String:
			if err := addColumnParam(ref, n.Str, n.Pos()); err != nil {
				return nil, err
			}

		case *ast.TypeCast:
			if n.TypeName == nil {
				return nil, fmt.Errorf("*ast.TypeCast has nil type name")
			}
			col := toColumn(n.TypeName)
			defaultP := named.NewInferredParam(col.Name, col.NotNull)
			p, _ := params.FetchMerge(ref.ref.Number, defaultP)

			col.Name = p.Name()
			col.NotNull = p.NotNull()
			if err := addParam(ref, Parameter{
				Number: ref.ref.Number,
				Column: col,
			}); err != nil {
				return nil, err
			}

		case *ast.ParamRef:
			if err := addParam(ref, Parameter{Number: ref.ref.Number}); err != nil {
				return nil, err
			}

		case *ast.In:
			if n == nil || n.List == nil {
				fmt.Println("ast.In is nil")
				continue
			}

			location := 0
			var key, alias string
			var items []string

			if left, ok := n.Expr.(*ast.ColumnRef); ok {
				location = left.Location
				items = stringSlice(left.Fields)
			} else if left, ok := n.Expr.(*ast.ParamRef); ok {
				if len(n.List) <= 0 {
					continue
				}
				if right, ok := n.List[0].(*ast.ColumnRef); ok {
					location = left.Location
					items = stringSlice(right.Fields)
				} else {
					continue
				}
			} else {
				continue
			}

			switch len(items) {
			case 1:
				key = items[0]
			case 2:
				alias = items[0]
				key = items[1]
			default:
				panic("too many field items: " + strconv.Itoa(len(items)))
			}

			var found int
			if n.Sel == nil {
				search := tables
				if alias != "" {
					if original, ok := aliasMap[alias]; ok {
						search = []*ast.TableName{original}
					} else {
						for _, fqn := range tables {
							if fqn.Name == alias {
								search = []*ast.TableName{fqn}
							}
						}
					}
				}

				for _, table := range search {
					schema := table.Schema
					if schema == "" {
						schema = c.DefaultSchema
					}
					if c, ok := typeMap[schema][table.Name][key]; ok {
						found += 1
						if ref.name != "" {
							key = ref.name
						}
						defaultP := named.NewInferredParam(key, c.IsNotNull)
						p, isNamed := params.FetchMerge(ref.ref.Number, defaultP)
						if err := addParam(ref, Parameter{
							Number: ref.ref.Number,
							Column: &Column{
								Name:         p.Name(),
								OriginalName: c.Name,
								DataType:     dataType(&c.Type),
								NotNull:      c.IsNotNull,
								Unsigned:     c.IsUnsigned,
								IsArray:      c.IsArray,
								ArrayDims:    c.ArrayDims,
								Table:        table,
								IsNamedParam: isNamed,
								IsSqlcSlice:  p.IsSqlcSlice(),
							},
						}); err != nil {
							return nil, err
						}
					}
				}
			}

			if found == 0 {
				return nil, &sqlerr.Error{
					Code:     "42703",
					Message:  fmt.Sprintf("396: column %q does not exist", key),
					Location: location,
				}
			}
			if found > 1 {
				return nil, &sqlerr.Error{
					Code:     "42703",
					Message:  fmt.Sprintf("in same name column reference %q is ambiguous", key),
					Location: location,
				}
			}

		default:
			slog.Debug("unsupported reference type", "type", fmt.Sprintf("%T", n))
			if err := addUnknownParam(ref); err != nil {
				return nil, err
			}
		}
	}
	return a, nil
}
