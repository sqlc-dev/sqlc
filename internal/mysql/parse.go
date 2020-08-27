package mysql

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"vitess.io/vitess/go/vt/sqlparser"

	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/metadata"
	"github.com/kyleconroy/sqlc/internal/migrations"
	"github.com/kyleconroy/sqlc/internal/multierr"
	"github.com/kyleconroy/sqlc/internal/sql/sqlpath"
)

// Query holds the data for walking and validating mysql querys
type Query struct {
	SQL              string // the string representation of the parsed query
	Columns          []Column
	Params           []*Param // "?" params in the query string
	Name             string   // the Go function name
	Cmd              string   // TODO: Pick a better name. One of: one, many, exec, execrows
	DefaultTableName string   // for columns that are not qualified

	Filename string
}

type Column struct {
	*sqlparser.ColumnDefinition
	Table string
}

func parsePath(sqlPath []string, generator PackageGenerator) (*Result, error) {
	files, err := sqlpath.Glob(sqlPath)
	if err != nil {
		return nil, err
	}

	parseErrors := multierr.New()
	parsedQueries := []*Query{}
	for _, filename := range files {
		blob, err := ioutil.ReadFile(filename)
		if err != nil {
			parseErrors.Add(filename, "", 0, err)
		}
		contents := migrations.RemoveRollbackStatements(string(blob))
		if err != nil {
			parseErrors.Add(filename, "", 0, err)
			continue
		}

		t := sqlparser.NewStringTokenizer(contents)
		var start int
		for {
			q, err := sqlparser.ParseNextStrictDDL(t)
			if err == io.EOF {
				break
			} else if err != nil {
				if posErr, ok := err.(sqlparser.PositionedErr); ok {
					message := fmt.Errorf(posErr.Err)
					if posErr.Near != nil {
						message = fmt.Errorf("%s at or near \"%s\"", posErr.Err, posErr.Near)
					}
					parseErrors.Add(filename, contents, posErr.Pos, message)
				} else {
					parseErrors.Add(filename, contents, start, err)
				}
				continue
			}
			query := contents[start : t.Position-1]
			result, err := generator.parseQueryString(q, query)
			if err != nil {
				parseErrors.Add(filename, contents, start, err)
				start = t.Position
				continue
			}
			start = t.Position
			if result == nil {
				continue
			}
			result.Filename = filepath.Base(filename)
			parsedQueries = append(parsedQueries, result)
		}
	}

	if len(parseErrors.Errs()) > 0 {
		return nil, parseErrors
	}

	return &Result{
		Queries:          parsedQueries,
		PackageGenerator: generator,
	}, nil
}

func (pGen PackageGenerator) parseQueryString(tree sqlparser.Statement, query string) (*Query, error) {
	var parsedQuery *Query
	switch tree := tree.(type) {
	case *sqlparser.Select:
		selectQuery, err := pGen.parseSelect(tree, query)
		if err != nil {
			return nil, err
		}
		parsedQuery = selectQuery
	case *sqlparser.Insert:
		insert, err := pGen.parseInsert(tree, query)
		if err != nil {
			return nil, err
		}
		parsedQuery = insert
	case *sqlparser.Update:
		update, err := pGen.parseUpdate(tree, query)
		if err != nil {
			return nil, err
		}
		parsedQuery = update
	case *sqlparser.Delete:
		delete, err := pGen.parseDelete(tree, query)
		if err != nil {
			return nil, err
		}
		parsedQuery = delete
	case *sqlparser.DDL:
		return nil, pGen.Schema.Add(tree)
	default:
		// panic("Unsupported SQL statement type")
		return nil, nil
	}
	paramsReplacedQuery, err := replaceParamStrs(sqlparser.String(tree), parsedQuery.Params)
	if err != nil {
		return nil, fmt.Errorf("failed to replace param variables in query string: %w", err)
	}
	parsedQuery.SQL = paramsReplacedQuery
	return parsedQuery, nil
}

func (q *Query) parseNameAndCmd() error {
	if q == nil {
		return fmt.Errorf("cannot parse name and cmd from null query")
	}
	_, comments := sqlparser.SplitMarginComments(q.SQL)
	name, cmd, err := metadata.Parse(comments.Leading, metadata.CommentSyntax{SlashStar: true})
	if err != nil {
		return err
	} else if name == "" || cmd == "" {
		return fmt.Errorf("failed to parse query leading comment")
	}
	q.Name = name
	q.Cmd = cmd
	return nil
}

