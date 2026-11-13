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

## Common-type inference

The shared core does not give every database the same coercion rules. SQL
Server supplies its own type combiner for set operations and
alternative-valued expressions such as `CASE`, `COALESCE`, `GREATEST`, `LEAST`
and `CHOOSE`. `ISNULL` instead keeps its first typed argument's type; both
arguments must be nullable for its result to be nullable.

Dialects without a combiner keep the existing core fallback: a set-operation
column takes its type from the left branch, and an alternative-valued
expression takes the first typed candidate. **This is incomplete static
inference, not those databases' common-type semantics.** It does not change
their normal `generate` paths.

The database rules differ in ways that need separate policies:

| Dialect | Database behavior and current core gap |
| --- | --- |
| SQL Server | Uses type precedence and combines decimal scale/integer digits with a precision limit of 38. The enabled policy models these rules, but is not a general implicit-cast validator. |
| PostgreSQL | Selects common types using categories, implicit casts, preferred types, unknown literals and domains. Modifiers survive only when the candidates have the same type and modifier: `numeric(5,2) UNION numeric(8,4)` has unconstrained `numeric`, not SQL Server-style decimal widening. |
| SQLite | Values are dynamically typed; compound-view/subquery affinity can come from any constituent `SELECT`. Declared decimal precision and string length do not impose SQL Server-style limits. |
| MySQL | Aggregates types and lengths across set-operation branches; alternative expressions also need signed/unsigned, numeric, binary/string, temporal and NULL-literal rules. |
| DuckDB | Uses combination casting for set operations, including casts not allowed by ordinary implicit casting. Integer widths and nested types need more than a precedence list. |
| ClickHouse | Computes compatible supertypes, including `Nullable`, signed/unsigned widths and nested types; decimal precision can reach 76. Some combinations have no compatible supertype. |
| GoogleSQL (Spanner) | Intersects supertype sets, distinguishes exact from inexact numeric inputs, and has separate literal and structured-type rules. Its type inventory is not identical to BigQuery's. |

These are reviewed database semantics, not a claim that sqlc has implemented
or live-tested the missing policies. See the upstream references:
[SQL Server](https://learn.microsoft.com/en-us/sql/t-sql/data-types/precision-scale-and-length-transact-sql),
[PostgreSQL common types](https://www.postgresql.org/docs/18/typeconv-union-case.html)
and [modifiers](https://github.com/postgres/postgres/blob/REL_18_STABLE/src/backend/parser/parse_coerce.c),
[SQLite](https://www.sqlite.org/datatype3.html#column_affinity_for_compound_views),
[MySQL set operations](https://docs.oracle.com/cd/E17952_01/mysql-8.0-en/set-operations.html)
and [alternative expressions](https://docs.oracle.com/cd/E17952_01/mysql-8.0-en/flow-control-functions.html),
[DuckDB](https://duckdb.org/docs/2.0/sql/data_types/typecasting),
[ClickHouse UNION](https://clickhouse.com/docs/reference/statements/select/union)
and [decimals](https://clickhouse.com/docs/reference/data-types/decimal),
and [Spanner](https://docs.cloud.google.com/spanner/docs/reference/standard-sql/conversion_rules).

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
