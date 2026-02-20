package testcases

// DefaultRegistry is the global registry of all standard test cases
var DefaultRegistry = NewRegistry()

func init() {
	registerSelectTests()
	registerInsertTests()
	registerUpdateTests()
	registerDeleteTests()
	registerJoinTests()
	registerCTETests()
	registerSubqueryTests()
	registerUnionTests()
	registerAggregateTests()
	registerOperatorTests()
	registerCaseTests()
	registerNullTests()
	registerCastTests()
	registerFunctionTests()
	registerDataTypeTests()
	registerDDLTests()
	registerViewTests()
	registerUpsertTests()
	registerParamTests()
	registerResultTests()
	registerErrorTests()

	// Extension tests
	registerEnumTests()
	registerSchemaTests()
	registerArrayTests()
	registerJSONTests()
}

func registerSelectTests() {
	tests := []*TestCase{
		{ID: "S01", Name: "select_star", Category: CategorySelect, Description: "Star expansion returns all columns", Required: true},
		{ID: "S02", Name: "select_columns", Category: CategorySelect, Description: "Specific column selection", Required: true},
		{ID: "S03", Name: "select_column_alias", Category: CategorySelect, Description: "Column aliasing with AS", Required: true},
		{ID: "S04", Name: "select_table_alias", Category: CategorySelect, Description: "Table aliasing", Required: true},
		{ID: "S05", Name: "select_distinct", Category: CategorySelect, Description: "DISTINCT keyword", Required: true},
		{ID: "S06", Name: "select_where", Category: CategorySelect, Description: "WHERE with parameter", Required: true},
		{ID: "S07", Name: "select_where_multiple", Category: CategorySelect, Description: "Multiple WHERE conditions", Required: true},
		{ID: "S08", Name: "select_order_by", Category: CategorySelect, Description: "ORDER BY clause", Required: true},
		{ID: "S09", Name: "select_order_by_desc", Category: CategorySelect, Description: "ORDER BY with DESC", Required: true},
		{ID: "S10", Name: "select_order_by_multiple", Category: CategorySelect, Description: "Multiple ORDER BY columns", Required: true},
		{ID: "S11", Name: "select_limit", Category: CategorySelect, Description: "LIMIT with parameter", Required: true},
		{ID: "S12", Name: "select_limit_offset", Category: CategorySelect, Description: "LIMIT and OFFSET", Required: true},
		{ID: "S13", Name: "select_qualified_star", Category: CategorySelect, Description: "Table-qualified star", Required: true},
	}
	for _, tc := range tests {
		DefaultRegistry.Register(tc)
	}
}

func registerInsertTests() {
	tests := []*TestCase{
		{ID: "I01", Name: "insert_single_row", Category: CategoryInsert, Description: "Single row insert", Required: true},
		{ID: "I02", Name: "insert_multiple_rows", Category: CategoryInsert, Description: "Multi-row insert", Required: true},
		{ID: "I03", Name: "insert_partial_columns", Category: CategoryInsert, Description: "Partial column insert", Required: true},
		{ID: "I04", Name: "insert_returning_id", Category: CategoryInsert, Description: "RETURNING specific column", Required: true},
		{ID: "I05", Name: "insert_returning_star", Category: CategoryInsert, Description: "RETURNING all columns", Required: true},
		{ID: "I06", Name: "insert_select", Category: CategoryInsert, Description: "INSERT...SELECT", Required: true},
	}
	for _, tc := range tests {
		DefaultRegistry.Register(tc)
	}
}

func registerUpdateTests() {
	tests := []*TestCase{
		{ID: "U01", Name: "update_single_column", Category: CategoryUpdate, Description: "Single column update", Required: true},
		{ID: "U02", Name: "update_multiple_columns", Category: CategoryUpdate, Description: "Multiple column update", Required: true},
		{ID: "U03", Name: "update_expression", Category: CategoryUpdate, Description: "Expression in SET", Required: true},
		{ID: "U04", Name: "update_all_rows", Category: CategoryUpdate, Description: "Update without WHERE", Required: true},
		{ID: "U05", Name: "update_returning", Category: CategoryUpdate, Description: "RETURNING with UPDATE", Required: true},
		{ID: "U06", Name: "update_returning_columns", Category: CategoryUpdate, Description: "RETURNING specific columns", Required: true},
	}
	for _, tc := range tests {
		DefaultRegistry.Register(tc)
	}
}