func (pGen PackageGenerator) parseSelect(tree *sqlparser.Select, query string) (*Query, error) {
	tableAliasMap, defaultTableName, err := parseFrom(tree.From, false)
	if err != nil {
		return nil, fmt.Errorf("failed to parse table name alias's: %w", err)
	}

	// handle * expressions first by expanding all columns of the default table
	_, ok := tree.SelectExprs[0].(*sqlparser.StarExpr)
	if ok {
		colNames := []sqlparser.SelectExpr{}
		colDfns := pGen.Schema.tables[defaultTableName]
		for _, col := range colDfns {
			colNames = append(colNames, &sqlparser.AliasedExpr{
				Expr: &sqlparser.ColName{
					Name: col.Name,
				}},
			)
		}
		tree.SelectExprs = colNames
	}

	parsedQuery := Query{
		SQL:              query,
		DefaultTableName: defaultTableName,
	}
	cols, err := pGen.parseSelectAliasExpr(tree.SelectExprs, tableAliasMap, defaultTableName)
	if err != nil {
		return nil, err
	}
	parsedQuery.Columns = cols

	whereParams, err := pGen.paramsInWhereExpr(tree.Where, tableAliasMap, defaultTableName)
	if err != nil {
		return nil, err
	}

	limitParams, err := pGen.paramsInLimitExpr(tree.Limit, tableAliasMap)
	if err != nil {
		return nil, err
	}
	parsedQuery.Params = append(whereParams, limitParams...)

	err = parsedQuery.parseNameAndCmd()
	if err != nil {
		return nil, err
	}

	return &parsedQuery, nil
}

// FromTable describes a table reference in the "FROM" clause of a query.
type FromTable struct {
	TrueName     string // the true table name as described in the schema
	IsLeftJoined bool   // which could result in null columns
}

// FromTables describes a map between table alias expressions and the
// proper table name
type FromTables map[string]FromTable

func parseFrom(from sqlparser.TableExprs, isLeftJoined bool) (FromTables, string, error) {
	tables := make(map[string]FromTable)
	var defaultTableName string
	for _, expr := range from {
		switch v := expr.(type) {
		case *sqlparser.AliasedTableExpr:
			name, ok := v.Expr.(sqlparser.TableName)
			if !ok {
				return nil, "", fmt.Errorf("failed to parse AliasedTableExpr name: %v", v)
			}
			t := FromTable{
				TrueName:     name.Name.String(),
				IsLeftJoined: isLeftJoined,
			}
			if v.As.String() != "" {
				tables[v.As.String()] = t
			} else {
				tables[name.Name.String()] = t
			}
			defaultTableName = name.Name.String()
		case *sqlparser.JoinTableExpr:
			isLeftJoin := v.Join == "left join"
			left, leftMostTableName, err := parseFrom([]sqlparser.TableExpr{v.LeftExpr}, false)
			if err != nil {
				return nil, "", err
			}
			right, _, err := parseFrom([]sqlparser.TableExpr{v.RightExpr}, isLeftJoin)
			if err != nil {
				return nil, "", err
			}
			// merge the left and right maps
			for k, v := range left {
				right[k] = v
			}
			return right, leftMostTableName, nil
		default:
			return nil, "", fmt.Errorf("failed to parse table expr: %v", v)
		}
	}
	return tables, defaultTableName, nil
}

func (pGen PackageGenerator) parseUpdate(node *sqlparser.Update, query string) (*Query, error) {
	tableAliasMap, defaultTable, err := parseFrom(node.TableExprs, false)
	if err != nil {
		return nil, fmt.Errorf("failed to parse table name alias's: %w", err)
	}

	params := []*Param{}
	for _, updateExpr := range node.Exprs {
		newValue, isValue := updateExpr.Expr.(*sqlparser.SQLVal)
		if !isValue {
			continue
		} else if isParam := newValue.Type == sqlparser.ValArg; !isParam {
			continue
		}
		col, err := pGen.getColType(updateExpr.Name, tableAliasMap, defaultTable)
		if err != nil {
			return nil, fmt.Errorf("failed to determine type of a parameter's column: %w", err)
		}
		originalParamName := string(newValue.Val)
		param := Param{
			OriginalName: originalParamName,
			Name:         paramName(col.Name, originalParamName),
			Typ:          pGen.goTypeCol(*col),
		}
		params = append(params, &param)
	}

	whereParams, err := pGen.paramsInWhereExpr(node.Where, tableAliasMap, defaultTable)
	if err != nil {
		return nil, fmt.Errorf("failed to parse params from WHERE expression: %w", err)
	}

	parsedQuery := Query{
		SQL:              query,
		Columns:          nil,
		Params:           append(params, whereParams...),
		DefaultTableName: defaultTable,
	}
	err = parsedQuery.parseNameAndCmd()
	if err != nil {
		return nil, err
	}

	return &parsedQuery, nil
}

