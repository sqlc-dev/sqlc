# Linting queries

*Added in v1.19.0*

`sqlc vet` runs queries through a set of lint rules.

Rules are defined in the `sqlc` [configuration](../reference/config) file. They consist
of a name, message, and a [Common Expression Language (CEL)](https://github.com/google/cel-spec)
expression. Expressions are evaluated using [cel-go](https://github.com/google/cel-go).
If an expression evaluates to `true`, `sqlc vet` will report an error using the given message.

## Defining lint rules

Each lint rule's CEL expression has access to variables from your sqlc configuration and queries,
defined in the following struct.

```proto
message Config
{
  string version = 1;
  string engine = 2 ;
  repeated string schema = 3;
  repeated string queries = 4;
}

message Query
{
  // SQL body
  string sql = 1;
  // Name of the query
  string name = 2; 
  // One of :many, :one, :exec, etc.
  string cmd = 3;
  // Query parameters, if any
  repeated Parameter params = 4;
}

message Parameter
{
  int32 number = 1;
}
```

This struct will likely expand in the future to include more query information.
We may also add information returned from a running database, such as the result from
`EXPLAIN ...`.

While these examples are simplistic, they give you a flavor of the types of
rules you can write.

```yaml
version: 2
sql:
  - schema: "query.sql"
    queries: "query.sql"
    engine: "postgresql"
    gen:
      go:
        package: "authors"
        out: "db"
    rules:
      - no-pg
      - no-delete
      - only-one-param
      - no-exec
rules:
  - name: no-pg
    message: "invalid engine: postgresql"
    rule: |
      config.engine == "postgresql"
  - name: no-delete
    message: "don't use delete statements"
    rule: |
      query.sql.contains("DELETE")
  - name: only-one-param
    message: "too many parameters"
    rule: |
      query.params.size() > 1
  - name: no-exec
    message: "don't use exec"
    rule: |
      query.cmd == "exec"
```

## Built-in rules

### sqlc/db-prepare

When a [database](../reference/config.html#database) connection is configured, you can
run the built-in `sqlc/db-prepare` rule. This rule will attempt to prepare
each of your queries against the connected database and report any failures.

```yaml
version: 2
sql:
  - schema: "query.sql"
    queries: "query.sql"
    engine: "postgresql"
    gen:
      go:
        package: "authors"
        out: "db"
    database:
      uri: "postgresql://postgres:password@localhost:5432/postgres"
    rules:
      - sqlc/db-prepare
```

To see this in action, check out the [authors
example](https://github.com/kyleconroy/sqlc/blob/main/examples/authors/sqlc.yaml).

Please note that `sqlc` does not manage or migrate your database. Use your
migration tool of choice to create the necessary database tables and objects
before running `sqlc vet` with the `sqlc/db-prepare` rule.

## Running lint rules

When you add the name of a defined rule to the rules list
for a [sql package](https://docs.sqlc.dev/en/stable/reference/config.html#sql),
`sqlc vet` will evaluate that rule against every query in the package.

In the example below, two rules are defined but only one is enabled.

```yaml
version: 2
sql:
  - schema: "query.sql"
    queries: "query.sql"
    engine: "postgresql"
    gen:
      go:
        package: "authors"
        out: "db"
    rules:
      - no-delete
rules:
  - name: no-pg
    message: "invalid engine: postgresql"
    rule: |
      config.engine == "postgresql"
  - name: no-delete
    message: "don't use delete statements"
    rule: |
      query.sql.contains("DELETE")
```

### Opting-out of lint rules

For any query, you can tell `sqlc vet` not to evaluate lint rules using the
`@sqlc-vet-disable` query annotation.

```sql
/* name: GetAuthor :one */
/* @sqlc-vet-disable */
SELECT * FROM authors
WHERE id = ? LIMIT 1;
```