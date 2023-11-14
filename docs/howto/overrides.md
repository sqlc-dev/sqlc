# Overriding types

In many cases it's useful to tell `sqlc` explicitly what Go type you want it to
use for a query input or output. For instance, a PostgreSQL UUID type will map
to `github.com/google/uuid` by default, but you may want `sqlc` to use
`github.com/gofrs/uuid` instead.

If you'd like `sqlc` to use a different Go type, specify the package and type
in the `overrides` list.

```yaml
version: "2"
sql:
- schema: "postgresql/schema.sql"
  queries: "postgresql/query.sql"
  engine: "postgresql"
  gen:
    go: 
      package: "authors"
      out: "postgresql"
      overrides:
        - db_type: "uuid"
          go_type: "github.com/gofrs/uuid.UUID"
```

## The `overrides` list

Each element in the `overrides` list has the following keys:

- `db_type`:
  - A database type to override. Find the full list of supported types in [postgresql_type.go](https://github.com/sqlc-dev/sqlc/blob/main/internal/codegen/golang/postgresql_type.go#L12) or [mysql_type.go](https://github.com/sqlc-dev/sqlc/blob/main/internal/codegen/golang/mysql_type.go#L12). Note that for Postgres you must use pg_catalog-prefixed names where available. And note that `db_type` and `column` are mutually exclusive.
- `column`:
  - A column name to override. The value should be of the form `table.column` but you can also specify `schema.table.column` or `catalog.schema.table.column`. Note that `column` and `db_type` are mutually exclusive.
- `go_type`:
  - The fully-qualified name of a Go type to use in generated code. This is usually a string but can also be [a map](#the-go_type-map) for more complex configurations.
- `go_struct_tag`:
  - A reflect-style struct tag to use in generated code, e.g. `a:"b" x:"y,z"`.
    If you want `json` or `db` tags for all fields, use `emit_json_tags` or `emit_db_tags` instead.
- `nullable`:
  - If `true`, sqlc will apply this override when a column is nullable.
    Otherwise `sqlc` will apply this override when a column is non-nullable.
    Note that this has no effect on `column` overrides. Defaults to `false`.

Note that a single `db_type` override configuration applies to either nullable or non-nullable
columns, but not both. If you want the same Go type to override in both cases, you'll
need to configure two overrides.

When generating code, entries using the `column` key will always take precedence over
entries using the `db_type` key.

### The `go_type` map

Some overrides may require more detailed configuration. If necessary, `go_type`
can also be an object with the following keys:

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
      out: "postgresql"
      overrides:
        - db_type: "uuid"
          go_type:
            import: "a/b/v2"
            package: "b"
            type: "MyType"
            pointer: true
```