func (pGen PackageGenerator) parseInsert(node *sqlparser.Insert, query string) (*Query, error) {
	params := []*Param{}
	cols := node.Columns
	tableName := node.Table.Name.String()

	switch rows := node.Rows.(type) {
	case *sqlparser.Select:
		selectQuery, err := pGen.parseSelect(rows, query)
		if err != nil {
			return nil, err
		}
		params = append(params, selectQuery.Params...)
	case sqlparser.Values:
		for _, row := range rows {
			for colIx, item := range row {
				switch v := item.(type) {
				case *sqlparser.SQLVal:
					if v.Type == sqlparser.ValArg {
						colName := cols[colIx].String()
						col, err := pGen.schemaLookup(tableName, colName)
						varName := string(v.Val)
						param := &Param{OriginalName: varName}
						if err == nil {
							param.Name = paramName(col.Name, varName)
							param.Typ = pGen.goTypeCol(*col)
						} else {
							param.Name = "Unknown"
							param.Typ = "interface{}"
						}
						params = append(params, param)
					}
				case *sqlparser.FuncExpr:
					name, raw, err := matchFuncExpr(v)

					if err != nil {
						return nil, err
					}
					if name == "" || raw == "" {
						continue
					}
					colName := cols[colIx].String()
					col, err := pGen.schemaLookup(tableName, colName)
					param := &Param{
						OriginalName: raw,
					}
					if err == nil {
						param.Name = name
						param.Typ = pGen.goTypeCol(*col)
					} else {
						param.Name = "Unknown"
						param.Typ = "interface{}"
					}
					params = append(params, param)
				default:
					return nil, fmt.Errorf("failed to parse insert query value")
				}
			}
		}
	default:
		return nil, fmt.Errorf("Unknown insert row type of %T", node.Rows)
	}

	parsedQuery := &Query{
		SQL:              query,
		Params:           params,
		Columns:          nil,
		DefaultTableName: tableName,
	}

	err := parsedQuery.parseNameAndCmd()
	if err != nil {
		return nil, err
	}
	return parsedQuery, nil
}

func (pGen PackageGenerator) parseDelete(node *sqlparser.Delete, query string) (*Query, error) {
	tableAliasMap, defaultTableName, err := parseFrom(node.TableExprs, false)
	if err != nil {
		return nil, fmt.Errorf("failed to parse table name alias's: %w", err)
	}

	whereParams, err := pGen.paramsInWhereExpr(node.Where, tableAliasMap, defaultTableName)
	if err != nil {
		return nil, err
	}

	limitParams, err := pGen.paramsInLimitExpr(node.Limit, tableAliasMap)
	if err != nil {
		return nil, err
	}
	parsedQuery := &Query{
		SQL:              query,
		Params:           append(whereParams, limitParams...),
		Columns:          nil,
		DefaultTableName: defaultTableName,
	}
	err = parsedQuery.parseNameAndCmd()
	if err != nil {
		return nil, err
	}

	return parsedQuery, nil
}

func (pGen PackageGenerator) parseSelectAliasExpr(exprs sqlparser.SelectExprs, tableAliasMap FromTables, defaultTable string) ([]Column, error) {
	cols := []Column{}
	for _, col := range exprs {
		switch expr := col.(type) {
		case *sqlparser.AliasedExpr:
			hasAlias := !expr.As.IsEmpty()

			switch v := expr.Expr.(type) {
			case *sqlparser.ColName:
				res, err := pGen.getColType(v, tableAliasMap, defaultTable)
				if err != nil {
					return nil, err
				}
				if hasAlias {
					res.Name = expr.As // applys the alias
				}

				cols = append(cols, *res)
			case *sqlparser.GroupConcatExpr:
				cols = append(cols, Column{
					ColumnDefinition: &sqlparser.ColumnDefinition{
						Name: sqlparser.NewColIdent(expr.As.String()),
						Type: sqlparser.ColumnType{
							Type:    "varchar",
							NotNull: true,
						},
					},
					Table: "", // group concat expressions don't originate from a table schema
				},
				)
			case *sqlparser.FuncExpr:
				funcName := v.Name.Lowered()
				funcType := functionReturnType(funcName)

				var returnVal sqlparser.ColIdent
				if hasAlias {
					returnVal = expr.As
				} else {
					returnVal = sqlparser.NewColIdent(funcName)
				}

				colDfn := &sqlparser.ColumnDefinition{
					Name: returnVal,
					Type: sqlparser.ColumnType{
						Type:    funcType,
						NotNull: sqlparser.BoolVal(!functionIsNullable(funcName)),
					},
				}
				cols = append(cols, Column{colDfn, ""}) // func returns types don't originate from a table schema
			}
		default:
			return nil, fmt.Errorf("Failed to handle select expr of type : %T", expr)
		}
	}
	return cols, nil
}

// GeneratePkg is the main entry to mysql generator package
func GeneratePkg(pkgName string, schemaPath, querysPath []string, settings config.CombinedSettings) (*Result, error) {
	s := NewSchema()
	generator := PackageGenerator{
		Schema:           s,
		CombinedSettings: settings,
		packageName:      pkgName,
	}
	_, err := parsePath(schemaPath, generator)
	if err != nil {
		return nil, err
	}
	result, err := parsePath(querysPath, generator)
	if err != nil {
		return nil, err
	}
	return result, nil
}
