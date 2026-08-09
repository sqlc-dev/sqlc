# Named Parameters Package - Claude Code Guide

This package models a query's parameters: their names, their nullability and
whether they came from `sqlc.slice()`.

It no longer knows anything about sqlc's *syntax*. `sqlc.arg()`, `sqlc.narg()`,
`sqlc.slice()` and `@name` are rewritten to native placeholders by
`internal/sql/preprocess` before any engine parser runs, and that package builds
the `ParamSet` while it rewrites.

## Param

A `Param` carries the user-facing name plus a nullability that combines what was
inferred from the schema with what the user asked for:

- `NewParam(name)` — unspecified nullability, from `sqlc.arg()` or `@name`
- `NewUserNullableParam(name)` — from `sqlc.narg()`, always nullable
- `NewSqlcSlice(name)` — from `sqlc.slice()`
- `NewInferredParam(name, notNull)` — what the compiler inferred

A user-specified nullability outranks an inferred one. `mergeParam` ORs the
nullability bits together, so the order the two sources arrive in does not
matter.

## ParamSet

`ParamSet` maps placeholder numbers to `Param` values for a single statement.

```go
ps := named.NewParamSet(numbersAlreadyUsed, hasNamedSupport)
n := ps.Add(named.NewParam("author_id")) // returns the placeholder number
```

`hasNamedSupport` is false for dialects that send one argument per `?`
(MySQL, ClickHouse). There, every occurrence of a name gets its own number. For
`$1`/`?1`/`@name` dialects repeated uses of a name share a number.

The compiler reads the set back with:

- `NameFor(number)` — the user-facing name for a placeholder
- `FetchMerge(number, inferred)` — the merged param and whether it was named

## MySQL @variable vs sqlc @param

`@name` is only sqlc syntax where the preprocessor's dialect says it is. MySQL
sets `AtSign: false`, so `@user_id` there is a user variable and reaches the
parser untouched (converted to `ast.VariableExpr`). PostgreSQL, SQLite and
GoogleSQL treat `@name` as a named parameter.

```sql
-- PostgreSQL with sqlc @param syntax:
SELECT * FROM users WHERE id = @user_id
-- Preprocessed to: SELECT * FROM users WHERE id = $1

-- MySQL with user variable:
SELECT * FROM users WHERE id != @user_id
-- Stays: SELECT * FROM users WHERE id != @user_id
```
