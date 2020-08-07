package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	sqllog "github.com/kyleconroy/sqlc/internal/sql/catalog"
)

// https://stackoverflow.com/questions/25308765/postgresql-how-can-i-inspect-which-arguments-to-a-procedure-have-a-default-valu
const catalogFuncs = `
SELECT p.proname as name,
  format_type(p.prorettype, NULL),
  array(select format_type(unnest(p.proargtypes), NULL)),
  p.proargnames,
  p.proargnames[p.pronargs-p.pronargdefaults+1:p.pronargs]
FROM pg_catalog.pg_proc p
LEFT JOIN pg_catalog.pg_namespace n ON n.oid = p.pronamespace
WHERE n.nspname OPERATOR(pg_catalog.~) '^(pg_catalog)$'
  AND p.proargmodes IS NULL
  AND pg_function_is_visible(p.oid)
ORDER BY 1;
`

// https://dba.stackexchange.com/questions/255412/how-to-select-functions-that-belong-in-a-given-extension-in-postgresql
//
// Extension functions are added to the public schema
const extensionFuncs = `
WITH extension_funcs AS (
  SELECT p.oid
  FROM pg_catalog.pg_extension AS e
      INNER JOIN pg_catalog.pg_depend AS d ON (d.refobjid = e.oid)
      INNER JOIN pg_catalog.pg_proc AS p ON (p.oid = d.objid)
      INNER JOIN pg_catalog.pg_namespace AS ne ON (ne.oid = e.extnamespace)
      INNER JOIN pg_catalog.pg_namespace AS np ON (np.oid = p.pronamespace)
  WHERE d.deptype = 'e' AND e.extname = $1
)
SELECT p.proname as name,
  format_type(p.prorettype, NULL),
  array(select format_type(unnest(p.proargtypes), NULL)),
  p.proargnames,
  p.proargnames[p.pronargs-p.pronargdefaults+1:p.pronargs]
FROM pg_catalog.pg_proc p
JOIN extension_funcs ef ON ef.oid = p.oid
WHERE p.proargmodes IS NULL
  AND pg_function_is_visible(p.oid)
ORDER BY 1;
`

// TODO: List the name of all installed extensions
type Proc struct {
	Name       string
	ReturnType string
	ArgTypes   []string
	ArgNames   []string
	HasDefault []string
}

func clean(arg string) string {
	arg = strings.TrimSpace(arg)
	arg = strings.Replace(arg, "\"any\"", "any", -1)
	arg = strings.Replace(arg, "\"char\"", "char", -1)
	arg = strings.Replace(arg, "\"timestamp\"", "char", -1)
	return arg
}

func (p Proc) Func() *sqllog.Function {
	return &sqllog.Function{
		Name:       p.Name,
		Args:       p.Args(),
		ReturnType: &ast.TypeName{Name: clean(p.ReturnType)},
	}
}

func (p Proc) Args() []*sqllog.Argument {
	defaults := map[string]bool{}
	var args []*sqllog.Argument
	if len(p.ArgTypes) == 0 {
		return args
	}
	for _, name := range p.HasDefault {
		defaults[name] = true
	}
	for i, arg := range p.ArgTypes {
		var name string
		if i < len(p.ArgNames) {
			name = p.ArgNames[i]
		}
		args = append(args, &sqllog.Argument{
			Name:       name,
			HasDefault: defaults[name],
			Type:       &ast.TypeName{Name: clean(arg)},
		})
	}
	return args
}

func scanFuncs(rows *sql.Rows) ([]*sqllog.Function, error) {
	defer rows.Close()
	// Iterate through the result set
	var funcs []*sqllog.Function
	for rows.Next() {
		var p Proc
		err := rows.Scan(
			&p.Name,
			&p.ReturnType,
			&p.ArgTypes,
			&p.ArgNames,
			&p.HasDefault,
		)
		if err != nil {
			return nil, err
		}

		// TODO: Filter these out in SQL
		if strings.HasPrefix(p.ReturnType, "SETOF") {
			continue
		}

		// The internal pseudo-type is used to declare functions that are meant
		// only to be called internally by the database system, and not by
		// direct invocation in an SQL query. If a function has at least one
		// internal-type argument then it cannot be called from SQL. To
		// preserve the type safety of this restriction it is important to
		// follow this coding rule: do not create any function that is declared
		// to return internal unless it has at least one internal argument
		//
		// https://www.postgresql.org/docs/current/datatype-pseudo.html
		var skip bool
		for i := range p.ArgTypes {
			if p.ArgTypes[i] == "internal" {
				skip = true
			}
		}
		if skip {
			continue
		}
		if p.ReturnType == "internal" {
			continue
		}

		funcs = append(funcs, p.Func())
	}
	return funcs, rows.Err()
}

func Load(ctx context.Context, db *sql.DB) (*sqllog.Catalog, error) {
	rows, err := db.QueryContext(ctx, catalogFuncs)
	if err != nil {
		return nil, err
	}
	funcs, err := scanFuncs(rows)
	if err != nil {
		return nil, err
	}
	// TODO: Load the list of installed extensions instead
	for _, extension := range contrib.SuppliedModules {
		rows, err := db.QueryContext(ctx, extensionFuncs, extension)
		if err != nil {
			return nil, fmt.Errorf("extension %s: %w", extension, err)
		}
		extFuncs, err := scanFuncs(rows)
		if err != nil {
			return nil, fmt.Errorf("extension %s: %w", extension, err)
		}
		// TODO: Add the Extension name to the function itself
		funcs = append(funcs, extFuncs...)
	}
	c := &sqllog.Catalog{
		DefaultSchema: "public",
		Schemas: []*sqllog.Schema{
			{
				Name:  "public",
				Funcs: funcs,
			},
		},
	}
	return c, nil
}