func registerDeleteTests() {
	tests := []*TestCase{
		{ID: "D01", Name: "delete_where", Category: CategoryDelete, Description: "Conditional delete", Required: true},
		{ID: "D02", Name: "delete_all", Category: CategoryDelete, Description: "Delete all rows", Required: true},
		{ID: "D03", Name: "delete_returning", Category: CategoryDelete, Description: "RETURNING with DELETE", Required: true},
		{ID: "D04", Name: "delete_returning_columns", Category: CategoryDelete, Description: "RETURNING specific columns", Required: true},
	}
	for _, tc := range tests {
		DefaultRegistry.Register(tc)
	}
}

func registerJoinTests() {
	tests := []*TestCase{
		{ID: "J01", Name: "join_inner", Category: CategoryJoin, Description: "INNER JOIN", Required: true},
		{ID: "J02", Name: "join_left", Category: CategoryJoin, Description: "LEFT JOIN with nullable result", Required: true},
		{ID: "J03", Name: "join_right", Category: CategoryJoin, Description: "RIGHT JOIN", Required: true},
		{ID: "J04", Name: "join_full", Category: CategoryJoin, Description: "FULL OUTER JOIN", Required: true},
		{ID: "J05", Name: "join_cross", Category: CategoryJoin, Description: "CROSS JOIN", Required: true},
		{ID: "J06", Name: "join_implicit", Category: CategoryJoin, Description: "Implicit join with comma", Required: true},
		{ID: "J07", Name: "join_self", Category: CategoryJoin, Description: "Self-join", Required: true},
		{ID: "J08", Name: "join_multiple", Category: CategoryJoin, Description: "Multiple joins", Required: true},
		{ID: "J09", Name: "join_star_expansion", Category: CategoryJoin, Description: "Star with JOIN", Required: true},
		{ID: "J10", Name: "join_qualified_star", Category: CategoryJoin, Description: "Qualified star with JOIN", Required: true},
		{ID: "J11", Name: "join_many_to_many", Category: CategoryJoin, Description: "Many-to-many JOIN", Required: true},
	}
	for _, tc := range tests {
		DefaultRegistry.Register(tc)
	}
}

func registerCTETests() {
	tests := []*TestCase{
		{ID: "C01", Name: "cte_basic", Category: CategoryCTE, Description: "Basic CTE", Required: true},
		{ID: "C02", Name: "cte_multiple", Category: CategoryCTE, Description: "Multiple CTEs", Required: true},
		{ID: "C03", Name: "cte_chained", Category: CategoryCTE, Description: "CTEs referencing CTEs", Required: true},
		{ID: "C04", Name: "cte_recursive", Category: CategoryCTE, Description: "Recursive CTE", Required: true},
		{ID: "C05", Name: "cte_insert", Category: CategoryCTE, Description: "CTE with INSERT", Required: true},
		{ID: "C06", Name: "cte_update", Category: CategoryCTE, Description: "CTE with UPDATE", Required: true},
		{ID: "C07", Name: "cte_delete", Category: CategoryCTE, Description: "CTE with DELETE", Required: true},
	}
	for _, tc := range tests {
		DefaultRegistry.Register(tc)
	}
}

func registerSubqueryTests() {
	tests := []*TestCase{
		{ID: "Q01", Name: "subquery_from", Category: CategorySubquery, Description: "Subquery in FROM", Required: true},
		{ID: "Q02", Name: "subquery_where_in", Category: CategorySubquery, Description: "Subquery in WHERE IN", Required: true},
		{ID: "Q03", Name: "subquery_scalar", Category: CategorySubquery, Description: "Scalar subquery", Required: true},
		{ID: "Q04", Name: "subquery_exists", Category: CategorySubquery, Description: "EXISTS subquery", Required: true},
		{ID: "Q05", Name: "subquery_not_exists", Category: CategorySubquery, Description: "NOT EXISTS", Required: true},
		{ID: "Q06", Name: "subquery_select_list", Category: CategorySubquery, Description: "Subquery in SELECT", Required: true},
		{ID: "Q07", Name: "subquery_correlated", Category: CategorySubquery, Description: "Correlated subquery", Required: true},
	}
	for _, tc := range tests {
		DefaultRegistry.Register(tc)
	}
}

