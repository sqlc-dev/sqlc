# Changelog
All notable changes to this project will be documented in this file.

## [1.17.2](https://github.com/kyleconroy/sqlc/releases/tag/1.17.2)
Released 2023-02-22

### Bug Fixes

- Fix build on Windows (#2102)

## [1.17.1](https://github.com/kyleconroy/sqlc/releases/tag/1.17.1)
Released 2023-02-22

### Bug Fixes

- Prefer to use []T over pgype.Array[T] (#2090)
- Revert changes to Dockerfile (#2091)
- Do not throw error when IF NOT EXISTS is used on ADD COLUMN (#2092)

### MySQL

- Add `float` support to MySQL (#2097)

### Build

- (deps) Bump golang from 1.20.0 to 1.20.1 (#2082)

## [1.17.0](https://github.com/kyleconroy/sqlc/releases/tag/1.17.0)
Released 2023-02-13

### Bug Fixes

- Initialize generated code outside function (#1850)
- (engine/mysql) Take into account column's charset to distinguish text/blob, (var)char/(var)binary (#776) (#1895)
- The enum Value method returns correct type (#1996)
- Documentation for Inserting Rows (#2034)
- Add import statements even if only pointer types exist (#2046)
- Search from Rexpr if not found from Lexpr (#2056)

### Documentation

- Change ENTRYPOINT to CMD (#1943)
- Update samples for HOW-TO GUIDES (#1953)

### Features

- Add the diff command (#1963)

### Build

- (deps) Bump github.com/mattn/go-sqlite3 from 1.14.15 to 1.14.16 (#1913)
- (deps) Bump github.com/spf13/cobra from 1.6.0 to 1.6.1 (#1909)
- Fix devcontainer (#1942)
- Run sqlc-pg-gen via GitHub Actions (#1944)
- Move large arrays out of functions (#1947)
- Fix conflicts from pointer configuration (#1950)
- (deps) Bump github.com/go-sql-driver/mysql from 1.6.0 to 1.7.0 (#1988)
- (deps) Bump github.com/jackc/pgtype from 1.12.0 to 1.13.0 (#1978)
- (deps) Bump golang from 1.19.3 to 1.19.4 (#1992)
- (deps) Bump certifi from 2020.12.5 to 2022.12.7 in /docs (#1993)
- (deps) Bump golang from 1.19.4 to 1.19.5 (#2016)
- (deps) Bump golang from 1.19.5 to 1.20.0 (#2045)
- (deps) Bump github.com/jackc/pgtype from 1.13.0 to 1.14.0 (#2062)
- (deps) Bump github.com/jackc/pgx/v4 from 4.17.2 to 4.18.0 (#2063)

### Cmd

- Generate packages in parallel (#2026)

### Cmd/sqlc

- Bump version to v1.17.0

### Codegen

- Remove built-in Kotlin support (#1935)
- Remove built-in Python support (#1936)

### Internal/codegen

- Cache pattern matching compilations (#2028)

### Mysql

- Add datatype tests (#1948)
- Fix blob tests (#1949)

### Plugins

- Upgrade to wasmtime 3.0.1 (#2009)

### Sqlite

- Supported between expr (#1958) (#1967)

### Tools

- Regenerate scripts skips dirs that contains diff exec command (#1987)

### Wasm

- Upgrade to wasmtime 5.0.0 (#2065)

## [1.16.0](https://github.com/kyleconroy/sqlc/releases/tag/1.16.0)
Released 2022-11-09


### Bug Fixes

- (validate) Sqlc.arg & sqlc.narg are not "missing" (#1814)
- Emit correct comment for nullable enums (#1819)
- 🐛 Correctly switch `coalesce()` result `.NotNull` value (#1664)
- Prevent batch infinite loop with arg length (#1794)
- Support version 2 in error message (#1839)
- Handle empty column list in postgresql (#1843)
- Batch imports filter queries, update cmds having ret type (#1842)
- Named params contribute to batch parameter count (#1841)

### Documentation

- Add a getting started guide for SQLite (#1798)
- Various readability improvements (#1854)
- Add documentation for codegen plugins (#1904)
- Update migration guides with links (#1933)

### Features

- Add HAVING support to MySQL (#1806)

### Miscellaneous Tasks

- Upgrade wasmtime version (#1827)
- Bump wasmtime version to v1.0.0 (#1869)

### Build

- (deps) Bump github.com/jackc/pgconn from 1.12.1 to 1.13.0 (#1785)
- (deps) Bump github.com/mattn/go-sqlite3 from 1.14.13 to 1.14.15 (#1799)
- (deps) Bump github.com/jackc/pgx/v4 from 4.16.1 to 4.17.0 (#1786)
- (deps) Bump github.com/jackc/pgx/v4 from 4.17.0 to 4.17.1 (#1825)
- (deps) Bump github.com/bytecodealliance/wasmtime-go (#1826)
- (deps) Bump github.com/jackc/pgx/v4 from 4.17.1 to 4.17.2 (#1831)
- (deps) Bump golang from 1.19.0 to 1.19.1 (#1834)
- (deps) Bump github.com/google/go-cmp from 0.5.8 to 0.5.9 (#1838)
- (deps) Bump github.com/lib/pq from 1.10.6 to 1.10.7 (#1835)
- (deps) Bump github.com/bytecodealliance/wasmtime-go (#1857)
- (deps) Bump github.com/spf13/cobra from 1.5.0 to 1.6.0 (#1893)
- (deps) Bump golang from 1.19.1 to 1.19.3 (#1920)

### Cmd/sqlc

- Bump to v1.16.0

### Codgen

- Include serialized codegen options (#1890)

### Compiler

- Move Kotlin parameter logic into codegen (#1910)

### Examples

- Port Python examples to WASM plugin (#1903)

### Pg-gen

- Make sqlc-pg-gen the complete source of truth for pg_catalog.go (#1809)
- Implement information_schema shema (#1815)

### Python

- Port all Python tests to sqlc-gen-python (#1907)
- Upgrade to sqlc-gen-python v1.0.0 (#1932)

## [1.15.0](https://github.com/kyleconroy/sqlc/releases/tag/1.15.0)
Released 2022-08-07

### Bug Fixes

- (mysql) Typo (#1700)
- (postgresql) Add quotes for CamelCase columns (#1729)
- Cannot parse SQLite upsert statement (#1732)
- (sqlite) Regenerate test output for builtins (#1735)
- (wasm) Version modules by wasmtime version (#1734)
- Missing imports (#1637)
- Missing slice import for querier (#1773)

### Documentation

- Add process-based plugin docs (#1669)
- Add links to downloads.sqlc.dev (#1681)
- Update transactions how to example (#1775)

### Features

- More SQL Syntax Support for SQLite (#1687)
- (sqlite) Promote SQLite support to beta (#1699)
- Codegen plugins, powered by WASM (#1684)
- Set user-agent for plugin downloads (#1707)
- Null enums types (#1485)
- (sqlite) Support stdlib functions (#1712)
- (sqlite) Add support for returning (#1741)

### Miscellaneous Tasks

- Add tests for quoting columns (#1733)
- Remove catalog tests (#1762)

### Testing

- Add tests for fixing slice imports (#1736)
- Add test cases for returning (#1737)

### Build

- Upgrade to Go 1.19 (#1780)
- Upgrade to go-wasmtime 0.39.0 (#1781)

### Plugins

- (wasm) Change default cache location (#1709)
- (wasm) Change the SHA-256 config key (#1710)

## [1.14.0](https://github.com/kyleconroy/sqlc/releases/tag/1.14.0)
Released 2022-06-09

### Bug Fixes

- (postgresql) Remove extra newline with db argument (#1417)
- (sqlite) Fix DROP TABLE   (#1443)
- (compiler) Fix left join nullability with table aliases (#1491)
- Regenerate testdata for CREATE TABLE AS (#1516)
- (bundler) Only close multipart writer once (#1528)
- (endtoend) Regenerate testdata for exex_lastid
- (pgx) Copyfrom imports (#1626)
- Validate sqlc function arguments (#1633)
- Fixed typo `sql.narg` in doc (#1668)

### Features

- (golang) Add Enum.Valid and AllEnumValues (#1613)
- (sqlite) Start expanding support (#1410)
- (pgx) Add support for batch operations (#1437)
- (sqlite) Add support for delete statements (#1447)
- (codegen) Insert comments in interfaces (#1458)
- (sdk) Add the plugin SDK package (#1463)
- Upload projects (#1436)
- Add sqlc version to generated Kotlin code (#1512)
- Add sqlc version to generated Go code (#1513)
- Pass sqlc version in codegen request (#1514)
- (postgresql) Add materialized view support (#1509)
- (python) Graduate Python support to beta (#1520)
- Run sqlc with docker on windows cmd (#1557)
- Add JSON "codegen" output (#1565)
- Add sqlc.narg() for nullable named params (#1536)
- Process-based codegen plugins (#1578)

### Miscellaneous Tasks

- Fix extra newline in comments for copyfrom (#1438)
- Generate marshal/unmarshal with vtprotobuf (#1467)

### Refactor

- (codegen) Port Kotlin codegen package to use plugin types (#1416)
- (codegen) Port Go to plugin types (#1460)
- (cmd) Simplify codegen selection logic (#1466)
- (sql/catalog) Improve Readability (#1595)
- Add basic fuzzing for config / overrides (#1500)

## [1.13.0](https://github.com/kyleconroy/sqlc/releases/tag/v1.13.0)
Released 2022-03-31

### Bug Fixes

- (compiler) Fix left join nullability with table aliases (#1491)
- (postgresql) Remove extra newline with db argument (#1417)
- (sqlite) Fix DROP TABLE (#1443)

### Features

- (cli) Upload projects (#1436)
- (codegen) Add sqlc version to generated Go code (#1513)
- (codegen) Add sqlc version to generated Kotlin code (#1512)
- (codegen) Insert comments in interfaces (#1458)
- (codegen) Pass sqlc version in codegen request (#1514)
- (pgx) Add support for batch operations (#1437)
- (postgresql) Add materialized view support (#1509)
- (python) Graduate Python support to beta (#1520)
- (sdk) Add the plugin SDK package (#1463)
- (sqlite) Add support for delete statements (#1447)
- (sqlite) Start expanding support (#1410)

### Miscellaneous Tasks

- Fix extra newline in comments for copyfrom (#1438)
- Generate marshal/unmarshal with vtprotobuf (#1467)

### Refactor

- (codegen) Port Kotlin codegen package to use plugin types (#1416)
- (codegen) Port Go to plugin types (#1460)
- (cmd) Simplify codegen selection logic (#1466)

### Config

- Add basic fuzzing for config / overrides (#1500)

## [1.12.0](https://github.com/kyleconroy/sqlc/releases/tag/v1.12.0)
Released 2022-02-05

### Bug

- ALTER TABLE SET SCHEMA (#1409)

### Bug Fixes

- Update ANTLR v4 go.mod entry (#1336)
- Check delete statements for CTEs (#1329)
- Fix validation of GROUP BY on field aliases (#1348)
- Fix imports when non-copyfrom queries needed imports that copyfrom queries didn't (#1386)
- Remove extra comment newline (#1395)
- Enable strict function checking (#1405)

### Documentation

- Bump version to 1.11.0 (#1308)

### Features

- Inheritance (#1339)
- Generate query code using ASTs instead of templates (#1338)
- Add support for CREATE TABLE a ( LIKE b ) (#1355)
- Add support for sql.NullInt16 (#1376)

### Miscellaneous Tasks

- Add tests for :exec{result,rows} (#1344)
- Delete template-based codegen (#1345)

### Build

- Bump github.com/jackc/pgx/v4 from 4.14.0 to 4.14.1 (#1316)
- Bump golang from 1.17.3 to 1.17.4 (#1331)
- Bump golang from 1.17.4 to 1.17.5 (#1337)
- Bump github.com/spf13/cobra from 1.2.1 to 1.3.0 (#1343)
- Remove devel Docker build
- Bump golang from 1.17.5 to 1.17.6 (#1369)
- Bump github.com/google/go-cmp from 0.5.6 to 0.5.7 (#1382)
- Format all Go code (#1387)

## [1.11.0](https://github.com/kyleconroy/sqlc/releases/tag/v1.11.0)
Released 2021-11-24


### Bug Fixes

- Update incorrect signatures (#1180)
- Correct aggregate func sig (#1182)
- Jsonb_build_object (#1211)
- Case-insensitive identifiers (#1216)
- Incorrect handling of meta (#1228)
- Detect invalid INSERT expression (#1231)
- Respect alias name for coalesce (#1232)
- Mark nullable when casting NULL (#1233)
- Support nullable fields in joins for MySQL engine (#1249)
- Fix between expression handling of table references (#1268)
- Support nullable fields in joins on same table (#1270)
- Fix missing binds in ORDER BY (#1273)
- Set RV for TargetList items on updates (#1252)
- Fix MySQL parser for query without trailing semicolon (#1282)
- Validate table alias references (#1283)
- Add support for MySQL ON DUPLICATE KEY UPDATE (#1286)
- Support references to columns in joined tables in UPDATE statements (#1289)
- Add validation for GROUP BY clause column references (#1285)
- Prevent variable redeclaration in single param conflict (#1298)
- Use common params struct field for same named params (#1296)

### Documentation

- Replace deprecated go get with go install (#1181)
- Fix package name referenced in tutorial (#1202)
- Add environment variables (#1264)
- Add go.17+ install instructions (#1280)
- Warn about golang-migrate file order (#1302)

### Features

- Instrument compiler via runtime/trace (#1258)
- Add MySQL support for BETWEEN arguments (#1265)

### Refactor

- Move from io/ioutil to io and os package (#1164)

### Styling

- Apply gofmt to sample code (#1261)

### Build

- Bump golang from 1.17.0 to 1.17.1 (#1173)
- Bump eskatos/gradle-command-action from 1 to 2 (#1220)
- Bump golang from 1.17.1 to 1.17.2 (#1227)
- Bump github.com/pganalyze/pg_query_go/v2 (#1234)
- Bump actions/checkout from 2.3.4 to 2.3.5 (#1238)
- Bump babel from 2.9.0 to 2.9.1 in /docs (#1245)
- Bump golang from 1.17.2 to 1.17.3 (#1272)
- Bump actions/checkout from 2.3.5 to 2.4.0 (#1267)
- Bump github.com/lib/pq from 1.10.3 to 1.10.4 (#1278)
- Bump github.com/jackc/pgx/v4 from 4.13.0 to 4.14.0 (#1303)

### Cmd/sqlc

- Bump version to v1.11.0

## [1.10.0](https://github.com/kyleconroy/sqlc/releases/tag/v1.10.0)
Released 2021-09-07


### Documentation

- Fix invalid language support table (#1161)
- Add a getting started guide for MySQL (#1163)

### Build

- Bump golang from 1.16.7 to 1.17.0 (#1129)
- Bump github.com/lib/pq from 1.10.2 to 1.10.3 (#1160)

### Ci

- Upgrade Go to 1.17 (#1130)

### Cmd/sqlc

- Bump version to v1.10.0 (#1165)

### Codegen/golang

- Consolidate import logic (#1139)
- Add pgx support for range types (#1146)
- Use pgtype for hstore when using pgx (#1156)

### Codgen/golang

- Use p[gq]type for network address types (#1142)

### Endtoend

- Run `go test` in CI (#1134)

### Engine/mysql

- Add support for LIKE (#1162)

### Golang

- Output NullUUID when necessary (#1137)

## [1.9.0](https://github.com/kyleconroy/sqlc/releases/tag/v1.9.0)
Released 2021-08-13


### Documentation

- Update documentation (a bit) for v1.9.0 (#1117)

### Build

- Bump golang from 1.16.6 to 1.16.7 (#1107)

### Cmd/sqlc

- Bump version to v1.9.0 (#1121)

### Compiler

- Add tests for COALESCE behavior (#1112)
- Handle subqueries in SELECT statements (#1113)

## [1.8.0](https://github.com/kyleconroy/sqlc/releases/tag/v1.8.0)
Released 2021-05-03


### Documentation

- Add language support Matrix (#920)

### Features

- Add case style config option (#905)

### Python

- Eliminate runtime package and use sqlalchemy (#939)

### Build

- Bump github.com/google/go-cmp from 0.5.4 to 0.5.5 (#926)
- Bump github.com/lib/pq from 1.9.0 to 1.10.0 (#931)
- Bump golang from 1.16.0 to 1.16.1 (#935)
- Bump golang from 1.16.1 to 1.16.2 (#942)
- Bump github.com/jackc/pgx/v4 from 4.10.1 to 4.11.0 (#956)
- Bump github.com/go-sql-driver/mysql from 1.5.0 to 1.6.0 (#961)
- Bump github.com/pganalyze/pg_query_go/v2 (#965)
- Bump urllib3 from 1.26.3 to 1.26.4 in /docs (#968)
- Bump golang from 1.16.2 to 1.16.3 (#963)
- Bump github.com/lib/pq from 1.10.0 to 1.10.1 (#980)

### Cmd

- Add the --experimental flag (#929)
- Fix sqlc init (#959)

### Cmd/sqlc

- Bump version to v1.7.1-devel (#913)
- Bump version to v1.8.0

### Codegen

- Generate valid enum names for symbols (#972)

### Postgresql

- Support generated columns
- Add test for PRIMARY KEY INCLUDE
- Add tests for CREATE TABLE PARTITION OF
- CREATE TRIGGER EXECUTE FUNCTION
- Add support for renaming types (#971)

### Sql/ast

- Resolve return values from functions (#964)

### Workflows

- Only run tests once (#924)

## [1.7.0](https://github.com/kyleconroy/sqlc/releases/tag/v1.7.0)
Released 2021-02-28


### Bug Fixes

- Struct tag formatting (#833)

### Documentation

- Include all the existing Markdown files (#877)
- Split docs into four sections (#882)
- Reorganize and consolidate documentation
- Add link to Windows download (#888)
- Shorten the README (#889)

### Features

- Adding support for pgx/v4
- Adding support for pgx/v4

### README

- Add Go Report Card badge (#891)

### Build

- Bump github.com/google/go-cmp from 0.5.3 to 0.5.4 (#813)
- Bump github.com/lib/pq from 1.8.0 to 1.9.0 (#820)
- Bump golang from 1.15.5 to 1.15.6 (#822)
- Bump github.com/jackc/pgx/v4 from 4.9.2 to 4.10.0 (#823)
- Bump github.com/jackc/pgx/v4 from 4.10.0 to 4.10.1 (#839)
- Bump golang from 1.15.6 to 1.15.7 (#855)
- Bump golang from 1.15.7 to 1.15.8 (#881)
- Bump github.com/spf13/cobra from 1.1.1 to 1.1.2 (#892)
- Bump golang from 1.15.8 to 1.16.0 (#897)
- Bump github.com/lfittl/pg_query_go from 1.0.1 to 1.0.2 (#901)
- Bump github.com/spf13/cobra from 1.1.2 to 1.1.3 (#893)

### Catalog

- Improve alter column type (#818)

### Ci

- Uprade to Go 1.15 (#887)

### Cmd

- Allow config file location to be specified (#863)

### Cmd/sqlc

- Bump to version v1.6.1-devel (#807)
- Bump version to v1.7.0 (#912)

### Codegen/golang

- Make sure to import net package (#858)

### Compiler

- Support UNION query

### Dolphin

- Generate bools for tinyint(1)
- Support joins in update statements (#883)
- Add support for union query

### Endtoend

- Add tests for INTERSECT and EXCEPT

### Go.mod

- Update to go 1.15 and run 'go mod tidy' (#808)

### Mysql

- Compile tinyint(1) to bool (#873)

### Sql/ast

- Add enum values for SetOperation

## [1.6.0](https://github.com/kyleconroy/sqlc/releases/tag/v1.6.0)
Released 2020-11-23


### Dolphin

- Implement Rename (#651)
- Skip processing view drops (#653)

### README

- Update language / database support (#698)

### Astutils

- Fix Params rewrite call (#674)

### Build

- Bump golang from 1.14 to 1.15.3 (#765)
- Bump docker/build-push-action from v1 to v2.1.0 (#764)
- Bump github.com/google/go-cmp from 0.4.0 to 0.5.2 (#766)
- Bump github.com/spf13/cobra from 1.0.0 to 1.1.1 (#767)
- Bump github.com/jackc/pgx/v4 from 4.6.0 to 4.9.2 (#768)
- Bump github.com/lfittl/pg_query_go from 1.0.0 to 1.0.1 (#773)
- Bump github.com/google/go-cmp from 0.5.2 to 0.5.3 (#783)
- Bump golang from 1.15.3 to 1.15.5 (#782)
- Bump github.com/lib/pq from 1.4.0 to 1.8.0 (#769)

### Catalog

- Improve variadic argument support (#804)

### Cmd/sqlc

- Bump to version v1.6.0 (#806)

### Codegen

- Fix errant database/sql imports (#789)

### Compiler

- Use engine-specific reserved keywords (#677)

### Dolphi

- Add list of builtin functions (#795)

### Dolphin

- Update to the latest MySQL parser (#665)
- Add ENUM() support (#676)
- Add test for table aliasing (#684)
- Add MySQL ddl_create_table test (#685)
- Implete TRUNCATE table (#697)
- Represent tinyint as int32 (#797)
- Add support for coalesce (#802)
- Add function signatures (#796)

### Endtoend

- Add MySQL json test (#692)
- Add MySQL update set multiple test (#696)

### Examples

- Use generated enum constants in db_test (#678)
- Port ondeck to MySQL (#680)
- Add MySQL authors example (#682)

### Internal/cmd

- Print correct config file on parse failure (#749)

### Kotlin

- Remove runtime dependency (#774)

### Metadata

- Support multiple comment prefixes (#683)

### Postgresql

- Support string concat operator (#701)

### Sql/catalog

- Add support for variadic functions (#798)

## [1.5.0](https://github.com/kyleconroy/sqlc/releases/tag/v1.5.0)
Released 2020-08-05


### Documentation

- Build sqlc using Go 1.14 (#549)

### Cmd

- Add debugging support (#573)

### Cmd/sqlc

- Bump version to v1.4.1-devel (#548)
- Bump version to v1.5.0

### Compiler

- Support calling functions with defaults (#635)
- Skip func args without a paramRef (#636)
- Return a single column from coalesce (#639)

### Config

- Add emit_empty_slices to version one (#552)

### Contrib

- Add generated code for contrib

### Dinosql

- Remove deprecated package (#554)

### Dolphin

- Add support for column aliasing (#566)
- Implement star expansion for subqueries (#619)
- Implement exapansion with reserved words (#620)
- Implement parameter refs (#621)
- Implement limit and offest (#622)
- Implement inserts (#623)
- Implement delete (#624)
- Implement simple update statements (#625)
- Implement INSERT ... SELECT (#626)
- Use test driver instead of TiDB driver (#629)
- Implement named parameters via sqlc.arg() (#632)

### Endtoend

- Add MySQL test for SELECT * JOIN (#565)
- Add MySQL test for inflection (#567)

### Engine

- Create engine package (#556)

### Equinox

- Use the new equinox-io/setup action (#586)

### Examples

- Run tests for MySQL booktest (#627)

### Golang

- Add support for the money type (#561)
- Generate correct types for int2 and int8 (#579)

### Internal

- Rm catalog, pg, postgres packages (#555)

### Mod

- Downgrade TiDB package to fix build (#603)

### Mysql

- Upgrade to the latest vitess commit (#562)
- Support to infer type of a duplicated arg (#615)
- Allow some builtin functions to be nullable (#616)

### Postgresql

- Generate all functions in pg_catalog (#550)
- Remove pg_catalog schema from tests (#638)
- Move contrib code to a package

### Sql/catalog

- Fix comparison of pg_catalog types (#637)

### Tools

- Generate functions for all of contrib

### Workflow

- Migrate to equinox-io/setup-release-tool (#614)

## [1.4.0](https://github.com/kyleconroy/sqlc/releases/tag/v1.4.0)
Released 2020-06-17


### Dockerfile

- Add version build argument (#487)

### MySQL

- Prevent Panic when WHERE clause contains parenthesis.  (#531)

### README

- Document emit_exact_table_names (#486)

### All

- Remove the exp build tag (#507)

### Catalog

- Support functions with table parameters (#541)

### Cmd

- Bump to version 1.3.1-devel (#485)

### Cmd/sqlc

- Bump version to v1.4.0 (#547)

### Codegen

- Add the new codegen packages (#513)
- Add the :execresult query annotation (#542)

### Compiler

- Validate function calls (#505)
- Port bottom of parseQuery (#510)
- Don't mutate table name (#517)
- Enable experimental parser by default (#518)
- Apply rename rules to enum constants (#523)
- Temp fix for typecast function parameters (#530)

### Endtoend

- Standardize JSON formatting (#490)
- Add per-test configuration files (#521)
- Read expected stderr failures from disk (#527)

### Internal/dinosql

- Check parameter style before ref (#488)
- Remove unneeded column suffix (#492)
- Support named function arguments (#494)

### Internal/postgresql

- Fix NamedArgExpr rewrite (#491)

### Multierr

- Move dinosql.ParserErr to a new package (#496)

### Named

- Port parameter style validation to SQL (#504)

### Parser

- Support columns from subselect statements (#489)

### Rewrite

- Move parameter rewrite to package (#499)

### Sqlite

- Use convert functions instead of the listener (#519)

### Sqlpath

- Move ReadSQLFiles into a separate package (#495)

### Validation

- Move query validation to separate package (#498)

## [1.3.0](https://github.com/kyleconroy/sqlc/releases/tag/v1.3.0)
Released 2020-05-12


### Makefile

- Update target (#449)

### README

- Add Myles as a sponsor (#469)

### Testing

- Make sure all Go examples build (#480)

### Cmd

- Bump version to v1.3.0 (#484)

### Cmd/sqlc

- Bump version to v1.2.1-devel (#442)

### Dinosql

- Inline addFile (#446)
- Add PostgreSQL support for TRUNCATE (#448)

### Gen

- Emit json.RawMessage for JSON columns (#461)

### Go.mod

- Use latest lib/pq (#471)

### Parser

- Use same function to load SQL files (#483)

### Postgresql

- Fix panic walking CreateTableAsStmt (#475)

## [1.2.0](https://github.com/kyleconroy/sqlc/releases/tag/v1.2.0)
Released 2020-04-07


### Documentation

- Publish to Docker Hub (#422)

### README

- Docker installation docs (#424)

### Cmd/sqlc

- Bump version to v1.1.1-devel (#407)
- Bump version to v1.2.0 (#441)

### Gen

- Add special case for "campus" (#435)
- Properly quote reserved keywords on expansion (#436)

### Migrations

- Move migration parsing to new package (#427)

### Parser

- Generate correct types for SELECT EXISTS (#411)

## [1.1.0](https://github.com/kyleconroy/sqlc/releases/tag/v1.1.0)
Released 2020-03-17


### README

- Add installation instructions (#350)
- Add section on running tests (#357)
- Fix typo (#371)

### Ast

- Add AST for ALTER TABLE ADD / DROP COLUMN (#376)
- Add support for CREATE TYPE as ENUM (#388)
- Add support for CREATE / DROP SCHEMA (#389)

### Astutils

- Apply changes to the ValuesList slice (#372)

### Cmd

- Return v1.0.0 (#348)
- Return next bug fix version (#349)

### Cmd/sqlc

- Bump version to v1.1.0 (#406)

### Compiler

- Wire up the experimental parsers

### Config

- Remove "emit_single_file" option (#367)

### Dolphin

- Add experimental parser for MySQL

### Gen

- Add option to emit single file for Go (#366)
- Add support for the ltree extension (#385)

### Go.mod

- Add packages for MySQL and SQLite parsers

### Internal/dinosql

- Support Postgres macaddr type in Go (#358)

### Internal/endtoend

- Remove %w (#354)

### Kotlin

- Add Query class to support timeout and cancellation (#368)

### Postgresql

- Add experimental parser for MySQL

### Sql

- Add generic SQL AST

### Sql/ast

- Port support for COMMENT ON (#391)
- Implement DROP TYPE (#397)
- Implement ALTER TABLE RENAME (#398)
- Implement ALTER TABLE RENAME column (#399)
- Implement ALTER TABLE SET SCHEMA (#400)

### Sql/catalog

- Port tests over from catalog pkg (#402)

### Sql/errors

- Add a new errors package (#390)

### Sqlite

- Add experimental parser for SQLite

## [1.0.0](https://github.com/kyleconroy/sqlc/releases/tag/v1.0.0)
Released 2020-02-18


### Documentation

- Add documentation for query commands (#270)
- Add named parameter documentation (#332)

### README

- Add sponsors section (#333)

### Cmd

- Remove parse subcommand (#322)

### Config

- Parse V2 config format
- Add support for YAML (#336)

### Examples

- Add the jets and booktest examples (#237)
- Move sqlc.json into examples folder (#238)
- Add the authors example (#241)
- Add build tag to authors tests (#319)

### Internal

- Allow CTE to be used with UPDATE (#268)
- Remove the PackageMap from settings (#295)

### Internal/config

- Create new config package (#313)

### Internal/dinosql

- Emit Querier interface (#240)
- Strip leading "go-" or trailing "-go" from import (#262)
- Overrides can now be basic types (#271)
- Import needed types for Querier (#285)
- Handle schema-scoped enums (#310)
- Ignore golang-migrate rollbacks (#320)

### Internal/endtoend

- Move more tests to the record/replay framework
- Add update test for named params (#329)

### Internal/mysql

- Fix flaky test (#242)
- Port tests to endtoend package (#315)

### Internal/parser

- Resolve nested CTEs (#324)
- Error if last query is missing (#325)
- Support joins with aliases (#326)
- Remove print statement (#327)

### Internal/sqlc

- Add support for composite types (#311)

### Kotlin

- Support primitives
- Arrays, enums, and dates
- Generate examples
- README for examples
- Factor out db setup extension
- Fix enums, use List instead of Array
- Port Go tests for examples
- Rewrite numbered params to positional params
- Always use use, fix indents
- Unbox query params

### Parser

- Attach range vars to insert params
- Attach range vars to insert params (#342)
- Remove dead code (#343)

## [0.1.0](https://github.com/kyleconroy/sqlc/releases/tag/v0.1.0)
Released 2020-01-07


### Documentation

- Replace remaining references to DinoSQL with sqlc (#149)

### README

- Fix download links (#66)
- Add LIMIT 1 to query that should return one (#99)

### Catalog

- Support "ALTER TABLE ... DROP CONSTRAINT ..." (#34)
- Differentiate functions with different argument types (#51)

### Ci

- Enable tests on pull requests

### Cmd

- Include filenames in error messages (#69)
- Do not output any changes on error (#72)

### Dinosql/internal

- Add lower and upper functions (#215)
- Ignore alter sequence commands (#219)

### Gen

- Add DO NOT EDIT comments to generated code (#50)
- Include all schemas when generating models (#90)
- Prefix structs with schema name (#91)
- Generate single import for uuid package (#98)
- Use same import logic for all Go files
- Pick correct struct to return for queries (#107)
- Create consistent JSON tags (#110)
- Add Close method to Queries struct (#127)
- Ignore empty override settings (#128)
- Turn SQL comments into Go comments (#136)

### Internal/catalog

- Parse unnamed function arguments (#166)

### Internal/dinosql

- Prepare() with no GoQueries still valid (#95)
- Fix multiline comment rendering (#142)
- Dereference alias nodes on walk (#158)
- Ignore sql-migrate rollbacks (#160)
- Sort imported packages (#165)
- Add support for timestamptz (#169)
- Error on missing queries (#180)
- Use more database/sql null types (#182)
- Support the pg_temp schema (#183)
- Override columns with array type (#184)
- Implement robust expansion
- Implement robust expansion (#186)
- Add COMMENT ON support (#191)
- Add DATE support
- Add DATE support (#196)
- Filter out invalid characters (#198)
- Quote reserved keywords (#205)
- Return parser errors first (#207)
- Implement advisory locks (#212)
- Error on duplicate query names (#221)
- Fix incorrect enum names (#223)
- Add support for numeric types
- Add support for numeric types (#228)

### Internal/dinosql/testdata/ondeck

- Add Makefile (#156)

### Ondeck

- Move all tests to GitHub CI (#58)

### ParseQuery

- Return either a query or an error (#178)

### Parser

- Use schema when resolving catalog refs (#82)
- Support function calls in expressions (#104)
- Correctly handle single files (#119)
- Return error if missing RETURNING (#131)
- Add support for mathmatical operators (#132)
- Add support for simple case expressions (#134)
- Error on mismatched INSERT input (#135)
- Set IsArray on joined columns (#139)

### Pg

- Store functions in the catalog (#41)
- Add location to errors (#73)

<!-- generated by git-cliff -->
