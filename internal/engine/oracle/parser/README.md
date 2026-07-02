# Oracle PL/SQL parser (ANTLR-generated)

This package contains the ANTLR-generated Go lexer/parser for Oracle SQL/PL-SQL,
following the same pattern as [`internal/engine/sqlite/parser`](../../sqlite/parser).

## Provenance

- **Grammar:** `PlSqlLexer.g4`, `PlSqlParser.g4` are vendored from
  [antlr/grammars-v4](https://github.com/antlr/grammars-v4) (`sql/plsql`). They have
  been transformed for the Go target: the embedded `this.` actions were rewritten to
  `p.` (parser) / `l.` (lexer) via `transformGrammar.py`.
- **Base classes:** `plsql_lexer_base.go` and `plsql_parser_base.go` are hand-written
  and supply the grammar `superClass` implementations
  (`PlSqlLexerBase` / `PlSqlParserBase`). They embed `antlr.BaseLexer` /
  `antlr.BaseParser` from `github.com/antlr4-go/antlr/v4`.
- **ANTLR version:** 4.13.1, matching the `github.com/antlr4-go/antlr/v4` runtime
  version pinned in `go.mod` (and the version used by the SQLite engine).

## Generated files (committed, do not edit)

- `plsql_lexer.go`
- `plsql_parser.go`
- `plsqlparser_listener.go`
- `plsqlparser_base_listener.go`
- `*.interp`, `*.tokens`

Because these are committed, a normal `go build` does **not** require Java. Java +
the ANTLR jar are only needed to regenerate.

## Regenerating

```sh
# from this directory
make fetch-grammar   # refresh .g4 from grammars-v4 and re-apply the Go transform
make                 # downloads the ANTLR jar (if missing) and regenerates the Go code
```

The ANTLR jar (`antlr-4.13.1-complete.jar`) is git-ignored.

## Entry point

The top-level parser rule is `Sql_script()` (grammar rule `sql_script`). Individual
statements are reachable via `Unit_statement()` /
`Data_manipulation_language_statements()`. The Phase 2 AST converter consumes these
to produce sqlc's shared AST in [`internal/sql/ast`](../../../sql/ast).
