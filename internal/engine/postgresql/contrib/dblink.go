// Code generated by sqlc-pg-gen. DO NOT EDIT.

package contrib

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
)

func Dblink() *catalog.Schema {
	s := &catalog.Schema{Name: "pg_catalog"}
	s.Funcs = []*catalog.Function{
		{
			Name: "dblink",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "record"},
		},
		{
			Name: "dblink",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "boolean"},
				},
			},
			ReturnType: &ast.TypeName{Name: "record"},
		},
		{
			Name: "dblink",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "record"},
		},
		{
			Name: "dblink",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "boolean"},
				},
			},
			ReturnType: &ast.TypeName{Name: "record"},
		},
		{
			Name: "dblink_build_sql_delete",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int2vector"},
				},
				{
					Type: &ast.TypeName{Name: "integer"},
				},
				{
					Type: &ast.TypeName{Name: "text[]"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "dblink_build_sql_insert",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int2vector"},
				},
				{
					Type: &ast.TypeName{Name: "integer"},
				},
				{
					Type: &ast.TypeName{Name: "text[]"},
				},
				{
					Type: &ast.TypeName{Name: "text[]"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "dblink_build_sql_update",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "int2vector"},
				},
				{
					Type: &ast.TypeName{Name: "integer"},
				},
				{
					Type: &ast.TypeName{Name: "text[]"},
				},
				{
					Type: &ast.TypeName{Name: "text[]"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "dblink_cancel_query",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "dblink_close",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "dblink_close",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "boolean"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "dblink_close",
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
			Name: "dblink_close",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "boolean"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "dblink_connect",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "dblink_connect",
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
			Name: "dblink_connect_u",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "dblink_connect_u",
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
			Name:       "dblink_current_query",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "dblink_disconnect",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name:       "dblink_disconnect",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "dblink_error_message",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "dblink_exec",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "dblink_exec",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "boolean"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "dblink_exec",
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
			Name: "dblink_exec",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "boolean"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "dblink_fdw_validator",
			Args: []*catalog.Argument{
				{
					Name: "options",
					Type: &ast.TypeName{Name: "text[]"},
				},
				{
					Name: "catalog",
					Type: &ast.TypeName{Name: "oid"},
				},
			},
			ReturnType: &ast.TypeName{Name: "void"},
		},
		{
			Name: "dblink_fetch",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "integer"},
				},
			},
			ReturnType: &ast.TypeName{Name: "record"},
		},
		{
			Name: "dblink_fetch",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "integer"},
				},
				{
					Type: &ast.TypeName{Name: "boolean"},
				},
			},
			ReturnType: &ast.TypeName{Name: "record"},
		},
		{
			Name: "dblink_fetch",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "integer"},
				},
			},
			ReturnType: &ast.TypeName{Name: "record"},
		},
		{
			Name: "dblink_fetch",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "integer"},
				},
				{
					Type: &ast.TypeName{Name: "boolean"},
				},
			},
			ReturnType: &ast.TypeName{Name: "record"},
		},
		{
			Name:       "dblink_get_connections",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "text[]"},
		},
		{
			Name: "dblink_get_pkey",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "dblink_pkey_results"},
		},
		{
			Name: "dblink_get_result",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "record"},
		},
		{
			Name: "dblink_get_result",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "boolean"},
				},
			},
			ReturnType: &ast.TypeName{Name: "record"},
		},
		{
			Name: "dblink_is_busy",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "integer"},
		},
		{
			Name: "dblink_open",
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
			Name: "dblink_open",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "boolean"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "dblink_open",
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
			Name: "dblink_open",
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
					Type: &ast.TypeName{Name: "boolean"},
				},
			},
			ReturnType: &ast.TypeName{Name: "text"},
		},
		{
			Name: "dblink_send_query",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "text"},
				},
				{
					Type: &ast.TypeName{Name: "text"},
				},
			},
			ReturnType: &ast.TypeName{Name: "integer"},
		},
	}
	return s
}
