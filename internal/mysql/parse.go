package mysql

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/kyleconroy/sqlc/internal/dinosql"
	"vitess.io/vitess/go/vt/sqlparser"
)

// Query holds the data for walking and validating mysql querys
type Query struct {
	SQL              string                        // the string representation of the parsed query
	Columns          []*sqlparser.ColumnDefinition // definitions for all columns returned by this query
	Params           []*Param                      // "?" params in the query string
	Name             string                        // the Go function name
	Cmd              string                        // TODO: Pick a better name. One of: one, many, exec, execrows
	DefaultTableName string                        // for columns that are not qualified
	SchemaLookup     *Schema                       // for validation and conversion to Go types

	Filename string
}

func parsePath(sqlPath string, inPkg string, s *Schema, settings dinosql.GenerateSettings) (*Result, error) {
	files, err := dinosql.ReadSQLFiles(sqlPath)
	if err != nil {
		return nil, err
	}

	parseErrors := dinosql.ParserErr{}

	parsedQueries := []*Query{}
	for _, filename := range files {
		blob, err := ioutil.ReadFile(filename)
		if err != nil {
			parseErrors.Add(filename, "", 0, err)
		}
		contents := dinosql.RemoveRollbackStatements(string(blob))
		if err != nil {
			parseErrors.Add(filename, "", 0, err)
			continue
		}
		queries, err := parseContents(filename, contents, s, settings)
		if err != nil {
			if positionedErr, ok := err.(PositionedErr); ok {
				parseErrors.Add(filename, contents, positionedErr.Pos, err)
			} else {
				parseErrors.Add(filename, contents, 0, err)
			}
			continue
		}
		parsedQueries = append(parsedQueries, queries...)
	}

	if len(parseErrors.Errs) > 0 {
		return nil, &parseErrors
	}

	return &Result{
		Queries:     parsedQueries,
		Schema:      s,
		packageName: inPkg,
	}, nil
}

func parseContents(filename, contents string, s *Schema, settings dinosql.GenerateSettings) ([]*Query, error) {
	t := sqlparser.NewStringTokenizer(contents)
	var queries []*Query
	var start int
	for {
		q, err := sqlparser.ParseNextStrictDDL(t)
		if err == io.EOF {
			break
		} else if err != nil {
			parsedLoc, locErr := locFromSyntaxErr(err)
			if locErr != nil {
				parsedLoc = start // next best guess of the error location
			}
			near, nearErr := nearStrFromSyntaxErr(err)
			if nearErr != nil {
				return nil, PositionedErr{parsedLoc, fmt.Errorf("syntax error")}
			}
			return nil, PositionedErr{parsedLoc, fmt.Errorf("syntax error at or near '%s'", near)}
		}
		query := contents[start : t.Position-1]
		result, err := parseQueryString(q, query, s, settings)
		if err != nil {
			return nil, PositionedErr{start, err}
		}
		start = t.Position
		if result == nil {
			continue
		}
		result.Filename = filepath.Base(filename)
		queries = append(queries, result)
	}
	return queries, nil
}