func registerUnionTests() {
	tests := []*TestCase{
		{ID: "N01", Name: "union_basic", Category: CategoryUnion, Description: "UNION deduplicates", Required: true},
		{ID: "N02", Name: "union_all", Category: CategoryUnion, Description: "UNION ALL keeps dupes", Required: true},
		{ID: "N03", Name: "union_order_by", Category: CategoryUnion, Description: "UNION with ORDER BY", Required: true},
		{ID: "N04", Name: "intersect", Category: CategoryUnion, Description: "INTERSECT", Required: true},
		{ID: "N05", Name: "except", Category: CategoryUnion, Description: "EXCEPT", Required: true},
	}
	for _, tc := range tests {
		DefaultRegistry.Register(tc)
	}
}

func registerAggregateTests() {
	tests := []*TestCase{
		{ID: "A01", Name: "agg_count_star", Category: CategoryAggregate, Description: "COUNT(*)", Required: true},
		{ID: "A02", Name: "agg_count_column", Category: CategoryAggregate, Description: "COUNT(column)", Required: true},
		{ID: "A03", Name: "agg_count_distinct", Category: CategoryAggregate, Description: "COUNT(DISTINCT)", Required: true},
		{ID: "A04", Name: "agg_sum", Category: CategoryAggregate, Description: "SUM", Required: true},
		{ID: "A05", Name: "agg_avg", Category: CategoryAggregate, Description: "AVG", Required: true},
		{ID: "A06", Name: "agg_min_max", Category: CategoryAggregate, Description: "MIN/MAX", Required: true},
		{ID: "A07", Name: "agg_group_by", Category: CategoryAggregate, Description: "GROUP BY", Required: true},
		{ID: "A08", Name: "agg_group_by_multiple", Category: CategoryAggregate, Description: "Multiple GROUP BY", Required: true},
		{ID: "A09", Name: "agg_having", Category: CategoryAggregate, Description: "HAVING clause", Required: true},
		{ID: "A10", Name: "agg_having_aggregate", Category: CategoryAggregate, Description: "HAVING with aggregate", Required: true},
		{ID: "A11", Name: "agg_mixed", Category: CategoryAggregate, Description: "Multiple aggregates", Required: true},
	}
	for _, tc := range tests {
		DefaultRegistry.Register(tc)
	}
}

func registerOperatorTests() {
	tests := []*TestCase{
		{ID: "O01", Name: "op_equal", Category: CategoryOperator, Description: "Equality", Required: true},
		{ID: "O02", Name: "op_not_equal", Category: CategoryOperator, Description: "Not equal", Required: true},
		{ID: "O03", Name: "op_less_than", Category: CategoryOperator, Description: "Less than", Required: true},
		{ID: "O04", Name: "op_less_equal", Category: CategoryOperator, Description: "Less or equal", Required: true},
		{ID: "O05", Name: "op_greater_than", Category: CategoryOperator, Description: "Greater than", Required: true},
		{ID: "O06", Name: "op_greater_equal", Category: CategoryOperator, Description: "Greater or equal", Required: true},
		{ID: "O07", Name: "op_between", Category: CategoryOperator, Description: "BETWEEN", Required: true},
		{ID: "O08", Name: "op_in_list", Category: CategoryOperator, Description: "IN with values", Required: true},
		{ID: "O09", Name: "op_is_null", Category: CategoryOperator, Description: "IS NULL", Required: true},
		{ID: "O10", Name: "op_is_not_null", Category: CategoryOperator, Description: "IS NOT NULL", Required: true},
		{ID: "O11", Name: "op_and", Category: CategoryOperator, Description: "AND", Required: true},
		{ID: "O12", Name: "op_or", Category: CategoryOperator, Description: "OR", Required: true},
		{ID: "O13", Name: "op_not", Category: CategoryOperator, Description: "NOT", Required: true},
		{ID: "O14", Name: "op_precedence", Category: CategoryOperator, Description: "Operator precedence", Required: true},
		{ID: "O15", Name: "op_like", Category: CategoryOperator, Description: "LIKE", Required: true},
		{ID: "O16", Name: "op_arithmetic", Category: CategoryOperator, Description: "Arithmetic operators", Required: true},
	}
	for _, tc := range tests {
		DefaultRegistry.Register(tc)
	}
}

