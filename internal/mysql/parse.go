package mysql

import (
	"fmt"
	"io/ioutil"
	"os"
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
	defaultTableName string                        // for columns that are not qualified
	schemaLookup     *Schema                       // for validation and conversion to Go types

	Filename string
}

func parsePath(sqlPath string, inPkg string, s *Schema, settings dinosql.GenerateSettings) (*Result, error) {
	files, err := dinosql.ReadSQLFiles(sqlPath)
	if err != nil {
		return nil, err
	}

	parsedQueries := []*Query{}
	for _, filename := range files {
		file, err := os.Open(filename)
		if err != nil {
			return nil, fmt.Errorf("Failed to open file [%v]: %v", filename, err)
		}
		contents, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("Failed to read contents of file [%v]: %v", filename, err)
		}
		rawQueries := strings.Split(string(contents), "\n\n")
		for _, query := range rawQueries {
			fmt.Println(query)
			if query == "" {
				continue
			}
			result, err := parseQueryString(query, s, settings)
			if err != nil {
				return nil, fmt.Errorf("Failed to parse query in filepath [%v]: %v", filename, err)
			}
			if result == nil {
				continue
			}
			result.Filename = filepath.Base(filename)
			parsedQueries = append(parsedQueries, result)
		}
	}

	return &Result{
		Queries:     parsedQueries,
		Schema:      s,
		packageName: inPkg,
	}, nil
}

func parseQueryString(query string, s *Schema, settings dinosql.GenerateSettings) (*Query, error) {
	tree, err := sqlparser.Parse(query)

	if err != nil {
		return nil, err
	}
	var parsedQuery *Query
	switch tree := tree.(type) {
	case *sqlparser.Select:
		selectQuery, err := parseSelect(tree, query, s, settings)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse SELECT query: %v", err)
		}
		parsedQuery = selectQuery
	case *sqlparser.Insert:
		insert, err := parseInsert(tree, query, s, settings)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse INSERT query: %v", err)
		}
		parsedQuery = insert
	case *sqlparser.Update:
		update, err := parseUpdate(tree, query, s, settings)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse UPDATE query: %v", err)
		}
		parsedQuery = update
	case *sqlparser.Delete:
		delete, err := parseDelete(tree, query, s, settings)
		delete.schemaLookup = nil
		if err != nil {
			return nil, fmt.Errorf("Failed to parse DELETE query: %v", err)
		}
		parsedQuery = delete
	case *sqlparser.DDL:
		s.Add(tree)
		return nil, nil
	default:
		panic("Unsupported SQL statement type")
	}
	paramsReplacedQuery, err := replaceParamStrs(sqlparser.String(tree), parsedQuery.Params)
	if err != nil {
		return nil, fmt.Errorf("Failed to replace param variables in query string: %v", err)
	}
	parsedQuery.SQL = paramsReplacedQuery
	return parsedQuery, nil
}

func (q *Query) parseNameAndCmd() error {
	if q == nil {
		return fmt.Errorf("Cannot parse name and cmd from null query")
	}
	_, comments := sqlparser.SplitMarginComments(q.SQL)
	err := q.parseLeadingComment(comments.Leading)
	if err != nil {
		return fmt.Errorf("Failed to parse leading comment %v", err)
	}
	return nil
}

func parseSelect(tree *sqlparser.Select, query string, s *Schema, settings dinosql.GenerateSettings) (*Query, error) {
	tableAliasMap, err := parseFrom(tree.From, false)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse table name alias's: %v", err)
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
		defaultTableName: defaultTableName,
		schemaLookup:     s,
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
				return nil, fmt.Errorf("Failed to parse AliasedTableExpr name: %v", spew.Sdump(v))
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
			return nil, fmt.Errorf("Failed to parse table expr: %v", spew.Sdump(v))
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
		return nil, fmt.Errorf("Failed to parse table name alias's: %v", err)
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
			return nil, fmt.Errorf("Failed to determine type of a parameter's column: %v", err)
		}
		originalParamName := string(newValue.Val)
		param := Param{
			originalName: originalParamName,
			name:         paramName(colDfn.Name, originalParamName),
			typ:          goTypeCol(colDfn, settings),
		}
		params = append(params, &param)
	}

	whereParams, err := paramsInWhereExpr(node.Where.Expr, s, tableAliasMap, defaultTable, settings)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse params from WHERE expression: %v", err)
	}

	parsedQuery := Query{
		SQL:              query,
		Columns:          nil,
		Params:           append(params, whereParams...),
		defaultTableName: defaultTable,
		schemaLookup:     s,
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
					colDfn, _ := s.schemaLookup(tableName, colName)
					varName := string(v.Val)
					p := &Param{
						originalName: varName,
						name:         paramName(colDfn.Name, varName),
						typ:          goTypeCol(colDfn, settings),
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
		defaultTableName: tableName,
		schemaLookup:     s,
	}
	parsedQuery.parseNameAndCmd()
	return parsedQuery, nil
}

func parseDelete(node *sqlparser.Delete, query string, s *Schema, settings dinosql.GenerateSettings) (*Query, error) {
	tableAliasMap, err := parseFrom(node.TableExprs, false)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse table name alias's: %v", err)
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
		defaultTableName: defaultTableName,
		schemaLookup:     s,
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
					panic(fmt.Sprintf("Column not found in schema: %v", err))
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
