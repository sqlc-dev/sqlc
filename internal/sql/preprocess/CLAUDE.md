# SQL Preprocess Package - Claude Code Guide

This package rewrites sqlc's query syntax into native SQL **before** any engine
parser sees a query. After it runs, no sqlc syntax remains in the text, so
parsers, converters and `astutils` traversals only ever deal with SQL.

## What it rewrites

| syntax | becomes |
| --- | --- |
| `sqlc.arg(name)` / `sqlc.narg(name)` | the dialect's native placeholder |
| `sqlc.slice(name)` | the placeholder wrapped in `/*SLICE:name*/` |
| `sqlc.embed(table)` | `table.*` |
| `@name` | the native placeholder, on dialects where `@` is sqlc syntax |

Native placeholders by dialect:

| engine | placeholder | `@name` |
| --- | --- | --- |
| postgresql | `$1` | sqlc syntax |
| mysql | `?` | user variable, left alone |
| sqlite | `?1` | sqlc syntax |

GoogleSQL and ClickHouse are deliberately absent. They handle their own
parameter syntax, so `File` returns their source unchanged and sqlc syntax is
not available to them — `sqlc.arg()` reaches the parser as the function call it
looks like.

## How it works

`lexer` is a dialect-parameterized scanner. It does not parse SQL — it knows
only enough to skip the regions where sqlc syntax must not be rewritten:
comments (`--`, `/* */`, `#`), string literals, quoted identifiers, backticks
and PostgreSQL dollar-quoted strings. `Dialect` (dialect.go) is the single
place those lexical rules live.

`File(engine, src)` looks up the dialect; an engine with no entry gets its
source back untouched. Otherwise it walks the whole file, splits it at
top-level semicolons and, for each statement:

1. `scan` collects every sqlc construct **and** every native placeholder, in
   source order.
2. `number` assigns placeholder numbers. Numbered dialects keep the numbers the
   user wrote and fill in the gaps; `?` dialects renumber everything in source
   order, because each `?` is its own argument.
3. The statement is rewritten into the output buffer and a `Statement` records
   what changed.

## The side table

`Result.Statement(offset)` returns the `Statement` covering an offset in the
**rewritten** text:

- `Params` — the `named.ParamSet` for the statement
- `Embeds` — each rewritten `table.*`, keyed by its location so the compiler can
  tell it apart from a star reference the user wrote
- `Slices` — the location of each `sqlc.slice()` placeholder
- `Numbers` — location → parameter number, used by the compiler to override the
  numbers an engine assigned in AST-conversion order
- `Dollar` / `ParamErr` — the placeholder-style validation that used to live in
  `validate.ParamRef`
- `Err` — a sqlc syntax error (unknown `sqlc.*` function, wrong arity, an
  argument that is not an identifier or string)

`Result.Origin(offset)` maps an offset in the rewritten text back to the
original source, so errors point at what the user wrote.

## Invariants

- **A rewrite never adds a line.** Replacements are single-line, so a line
  number taken from the rewritten text is never past the end of the original. It
  can *lose* lines, when the sqlc call itself spanned several — map offsets
  through `Origin` rather than trusting line numbers.
- **An invalid statement is copied through untouched.** The engine still parses
  what the user wrote and the error is reported per statement, not per file.
- **Nothing inside a comment or literal is rewritten.** Query annotations like
  `-- name: GetAuthor :one` are comments, so they are always safe.

## Tests

`testdata/<engine>/<case>/` holds the input and the expected results:

| file | contents |
| --- | --- |
| `input.sql` | the query file |
| `output.sql` | the rewritten SQL |
| `side_table.json` | everything the preprocessor recorded |
| `stderr.txt` | reported errors, only when the input is invalid |

`TestRewrite` runs one subtest per directory. `side_table.json` is rendered by
reading the result back through the same API the compiler uses, so it covers
parameter names and nullability, embed and slice spans, placeholder numbering
and the offset map — a parameter's `location` is its offset in the rewritten
text and its `origin` is where it came from.

Regenerate every golden with:

```bash
go test ./internal/sql/preprocess -update
```

The cases cover every shape of sqlc syntax that appears in
`internal/endtoend`, per preprocessed engine — each function against a bare
reference, a string constant and a quoted identifier, the placeholder styles,
and the constructs that must be left alone (comments, literals, MySQL user
variables and `$1`, PostgreSQL's `@>` and `?` operators). When you add a shape
to the corpus, add it here too.

## Adding a dialect

Add an entry to `dialects` in dialect.go. Nothing else in the codebase needs to
change for `sqlc.arg`/`narg`/`slice`/`embed` to work for a new engine — but only
add one if the engine should support sqlc syntax at all. Leaving an engine out
is a deliberate choice, not an oversight.
