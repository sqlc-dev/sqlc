# Linting queries

*Added in v1.19.0*

`sqlc vet` runs queries through a set of lint rules.

Rules are defined in the `sqlc` [configuration](../reference/config) file. They consist
of a name, message, and an expression. If the expression evaluates to `true`, an
error is reported. These expressions are evaluated using
[cel-go](https://github.com/google/cel-go).

Each expression has access to a query object, which is defined as the following
struct:

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

This struct may be expanded in the future to include more query information.
We may also add information from a running database, such as the result from
`EXPLAIN`.

While these examples are simplistic, they give you an idea on what types of
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

When a [database](../reference/config.html#database) in configured, the `sqlc/db-preapre`
rule will attempt to prepare each of your queries against the connected
database. Any failures will be reported to standard error.

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

Please note that `sqlc` does not manage or migrate the database. Use your
migration tool of choice to create the necessary database tables and objects.