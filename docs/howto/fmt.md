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

Comments are never deleted. The comments above each query are kept, including
the [`-- name:`](../reference/query-annotations.md) annotation and multi-line
`/* */` blocks, and a comment on the same line as a statement's closing
semicolon stays attached to it.

Comments inside a statement are formatted along with it: each comment is
anchored to the code around it by source position and printed back there —
a comment trailing a select-list item stays with that item, a comment above
a clause stays above its keyword — and a statement carrying comments keeps
its multi-line shape. Any statement that cannot be proven to survive
formatting unchanged is left exactly as written.

Formatting currently supports the `sqlite` engine, whose parser surfaces the
comments its lexer sees; other engines' query files are left untouched and
gain support as their parsers are updated. Files that cannot be parsed
without the compiler's preprocessing (for example, files using
`sqlc.slice()` on engines whose parser rejects it) are also skipped; a
skipped file is reported on standard error and left unchanged.

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
where id =  ? limit 1;

-- name: SearchAuthors :many
SELECT id, -- the primary key
       name, bio, created_at FROM authors WHERE name LIKE ? AND bio IS NOT NULL AND id > ? AND created_at > ? AND name <> ? ORDER BY name;
```

running `sqlc fmt` rewrites it to:

```sql
-- name: GetAuthor :one
SELECT id, name, bio FROM authors WHERE id = ? LIMIT 1;

-- name: SearchAuthors :many
SELECT
  id, -- the primary key
  name,
  bio,
  created_at
FROM authors
WHERE name LIKE ?
  AND bio IS NOT NULL
  AND id > ?
  AND created_at > ?
  AND name <> ?
ORDER BY name;
```
