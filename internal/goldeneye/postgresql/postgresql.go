// Package postgresql generates the PostgreSQL dialect seed under
// internal/engine/postgresql/dialect from a live server: pg_catalog's
// functions, the relations of pg_catalog and information_schema, and one
// directory per contrib extension holding the types and functions CREATE
// EXTENSION adds. types.jsonl and operators.jsonl at the top of the dialect
// are written by hand and are not this package's business.
//
// The server is named by POSTGRESQL_SERVER_URI and has to be the major
// release in Major, since each release adds to its catalogs.
package postgresql

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/dialect"
)

// Engine is the name of the engine directory the dialect lives under.
const Engine = "postgresql"

// Major is the PostgreSQL major release the dialect is generated from.
// Bumping it is a deliberate change: every release adds functions and
// catalog columns, so regenerate and review the dialect after changing it.
const Major = 16

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
  p.proargnames[p.pronargs-p.pronargdefaults+1:p.pronargs],
  p.proargmodes::text[]
FROM pg_catalog.pg_proc p
JOIN extension_funcs ef ON ef.oid = p.oid
WHERE pg_function_is_visible(p.oid)
-- simply order all columns to keep subsequent runs stable
ORDER BY 1, 2, 3, 4, 5;
`

// The types an extension defines, in the order they were declared so that a
// type appears after the types it is built on. Array types are created
// implicitly rather than by the extension's script, so they carry an internal
// dependency and never show up here.
const extensionTypes = `
SELECT t.typname, t.typcategory
FROM pg_catalog.pg_extension AS e
    INNER JOIN pg_catalog.pg_depend AS d ON (d.refobjid = e.oid)
    INNER JOIN pg_catalog.pg_type AS t ON (t.oid = d.objid)
WHERE d.deptype = 'e'
  AND d.classid = 'pg_catalog.pg_type'::regclass
  AND e.extname = $1
  AND t.typtype <> 'p'
