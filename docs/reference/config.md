# Configuration file (version 1)

The `sqlc` tool is configured via a `sqlc.yaml` or `sqlc.json` file. This file must be
in the directory where the `sqlc` command is run.

```yaml
version: "1"
packages:
  - name: "db"
    path: "internal/db"
    queries: "./sql/query/"
    schema: "./sql/schema/"
    engine: "postgresql"
    emit_json_tags: true
    emit_prepared_queries: true
    emit_interface: false
    emit_exact_table_names: false
    emit_empty_slices: false
```

Each package document has the following keys:
- `name`:
  - The package name to use for the generated code. Defaults to `path` basename
- `path`:
  - Output directory for generated code
- `queries`:
  - Directory of SQL queries or path to single SQL file; or a list of paths
- `schema`:
  - Directory of SQL migrations or path to single SQL file; or a list of paths
- `engine`:
  - Either `postgresql` or `mysql`. Defaults to `postgresql`. MySQL support is experimental
- `emit_json_tags`:
  - If true, add JSON tags to generated structs. Defaults to `false`.
- `emit_db_tags`:
  - If true, add DB tags to generated structs. Defaults to `false`.
- `emit_prepared_queries`:
  - If true, include support for prepared queries. Defaults to `false`.
- `emit_interface`:
  - If true, output a `Querier` interface in the generated package. Defaults to `false`.
- `emit_exact_table_names`:
  - If true, struct names will mirror table names. Otherwise, sqlc attempts to singularize plural table names. Defaults to `false`.
- `emit_empty_slices`:
  - If true, slices returned by `:many` queries will be empty instead of `nil`. Defaults to `false`.

## Type Overrides

The default mapping of PostgreSQL types to Go types only uses packages outside
the standard library when it must.

For example, the `uuid` PostgreSQL type is mapped to `github.com/google/uuid`.
If a different Go package for UUIDs is required, specify the package in the
`overrides` array. In this case, I'm going to use the `github.com/gofrs/uuid`
instead.

```yaml
version: "1"
packages: [...]
overrides:
  - go_type: "github.com/gofrs/uuid.UUID"
    db_type: "uuid"
```

Each override document has the following keys:
- `db_type`:
  - The PostgreSQL type to override. Find the full list of supported types in [postgresql_type.go](https://github.com/kyleconroy/sqlc/blob/master/internal/codegen/golang/postgresql_type.go#L12).
- `go_type`:
  - A fully qualified name to a Go type to use in the generated code.
- `nullable`:
  - If true, use this type when a column is nullable. Defaults to `false`.

## Per-Column Type Overrides

Sometimes you would like to override the Go type used in model or query generation for
a specific field of a table and not on a type basis as described in the previous section.

This may be configured by specifying the `column` property in the override definition. `column`
should be of the form `table.column` buy you may be even more specify by specifying `schema.table.column`
or `catalog.schema.table.column`.

```yaml
version: "1"
packages: [...]
overrides:
  - column: "authors.id"
    go_type: "github.com/segmentio/ksuid.KSUID"
```

## Package Level Overrides

Overrides can be configured globally, as demonstrated in the previous sections, or they can be configured on a per-package which
scopes the override behavior to just a single package:

```yaml
version: "1"
packages:
  - overrides: [...]
```

## Renaming Struct Fields

Struct field names are generated from column names using a simple algorithm:
split the column name on underscores and capitalize the first letter of each
part.

```
account     -> Account
spotify_url -> SpotifyUrl
app_id      -> AppID
```

If you're not happy with a field's generated name, use the `rename` dictionary
to pick a new name. The keys are column names and the values are the struct
field name to use.

```yaml
version: "1"
packages: [...]
rename:
  spotify_url: "SpotifyURL"
```