func registerCaseTests() {
	tests := []*TestCase{
		{ID: "E01", Name: "case_simple", Category: CategoryCase, Description: "Simple CASE", Required: true},
		{ID: "E02", Name: "case_searched", Category: CategoryCase, Description: "Searched CASE", Required: true},
		{ID: "E03", Name: "case_no_else", Category: CategoryCase, Description: "CASE without ELSE", Required: true},
		{ID: "E04", Name: "case_in_where", Category: CategoryCase, Description: "CASE in WHERE", Required: true},
	}
	for _, tc := range tests {
		DefaultRegistry.Register(tc)
	}
}

func registerNullTests() {
	tests := []*TestCase{
		{ID: "F01", Name: "null_coalesce_two", Category: CategoryNull, Description: "COALESCE with 2 args", Required: true},
		{ID: "F02", Name: "null_coalesce_multiple", Category: CategoryNull, Description: "COALESCE with multiple", Required: true},
		{ID: "F03", Name: "null_coalesce_literal", Category: CategoryNull, Description: "COALESCE with literal", Required: true},
		{ID: "F04", Name: "null_nullif", Category: CategoryNull, Description: "NULLIF", Required: true},
	}
	for _, tc := range tests {
		DefaultRegistry.Register(tc)
	}
}

func registerCastTests() {
	tests := []*TestCase{
		{ID: "T01", Name: "cast_to_int", Category: CategoryCast, Description: "CAST to integer", Required: true},
		{ID: "T02", Name: "cast_to_text", Category: CategoryCast, Description: "CAST to text", Required: true},
		{ID: "T03", Name: "cast_param", Category: CategoryCast, Description: "CAST parameter", Required: true},
	}
	for _, tc := range tests {
		DefaultRegistry.Register(tc)
	}
}

func registerFunctionTests() {
	tests := []*TestCase{
		{ID: "B01", Name: "func_upper", Category: CategoryFunction, Description: "UPPER", Required: true},
		{ID: "B02", Name: "func_lower", Category: CategoryFunction, Description: "LOWER", Required: true},
		{ID: "B03", Name: "func_length", Category: CategoryFunction, Description: "LENGTH", Required: true},
		{ID: "B04", Name: "func_trim", Category: CategoryFunction, Description: "TRIM", Required: true},
		{ID: "B05", Name: "func_substring", Category: CategoryFunction, Description: "SUBSTRING", Required: true},
		{ID: "B06", Name: "func_replace", Category: CategoryFunction, Description: "REPLACE", Required: true},
		{ID: "B07", Name: "func_abs", Category: CategoryFunction, Description: "ABS", Required: true},
		{ID: "B08", Name: "func_round", Category: CategoryFunction, Description: "ROUND", Required: true},
		{ID: "B09", Name: "func_now", Category: CategoryFunction, Description: "NOW", Required: true},
		{ID: "B10", Name: "func_current_date", Category: CategoryFunction, Description: "CURRENT_DATE", Required: true},
	}
	for _, tc := range tests {
		DefaultRegistry.Register(tc)
	}
}