ORDER BY t.oid;
`

// Locate returns the server to generate from, named by POSTGRESQL_SERVER_URI.
func Locate() (string, error) {
	if url := os.Getenv("POSTGRESQL_SERVER_URI"); url != "" {
		return url, nil
	}
	return "", errors.New("POSTGRESQL_SERVER_URI is not set")
}

// Version reports the release a server is.
func Version(ctx context.Context, url string) (string, error) {
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return "", err
	}
	defer conn.Close(ctx)
	var version string
	if err := conn.QueryRow(ctx, "SELECT version()").Scan(&version); err != nil {
		return "", err
	}
	return version, nil
}

// checkVersion refuses a server of another major release than the dialect
// is generated from, whose catalogs would differ from the committed ones
// without anything being wrong.
func checkVersion(ctx context.Context, conn *pgx.Conn) error {
	var num string
	if err := conn.QueryRow(ctx, "SELECT current_setting('server_version_num')").Scan(&num); err != nil {
		return err
	}
	n, err := strconv.Atoi(num)
	if err != nil {
		return fmt.Errorf("server_version_num %q: %w", num, err)
	}
	if major := n / 10000; major != Major {
		return fmt.Errorf("the dialect is generated from PostgreSQL %d, but the server is PostgreSQL %d", Major, major)
	}
	return nil
}

func clean(arg string) string {
	arg = strings.TrimSpace(arg)
	arg = strings.Replace(arg, "\"any\"", "any", -1)
	arg = strings.Replace(arg, "\"char\"", "char", -1)
	arg = strings.Replace(arg, "\"timestamp\"", "char", -1)
	return arg
}

// encodeFunctions encodes procs as JSONL, one function per line, in the form
// the seed package reads. Both the analysis core and the catalog the legacy
// compiler builds load the result.
func encodeFunctions(procs []Proc) ([]byte, error) {
	funcs := make([]dialect.Function, 0, len(procs))
	for _, proc := range procs {
		fn := dialect.Function{Name: proc.Name, Returns: proc.ReturnTypeName()}
		for _, arg := range proc.Args() {
			a := dialect.Arg{Name: arg.Name, Type: arg.TypeName(), HasDefault: arg.HasDefault}
			if arg.Mode != "" && arg.Mode != "i" {
				a.Mode = arg.Mode
			}
			fn.Args = append(fn.Args, a)
		}
		funcs = append(funcs, fn)
	}
	return dialect.JSONL(funcs)
}

// readExtensionTypes reads the types an extension defines.
func readExtensionTypes(ctx context.Context, conn *pgx.Conn, extension string) ([]dialect.Type, error) {
	rows, err := conn.Query(ctx, extensionTypes, extension)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var types []dialect.Type
	for rows.Next() {
		var t dialect.Type
		if err := rows.Scan(&t.Name, &t.Category); err != nil {
			return nil, err
		}
		types = append(types, t)
	}
	return types, rows.Err()
}

// writeRelations appends a schema's relations to out as JSONL, one relation
// per line, in the form the seed package reads.
func writeRelations(out *bytes.Buffer, schemaName string, relations []Relation) error {
	enc := json.NewEncoder(out)
	for _, relation := range relations {
		rec := dialect.Relation{
			Catalog: relation.Catalog,
			Schema:  schemaName,
			Name:    relation.Name,
		}
		for _, col := range relation.Columns {
			c := dialect.Column{
				Name:    col.Name,
				Type:    col.Type,
				NotNull: col.IsNotNull,
				Array:   col.IsArray,
			}
			if col.Length != nil {
				c.Length = *col.Length
			}
			rec.Columns = append(rec.Columns, c)
		}
		if err := enc.Encode(rec); err != nil {
			return err
		}
	}
	return nil
}

// preserveLegacyCatalogBehavior maintain previous ordering and filtering
// that was manually done to the generated file pg_catalog.go.
// Some of the test depend on this ordering - in particular, function lookups
// where there might be multiple matching functions (due to type overloads)
// Until sqlc supports "smarter" looking up of these functions,
// preserveLegacyCatalogBehavior ensures there are no accidental test breakages
func preserveLegacyCatalogBehavior(allProcs []Proc) []Proc {
	// Preserve the legacy sort order of the end-to-end tests
	sort.SliceStable(allProcs, func(i, j int) bool {
		fnA := allProcs[i]
		fnB := allProcs[j]

		if fnA.Name == "lower" && fnB.Name == "lower" && len(fnA.ArgTypes) == 1 && fnA.ArgTypes[0] == "text" {
			return true
		}

		if fnA.Name == "generate_series" && fnB.Name == "generate_series" && len(fnA.ArgTypes) == 2 && fnA.ArgTypes[0] == "numeric" {
			return true
		}

		return false
	})

	procs := make([]Proc, 0, len(allProcs))
	for _, p := range allProcs {
		// Skip generating pg_catalog.concat to preserve legacy behavior
		if p.Name == "concat" {
			continue
		}

		procs = append(procs, p)
	}

	return procs
}

// Generate reads the dialect from the server. It creates every contrib
// extension it describes, so the server's contrib package has to be
// installed.
func Generate(ctx context.Context, url string) (dialect.Files, error) {
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)
	if err := checkVersion(ctx, conn); err != nil {
		return nil, err
	}

	files := dialect.Files{}

	// The two schemas sqlc knows PostgreSQL by. pg_catalog's functions are
	// the dialect's standard library; the relations of both share one file,
	// keyed by schema.
	procs, err := readProcs(ctx, conn, "pg_catalog", extensions)
	if err != nil {
		return nil, err
	}
	if files[dialect.FunctionsFile], err = encodeFunctions(preserveLegacyCatalogBehavior(procs)); err != nil {
		return nil, err
	}
	var relations bytes.Buffer
	for _, schema := range []string{"pg_catalog", "information_schema"} {
		rels, err := readRelations(ctx, conn, schema)
		if err != nil {
			return nil, err
		}
		if err := writeRelations(&relations, schema, rels); err != nil {
			return nil, err
		}
	}
	files[dialect.RelationsFile] = relations.Bytes()

	// Each extension is a directory of its own under the dialect, holding the
	// functions CREATE EXTENSION adds to the catalog.
	for _, extension := range extensions {
		if _, err := conn.Exec(ctx, fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS %q", extension)); err != nil {
			return nil, fmt.Errorf("error creating %s: %s", extension, err)
		}

		rows, err := conn.Query(ctx, extensionFuncs, extension)
		if err != nil {
			return nil, err
		}
		procs, err := scanProcs(rows)
		if err != nil {
			return nil, err
		}
		types, err := readExtensionTypes(ctx, conn, extension)
		if err != nil {
			return nil, err
		}
		if len(procs) == 0 && len(types) == 0 {
			continue
		}

		// Preserve the legacy sort order of the end-to-end tests
		sort.SliceStable(procs, func(i, j int) bool {
			fnA := procs[i]
			fnB := procs[j]

			if extension == "pgcrypto" {
				if fnA.Name == "digest" && fnB.Name == "digest" && len(fnA.ArgTypes) == 2 && fnA.ArgTypes[0] == "text" {
					return true
				}
			}

			return false
		})

		extensionDir := path.Join(dialect.ExtensionsDir, extension)
		if len(types) > 0 {
			if files[path.Join(extensionDir, dialect.TypesFile)], err = dialect.JSONL(types); err != nil {
				return nil, fmt.Errorf("error generating extension %s: %w", extension, err)
			}
		}
		if files[path.Join(extensionDir, dialect.FunctionsFile)], err = encodeFunctions(procs); err != nil {
			return nil, fmt.Errorf("error generating extension %s: %w", extension, err)
		}
	}

	return files, nil
}

// https://www.postgresql.org/docs/current/contrib.html
var extensions = []string{
	"adminpack",
	"amcheck",
	// "auth_delay",
	// "auto_explain",
	// "bloom",
	"btree_gin",
	"btree_gist",
	"citext",
	"cube",
	"dblink",
	// "dict_int",
	// "dict_xsyn",
	"earthdistance",
	"file_fdw",
	"fuzzystrmatch",
	"hstore",
	"intagg",
	"intarray",
	"isn",
	"lo",
	"ltree",
	"pageinspect",
	// "passwordcheck",
	"pg_buffercache",
	"pg_freespacemap",
	"pg_prewarm",
	"pg_stat_statements",
	"pg_trgm",
	"pg_visibility",
	"pgcrypto",
	"pgrowlocks",
	"pgstattuple",
	"postgres_fdw",
	"seg",
	// "sepgsql",
	// "spi",
	"sslinfo",
	"tablefunc",
	"tcn",
	// "test_decoding",
	// "tsm_system_rows",
	// "tsm_system_time",
	"unaccent",
	"uuid-ossp",
	"xml2",
}
