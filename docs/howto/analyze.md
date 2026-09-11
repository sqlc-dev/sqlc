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
          "name": "bigint"
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
            "name": "bigint"
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
`int`, `bool`, `string` or `ident`, with `nullable` set at whatever depth it
applies. A `numeric(10,2)` column is `numeric` applied to `10` and `2`; an
array of text is `array` applied to `text`, and an array of arrays nests
one `array` per dimension; a `Map(String, Nullable(UInt8))` in ClickHouse is
`map` applied to `string` and a nullable `uint8`; a `STRUCT<a INT64>` in
GoogleSQL is `struct` applied to an `int64` labelled `a`; the `MAX` of SQL
Server's `nvarchar(max)` is the identifier `max`.

Types are reported the way the engine itself stores and reports them rather
than the way the schema spelled them: PostgreSQL's `int` and `bigserial` are
`integer` and `bigint`, as `format_type` prints them; MySQL's `BOOLEAN` is
`tinyint(1)`; ClickHouse's `Decimal32(4)` is `decimal(9, 4)`; DuckDB's
`TEXT` is `varchar`; SQL Server's `FLOAT(24)` is `real`. SQLite, which
keeps a declared type as written, is reported as written.

Pass `--ast` to also include each statement's parsed AST under an `ast` key. It
has the same shape as the output of [`parse`](parse.md), with every node tagged
by type.