func registerDataTypeTests() {
	tests := []*TestCase{
		{ID: "DT01", Name: "datatype_int", Category: CategoryDataType, Description: "INT type", Required: true},
		{ID: "DT02", Name: "datatype_int_nullable", Category: CategoryDataType, Description: "Nullable INT", Required: true},
		{ID: "DT03", Name: "datatype_bigint", Category: CategoryDataType, Description: "BIGINT type", Required: true},
		{ID: "DT04", Name: "datatype_float", Category: CategoryDataType, Description: "FLOAT type", Required: true},
		{ID: "DT05", Name: "datatype_decimal", Category: CategoryDataType, Description: "DECIMAL type", Required: true},
		{ID: "DT06", Name: "datatype_text", Category: CategoryDataType, Description: "TEXT type", Required: true},
		{ID: "DT07", Name: "datatype_varchar", Category: CategoryDataType, Description: "VARCHAR type", Required: true},
		{ID: "DT08", Name: "datatype_boolean", Category: CategoryDataType, Description: "BOOLEAN type", Required: true},
		{ID: "DT09", Name: "datatype_timestamp", Category: CategoryDataType, Description: "TIMESTAMP type", Required: true},
		{ID: "DT10", Name: "datatype_date", Category: CategoryDataType, Description: "DATE type", Required: true},
		{ID: "DT11", Name: "datatype_blob", Category: CategoryDataType, Description: "BLOB type", Required: true},
	}
	for _, tc := range tests {
		DefaultRegistry.Register(tc)
	}
}

func registerDDLTests() {
	tests := []*TestCase{
		{ID: "DDL01", Name: "ddl_create_basic", Category: CategoryDDL, Description: "Basic table creation", Required: true},
		{ID: "DDL02", Name: "ddl_create_not_null", Category: CategoryDDL, Description: "NOT NULL constraint", Required: true},
		{ID: "DDL03", Name: "ddl_create_primary_key", Category: CategoryDDL, Description: "Primary key", Required: true},
		{ID: "DDL04", Name: "ddl_create_composite_pk", Category: CategoryDDL, Description: "Composite primary key", Required: true},
		{ID: "DDL05", Name: "ddl_create_unique", Category: CategoryDDL, Description: "UNIQUE constraint", Required: true},
		{ID: "DDL06", Name: "ddl_create_default", Category: CategoryDDL, Description: "Default value", Required: true},
		{ID: "DDL07", Name: "ddl_create_foreign_key", Category: CategoryDDL, Description: "Foreign key", Required: true},
		{ID: "DDL08", Name: "ddl_alter_add_column", Category: CategoryDDL, Description: "Add column", Required: true},
		{ID: "DDL09", Name: "ddl_alter_drop_column", Category: CategoryDDL, Description: "Drop column", Required: true},
		{ID: "DDL10", Name: "ddl_drop_table", Category: CategoryDDL, Description: "Drop table", Required: true},
	}
	for _, tc := range tests {
		DefaultRegistry.Register(tc)
	}
}

func registerViewTests() {
	tests := []*TestCase{
		{ID: "V01", Name: "view_select", Category: CategoryView, Description: "Query view", Required: true},
		{ID: "V02", Name: "view_filter", Category: CategoryView, Description: "Filter view", Required: true},
		{ID: "V03", Name: "view_complex", Category: CategoryView, Description: "Complex view with joins", Required: true},
	}
	for _, tc := range tests {
		DefaultRegistry.Register(tc)
	}
}

func registerUpsertTests() {
	tests := []*TestCase{
		{ID: "UP01", Name: "upsert_do_nothing", Category: CategoryUpsert, Description: "ON CONFLICT DO NOTHING", Required: true},
		{ID: "UP02", Name: "upsert_do_update", Category: CategoryUpsert, Description: "ON CONFLICT DO UPDATE", Required: true},
		{ID: "UP03", Name: "upsert_excluded", Category: CategoryUpsert, Description: "excluded pseudo-table", Required: true},
	}
	for _, tc := range tests {
		DefaultRegistry.Register(tc)
	}
}

func registerParamTests() {
	tests := []*TestCase{
		{ID: "P01", Name: "param_single", Category: CategoryParam, Description: "Single parameter", Required: true},
		{ID: "P02", Name: "param_multiple_same_type", Category: CategoryParam, Description: "Multiple same-type params", Required: true},
		{ID: "P03", Name: "param_multiple_different_type", Category: CategoryParam, Description: "Different type params", Required: true},
		{ID: "P04", Name: "param_repeated", Category: CategoryParam, Description: "Same param used twice", Required: true},
		{ID: "P05", Name: "param_in_function", Category: CategoryParam, Description: "Param inside function", Required: true},
	}
	for _, tc := range tests {
		DefaultRegistry.Register(tc)
	}
}

