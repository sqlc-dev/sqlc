# `fmt` - Formatting queries

`sqlc fmt` rewrites the query files referenced by your configuration file in a
canonical format. Each query is parsed with the engine's parser and printed
back from the syntax tree, so formatting never depends on how the query was
written — only on what it means.

Like `gofmt`, the formatter does not impose a maximum line width. A statement
written on a single line stays on a single line, and a statement the author
broke across lines keeps its breaks: the printer notices which clause
boundaries (`FROM`, `WHERE`, `ORDER BY`, ...) and list boundaries the author
broke at and preserves them, normalizing indentation and spacing around them.

Comments inside a statement are formatted along with it: each comment is
anchored to the code around it by source position and printed back there —
a comment trailing a select-list item stays with that item, a comment above
a clause stays above its keyword — and a comment that runs to the end of its
line breaks the statement open around it. Any statement that cannot be
proven to survive formatting unchanged is left exactly as written.

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
SELECT id, name, bio
FROM authors
WHERE id = ?
LIMIT 1;

-- name: SearchAuthors :many
SELECT
  id, -- the primary key
  name,
  bio,
  created_at
FROM authors
WHERE name LIKE ? AND bio IS NOT NULL AND id > ? AND created_at > ? AND name <> ?
ORDER BY name;
```

`GetAuthor` keeps the line breaks its author wrote; had it been written on
one line, it would stay on one line. In `SearchAuthors`, the line comment
cannot share a line with the code after it, so the statement breaks open
around it, while the `WHERE` chain — written on one line — stays on one.
