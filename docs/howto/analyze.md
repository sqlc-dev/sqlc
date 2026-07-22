# `analyze` - Analyzing query result types

`sqlc analyze` analyzes a query against a schema and prints the inferred result
columns and parameters as a single JSON document.

Unlike [`generate`](generate.md), this command does not require a configuration
file and does not connect to a database. It uses sqlc's native static analysis
to infer types directly from the provided schema.

## Usage

```sh
sqlc analyze --dialect <dialect> --schema <schema-file> [query-file]
```

The query is read from the given file, or from standard input when no file is
provided. The schema is always read from the `--schema` file.

## Flags

- `--dialect`, `-d` - The SQL dialect to use. One of `postgresql`, `mysql`,
  `sqlite`, or `googlesql`. Required.
- `--schema`, `-s` - Path to the schema (DDL) file. Required.
- `--ast` - Include each statement's AST in the output. Defaults to `false`.

## Examples

Given a schema in `schema.sql`:

```sql
CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  bio  text
);
```

and a query in `query.sql`:

```sql
-- name: GetAuthor :one
SELECT * FROM authors WHERE id = $1;
```

Running:

```sh
sqlc analyze --dialect postgresql --schema schema.sql query.sql
```

reports the result columns and parameters:

```json
[
  {
    "name": "GetAuthor",
    "cmd": ":one",
    "columns": [
      {
        "name": "id",
        "data_type": "bigserial",
        "not_null": true,
        "is_array": false,
        "table": "authors"
      },
      {
        "name": "name",
        "data_type": "text",
        "not_null": true,
        "is_array": false,
        "table": "authors"
      },
      {
        "name": "bio",
        "data_type": "text",
        "not_null": false,
        "is_array": false,
        "table": "authors"
      }
    ],
    "params": [
      {
        "number": 1,
        "column": {
          "name": "id",
          "data_type": "bigserial",
          "not_null": true,
          "is_array": false,
          "table": "authors"
        }
      }
    ]
  }
]
```

Pass `--ast` to also include each statement's parsed AST under an `ast` key.