func registerResultTests() {
	tests := []*TestCase{
		{ID: "R01", Name: "result_one", Category: CategoryResult, Description: ":one annotation", Required: true},
		{ID: "R02", Name: "result_many", Category: CategoryResult, Description: ":many annotation", Required: true},
		{ID: "R03", Name: "result_exec", Category: CategoryResult, Description: ":exec annotation", Required: true},
		{ID: "R04", Name: "result_execrows", Category: CategoryResult, Description: ":execrows annotation", Required: true},
	}
	for _, tc := range tests {
		DefaultRegistry.Register(tc)
	}
}

func registerErrorTests() {
	tests := []*TestCase{
		{ID: "ER01", Name: "error_unknown_table", Category: CategoryError, Description: "Unknown table error", Required: true},
		{ID: "ER02", Name: "error_unknown_column_select", Category: CategoryError, Description: "Unknown column in SELECT", Required: true},
		{ID: "ER03", Name: "error_unknown_column_where", Category: CategoryError, Description: "Unknown column in WHERE", Required: true},
		{ID: "ER04", Name: "error_unknown_column_insert", Category: CategoryError, Description: "Unknown column in INSERT", Required: true},
		{ID: "ER05", Name: "error_unknown_column_update", Category: CategoryError, Description: "Unknown column in UPDATE", Required: true},
		{ID: "ER06", Name: "error_syntax", Category: CategoryError, Description: "Syntax error", Required: true},
	}
	for _, tc := range tests {
		DefaultRegistry.Register(tc)
	}
}

// Extension tests - not required for all engines

func registerEnumTests() {
	tests := []*TestCase{
		{ID: "EN01", Name: "enum_select", Category: CategoryEnum, Description: "Select with enum column", Required: false},
		{ID: "EN02", Name: "enum_filter", Category: CategoryEnum, Description: "Filter by enum", Required: false},
		{ID: "EN03", Name: "enum_insert", Category: CategoryEnum, Description: "Insert enum value", Required: false},
		{ID: "EN04", Name: "enum_update", Category: CategoryEnum, Description: "Update enum value", Required: false},
		{ID: "EN05", Name: "enum_nullable", Category: CategoryEnum, Description: "Nullable enum", Required: false},
	}
	for _, tc := range tests {
		DefaultRegistry.Register(tc)
	}
}

func registerSchemaTests() {
	tests := []*TestCase{
		{ID: "SC01", Name: "schema_select", Category: CategorySchema, Description: "Schema-qualified SELECT", Required: false},
		{ID: "SC02", Name: "schema_insert", Category: CategorySchema, Description: "Schema-qualified INSERT", Required: false},
		{ID: "SC03", Name: "schema_update", Category: CategorySchema, Description: "Schema-qualified UPDATE", Required: false},
		{ID: "SC04", Name: "schema_delete", Category: CategorySchema, Description: "Schema-qualified DELETE", Required: false},
		{ID: "SC05", Name: "schema_join", Category: CategorySchema, Description: "Cross-schema JOIN", Required: false},
	}
	for _, tc := range tests {
		DefaultRegistry.Register(tc)
	}
}

func registerArrayTests() {
	tests := []*TestCase{
		{ID: "AR01", Name: "array_select", Category: CategoryArray, Description: "Select with array columns", Required: false},
		{ID: "AR02", Name: "array_insert", Category: CategoryArray, Description: "Insert array value", Required: false},
		{ID: "AR03", Name: "array_any", Category: CategoryArray, Description: "ANY with array", Required: false},
	}
	for _, tc := range tests {
		DefaultRegistry.Register(tc)
	}
}

func registerJSONTests() {
	tests := []*TestCase{
		{ID: "JS01", Name: "json_select", Category: CategoryJSON, Description: "Select with JSON columns", Required: false},
		{ID: "JS02", Name: "json_insert", Category: CategoryJSON, Description: "Insert JSON value", Required: false},
		{ID: "JS03", Name: "json_nullable", Category: CategoryJSON, Description: "Nullable JSON", Required: false},
	}
	for _, tc := range tests {
		DefaultRegistry.Register(tc)
	}
}
