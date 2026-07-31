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
| googlesql | `@name` | native, left in place |
| clickhouse | `?` | not a parameter |

## How it works

`lexer` is a dialect-parameterized scanner. It does not parse SQL — it knows
only enough to skip the regions where sqlc syntax must not be rewritten:
comments (`--`, `/* */`, `#`), string literals, quoted identifiers, backticks
and PostgreSQL dollar-quoted strings. `Dialect` (dialect.go) is the single
place those lexical rules live.

`File(engine, src)` walks the whole file, splits it at top-level semicolons and,
for each statement:

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

- **Line structure is preserved.** A rewrite never adds or removes a newline, so
  line numbers survive even before `Origin` is applied.
- **An invalid statement is copied through untouched.** The engine still parses
  what the user wrote and the error is reported per statement, not per file.
- **Nothing inside a comment or literal is rewritten.** Query annotations like
  `-- name: GetAuthor :one` are comments, so they are always safe.

## Tests

`testdata/<engine>/<case>/` holds an `input.sql`, the expected `output.sql` and,
for invalid input, a `stderr.txt`. `TestRewrite` runs one subtest per directory.
Regenerate the goldens with:

```bash
go test ./internal/sql/preprocess -update
```

## Adding a dialect

Add an entry to `dialects` in dialect.go. Nothing else in the codebase needs to
change for `sqlc.arg`/`narg`/`slice`/`embed` to work for a new engine.
