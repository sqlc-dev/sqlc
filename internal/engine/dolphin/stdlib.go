package dolphin

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
)

func defaultSchema(name string) *catalog.Schema {
	s := &catalog.Schema{Name: name}
	s.Funcs = []*catalog.Function{

		{
			Name:       "ABS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ACOS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ADDDATE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ADDTIME",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "AES_DECRYPT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "AES_ENCRYPT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ANY_VALUE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ASCII",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ASIN",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ATAN",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ATAN2",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "AVG",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "BENCHMARK",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "BIN",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "BIN_TO_UUID",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "BIT_AND",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "BIT_COUNT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "BIT_LENGTH",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "BIT_OR",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "BIT_XOR",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "CAST",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "CEIL",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "CEILING",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "CHAR",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "CHARACTER_LENGTH",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "CHARSET",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "CHAR_LENGTH",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "COALESCE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "COERCIBILITY",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "COLLATION",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "COMPRESS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "CONCAT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "CONCAT_WS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "CONNECTION_ID",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "CONV",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "CONVERT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "CONVERT_TZ",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "COS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "COT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "COUNT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "bigint"},
		},
		{
			Name: "COUNT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bigint"},
		},
		{
			Name:       "CRC32",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "CUME_DIST",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "CURDATE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "CURDATE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "CURRENT_DATE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "CURRENT_ROLE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "CURRENT_TIME",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "CURRENT_TIMESTAMP",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "CURRENT_USER",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "CURTIME",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "CURTIME",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "DATABASE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "DATE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "DATEDIFF",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "DATE_ADD",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "DATE_ADD_INTERVAL",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "DATE_FORMAT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "DATE_SUB",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "DATE_SUB_INTERVAL",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "DAY",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "DAYNAME",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "DAYOFMONTH",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "DAYOFWEEK",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "DAYOFYEAR",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "DEFAULT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "DEGREES",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "DENSE_RANK",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "DISTINCT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ELT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "EXP",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "EXPORT_SET",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "EXTRACT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "EXTRACT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "EXTRACTVALUE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "FIELD",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "FIND_IN_SET",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "FIRST_VALUE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "FLOOR",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "FORMAT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "FORMAT_BYTES",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "FORMAT_PICO_TIME",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "FOUND_ROWS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "FROM_BASE64",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "FROM_DAYS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "FROM_UNIXTIME",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "GEOMCOLLECTION",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "GEOMETRYCOLLECTION",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "GET_FORMAT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "GET_LOCK",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "GREATEST",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "GROUPING",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "GROUP_CONCAT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "GTID_SUBSET",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "GTID_SUBTRACT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "HEX",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "HOUR",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ICU_VERSION",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "IF",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "IFNULL",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "INET6_ATON",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "INET6_NTOA",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "INET_ATON",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "INET_NTOA",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "INSERT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "INSTR",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "INTERVAL",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ISNULL",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "IS_FREE_LOCK",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "IS_IPV4",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "IS_IPV4_COMPAT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "IS_IPV4_MAPPED",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "IS_IPV6",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "IS_USED_LOCK",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "IS_UUID",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_ARRAY",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_ARRAYAGG",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_ARRAY_APPEND",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_ARRAY_INSERT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_CONTAINS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_CONTAINS_PATH",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_DEPTH",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_EXTRACT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_INSERT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_KEYS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_LENGTH",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_MERGE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_MERGE_PATCH",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_MERGE_PRESERVE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_OBJECT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_OBJECTAGG",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_OVERLAPS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_PRETTY",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_QUOTE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_REMOVE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_REPLACE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_SCHEMA_VALID",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_SCHEMA_VALIDATION_REPORT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_SCHEMA_VALIDATION_REPORT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_SEARCH",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_SET",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_STORAGE_FREE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_STORAGE_SIZE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_TYPE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_UNQUOTE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_VALID",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "JSON_VALUE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "LAG",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "LAST_DAY",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "LAST_INSERT_ID",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "LAST_VALUE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "LCASE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "LEAD",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "LEAST",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "LEFT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "LENGTH",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "LIKE_RANGE_MAX",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "LIKE_RANGE_MIN",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "LINESTRING",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "LN",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "LOAD_FILE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "LOCALTIME",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "LOCALTIMESTAMP",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "LOCATE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "LOG",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "LOG10",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "LOG2",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "LOWER",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "LPAD",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "LTRIM",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "MAKEDATE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "MAKETIME",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "MAKE_SET",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "MASTER_POS_WAIT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "MAX",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "MBRCONTAINS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "MBRCOVEREDBY",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "MBRCOVERS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "MBRDISJOINT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "MBREQUALS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "MBRINTERSECTS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "MBROVERLAPS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "MBRTOUCHES",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "MBRWITHIN",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "MD5",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "MICROSECOND",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "MID",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "MIN",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "MINUTE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "MOD",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "MONTH",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "MONTHNAME",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "MULTILINESTRING",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "MULTIPOINT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "MULTIPOLYGON",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "NAME_CONST",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "NOW",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "NTH_VALUE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "NTILE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "NULLIF",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "OCT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "OCTET_LENGTH",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ORD",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "PASSWORD",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "PERCENT_RANK",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "PERIOD_ADD",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "PERIOD_DIFF",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "PI",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "POINT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "POLYGON",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "POSITION",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "POSITION",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "POW",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "POWER",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "PS_CURRENT_THREAD_ID",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "PS_THREAD_ID",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "QUARTER",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "QUOTE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "RADIANS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "RAND",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "RANDOM_BYTES",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "RANK",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "REGEXP_INSTR",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "REGEXP_LIKE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "REGEXP_REPLACE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "REGEXP_SUBSTR",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "RELEASE_ALL_LOCKS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "RELEASE_LOCK",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "REPEAT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "REPLACE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "REVERSE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "REVERSE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "RIGHT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ROLES_GRAPHML",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ROUND",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ROW_COUNT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ROW_NUMBER",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "RPAD",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "RTRIM",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "SCHEMA",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "SECOND",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "SEC_TO_TIME",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "SESSION_USER",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "SHA",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "SHA1",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "SHA2",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "SIGN",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "SIN",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "SLEEP",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "SOUNDEX",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "SPACE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "SQRT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "STATEMENT_DIGEST",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "STATEMENT_DIGEST_TEXT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "STD",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "STDDEV",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "STDDEV_POP",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "STDDEV_SAMP",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "STRCMP",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "STR_TO_DATE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_AREA",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_ASBINARY",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_ASGEOJSON",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_ASTEXT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_ASWKB",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_ASWKT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_BUFFER",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_BUFFER_STRATEGY",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_CENTROID",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_CONTAINS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_CONVEXHULL",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_CROSSES",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_DIFFERENCE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_DIMENSION",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_DISJOINT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_DISTANCE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_DISTANCE_SPHERE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_ENDPOINT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_ENVELOPE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_EQUALS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_EXTERIORRING",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_GEOHASH",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_GEOMCOLLFROMTEXT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_GEOMCOLLFROMTXT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_GEOMCOLLFROMWKB",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_GEOMETRYCOLLECTIONFROMTEXT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_GEOMETRYCOLLECTIONFROMWKB",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_GEOMETRYFROMTEXT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_GEOMETRYFROMWKB",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_GEOMETRYN",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_GEOMETRYTYPE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_GEOMFROMGEOJSON",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_GEOMFROMTEXT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_GEOMFROMWKB",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_INTERIORRINGN",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_INTERSECTION",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_INTERSECTS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_ISCLOSED",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_ISEMPTY",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_ISSIMPLE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_ISVALID",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_LATFROMGEOHASH",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_LATITUDE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_LENGTH",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_LINEFROMTEXT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_LINEFROMWKB",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_LINESTRINGFROMTEXT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_LINESTRINGFROMWKB",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_LONGFROMGEOHASH",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_LONGITUDE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_MAKEENVELOPE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_MLINEFROMTEXT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_MLINEFROMWKB",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_MPOINTFROMTEXT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_MPOINTFROMWKB",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_MPOLYFROMTEXT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_MPOLYFROMWKB",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_MULTILINESTRINGFROMTEXT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_MULTILINESTRINGFROMWKB",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_MULTIPOINTFROMTEXT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_MULTIPOINTFROMWKB",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_MULTIPOLYGONFROMTEXT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_MULTIPOLYGONFROMWKB",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_NUMGEOMETRIES",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_NUMINTERIORRING",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_NUMINTERIORRINGS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_NUMPOINTS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_OVERLAPS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_POINTFROMGEOHASH",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_POINTFROMTEXT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_POINTFROMWKB",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_POINTN",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_POLYFROMTEXT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_POLYFROMWKB",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_POLYGONFROMTEXT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_POLYGONFROMWKB",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_SIMPLIFY",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_SRID",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_STARTPOINT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_SWAPXY",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_SYMDIFFERENCE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_TOUCHES",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_TRANSFORM",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_UNION",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_VALIDATE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_WITHIN",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_X",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "ST_Y",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "SUBDATE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "SUBDATE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "SUBSTR",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "SUBSTRING",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "SUBSTRING",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "SUBSTRING_INDEX",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "SUBTIME",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "SUM",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "SYSDATE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "SYSDATE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "SYSTEM_USER",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "TAN",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "TIME",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "TIMEDIFF",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "TIMESTAMP",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "TIMESTAMPADD",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "TIMESTAMPDIFF",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "TIME_FORMAT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "TIME_TO_SEC",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "TO_BASE64",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "TO_DAYS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "TO_SECONDS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "TRIM",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "TRIM",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "TRUNCATE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "UCASE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "UNCOMPRESS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "UNCOMPRESSED_LENGTH",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "UNHEX",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "UNIX_TIMESTAMP",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "UPDATEXML",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "UPPER",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "USER",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "UTC_DATE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "UTC_TIME",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "UTC_TIMESTAMP",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "UUID",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "UUID_SHORT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "UUID_TO_BIN",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "VALIDATE_PASSWORD_STRENGTH",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "VALUES",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "VARIANCE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "VAR_POP",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "VAR_SAMP",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "VERSION",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "WAIT_FOR_EXECUTED_GTID_SET",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "WAIT_UNTIL_SQL_THREAD_AFTER_GTIDS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "WEEK",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "WEEKDAY",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "WEEKOFYEAR",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "WEIGHT_STRING",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "YEAR",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "YEARWEEK",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
	return s
}
