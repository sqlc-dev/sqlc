# Separating models file

By default, sqlc uses a single package to place all the generated code. But you may want to separate
the generated models file into another package for loose coupling purposes in your project.

To do this, you can use the following configuration:

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "queries.sql"
    schema: "schema.sql"
    gen:
      go:
        out: "internal/"  # Base directory for the generated files. You can also just use "."
        sql_package: "pgx/v5"
        package: "sqlcrepo"
        output_batch_file_name: "db/sqlcrepo/batch.go"
        output_db_file_name: "db/sqlcrepo/db.go"
        output_querier_file_name: "db/sqlcrepo/querier.go"
        output_copyfrom_file_name: "db/sqlcrepo/copyfrom.go"
        output_query_files_directory: "db/sqlcrepo/"
        output_models_file_name: "business/entities/models.go"
        output_models_package: "entities"
        models_package_import_path: "example.com/project/module-path/internal/business/entities"
```

This configuration will generate files in the `internal/db/sqlcrepo` directory with `sqlcrepo`
package name, except for the models file which will be generated in the `internal/business/entities`
directory. The generated models file will use the package name `entities` and it will be imported in
the other generated files using the given
`"example.com/project/module-path/internal/business/entities"` import path when needed.

The generated files will look like this:

```
my-app/
├── internal/
│   ├── db/
│   │   └── sqlcrepo/
│   │       ├── db.go
│   │       └── queries.sql.go
│   └── business/
│       └── entities/
│           └── models.go
├── queries.sql
├── schema.sql
└── sqlc.yaml
```