func parseQueryString(tree sqlparser.Statement, query string, s *Schema, settings dinosql.GenerateSettings) (*Query, error) {
	var parsedQuery *Query
	switch tree := tree.(type) {
	case *sqlparser.Select:
		selectQuery, err := parseSelect(tree, query, s, settings)
		if err != nil {
			return nil, err
		}
		parsedQuery = selectQuery
	case *sqlparser.Insert:
		insert, err := parseInsert(tree, query, s, settings)
		if err != nil {
			return nil, err
		}
		parsedQuery = insert
	case *sqlparser.Update:
		update, err := parseUpdate(tree, query, s, settings)
		if err != nil {
			return nil, err
		}
		parsedQuery = update
	case *sqlparser.Delete:
		delete, err := parseDelete(tree, query, s, settings)
		if err != nil {
			return nil, err
		}
		parsedQuery = delete
	case *sqlparser.DDL:
		s.Add(tree)
		return nil, nil
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
	err := q.parseLeadingComment(comments.Leading)
	if err != nil {
		return fmt.Errorf("failed to parse leading comment %w", err)
	}
	return nil
}

func parseSelect(tree *sqlparser.Select, query string, s *Schema, settings dinosql.GenerateSettings) (*Query, error) {
	tableAliasMap, err := parseFrom(tree.From, false)
	if err != nil {
		return nil, fmt.Errorf("failed to parse table name alias's: %w", err)
	}
	defaultTableName := getDefaultTable(tableAliasMap)

	// handle * expressions first by expanding all columns of the default table
	_, ok := tree.SelectExprs[0].(*sqlparser.StarExpr)
	if ok {
		colNames := []sqlparser.SelectExpr{}
		colDfns := s.tables[defaultTableName]
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
		SchemaLookup:     s,
	}
	cols, err := parseSelectAliasExpr(tree.SelectExprs, s, tableAliasMap, defaultTableName)
	if err != nil {
		return nil, err
	}
	parsedQuery.Columns = cols

	whereParams, err := paramsInWhereExpr(tree.Where, s, tableAliasMap, defaultTableName, settings)
	if err != nil {
		return nil, err
	}

	limitParams, err := paramsInLimitExpr(tree.Limit, s, tableAliasMap, settings)
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

func parseFrom(from sqlparser.TableExprs, isLeftJoined bool) (FromTables, error) {
	tables := make(map[string]FromTable)
	for _, expr := range from {
		switch v := expr.(type) {
		case *sqlparser.AliasedTableExpr:
			name, ok := v.Expr.(sqlparser.TableName)
			if !ok {
				return nil, fmt.Errorf("failed to parse AliasedTableExpr name: %v", spew.Sdump(v))
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
		case *sqlparser.JoinTableExpr:
			isLeftJoin := v.Join == "left join"
			left, err := parseFrom([]sqlparser.TableExpr{v.LeftExpr}, false)
			if err != nil {
				return nil, err
			}
			right, err := parseFrom([]sqlparser.TableExpr{v.RightExpr}, isLeftJoin)
			if err != nil {
				return nil, err
			}
			// merge the left and right maps
			for k, v := range left {
				right[k] = v
			}
			return right, nil
		default:
			return nil, fmt.Errorf("failed to parse table expr: %v", spew.Sdump(v))
		}
	}
	return tables, nil
}

func getDefaultTable(tableAliasMap FromTables) string {
	if len(tableAliasMap) != 1 {
		return ""
	}
	for _, val := range tableAliasMap {
		return val.TrueName
	}
	panic("Should never be reached.")
}

func parseUpdate(node *sqlparser.Update, query string, s *Schema, settings dinosql.GenerateSettings) (*Query, error) {
	tableAliasMap, err := parseFrom(node.TableExprs, false)
	if err != nil {
		return nil, fmt.Errorf("failed to parse table name alias's: %w", err)
	}
	defaultTable := getDefaultTable(tableAliasMap)
	if err != nil {
		return nil, err
	}

	params := []*Param{}
	for _, updateExpr := range node.Exprs {
		col := updateExpr.Name
		newValue, isParam := updateExpr.Expr.(*sqlparser.SQLVal)
		if !isParam {
			continue
		}
		colDfn, err := s.getColType(col, tableAliasMap, defaultTable)
		if err != nil {
			return nil, fmt.Errorf("failed to determine type of a parameter's column: %w", err)
		}
		originalParamName := string(newValue.Val)
		param := Param{
			OriginalName: originalParamName,
			Name:         paramName(colDfn.Name, originalParamName),
			Typ:          goTypeCol(colDfn, settings),
		}
		params = append(params, &param)
	}

	whereParams, err := paramsInWhereExpr(node.Where.Expr, s, tableAliasMap, defaultTable, settings)
	if err != nil {
		return nil, fmt.Errorf("failed to parse params from WHERE expression: %w", err)
	}

	parsedQuery := Query{
		SQL:              query,
		Columns:          nil,
		Params:           append(params, whereParams...),
		DefaultTableName: defaultTable,
		SchemaLookup:     s,
	}
	parsedQuery.parseNameAndCmd()

	return &parsedQuery, nil
}

func parseInsert(node *sqlparser.Insert, query string, s *Schema, settings dinosql.GenerateSettings) (*Query, error) {
	cols := node.Columns
	tableName := node.Table.Name.String()
	rows, ok := node.Rows.(sqlparser.Values)
	if !ok {
		return nil, fmt.Errorf("Unknown insert row type of %T", node.Rows)
	}

	params := []*Param{}

	for _, row := range rows {
		for colIx, item := range row {
			switch v := item.(type) {
			case *sqlparser.SQLVal:
				if v.Type == sqlparser.ValArg {
					colName := cols[colIx].String()
					colDfn, err := s.schemaLookup(tableName, colName)
					varName := string(v.Val)
					p := &Param{OriginalName: varName}
					if err == nil {
						p.Name = paramName(colDfn.Name, varName)
						p.Typ = goTypeCol(colDfn, settings)
					} else {
						p.Name = "Unknown"
						p.Typ = "interface{}"
					}
					params = append(params, p)
				}
			default:
				panic("Error occurred in parsing INSERT statement")
			}
		}
	}
	parsedQuery := &Query{
		SQL:              query,
		Params:           params,
		Columns:          nil,
		DefaultTableName: tableName,
		SchemaLookup:     s,
	}
	parsedQuery.parseNameAndCmd()
	return parsedQuery, nil
}

func parseDelete(node *sqlparser.Delete, query string, s *Schema, settings dinosql.GenerateSettings) (*Query, error) {
	tableAliasMap, err := parseFrom(node.TableExprs, false)
	if err != nil {
		return nil, fmt.Errorf("failed to parse table name alias's: %w", err)
	}
	defaultTableName := getDefaultTable(tableAliasMap)
	if err != nil {
		return nil, err
	}

	whereParams, err := paramsInWhereExpr(node.Where, s, tableAliasMap, defaultTableName, settings)
	if err != nil {
		return nil, err
	}

	limitParams, err := paramsInLimitExpr(node.Limit, s, tableAliasMap, settings)
	if err != nil {
		return nil, err
	}
	parsedQuery := &Query{
		SQL:              query,
		Params:           append(whereParams, limitParams...),
		Columns:          nil,
		DefaultTableName: defaultTableName,
		SchemaLookup:     s,
	}
	err = parsedQuery.parseNameAndCmd()
	if err != nil {
		return nil, err
	}

	return parsedQuery, nil
}

func (q *Query) parseLeadingComment(comment string) error {
	for _, line := range strings.Split(comment, "\n") {
		if !strings.HasPrefix(line, "/* name:") {
			continue
		}
		part := strings.Split(strings.TrimSpace(line), " ")
		if len(part) == 3 {
			return fmt.Errorf("missing query type [':one', ':many', ':exec', ':execrows']: %s", line)
		}
		if len(part) != 5 {
			return fmt.Errorf("invalid query comment: %s", line)
		}
		queryName := part[2]
		queryType := strings.TrimSpace(part[3])
		switch queryType {
		case ":one", ":many", ":exec", ":execrows":
		default:
			return fmt.Errorf("invalid query type: %s", queryType)
		}
		// if err := validateQueryName(queryName); err != nil {
		// 	return err
		// }
		q.Name = queryName
		q.Cmd = queryType
	}
	return nil
}

func parseSelectAliasExpr(exprs sqlparser.SelectExprs, s *Schema, tableAliasMap FromTables, defaultTable string) ([]*sqlparser.ColumnDefinition, error) {
	colDfns := []*sqlparser.ColumnDefinition{}
	for _, col := range exprs {
		switch expr := col.(type) {
		case *sqlparser.AliasedExpr:
			hasAlias := !expr.As.IsEmpty()

			switch v := expr.Expr.(type) {
			case *sqlparser.ColName:
				res, err := s.getColType(v, tableAliasMap, defaultTable)
				if err != nil {
					return nil, err
				}
				if hasAlias {
					res.Name = expr.As // applys the alias
				}
				colDfns = append(colDfns, res)
			case *sqlparser.GroupConcatExpr:
				colDfns = append(colDfns, &sqlparser.ColumnDefinition{
					Name: sqlparser.NewColIdent(expr.As.String()),
					Type: sqlparser.ColumnType{
						Type:    "varchar",
						NotNull: true,
					},
				})
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
						NotNull: true,
					},
				}
				colDfns = append(colDfns, colDfn)
			}
		default:
			return nil, fmt.Errorf("Failed to handle select expr of type : %T", expr)
		}
	}
	return colDfns, nil
}

// GeneratePkg is the main entry to mysql generator package
func GeneratePkg(pkgName, schemaPath, querysPath string, settings dinosql.GenerateSettings) (*Result, error) {
	s := NewSchema()
	_, err := parsePath(schemaPath, pkgName, s, settings)
	if err != nil {
		return nil, err
	}
	result, err := parsePath(querysPath, pkgName, s, settings)
	if err != nil {
		return nil, err
	}
	return result, nil
}
