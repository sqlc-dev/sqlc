# ClickHouse testgen

`testgen` records what ClickHouse itself reports about a set of sqlc queries,
so the answer can be committed as a golden file and held against what sqlc's
own analysis produces for the same schema and queries.

It is a nested Go module with no dependencies beyond the standard library.
Run it from this directory:

```bash
# Download the pinned clickhouse binary into the user cache directory.
go run . install

# Analyze every query in query.sql against schema.sql with fixture.sql loaded.
go run . analyze --schema schema.sql --fixture fixture.sql query.sql
```

The binary is looked up in the `CLICKHOUSE` environment variable first, then
in the cache populated by `install`. The version is pinned in `install.go`;
bumping it can change the query tree format and type inference, so regenerate
and review the goldens afterwards.

## How it works

Each query runs in its own `clickhouse local` process, which needs no server,
no network and no configuration: the schema and fixture are loaded fresh, the
query is explained and then executed, and the process exits.

- **Types and nullability** come from the executed query's result header,
  exactly as a driver would see them.
- **Provenance** comes from `EXPLAIN QUERY TREE`, whose resolved columns point
  at the table expression they read from. References are followed through
  subqueries, CTEs and unions to the base table.
- **Parameters** are invisible to ClickHouse, which substitutes `?` on the
  client. Each `?`, `sqlc.arg()` and `sqlc.narg()` is replaced by a constant
  expression that carries its ordinal, `(NULL + k)`, or `toUInt64(4294967295 + k)`
  after `LIMIT` and `OFFSET`. The query tree prints the expression each folded
  constant came from, so the placeholder is found again and described by the
  operand it is compared with. Parameters of `INSERT ... VALUES` map onto the
  target columns reported by `DESCRIBE TABLE`.

The output has the shape of `sqlc analyze`, except that each column's type is
one expression rather than a name and flags. A type is written as a call: a
lowercased name applied to arguments that are numbers, quoted strings,
identifiers, other calls, or any of those with a label. `Nullable`, `Array`
and `LowCardinality` are ordinary names in that grammar, so nothing about
nesting is lost:

| ClickHouse                        | testgen                            |
|-----------------------------------|------------------------------------|
| `Array(Nullable(String))`         | `array(nullable(string))`          |
| `Map(String, UInt32)`             | `map(string, uint32)`              |
| `Tuple(lat Float64, lon Float64)` | `tuple(lat: float64, lon: float64)`|
| `Enum8('a' = 1, 'b' = 2)`         | `enum8('a': 1, 'b': 2)`            |
| `DateTime64(3, 'UTC')`            | `datetime64(3, 'UTC')`             |
| `LowCardinality(Nullable(String))`| `lowcardinality(string)`, `not_null: false` |

An outer `Nullable` is the column's nullability rather than part of its type,
so it is lifted into `not_null`; a `Nullable` at any greater depth stays in
the expression. Resolving the names is left to whoever reads the output.

## Tests

`testdata/<case>/` holds a `schema.sql`, `query.sql`, an optional
`fixture.sql` and the expected `stdout.txt`. The test skips unless a binary is
available.

```bash
go run . install
go test .
go test . -update   # rewrite every stdout.txt
```
