package main

import (
	"bytes"
	"context"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	pgx "github.com/jackc/pgx/v4"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
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

const catalogTmpl = `
package postgresql

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
)

func gen{{.Name}}() *catalog.Schema {
	s := &catalog.Schema{Name: "pg_catalog"}
	s.Funcs = []*catalog.Function{
	    {{- range .Funcs}}
		{
			Name: "{{.Name}}",
			Args: []*catalog.Argument{
				{{range .Args}}{
				{{- if .Name}}
				Name: "{{.Name}}",
				{{- end}}
				{{- if .HasDefault}}
				HasDefault: true,
				{{- end}}
				Type: &ast.TypeName{Name: "{{.Type.Name}}"},
				},
				{{end}}
			},
			ReturnType: &ast.TypeName{Name: "{{.ReturnType.Name}}"},
		},
		{{- end}}
	}
	return s
}
`

type tmplCtx struct {
	Name  string
	Funcs []catalog.Function
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

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

func (p Proc) Func() catalog.Function {
	return catalog.Function{
		Name:       p.Name,
		Args:       p.Args(),
		ReturnType: &ast.TypeName{Name: clean(p.ReturnType)},
	}
}

func (p Proc) Args() []*catalog.Argument {
	defaults := map[string]bool{}
	var args []*catalog.Argument
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
		args = append(args, &catalog.Argument{
			Name:       name,
			HasDefault: defaults[name],
			Type:       &ast.TypeName{Name: clean(arg)},
		})
	}
	return args
}

func scanFuncs(rows pgx.Rows) ([]catalog.Function, error) {
	defer rows.Close()
	// Iterate through the result set
	var funcs []catalog.Function
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

func run(ctx context.Context) error {
	tmpl, err := template.New("").Parse(catalogTmpl)
	if err != nil {
		return err
	}
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	// Generate internal/engine/postgresql/pg_catalog.gen.go
	rows, err := conn.Query(ctx, catalogFuncs)
	if err != nil {
		return err
	}
	funcs, err := scanFuncs(rows)
	if err != nil {
		return err
	}
	out := bytes.NewBuffer([]byte{})
	if err := tmpl.Execute(out, tmplCtx{Name: "PGCatalog", Funcs: funcs}); err != nil {
		return err
	}
	code, err := format.Source(out.Bytes())
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join("internal", "engine", "postgresql", "pg_catalog.gen.go"), code, 0644)
	if err != nil {
		return err
	}

	// https://www.postgresql.org/docs/current/contrib.html
	extensions := map[string]string{
		"citext":    "CIText",
		"pg_trgm":   "PGTrigram",
		"pgcrypto":  "PGCrypto",
		"uuid-ossp": "UUIDOSSP",
	}

	for extension, name := range extensions {
		_, err := conn.Exec(ctx, fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS \"%s\"", extension))
		if err != nil {
			return err
		}
		rows, err := conn.Query(ctx, extensionFuncs, extension)
		funcs, err := scanFuncs(rows)
		if err != nil {
			return err
		}
		out := bytes.NewBuffer([]byte{})
		if err := tmpl.Execute(out, tmplCtx{Name: name, Funcs: funcs}); err != nil {
			return err
		}
		code, err := format.Source(out.Bytes())
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filepath.Join("internal", "engine", "postgresql", "extension_"+strings.Replace(extension, "-", "_", -1)+".gen.go"), code, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
