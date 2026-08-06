package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/core/seed"

	"github.com/jackc/pgx/v5"
)

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

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func clean(arg string) string {
	arg = strings.TrimSpace(arg)
	arg = strings.Replace(arg, "\"any\"", "any", -1)
	arg = strings.Replace(arg, "\"char\"", "char", -1)
	arg = strings.Replace(arg, "\"timestamp\"", "char", -1)
	return arg
}

// writeFunctions writes procs to destPath as JSONL, one function per line, in
// the form the seed package reads. Both the analysis core and the catalog the
// legacy compiler builds load the result.
func writeFunctions(procs []Proc, destPath string) error {
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	enc := json.NewEncoder(out)
	for _, proc := range procs {
		fn := seed.Function{Name: proc.Name, Returns: proc.ReturnTypeName()}
		for _, arg := range proc.Args() {
			a := seed.Arg{Name: arg.Name, Type: arg.TypeName(), HasDefault: arg.HasDefault}
			if arg.Mode != "" && arg.Mode != "i" {
				a.Mode = arg.Mode
			}
			fn.Args = append(fn.Args, a)
		}
		if err := enc.Encode(fn); err != nil {
			return err
		}
	}
	return nil
}

// writeRelations appends a schema's relations to out as JSONL, one relation
// per line, in the form the seed package reads.
func writeRelations(out io.Writer, schemaName string, relations []Relation) error {
	enc := json.NewEncoder(out)
	for _, relation := range relations {
		rec := seed.Relation{
			Catalog: relation.Catalog,
			Schema:  schemaName,
			Name:    relation.Name,
		}
		for _, col := range relation.Columns {
			c := seed.Column{
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

func databaseURL() string {
	dburl := os.Getenv("DATABASE_URL")
	if dburl != "" {
		return dburl
	}
	pgUser := os.Getenv("PG_USER")
	pgHost := os.Getenv("PG_HOST")
	pgPort := os.Getenv("PG_PORT")
	pgPass := os.Getenv("PG_PASSWORD")
	pgDB := os.Getenv("PG_DATABASE")
	if pgUser == "" {
		pgUser = "postgres"
	}
	if pgPass == "" {
		pgPass = "mysecretpassword"
	}
	if pgPort == "" {
		pgPort = "5432"
	}
	if pgHost == "" {
		pgHost = "127.0.0.1"
	}
	if pgDB == "" {
		pgDB = "dinotest"
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", pgUser, pgPass, pgHost, pgPort, pgDB)
}

func run(ctx context.Context) error {
	flag.Parse()

	dir := flag.Arg(0)
	if dir == "" {
		dir = filepath.Join("internal", "engine", "postgresql")
	}

	conn, err := pgx.Connect(ctx, databaseURL())
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	// The two schemas sqlc knows PostgreSQL by are written to the dialect
	// directory rather than to Go: the engine and the analysis core both read
	// them from there. Their relations share one file, keyed by schema.
	dialectDir := filepath.Join(dir, "dialect")
	relationsPath := filepath.Join(dialectDir, seed.RelationsFile)
	schemas := []schemaToLoad{
		{
			Name:      "pg_catalog",
			FuncsPath: filepath.Join(dialectDir, seed.FunctionsFile),
		},
		{
			Name: "information_schema",
		},
	}

	relationsFile, err := os.Create(relationsPath)
	if err != nil {
		return err
	}
	defer relationsFile.Close()

	for _, schema := range schemas {
		procs, err := readProcs(ctx, conn, schema.Name)
		if err != nil {
			return err
		}

		if schema.Name == "pg_catalog" {
			procs = preserveLegacyCatalogBehavior(procs)
		}

		relations, err := readRelations(ctx, conn, schema.Name)
		if err != nil {
			return err
		}

		if schema.FuncsPath != "" {
			if err := writeFunctions(procs, schema.FuncsPath); err != nil {
				return err
			}
		}
		if err := writeRelations(relationsFile, schema.Name, relations); err != nil {
			return err
		}
	}

	// Each extension is a directory of its own under the dialect, holding the
	// functions CREATE EXTENSION adds to the catalog.
	for _, extension := range extensions {
		if _, err := conn.Exec(ctx, fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS %q", extension)); err != nil {
			return fmt.Errorf("error creating %s: %s", extension, err)
		}

		rows, err := conn.Query(ctx, extensionFuncs, extension)
		if err != nil {
			return err
		}
		procs, err := scanProcs(rows)
		if err != nil {
			return err
		}
		if len(procs) == 0 {
			log.Printf("no functions in %s, skipping", extension)
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

		extensionDir := filepath.Join(dialectDir, "extensions", extension)
		if err := os.MkdirAll(extensionDir, 0o755); err != nil {
			return err
		}
		if err := writeFunctions(procs, filepath.Join(extensionDir, seed.FunctionsFile)); err != nil {
			return fmt.Errorf("error generating extension %s: %w", extension, err)
		}
	}

	return nil
}

type schemaToLoad struct {
	// name is the name of a schema to load
	Name string
	// FuncsPath, when set, is the JSONL file this schema's functions are
	// written to.
	FuncsPath string
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
