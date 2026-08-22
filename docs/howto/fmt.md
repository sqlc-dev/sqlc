# `fmt` - Formatting queries

`sqlc fmt` rewrites the query files referenced by your configuration file in a
canonical format. Each query is parsed with the engine's parser and printed
back from the syntax tree, so formatting never depends on how the query was
written — only on what it means.

A statement that fits within 80 columns is printed on a single line. A longer
statement breaks at clause boundaries (`FROM`, `WHERE`, `ORDER BY`, ...), and
any part that still does not fit — a column list, an `AND`/`OR` chain, a
parenthesized subquery — breaks again, indented one level deeper. This is the
same layout model used by Prettier and by ruff's Python formatter: the printer
lays each group of the statement out flat when it fits and breaks it when it
does not.

The comments above each query are kept, including the
[`-- name:`](../reference/query-annotations.md) annotation. A statement that
cannot be proven to survive formatting unchanged — for example one using
syntax the formatter does not support yet — is left exactly as written.

Formatting is supported for the `postgresql`, `mysql`, `sqlite` and
`clickhouse` engines. Query files for other engines are skipped, as are files
that cannot be parsed without the compiler's preprocessing (for example, files
using `sqlc.slice()` on engines whose parser rejects it); a skipped file is
reported on standard error and left unchanged.

## Usage

```sh
sqlc fmt [--diff]
```

Without flags, the query files are rewritten in place. With `--diff`, the
changes are printed to standard output instead and no files are modified.

## Examples

Given this query file:

```sql
-- name: GetAuthor :one
select  id,name , bio
from   authors
where id =  $1 limit 1;

-- name: SearchAuthors :many
SELECT id, name, bio, created_at FROM authors WHERE name LIKE $1 AND bio IS NOT NULL AND id > $2 AND created_at > $3 AND name <> $4 ORDER BY name;
```

running `sqlc fmt` rewrites it to:

```sql
-- name: GetAuthor :one
SELECT id, name, bio FROM authors WHERE id = $1 LIMIT 1;

-- name: SearchAuthors :many
SELECT id, name, bio, created_at
FROM authors
WHERE name LIKE $1
  AND bio IS NOT NULL
  AND id > $2
  AND created_at > $3
  AND name <> $4
ORDER BY name;
```
