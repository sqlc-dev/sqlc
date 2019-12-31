package mysql

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/kyleconroy/sqlc/internal/dinosql"
	"vitess.io/vitess/go/vt/sqlparser"
)

// Query holds the data for walking and validating mysql querys
type Query struct {
	SQL              string
	Columns          []*sqlparser.ColumnDefinition
	Params           []*Param
	Name             string
	Cmd              string // TODO: Pick a better name. One of: one, many, exec, execrows
	defaultTableName string // for columns that are not qualified
	schemaLookup     *Schema
}

func parseFile(filepath string, inPkg string, s *Schema, settings dinosql.GenerateSettings) (*Result, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open file [%v]: %v", filepath, err)
	}
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("Failed to read contents of file [%v]: %v", filepath, err)
	}
	rawQueries := strings.Split(string(contents), "\n\n")

	parsedQueries := []*Query{}

	for _, query := range rawQueries {
		result, err := parseQueryString(query, s, settings)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse query in filepath [%v]: %v", filepath, err)
		}
		if result == nil {
			continue
		}
		parsedQueries = append(parsedQueries, result)
	}

	r := Result{
		Queries:     parsedQueries,
		Schema:      s,
		packageName: inPkg,
	}
	return &r, nil
}

func parseQueryString(query string, s *Schema, settings dinosql.GenerateSettings) (*Query, error) {
	tree, err := sqlparser.Parse(query)

	if err != nil {
		return nil, err
	}

	switch tree := tree.(type) {
	case *sqlparser.Select:
		defaultTableName := getDefaultTable(&tree.From)
		res, err := parseSelect(tree, query, s, defaultTableName, settings)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse SELECT query: %v", err)
		}
		return res, nil
	case *sqlparser.Insert:
		insert, err := parseInsert(tree, query, s, settings)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse INSERT query: %v", err)
		}
		return insert, nil
	case *sqlparser.Update:
		update, err := parseUpdate(tree, query, s, settings)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse UPDATE query: %v", err)
		}
		return update, nil
	case *sqlparser.DDL:
		s.Add(tree)
		return nil, nil
	default:
		panic("Unsupported SQL statement type")
		// return &Query{}, nil
	}
	return nil, fmt.Errorf("Failed to parse query statement: %v", query)
}

func (q *Query) parseNameAndCmd() error {
	if q == nil {
		return errors.New("Cannot parse name and cmd from null query")
	}
	_, comments := sqlparser.SplitMarginComments(q.SQL)
	err := q.parseLeadingComment(comments.Leading)
	if err != nil {
		return fmt.Errorf("Failed to parse leading comment %v", err)
	}
	return nil
}

func parseSelect(tree *sqlparser.Select, query string, s *Schema, defaultTableName string, settings dinosql.GenerateSettings) (*Query, error) {
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

	parsedQuery.Columns = parseSelectAliasExpr(tree.SelectExprs, s, defaultTableName)

	whereParams, err := paramsInWhereExpr(tree.Where, s, defaultTableName, settings)
	if err != nil {
		return nil, err
	}

	limitParams, err := paramsInLimitExpr(tree.Limit, s, settings)
	if err != nil {
		return nil, err
	}
	parsedQuery.Params = append(whereParams, limitParams...)

	err = parsedQuery.parseNameAndCmd()
	if err != nil {
		return nil, err
	}
	parsedQuery.SQL = sqlparser.String(tree)

	return &parsedQuery, nil
}

func getDefaultTable(node *sqlparser.TableExprs) string {
	// TODO: improve this
	var tableName string
	visit := func(node sqlparser.SQLNode) (bool, error) {
		switch v := node.(type) {
		case sqlparser.TableName:
			if name := v.Name.String(); name != "" {
				tableName = name
				return false, nil
			}
		}
		return true, nil
	}
	sqlparser.Walk(visit, node)
	return tableName
}

func parseUpdate(node *sqlparser.Update, query string, s *Schema, settings dinosql.GenerateSettings) (*Query, error) {
	defaultTable := getDefaultTable(&node.TableExprs)

	params := []*Param{}
	for _, updateExpr := range node.Exprs {
		col := updateExpr.Name
		newValue, isParam := updateExpr.Expr.(*sqlparser.SQLVal)
		if !isParam {
			continue
		}
		colDfn, err := s.getColType(col, defaultTable)
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

	whereParams, err := paramsInWhereExpr(node.Where.Expr, s, defaultTable, settings)
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

func parseSelectAliasExpr(exprs sqlparser.SelectExprs, s *Schema, defaultTable string) []*sqlparser.ColumnDefinition {
	colDfns := []*sqlparser.ColumnDefinition{}
	for _, col := range exprs {
		switch expr := col.(type) {
		case *sqlparser.AliasedExpr:
			hasAlias := !expr.As.IsEmpty()

			switch v := expr.Expr.(type) {
			case *sqlparser.ColName:
				res, err := s.getColType(v, defaultTable)
				if err != nil {
					panic(fmt.Sprintf("Column not found in schema: %v", err))
				}
				if hasAlias {
					res.Name = expr.As // applys the alias
				}
				colDfns = append(colDfns, res)
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
			panic(fmt.Sprintf("Failed to handle select expr of type : %T\n", expr))
		}
	}
	return colDfns
}

// GeneratePkg is the main entry to mysql generator package
func GeneratePkg(pkgName string, querysPath string, settings dinosql.GenerateSettings) (map[string]string, error) {
	s := NewSchema()
	result, err := parseFile(querysPath, pkgName, s, settings)
	if err != nil {
		return nil, err
	}
	output, err := dinosql.Generate(result, settings)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate output: %v", err)
	}

	return output, nil
}
