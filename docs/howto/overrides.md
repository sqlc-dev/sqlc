# Overriding types

:::{note}
Type overrides and field renaming are only fully-supported for Go.
:::

In many cases it's useful to tell `sqlc` explicitly what Go type you want it to
use for a query input or output. For instance, by default when you use
`pgx/v5`, `sqlc` will map a PostgreSQL UUID type to `UUID` from `github.com/jackc/pgx/pgtype`.
But you may want `sqlc` to use `UUID` from `github.com/google/uuid` instead.

To tell `sqlc` to use a different Go type, add an entry to the `overrides` list in your
configuration.

`sqlc` offers two kinds of Go type overrides:
* `db_type` overrides, which override the Go type for a specific database type.
* `column` overrides, which override the Go type for a column or columns by name.

Here's an example including one of each kind:

```yaml
version: "2"
sql:
- schema: "postgresql/schema.sql"
  queries: "postgresql/query.sql"
  engine: "postgresql"
  gen:
    go: 
      package: "authors"
      out: "db"
      sql_package: "pgx/v5"
      overrides:
        - db_type: "uuid"
          nullable: true
          go_type:
            import: "github.com/google/uuid"
            type: "UUID"
        - column: "users.birthday"
          go_type: "time.Time"
```

:::{tip}
  A single `db_type` override configuration applies to either nullable or non-nullable
  columns, but not both. If you want the same Go type to override regardless of
  nullability, you'll need to configure two overrides: one with `nullable: true` and one without.
:::

## The `overrides` list

Each element in the `overrides` list has the following keys:

- `db_type`:
  - A database type to override. Find the full list of supported types in [postgresql_type.go](https://github.com/sqlc-dev/sqlc/blob/main/internal/codegen/golang/postgresql_type.go#L12) or [mysql_type.go](https://github.com/sqlc-dev/sqlc/blob/main/internal/codegen/golang/mysql_type.go#L12). Note that for Postgres you must use pg_catalog-prefixed names where available. `db_type` and `column` are mutually exclusive.
- `column`:
  - A column name to override. The value should be of the form `table.column` but you can also specify `schema.table.column` or `catalog.schema.table.column`. `column` and `db_type` are mutually exclusive.
- `go_type`:
  - The fully-qualified name of a Go type to use in generated code. This is usually a string but can also be [a map](#the-go-type-map) for more complex configurations.
- `go_struct_tag`:
  - A reflect-style struct tag to use in generated code, e.g. `a:"b" x:"y,z"`.
    If you want `json` or `db` tags for all fields, configure `emit_json_tags` or `emit_db_tags` instead.
- `unsigned`:
  - If `true`, sqlc will apply this override when a numeric column is unsigned.
    Note that this only applies to `db_type` overrides and has no effect on `column` overrides.
    Defaults to `false`.
- `nullable`:
  - If `true`, sqlc will apply this override when a column is nullable.
    Otherwise `sqlc` will apply this override when a column is non-nullable.
    Note that this only applies to `db_type` overrides and has no effect on `column` overrides.
    Defaults to `false`.

:::{tip}
  A single `db_type` override configuration applies to either nullable or non-nullable
  columns, but not both. If you want the same Go type to override regardless of nullability, you'll
  need to configure two overrides: one with `nullable: true` and one without.
:::

:::{note}
When generating code, `column` override configurations take precedence over `db_type` configurations.
:::

### The `go_type` map

Some overrides may require more detailed configuration. If necessary, `go_type`
can be a map with the following keys:

- `import`:
  - The import path for the package where the type is defined.
- `package`:
  - The package name where the type is defined. This should only be necessary when your import path doesn't end with the desired package name.
- `type`:
  - The type name itself, without any package prefix.
- `pointer`:
  - If `true`, generated code will use a pointer to the type rather than the type itself.
- `slice`:
  - If `true`, generated code will use a slice of the type rather than the type itself.

An example:

```yaml
version: "2"
sql:
- schema: "postgresql/schema.sql"
  queries: "postgresql/query.sql"
  engine: "postgresql"
  gen:
    go:
      package: "authors"
      out: "db"
      sql_package: "pgx/v5"
      overrides:
        - db_type: "uuid"
          go_type:
            import: "a/b/v2"
            package: "b"
            type: "MyType"
            pointer: true
```

## Global overrides

To override types in all packages that `sqlc` generates, add an override
configuration to the top-level `overrides` section of your `sqlc` config:

```yaml
version: "2"
overrides:
  go:
    overrides:
      - db_type: "pg_catalog.timestamptz"
        nullable: true
        engine: "postgresql"
        go_type:
          import: "gopkg.in/guregu/null.v4"
          package: "null"
          type: "Time"
sql:
- schema: "service1/schema.sql"
  queries: "service1/query.sql"
  engine: "postgresql"
  gen:
    go: 
      package: "service1"
      out: "service1"
- schema: "service2/schema.sql"
  queries: "service2/query.sql"
  engine: "postgresql"
  gen:
    go:
      package: "service2"
      out: "service2"
```

Using this configuration, whenever there is a nullable `timestamp with time zone`
column in a Postgres table, `sqlc` will generate Go code using `null.Time`.

Note that the mapping for global type overrides has a field called `engine` that
is absent in per-package type overrides. This field is only used when there are
multiple `sql` sections using different engines. If you're only generating code
for a single database engine you can omit it.

#### Version 1 configuration

If you are using the older version 1 of the `sqlc` configuration format, override
configurations themselves are unchanged but are nested differently.

Per-package configurations are nested under the `overrides` key within an item
in the `packages` list:

```yaml
version: "1"
packages:
  - name: "db"
    path: "internal/db"
    queries: "./sql/query/"
    schema: "./sql/schema/"
    engine: "postgresql"
    overrides: [...]
```

And global configurations are nested under the top-level `overrides` key:

```yaml
version: "1"
packages: [...]
overrides:
  - db_type: "uuid"
    go_type: "github.com/gofrs/uuid.UUID"
```
