# `parse` - Parsing SQL into an AST

> [!NOTE]
> `parse` is in beta. Its flags and the shape of the JSON AST it
> prints may change in a future release.

`sqlc parse` parses SQL from a file or standard input and prints the abstract
syntax tree (AST) as a single JSON document. It does not require a configuration
file or a database connection.

Each statement is reported with its sqlc query name and command (when the
statement carries a [`-- name:`](../reference/query-annotations.md) annotation)
alongside its AST.

## Usage

```sh
sqlc parse --dialect <dialect> [file]
```

The SQL is read from the given file, or from standard input when no file is
provided.

## Flags

- `--dialect`, `-d` - The SQL dialect to use. One of `postgresql`, `mysql`,
  `sqlite`, `clickhouse`, `googlesql`, `mssql`, or `duckdb`. Required.

## Examples

Parse a query file:

```sh
sqlc parse --dialect postgresql query.sql
```

Parse SQL piped via standard input:

```sh
echo "SELECT 1;" | sqlc parse --dialect mysql
```

The output is a JSON array with one object per statement:

```json
[
  {
    "name": "GetAuthor",
    "cmd": ":one",
    "ast": {
      "tag": "RawStmt",
      "stmt": {
        "...": "..."
      },
      "stmt_location": 0,
      "stmt_len": 42
    }
  }
]
```

Statements without a `-- name:` annotation (for example schema DDL) omit the
`name` and `cmd` fields. Field names are `snake_case` versions of the AST node
field names.

## Node types

Every node in the AST carries a `tag` naming its type. Some nodes have no
fields of their own, so without it a star, a null literal and an untranslated
clause would all print as `{}`.

```json
"val": {
  "tag": "ColumnRef",
  "name": "",
  "fields": {
    "tag": "List",
    "items": [
      {
        "tag": "A_Star"
      }
    ]
  },
  "location": 93
}
```

A `tag` of `TODO` marks a clause the dialect's converter does not translate
yet. It means the clause was parsed but is not represented in the AST, not that
the clause was absent from the query.

## Absent fields

A field the statement does not use is left out rather than printed as `null`.
An `A_Const` carrying an integer reports only what it has:

```json
"val": {
  "tag": "A_Const",
  "val": {
    "tag": "Integer",
    "ival": 1
  },
  "location": 30
}
```

Only absent fields — and empty lists, which engines construct differently for
the same SQL — are omitted. A zero keeps its place, because zero is a value the
parser can find: `stmt_location` is 0 for the first statement in a file, and
`LIMIT 0` parses to an `ival` of 0.
