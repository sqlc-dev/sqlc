# Database and language support

| Language   | Plugin                                                                 | MySQL  | PostgreSQL | SQLite          | ClickHouse      | DuckDB          | Spanner         | SQL Server      |
| ---------- | ---------------------------------------------------------------------- | ------ | ---------- | --------------- | --------------- | --------------- | --------------- | --------------- |
| Go         | (built-in)                                                             | Stable | Stable     | Beta            | Beta            | Beta            | Beta            | Beta            |
| Go         | [sqlc-gen-go](https://github.com/sqlc-dev/sqlc-gen-go)                 | Stable | Stable     | Beta            | Not implemented | Not implemented | Not implemented | Not implemented |
| Kotlin     | [sqlc-gen-kotlin](https://github.com/sqlc-dev/sqlc-gen-kotlin)         | Beta   | Beta       | Not implemented | Not implemented | Not implemented | Not implemented | Not implemented |
| Python     | [sqlc-gen-python](https://github.com/sqlc-dev/sqlc-gen-python)         | Beta   | Beta       | Not implemented | Not implemented | Not implemented | Not implemented | Not implemented |
| TypeScript | [sqlc-gen-typescript](https://github.com/sqlc-dev/sqlc-gen-typescript) | Beta   | Beta       | Not implemented | Not implemented | Not implemented | Not implemented | Not implemented |

The built-in Go generator targets `database/sql` for ClickHouse, DuckDB,
Spanner and SQL Server, with the types their drivers hand back:
[clickhouse-go](https://github.com/ClickHouse/clickhouse-go),
[duckdb-go](https://github.com/duckdb/duckdb-go),
[go-sql-spanner](https://github.com/googleapis/go-sql-spanner) and
[go-mssqldb](https://github.com/microsoft/go-mssqldb). Enums and alias types
are not yet resolved for these engines and come out as `any`.

## Community language support

New languages can be added via [plugins](../guides/plugins.md).

| Language | Plugin                                                                                | MySQL  | PostgreSQL | SQLite |
| -------- | ------------------------------------------------------------------------------------- | ------ | ---------- | ------ |
| C#       | [DaredevilOSS/sqlc-gen-csharp](https://github.com/DaredevilOSS/sqlc-gen-csharp)        | Stable | Stable     | Stable |
| F#       | [kaashyapan/sqlc-gen-fsharp](https://github.com/kaashyapan/sqlc-gen-fsharp)            | N/A    | Beta       | Beta   |
| Java     | [tandemdude/sqlc-gen-java](https://github.com/tandemdude/sqlc-gen-java)                | Beta   | Beta       | N/A    |
| PHP      | [lcarilla/sqlc-plugin-php-dbal](https://github.com/lcarilla/sqlc-plugin-php-dbal)      | Beta   | N/A        | N/A    |
| Ruby     | [DaredevilOSS/sqlc-gen-ruby](https://github.com/DaredevilOSS/sqlc-gen-ruby)            | Beta   | Beta       | Beta   |
| Zig      | [tinyzimmer/sqlc-gen-zig](https://github.com/tinyzimmer/sqlc-gen-zig)                  | N/A    | Beta       | Beta   |
| Python   | [rayakame/sqlc-gen-better-python](https://github.com/rayakame/sqlc-gen-better-python)  | Beta   | Stable     | Stable |
| Rust     | [mathematic-inc/sqlc-gen-sqlx](https://github.com/mathematic-inc/sqlc-gen-sqlx)        | N/A    | Beta       | N/A    |
| \[Any\]  | [fdietze/sqlc-gen-from-template](https://github.com/fdietze/sqlc-gen-from-template)    | Stable | Stable     | Stable |

Plugins developed by our Community can also be found using our
[github topic](https://github.com/topics/sqlc-plugin).

## Community projects

| Language | Project                                                           | MySQL  | PostgreSQL | SQLite |
| -------- | ----------------------------------------------------------------- | ------ | ---------- | ------ |
| Gleam    | [daniellionel01/parrot](https://github.com/daniellionel01/parrot) | Stable | Stable     | Stable |
