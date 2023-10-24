# `vet` - Linting queries

*Added in v1.19.0*

`sqlc vet` runs queries through a set of lint rules.

Rules are defined in the `sqlc` [configuration](../reference/config) file. They
consist of a name, message, and a [Common Expression Language
(CEL)](https://github.com/google/cel-spec) expression. Expressions are evaluated
using [cel-go](https://github.com/google/cel-go).  If an expression evaluates to
`true`, `sqlc vet` will report an error using the given message.

## Defining lint rules

Each lint rule's CEL expression has access to information from your sqlc
configuration and queries via variables defined in the following proto messages.

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
  // One of "many", "one", "exec", etc.
  string cmd = 3;
  // Query parameters, if any
  repeated Parameter params = 4;
}

message Parameter
{
  int32 number = 1;
}
```

In addition to this basic information, when you have a PostgreSQL or MySQL
[database connection configured](../reference/config.md#database)
each CEL expression has access to the output from running `EXPLAIN ...` on your query
via the `postgresql.explain` and `mysql.explain` variables.
This output is quite complex and depends on the structure of your query but sqlc attempts
to parse and provide as much information as it can. See
[Rules using `EXPLAIN ...` output](#rules-using-explain-output) for more information.

Here are a few example rules just using the basic configuration and query information available
to the CEL expression environment. While these examples are simplistic, they give you a flavor
of the types of rules you can write.

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

### Rules using `EXPLAIN ...` output

*Added in v1.20.0*

The CEL expression environment has two variables containing `EXPLAIN ...` output,
`postgresql.explain` and `mysql.explain`. `sqlc` only populates the variable associated with
your configured database engine, and only when you have a
[database connection configured](../reference/config.md#database).

For the `postgresql` engine, `sqlc` runs

```sql
EXPLAIN (ANALYZE false, VERBOSE, COSTS, SETTINGS, BUFFERS, FORMAT JSON) ...
```

where `"..."` is your query string, and parses the output into a [`PostgreSQLExplain`](https://buf.build/sqlc/sqlc/docs/v1.20.0:vet#vet.PostgreSQLExplain) proto message.

For the `mysql` engine, `sqlc` runs

```sql
EXPLAIN FORMAT=JSON ...
```

where `"..."` is your query string, and parses the output into a [`MySQLExplain`](https://buf.build/sqlc/sqlc/docs/v1.20.0:vet#vet.MySQLExplain) proto message.

These proto message definitions are too long to include here, but you can find them in the `protos`
directory within the `sqlc` source tree.

The output from `EXPLAIN ...` depends on the structure of your query so it's a bit difficult
to offer generic examples. Refer to the
[PostgreSQL documentation](https://www.postgresql.org/docs/current/using-explain.html) and
[MySQL documentation](https://dev.mysql.com/doc/refman/en/explain-output.html) for more
information.

```yaml
...
rules:
- name: postgresql-query-too-costly
  message: "Query cost estimate is too high"
  rule: "postgresql.explain.plan.total_cost > 1.0"
- name: postgresql-no-seq-scan
  message: "Query plan results in a sequential scan"
  rule: "postgresql.explain.plan.node_type == 'Seq Scan'"
- name: mysql-query-too-costly
  message: "Query cost estimate is too high"
  rule: "has(mysql.explain.query_block.cost_info) && double(mysql.explain.query_block.cost_info.query_cost) > 2.0"
- name: mysql-must-use-primary-key
  message: "Query plan doesn't use primary key"
  rule: "has(mysql.explain.query_block.table.key) && mysql.explain.query_block.table.key != 'PRIMARY'"
```

When building rules that depend on `EXPLAIN ...` output, it may be helpful to see the actual JSON
returned from the database. `sqlc` will print it When you set the environment variable
`SQLCDEBUG=dumpexplain=1`. Use this environment variable together with a dummy rule to see
`EXPLAIN ...` output for all of your queries.

```yaml
version: 2
sql:
  - schema: "query.sql"
    queries: "query.sql"
    engine: "postgresql"
    database:
      uri: "postgresql://postgres:postgres@localhost:5432/postgres"
    gen:
      go:
        package: "db"
        out: "db"
    rules:
      - debug
rules:
- name: debug
  rule: "!has(postgresql.explain)" # A dummy rule to trigger explain
```

Please note that databases configured with a `uri` must have an up-to-date
schema for `vet` to work correctly, and `sqlc` does not apply schema migrations
to your database. Use your migration tool of choice to create the necessary
tables and objects before running `sqlc vet` with rules that depend on
`EXPLAIN ...` output.

Alternatively, configure [managed databases](managed-databases.md) to have
`sqlc` create hosted ephemeral databases with the correct schema automatically.

## Built-in rules

### sqlc/db-prepare

When a [database](../reference/config.md#database) connection is configured, you can
run the built-in `sqlc/db-prepare` rule. This rule will attempt to prepare
each of your queries against the connected database and report any failures.

```yaml
version: 2
sql:
  - schema: "schema.sql"
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

Please note that databases configured with a `uri` must have an up-to-date
schema for `vet` to work correctly, and `sqlc` does not apply schema migrations
to your database. Use your migration tool of choice to create the necessary
tables and objects before running `sqlc vet` with the `sqlc/db-prepare` rule.

Alternatively, configure [managed databases](managed-databases.md) to have
`sqlc` create hosted ephemeral databases with the correct schema automatically.

```yaml
version: 2
cloud:
  project: "<PROJECT_ID>"
sql:
  - schema: "schema.sql"
    queries: "query.sql"
    engine: "postgresql"
    gen:
      go:
        package: "authors"
        out: "db"
    database:
      managed: true
    rules:
      - sqlc/db-prepare
```

To see this in action, check out the [authors
example](https://github.com/sqlc-dev/sqlc/blob/main/examples/authors/sqlc.yaml).

## Running lint rules

When you add the name of a defined rule to the rules list
for a [sql package](../reference/config.md#sql),
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
