# `analyze` - Analyzing query result types

> [!NOTE]
> `analyze` is in beta. Its flags and JSON output may change in a
> future release, and the types it reports can differ from the ones
> `generate` produces.

`sqlc analyze` analyzes a query against a schema and prints the inferred result
columns and parameters as a single JSON document.

Unlike [`generate`](generate.md), this command does not require a configuration
file and does not connect to a database. It uses sqlc's native static analysis
to infer types directly from the provided schema.

Every dialect is analyzed by the same engine-neutral analysis core: the schema
is loaded into a catalog seeded with the dialect's types, operators and
functions, and each query is resolved against it. `generate` still uses each
engine's own analysis path, so the two can report a type differently — most
visibly, `analyze` reports type names as the catalog stores them, in lower
case.

## Usage

```sh
sqlc analyze --dialect <dialect> --schema <schema-file> [query-file]
```

The query is read from the given file, or from standard input when no file is
provided. The schema is always read from the `--schema` file.

## Flags

- `--dialect`, `-d` - The SQL dialect to use. One of `postgresql`, `mysql`,
  `sqlite`, `clickhouse`, `googlesql`, `mssql`, or `duckdb`. Required.
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
        "type": {
          "name": "bigserial"
        },
        "table": "authors"
      },
      {
        "name": "name",
        "type": {
          "name": "text"
        },
        "table": "authors"
      },
      {
        "name": "bio",
        "type": {
          "name": "text",
          "nullable": true
        },
        "table": "authors"
      }
    ],
    "params": [
      {
        "number": 1,
        "column": {
          "name": "id",
          "type": {
            "name": "bigserial"
          },
          "table": "authors"
        }
      }
    ]
  }
]
```

A column's `type` is written as a call expression: a `name` applied to
`args`, each of which carries an optional `label` and exactly one of `type`,
`int`, `bool` or `string`, with `nullable` set at whatever depth it applies.
An array of text is `array` applied to `text`; a `Map(String, Nullable(UInt8))`
in ClickHouse is `map` applied to `string` and a nullable `uint8`. Names are
recorded as the engine reports them.

Pass `--ast` to also include each statement's parsed AST under an `ast` key. It
has the same shape as the output of [`parse`](parse.md), with every node tagged
by type.
