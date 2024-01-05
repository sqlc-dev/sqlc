package compiler

import (
	"fmt"
	"log/slog"
	"strconv"

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

	addUnknownParam := func(ref paramRef) {
		defaultP := named.NewInferredParam(ref.name, false)
		p, isNamed := params.FetchMerge(ref.ref.Number, defaultP)
		a = append(a, Parameter{
			Number: ref.ref.Number,
			Column: &Column{
				Name:         p.Name(),
				DataType:     "any",
				IsNamedParam: isNamed,
			},
		})
	}

	for _, ref := range args {
		switch n := ref.parent.(type) {

		case *limitOffset:
			defaultP := named.NewInferredParam("offset", true)
			p, isNamed := params.FetchMerge(ref.ref.Number, defaultP)
			a = append(a, Parameter{
				Number: ref.ref.Number,
				Column: &Column{
					Name:         p.Name(),
					DataType:     "integer",
					NotNull:      p.NotNull(),
					IsNamedParam: isNamed,
				},
			})

		case *limitCount:
			defaultP := named.NewInferredParam("limit", true)
			p, isNamed := params.FetchMerge(ref.ref.Number, defaultP)
			a = append(a, Parameter{
				Number: ref.ref.Number,
				Column: &Column{
					Name:         p.Name(),
					DataType:     "integer",
					NotNull:      p.NotNull(),
					IsNamedParam: isNamed,
				},
			})

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
				a = append(a, Parameter{
					Number: ref.ref.Number,
					Column: &Column{
						Name:         p.Name(),
						DataType:     dataType,
						IsNamedParam: isNamed,
						NotNull:      p.NotNull(),
						IsSqlcSlice:  p.IsSqlcSlice(),
					},
				})
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
				if alias != "" {
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
						a = append(a, Parameter{
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
						})
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
						if ref.ref == n.Left {
							namePrefix = "from_"
						} else if ref.ref == n.Right {
							namePrefix = "to_"
						}
					}

					a = append(a, Parameter{
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
					})
				}
			}

		case *ast.FuncCall:
			fun, err := c.ResolveFuncCall(n)
			if err != nil {
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
					a = append(a, Parameter{
						Number: ref.ref.Number,
						Column: &Column{
							Name:         p.Name(),
							DataType:     "any",
							IsNamedParam: isNamed,
							NotNull:      p.NotNull(),
							IsSqlcSlice:  p.IsSqlcSlice(),
						},
					})
					continue
				}

				var paramName string
				var paramType *ast.TypeName

				if argName == "" {
					if i < len(fun.Args) {
						paramName = fun.Args[i].Name
						paramType = fun.Args[i].Type
					}
				} else {
					paramName = argName
					for _, arg := range fun.Args {
						if arg.Name == argName {
							paramType = arg.Type
						}
					}
					if paramType == nil {
						panic(fmt.Sprintf("named argument %s has no type", paramName))
					}
				}
				if paramName == "" {
					paramName = funcName
				}
				if paramType == nil {
					paramType = &ast.TypeName{Name: ""}
				}

				defaultP := named.NewInferredParam(paramName, true)
				p, isNamed := params.FetchMerge(ref.ref.Number, defaultP)
				added = true
				a = append(a, Parameter{
					Number: ref.ref.Number,
					Column: &Column{
						Name:         p.Name(),
						DataType:     dataType(paramType),
						NotNull:      p.NotNull(),
						IsNamedParam: isNamed,
						IsSqlcSlice:  p.IsSqlcSlice(),
					},
				})
			}

			if fun.ReturnType == nil {
				if !added {
					addUnknownParam(ref)
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
					addUnknownParam(ref)
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
			key := *n.Name

			var schema, rel string
			// TODO: Deprecate defaultTable
			if defaultTable != nil {
				schema = defaultTable.Schema
				rel = defaultTable.Name
			}
			if ref.rv != nil {
				fqn, err := ParseTableName(ref.rv)
				if err != nil {
					return nil, err
				}
				schema = fqn.Schema
				rel = fqn.Name
			}
			if schema == "" {
				schema = c.DefaultSchema
			}

			tableMap, ok := typeMap[schema][rel]
			if !ok {
				return nil, sqlerr.RelationNotFound(rel)
			}

			if c, ok := tableMap[key]; ok {
				defaultP := named.NewInferredParam(key, c.IsNotNull)
				p, isNamed := params.FetchMerge(ref.ref.Number, defaultP)
				a = append(a, Parameter{
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
			} else {
				return nil, &sqlerr.Error{
					Code:     "42703",
					Message:  fmt.Sprintf("column %q does not exist", key),
					Location: n.Location,
				}
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
			a = append(a, Parameter{
				Number: ref.ref.Number,
				Column: col,
			})

		case *ast.ParamRef:
			a = append(a, Parameter{Number: ref.ref.Number})

		case *ast.In:
			if n == nil || n.List == nil {
				fmt.Println("ast.In is nil")
				continue
			}

			number := 0
			if pr, ok := n.List[0].(*ast.ParamRef); ok {
				number = pr.Number
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
						a = append(a, Parameter{
							Number: number,
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
						})
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
			addUnknownParam(ref)
		}
	}
	return a, nil
}
