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
      "Stmt": {
        "...": "..."
      },
      "StmtLocation": 0,
      "StmtLen": 42
    }
  }
]
```

Statements without a `-- name:` annotation (for example schema DDL) omit the
`name` and `cmd` fields.
