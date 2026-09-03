package postgresql

import (
	"context"
	"strings"

	pgx "github.com/jackc/pgx/v5"
)

// https://stackoverflow.com/questions/25308765/postgresql-how-can-i-inspect-which-arguments-to-a-procedure-have-a-default-valu
const catalogFuncs = `
SELECT p.proname as name,
  format_type(p.prorettype, NULL),
  array(select format_type(unnest(p.proargtypes), NULL)),
  p.proargnames,
  p.proargnames[p.pronargs-p.pronargdefaults+1:p.pronargs],
  p.proargmodes::text[]
FROM pg_catalog.pg_proc p
LEFT JOIN pg_catalog.pg_namespace n ON n.oid = p.pronamespace
WHERE n.nspname::text = $1
  AND pg_function_is_visible(p.oid)
  -- A function that one of the extensions the dialect describes separately
  -- put in the schema, as adminpack puts its own in pg_catalog, belongs to
  -- that extension's directory: leaving it out here keeps the output the
  -- same whether or not the extension is installed.
  AND NOT EXISTS (
    SELECT 1
    FROM pg_catalog.pg_depend d
    INNER JOIN pg_catalog.pg_extension e ON e.oid = d.refobjid
    WHERE d.classid = 'pg_catalog.pg_proc'::regclass
      AND d.objid = p.oid
      AND d.deptype = 'e'
      AND e.extname = ANY($2))
-- simply order all columns to keep subsequent runs stable
ORDER BY 1, 2, 3, 4, 5;
`

type Proc struct {
	Name       string
	ReturnType string
	ArgTypes   []string
	ArgNames   []string
	HasDefault []string
	ArgModes   []string
}

func (p *Proc) ReturnTypeName() string {
	return clean(p.ReturnType)
}

func (p *Proc) Args() []Arg {
	var args []Arg
	defaults := map[string]bool{}
	for _, name := range p.HasDefault {
		defaults[name] = true
	}

	for i, argType := range p.ArgTypes {
		mode := "i"
		name := ""
		if i < len(p.ArgModes) {
			mode = p.ArgModes[i]
		}
		if i < len(p.ArgNames) {
			name = p.ArgNames[i]
		}

		args = append(args, Arg{
			Name:       name,
			Type:       argType,
			Mode:       mode,
			HasDefault: defaults[name],
		})
	}

	// Some manual changes until https://github.com/sqlc-dev/sqlc/pull/1748
	// can be completely implmented
	if p.Name == "mode" {
		return nil
	}

	if p.Name == "percentile_cont" && len(args) == 2 {
		args = args[:1]
	}

	if p.Name == "percentile_disc" && len(args) == 2 {
		args = args[:1]
	}

	return args
}

type Arg struct {
	Name       string
	Mode       string
	Type       string
	HasDefault bool
}

func (a *Arg) TypeName() string {
	return clean(a.Type)
}

func scanProcs(rows pgx.Rows) ([]Proc, error) {
	defer rows.Close()
	// Iterate through the result set
	var procs []Proc
	for rows.Next() {
		var p Proc
		err := rows.Scan(
			&p.Name,
			&p.ReturnType,
			&p.ArgTypes,
			&p.ArgNames,
			&p.HasDefault,
			&p.ArgModes,
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

		procs = append(procs, p)
	}
	return procs, rows.Err()
}

// readProcs reads the functions of a schema, leaving out the ones that
// belong to any of the named extensions.
func readProcs(ctx context.Context, conn *pgx.Conn, schemaName string, extensions []string) ([]Proc, error) {
	rows, err := conn.Query(ctx, catalogFuncs, schemaName, extensions)
	if err != nil {
		return nil, err
	}

	return scanProcs(rows)
}
