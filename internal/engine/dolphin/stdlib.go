package dolphin

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func defaultSchema(name string) *catalog.Schema {
	s := &catalog.Schema{Name: name}
	s.Funcs = []*catalog.Function{
		{
			Name: "ABS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "tinyint"},
				},
			},
			ReturnType: &ast.TypeName{Name: "tinyint"},
		},
		{
			Name: "ABS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "smallint"},
				},
			},
			ReturnType: &ast.TypeName{Name: "smallint"},
		},
		{
			Name: "ABS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "mediumint"},
				},
			},
			ReturnType: &ast.TypeName{Name: "mediumint"},
		},
		{
			Name: "ABS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "ABS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "bigint"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bigint"},
		},
		{
			Name: "ABS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double"},
		},
		{
			Name: "ABS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "ACOS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "ADDDATE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "date"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "date"},
		},
		{
			Name: "ADDTIME",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "time"},
				},
				{
					Type: &ast.TypeName{Name: "time"},
				},
			},
			ReturnType: &ast.TypeName{Name: "time"},
		},
		{
			Name: "AES_DECRYPT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "AES_DECRYPT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "binary"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "AES_ENCRYPT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "binary"},
		},
		{
			Name: "AES_ENCRYPT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "binary"},
				},
			},
			ReturnType: &ast.TypeName{Name: "binary"},
		},
		{
			Name: "ANY_VALUE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ASCII",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "ASIN",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "ATAN",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "ATAN2",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "AVG",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "BENCHMARK",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "BIN",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "BIN_TO_UUID",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "BIN_TO_UUID",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "tinyint"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "BIT_AND",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "BIT_COUNT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "bigint"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bigint"},
		},
		{
			Name: "BIT_LENGTH",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "BIT_OR",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "BIT_XOR",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "CAST",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "CEIL",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "CEIL",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "CEILING",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "CEILING",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "CHAR",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "CHARACTER_LENGTH",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "CHARSET",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "CHAR_LENGTH",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "COERCIBILITY",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "COLLATION",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "COMPRESS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "binary"},
		},
		{
			Name: "CONCAT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "CONCAT_WS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name:       "CONNECTION_ID",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "CONV",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "CONVERT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "CONVERT_TZ",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "datetime"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "datetime"},
		},
		{
			Name: "COS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "COT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
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
			Name: "CRC32",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name:       "CUME_DIST",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name:       "CURDATE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "date"},
		},
		{
			Name:       "CURRENT_DATE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "date"},
		},
		{
			Name:       "CURRENT_ROLE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name:       "CURRENT_TIME",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "time"},
		},
		{
			Name: "CURRENT_TIME",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "time"},
		},
		{
			Name:       "CURRENT_TIMESTAMP",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "datetime"},
		},
		{
			Name: "CURRENT_TIMESTAMP",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "datetime"},
		},
		{
			Name:       "CURRENT_USER",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name:       "CURTIME",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "time"},
		},
		{
			Name: "CURTIME",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "time"},
		},
		{
			Name:       "DATABASE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "DATE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "datetime"},
				},
			},
			ReturnType: &ast.TypeName{Name: "date"},
		},
		{
			Name: "DATEDIFF",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "datetime"},
				},
				{
					Type: &ast.TypeName{Name: "datetime"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "DATE_ADD",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "date"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "date"},
		},
		{
			Name: "DATE_ADD_INTERVAL",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "date"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "date"},
		},
		{
			Name: "DATE_FORMAT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "date"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "DATE_SUB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "date"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "date"},
		},
		{
			Name: "DATE_SUB_INTERVAL",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "date"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "date"},
		},
		{
			Name: "DAY",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "date"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "DAYNAME",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "date"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "DAYOFMONTH",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "date"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "DAYOFWEEK",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "date"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "DAYOFYEAR",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "date"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "DEFAULT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "DEGREES",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name:       "DENSE_RANK",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "DISTINCT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ELT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "EXP",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "EXPORT_SET",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "EXPORT_SET",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "EXPORT_SET",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "EXTRACT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "date"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "EXTRACTVALUE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "FIELD",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "FIND_IN_SET",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "FIRST_VALUE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "FLOOR",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "FORMAT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "FORMAT_BYTES",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "FORMAT_PICO_TIME",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "time"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name:       "FOUND_ROWS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "FROM_BASE64",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "FROM_DAYS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "date"},
		},
		{
			Name: "FROM_UNIXTIME",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "bigint"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "FROM_UNIXTIME",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "bigint"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "GEOMCOLLECTION",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "GEOMETRYCOLLECTION",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "GET_FORMAT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "GET_LOCK",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "GREATEST",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "GROUPING",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "GROUP_CONCAT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType:         &ast.TypeName{Name: "text"},
			ReturnTypeNullable: true,
		},
		{
			Name: "GTID_SUBSET",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "GTID_SUBTRACT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "HEX",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "HEX",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "HOUR",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "time"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name:       "ICU_VERSION",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "IF",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "IFNULL",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "INET6_ATON",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "binary"},
		},
		{
			Name: "INET6_NTOA",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "INET_ATON",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "INET_NTOA",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "INSERT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "INSTR",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "INTERVAL",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ISNULL",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "IS_FREE_LOCK",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "IS_IPV4",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "IS_IPV4_COMPAT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "IS_IPV4_MAPPED",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "IS_IPV6",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "IS_USED_LOCK",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "IS_UUID",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "JSON_ARRAY",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "json"},
		},
		{
			Name: "JSON_ARRAYAGG",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "json"},
		},
		{
			Name: "JSON_ARRAY_APPEND",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "json"},
		},
		{
			Name: "JSON_ARRAY_INSERT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "json"},
		},
		{
			Name: "JSON_CONTAINS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "JSON_CONTAINS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "JSON_CONTAINS_PATH",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "JSON_DEPTH",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "JSON_EXTRACT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "JSON_INSERT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "json"},
		},
		{
			Name: "JSON_KEYS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "json"},
		},
		{
			Name: "JSON_KEYS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "json"},
		},
		{
			Name: "JSON_LENGTH",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "JSON_LENGTH",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "JSON_MERGE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "json"},
		},
		{
			Name: "JSON_MERGE_PATCH",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "json"},
		},
		{
			Name: "JSON_MERGE_PRESERVE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "json"},
		},
		{
			Name: "JSON_OBJECT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "json"},
		},
		{
			Name: "JSON_OBJECTAGG",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "json"},
		},
		{
			Name: "JSON_OVERLAPS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "JSON_PRETTY",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "JSON_QUOTE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "JSON_REMOVE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "json"},
		},
		{
			Name: "JSON_REPLACE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "json"},
		},
		{
			Name: "JSON_SCHEMA_VALID",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "JSON_SCHEMA_VALIDATION_REPORT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "json"},
		},
		{
			Name: "JSON_SEARCH",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "JSON_SEARCH",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "JSON_SET",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "json"},
		},
		{
			Name: "JSON_STORAGE_FREE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "JSON_STORAGE_SIZE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "JSON_TYPE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "JSON_UNQUOTE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "json"},
		},
		{
			Name: "JSON_VALID",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "JSON_VALUE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "LAG",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type:       &ast.TypeName{Name: "int"},
					HasDefault: true,
				},
				{
					Type:       &ast.TypeName{Name: "any"},
					HasDefault: true,
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "LAST_DAY",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "date"},
				},
			},
			ReturnType: &ast.TypeName{Name: "date"},
		},
		{
			Name:       "LAST_INSERT_ID",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "bigint"},
		},
		{
			Name: "LAST_INSERT_ID",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bigint"},
		},
		{
			Name: "LAST_VALUE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "LCASE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "LEAD",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type:       &ast.TypeName{Name: "int"},
					HasDefault: true,
				},
				{
					Type:       &ast.TypeName{Name: "any"},
					HasDefault: true,
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "LEAST",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "LEFT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "LENGTH",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "LINESTRING",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "LN",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "LOAD_FILE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name:       "LOCALTIME",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "datetime"},
		},
		{
			Name: "LOCALTIME",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "datetime"},
		},

		{
			Name:       "LOCALTIMESTAMP",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "datetime"},
		},
		{
			Name: "LOCALTIMESTAMP",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "datetime"},
		},
		{
			Name: "LOCATE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "LOCATE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "LOG",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "LOG",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "LOG10",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "LOG2",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "LOWER",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "LPAD",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "LTRIM",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "MAKEDATE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "date"},
		},
		{
			Name: "MAKETIME",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "time"},
		},
		{
			Name: "MAKE_SET",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "MASTER_POS_WAIT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "MASTER_POS_WAIT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "MASTER_POS_WAIT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "MAX",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "MBRCONTAINS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "MBRCOVEREDBY",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "MBRCOVERS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "MBRDISJOINT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "MBREQUALS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "MBRINTERSECTS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "MBROVERLAPS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "MBRTOUCHES",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "MBRWITHIN",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "MD5",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "MICROSECOND",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "time"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "MID",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "MIN",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "MINUTE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "time"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "MOD",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "MONTH",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "date"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "MONTHNAME",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "date"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "MULTILINESTRING",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "MULTIPOINT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "MULTIPOLYGON",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "NAME_CONST",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "NOW",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "datetime"},
		},
		{
			Name: "NEXTVAL",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "NOW",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "datetime"},
		},
		{
			Name: "NTH_VALUE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "NTILE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "NULLIF",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "OCT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "bigint"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "OCTET_LENGTH",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "ORD",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name:       "PERCENT_RANK",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "PERIOD_ADD",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "PERIOD_DIFF",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name:       "PI",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "POINT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "POLYGON",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "POSITION",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "POW",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "POWER",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name:       "PS_CURRENT_THREAD_ID",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "PS_THREAD_ID",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "QUARTER",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "date"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "QUOTE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "RADIANS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name:       "RAND",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "RAND",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "RANDOM_BYTES",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "binary"},
		},
		{
			Name:       "RANK",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "REGEXP_INSTR",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type:       &ast.TypeName{Name: "int"},
					HasDefault: true,
				},
				{
					Type:       &ast.TypeName{Name: "int"},
					HasDefault: true,
				},
				{
					Type:       &ast.TypeName{Name: "int"},
					HasDefault: true,
				},
				{
					Type:       &ast.TypeName{Name: "text"},
					HasDefault: true,
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "REGEXP_LIKE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type:       &ast.TypeName{Name: "text"},
					HasDefault: true,
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "REGEXP_REPLACE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type:       &ast.TypeName{Name: "int"},
					HasDefault: true,
				},
				{
					Type:       &ast.TypeName{Name: "int"},
					HasDefault: true,
				},
				{
					Type:       &ast.TypeName{Name: "text"},
					HasDefault: true,
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "REGEXP_SUBSTR",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type:       &ast.TypeName{Name: "int"},
					HasDefault: true,
				},
				{
					Type:       &ast.TypeName{Name: "int"},
					HasDefault: true,
				},
				{
					Type:       &ast.TypeName{Name: "text"},
					HasDefault: true,
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name:       "RELEASE_ALL_LOCKS",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "RELEASE_LOCK",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "REPEAT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "REPLACE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "REVERSE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "RIGHT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name:       "ROLES_GRAPHML",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "ROUND",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "ROUND",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name:       "ROW_COUNT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "bigint"},
		},
		{
			Name:       "ROW_NUMBER",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "RPAD",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "RTRIM",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name:       "SCHEMA",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "SECOND",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "time"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "SEC_TO_TIME",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "time"},
		},
		{
			Name:       "SESSION_USER",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "SHA",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "SHA1",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "SHA2",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "SIGN",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "SIN",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "SLEEP",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "SOUNDEX",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "SPACE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "SQRT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "STATEMENT_DIGEST",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "STATEMENT_DIGEST_TEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "STD",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "STDDEV",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "STDDEV_POP",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "STDDEV_SAMP",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "STRCMP",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "tinyint"},
		},
		{
			Name: "STR_TO_DATE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "datetime"},
		},
		{
			Name: "ST_AREA",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "ST_ASBINARY",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "binary"},
		},
		{
			Name: "ST_ASBINARY",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "binary"},
		},
		{
			Name: "ST_ASGEOJSON",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "json"},
		},
		{
			Name: "ST_ASGEOJSON",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "json"},
		},
		{
			Name: "ST_ASGEOJSON",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "tinyint"},
				},
			},
			ReturnType: &ast.TypeName{Name: "json"},
		},
		{
			Name: "ST_ASTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "ST_ASTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "ST_ASWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "binary"},
		},
		{
			Name: "ST_ASWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "binary"},
		},
		{
			Name: "ST_ASWKT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "ST_ASWKT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "ST_BUFFER",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_BUFFER",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
				{
					Type: &ast.TypeName{Name: "binary"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_BUFFER",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "binary"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_BUFFER",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "binary"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_BUFFER_STRATEGY",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "binary"},
		},
		{
			Name: "ST_BUFFER_STRATEGY",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "binary"},
		},
		{
			Name: "ST_CENTROID",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_CONTAINS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "ST_CONVEXHULL",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_CROSSES",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "ST_DIFFERENCE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_DIMENSION",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "tinyint"},
		},
		{
			Name: "ST_DISJOINT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "ST_DISTANCE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "ST_DISTANCE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "ST_DISTANCE_SPHERE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "ST_DISTANCE_SPHERE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "ST_ENDPOINT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_ENVELOPE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_EQUALS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "ST_EXTERIORRING",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_FRECHET_DISTANCE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "ST_FRECHET_DISTANCE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "ST_GEOHASH",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "ST_GEOHASH",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "ST_GEOMCOLLFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMCOLLFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMCOLLFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMCOLLFROMTXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMCOLLFROMTXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMCOLLFROMTXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMCOLLFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMCOLLFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMCOLLFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMETRYCOLLECTIONFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMETRYCOLLECTIONFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMETRYCOLLECTIONFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMETRYCOLLECTIONFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMETRYCOLLECTIONFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMETRYCOLLECTIONFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMETRYFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMETRYFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMETRYFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMETRYFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMETRYFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMETRYFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMETRYN",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMETRYTYPE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "ST_GEOMFROMGEOJSON",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMFROMGEOJSON",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "tinyint"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMFROMGEOJSON",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "tinyint"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_GEOMFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_HAUSDORFF_DISTANCE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "ST_HAUSDORFF_DISTANCE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},

		{
			Name: "ST_INTERIORRINGN",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_INTERSECTION",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_INTERSECTS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "ST_ISCLOSED",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "ST_ISEMPTY",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "ST_ISSIMPLE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "ST_ISVALID",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "ST_LATFROMGEOHASH",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "ST_LATITUDE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "ST_LATITUDE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_LENGTH",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "ST_LENGTH",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "ST_LINEFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_LINEFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_LINEFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_LINEFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_LINEFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_LINEFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_LINEINTERPOLATEPOINT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_LINEINTERPOLATEPOINTS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_LINESTRINGFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_LINESTRINGFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_LINESTRINGFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_LINESTRINGFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_LINESTRINGFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_LINESTRINGFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_LONGFROMGEOHASH",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "ST_LONGITUDE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "ST_LONGITUDE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MAKEENVELOPE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MLINEFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MLINEFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MLINEFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MLINEFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MLINEFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MLINEFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MPOINTFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MPOINTFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MPOINTFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MPOINTFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MPOINTFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MPOINTFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MPOLYFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MPOLYFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MPOLYFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MPOLYFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MPOLYFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MPOLYFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MULTILINESTRINGFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MULTILINESTRINGFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MULTILINESTRINGFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MULTILINESTRINGFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MULTILINESTRINGFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MULTILINESTRINGFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MULTIPOINTFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MULTIPOINTFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MULTIPOINTFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MULTIPOINTFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MULTIPOINTFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MULTIPOINTFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MULTIPOLYGONFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MULTIPOLYGONFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MULTIPOLYGONFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MULTIPOLYGONFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MULTIPOLYGONFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_MULTIPOLYGONFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_NUMGEOMETRIES",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bigint"},
		},
		{
			Name: "ST_NUMINTERIORRING",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "ST_NUMINTERIORRINGS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "ST_NUMPOINTS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bigint"},
		},
		{
			Name: "ST_OVERLAPS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "ST_POINTATDISTANCE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_POINTFROMGEOHASH",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_POINTFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_POINTFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_POINTFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_POINTFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_POINTFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_POINTFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_POINTN",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_POLYFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_POLYFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_POLYFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_POLYFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_POLYFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_POLYFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_POLYGONFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_POLYGONFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_POLYGONFROMTEXT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_POLYGONFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_POLYGONFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_POLYGONFROMWKB",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_SIMPLIFY",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_SRID",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "ST_SRID",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "ST_STARTPOINT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_SWAPXY",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_SYMDIFFERENCE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_TOUCHES",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "ST_TRANSFORM",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_UNION",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_VALIDATE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_WITHIN",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "ST_X",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "ST_X",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "ST_Y",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "ST_Y",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "SUBDATE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "date"},
				},
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "date"},
		},
		{
			Name: "SUBSTR",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "SUBSTR",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "SUBSTRING",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "SUBSTRING",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "SUBSTRING_INDEX",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "SUBTIME",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "time"},
				},
				{
					Type: &ast.TypeName{Name: "time"},
				},
			},
			ReturnType: &ast.TypeName{Name: "time"},
		},
		{
			Name: "SUM",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "SYSDATE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "datetime"},
		},
		{
			Name: "SYSDATE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "datetime"},
		},
		{
			Name:       "SYSTEM_USER",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "TAN",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
			},
			ReturnType: &ast.TypeName{Name: "double precision"},
		},
		{
			Name: "TIME",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "datetime"},
				},
			},
			ReturnType: &ast.TypeName{Name: "time"},
		},
		{
			Name: "TIMEDIFF",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "time"},
				},
				{
					Type: &ast.TypeName{Name: "time"},
				},
			},
			ReturnType: &ast.TypeName{Name: "time"},
		},
		{
			Name: "TIMESTAMP",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "datetime"},
				},
			},
			ReturnType: &ast.TypeName{Name: "datetime"},
		},
		{
			Name: "TIMESTAMP",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "datetime"},
				},
				{
					Type: &ast.TypeName{Name: "datetime"},
				},
			},
			ReturnType: &ast.TypeName{Name: "datetime"},
		},
		{
			Name: "TIMESTAMPADD",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
				{
					Type: &ast.TypeName{Name: "datetime"},
				},
			},
			ReturnType: &ast.TypeName{Name: "datetime"},
		},
		{
			Name: "TIMESTAMPDIFF",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
				{
					Type: &ast.TypeName{Name: "datetime"},
				},
				{
					Type: &ast.TypeName{Name: "datetime"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "TIME_FORMAT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "time"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "TIME_TO_SEC",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "time"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "TO_BASE64",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "TO_DAYS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "date"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "TO_SECONDS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "datetime"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bigint"},
		},
		{
			Name: "TRIM",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "TRIM",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "TRUNCATE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "double precision"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "UCASE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "UNCOMPRESS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "UNCOMPRESSED_LENGTH",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "binary"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "UNHEX",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
		},
		{
			Name:       "UNIX_TIMESTAMP",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "bigint"},
		},
		{
			Name: "UNIX_TIMESTAMP",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "datetime"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bigint"},
		},
		{
			Name: "UPDATEXML",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "UPPER",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name:       "USER",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name:       "UTC_DATE",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "date"},
		},
		{
			Name:       "UTC_TIME",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "time"},
		},
		{
			Name:       "UTC_TIMESTAMP",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "datetime"},
		},
		{
			Name:       "UUID",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name:       "UUID_SHORT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "bigint"},
		},
		{
			Name: "UUID_TO_BIN",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "binary"},
		},
		{
			Name: "UUID_TO_BIN",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "tinyint"},
				},
			},
			ReturnType: &ast.TypeName{Name: "binary"},
		},
		{
			Name: "VALIDATE_PASSWORD_STRENGTH",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "VALUES",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "VARIANCE",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "VAR_POP",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "VAR_SAMP",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "VERSION",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "WAIT_FOR_EXECUTED_GTID_SET",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type:       &ast.TypeName{Name: "int"},
					HasDefault: true,
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "WAIT_UNTIL_SQL_THREAD_AFTER_GTIDS",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type:       &ast.TypeName{Name: "int"},
					HasDefault: true,
				},
				{
					Type:       &ast.TypeName{Name: "text"},
					HasDefault: true,
				},
			},
			ReturnType: &ast.TypeName{Name: "bool"},
		},
		{
			Name: "WEEK",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "date"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "WEEK",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "date"},
				},
				{
					Type: &ast.TypeName{Name: "int"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "WEEKDAY",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "date"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "WEEKOFYEAR",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "date"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "WEIGHT_STRING",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "YEAR",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "date"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
		{
			Name: "YEARWEEK",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "date"},
				},
			},
			ReturnType: &ast.TypeName{Name: "int"},
		},
	}
	return s
}
