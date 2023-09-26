# Overriding types

The default mapping of PostgreSQL/MySQL types to Go types only uses packages outside
the standard library when it must.

For example, the `uuid` PostgreSQL type is mapped to `github.com/google/uuid`.
If a different Go package for UUIDs is required, specify the package in the
`overrides` array. In this case, I'm going to use the `github.com/gofrs/uuid`
instead.

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

Each mapping of the `overrides` collection has the following keys:

- `db_type`:
  - The PostgreSQL or MySQL type to override. Find the full list of supported types in [postgresql_type.go](https://github.com/sqlc-dev/sqlc/blob/main/internal/codegen/golang/postgresql_type.go#L12) or [mysql_type.go](https://github.com/sqlc-dev/sqlc/blob/main/internal/codegen/golang/mysql_type.go#L12). Note that for Postgres you must use the pg_catalog prefixed names where available. Can't be used if the `column` key is defined.
- `column`:
  - In case the type overriding should be done on specific a column of a table instead of a type. `column` should be of the form `table.column` but you can be even more specific by specifying `schema.table.column` or `catalog.schema.table.column`. Can't be used if the `db_type` key is defined.
- `go_type`:
  - A fully qualified name to a Go type to use in the generated code.
- `go_struct_tag`:
  - A reflect-style struct tag to use in the generated code, e.g. `a:"b" x:"y,z"`.
    If you want general json/db tags for all fields, use `emit_db_tags` and/or `emit_json_tags` instead.
- `nullable`:
  - If `true`, use this type when a column is nullable. Defaults to `false`.

Note that a single `db_type` override configuration applies to either nullable or non-nullable
columns, but not both. If you want a single `go_type` to override in both cases, you'll
need to specify two overrides.

When generating code, entries using the `column` key will always have preference over
entries using the `db_type` key in order to generate the struct.

For more complicated import paths, the `go_type` can also be an object with the following keys:

- `import`:
  - The import path for the package where the type is defined.
- `package`:
  - The package name where the type is defined. This should only be necessary when your import path doesn't end with the desired package name.
- `type`:
  - The type name itself, without any package prefix.
- `pointer`:
  - If set to `true`, generated code will use pointers to the type rather than the type itself.
- `slice`:
  - If set to `true`, generated code will use a slice of the type rather than the type itself.

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