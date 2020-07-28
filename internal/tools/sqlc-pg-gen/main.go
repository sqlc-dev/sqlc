package main

import (
	"bytes"
	"context"
	"fmt"
	"go/format"
	"log"
	"os"
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

const catalogTmpl = `
package postgresql

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
)

func genPGCatalog() *catalog.Schema {
	s := &catalog.Schema{Name: "pg_catalog"}
	s.Funcs = []*catalog.Function{
	    {{- range .}}
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

	rows, err := conn.Query(ctx, catalogFuncs)
	if err != nil {
		return err
	}

	defer rows.Close()

	// Iterate through the result set
	var funcs []catalog.Function
	for rows.Next() {
		var p Proc
		err = rows.Scan(
			&p.Name,
			&p.ReturnType,
			&p.ArgTypes,
			&p.ArgNames,
			&p.HasDefault,
		)
		if err != nil {
			return err
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
		for i := range p.ArgTypes {
			if p.ArgTypes[i] == "internal" {
				continue
			}
		}

		funcs = append(funcs, p.Func())
	}

	if rows.Err() != nil {
		return err
	}

	out := bytes.NewBuffer([]byte{})
	if err := tmpl.Execute(out, funcs); err != nil {
		return err
	}
	code, err := format.Source(out.Bytes())
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(os.Stdout, string(code))
	return err
}
