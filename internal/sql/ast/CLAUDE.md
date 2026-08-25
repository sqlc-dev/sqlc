# AST Package - Claude Code Guide

This package defines the Abstract Syntax Tree (AST) nodes used by sqlc to represent SQL statements across all supported databases (PostgreSQL, MySQL, SQLite).

## Key Concepts

### Node Interface
All AST nodes implement the `Node` interface with:
- `Pos() int` - returns the source position
- `Format(buf *TrackedBuffer)` - formats the node back to SQL

### TrackedBuffer
The `TrackedBuffer` type (`print.go`) handles SQL formatting with dialect-specific behavior:
- `astFormat(node Node)` - formats any AST node
- `join(list *List, sep string)` - joins list items with separator
- `WriteString(s string)` - writes raw SQL
- `QuoteIdent(name string)` - quotes identifiers (dialect-specific)
- `TypeName(ns, name string)` - formats type names (dialect-specific)

### Pretty printing
`TrackedBuffer` records a Wadler-style document (the model behind Prettier
and ruff): Format methods emit text plus layout tokens, and the renderer
decides which break opportunities become newlines.

- `line()` - a space when the group fits on one line, a newline otherwise
- `softline()` - nothing when the group fits, a newline otherwise
- `group()` / `endGroup()` - a region laid out flat when its width fits
- `indent()` / `endIndent()` - one level deeper (2 spaces) after breaks
- `joinComma(list)` - joins items with `,` + `line()`
- `condition(node)` - clause-level AND/OR chain without outer parentheses,
  one branch per line when broken

`ast.Format(n, d)` renders on a single line (all breaks collapse);
`ast.Pretty(n, d, width)` breaks lines to fit `width` columns. Statement
Format methods open a group and put `line()` before each clause keyword
(`FROM`, `WHERE`, ...), so a statement that fits stays on one line and a
long one breaks at clause boundaries. When adding a Format method, write
tokens so the flat rendering is correct SQL; layout tokens are optional.

### Comments
`ast.File{Stmts, Comments}` is what a comment-surfacing parser returns
(SQLite via meyer's ParseFile). `AttachComments(raw, d, comments, src)`
classifies each comment once, against a dry-run of the printer: source
positions and line numbers decide trailing (same line as the code
before) vs leading, and each comment is attached to the emission point
— node or clause/list boundary — where the printer will reach it. From
then on positions are never consulted: `PrettyWithComments(n, d, width,
table)` emits by node identity, which is what lets edited or synthetic
trees print their comments correctly (the dave/dst model; each record
also keeps the node the comment followed, for future rewriting tools).
A line comment forces every enclosing group to break — `hardline` and
`breaker` tokens measure as infinitely wide — so commented statements
format instead of collapsing. Emission points (`beforeClause`,
`boundary` in joinComma/condition) double as classification markers on
the dry run, guaranteeing attach-time decisions and print-time emission
agree.

### Dialect Interface
Dialect-specific formatting is handled via the `Dialect` interface:
```go
type Dialect interface {
    QuoteIdent(string) string
    TypeName(ns, name string) string
    Param(int) string      // $1 for PostgreSQL, ? for MySQL
    NamedParam(string) string // @name for PostgreSQL, :name for SQLite
    Cast(string) string
}
```

## Adding New AST Nodes

When adding a new AST node type:

1. **Create the node file** (e.g., `variable_expr.go`):
```go
package ast

type VariableExpr struct {
    Name     string
    Location int
}

func (n *VariableExpr) Pos() int {
    return n.Location
}

func (n *VariableExpr) Format(buf *TrackedBuffer) {
    if n == nil {
        return
    }
    buf.WriteString("@")
    buf.WriteString(n.Name)
}
```

2. **Add to `astutils/walk.go`** - Add a case in the Walk function:
```go
case *ast.VariableExpr:
    // Leaf node - no children to traverse
```

3. **Add to `astutils/rewrite.go`** - Add a case in the Apply function:
```go
case *ast.VariableExpr:
    // Leaf node - no children to traverse
```

4. **Update the parser/converter** - In the relevant engine (e.g., `dolphin/convert.go` for MySQL)

## Helper Functions for Format Methods

- `set(node Node) bool` - returns true if node is non-nil and not an empty List
- `items(list *List) bool` - returns true if list has items
- `todo(node) Node` - placeholder for unimplemented conversions (returns nil)

## Common Node Types

### Statements
- `SelectStmt` - SELECT queries with FromClause, WhereClause, etc.
- `InsertStmt` - INSERT with Relation, Cols, SelectStmt, OnConflictClause
- `UpdateStmt` - UPDATE with Relations, TargetList, WhereClause
- `DeleteStmt` - DELETE with Relations, FromClause (for JOINs), Targets

### Expressions
- `A_Expr` - General expression with operator (e.g., `a + b`, `@param`)
- `ColumnRef` - Column reference with Fields list
- `FuncCall` - Function call with Func, Args, aggregation options
- `TypeCast` - Type cast with Arg and TypeName
- `ParenExpr` - Parenthesized expression
- `VariableExpr` - MySQL user variable (e.g., `@user_id`)

### Table References
- `RangeVar` - Table reference with schema, name, alias
- `JoinExpr` - JOIN with Larg, Rarg, Jointype, Quals/UsingClause

## MySQL-Specific Nodes

- `VariableExpr` - User variables (`@var`), distinct from sqlc's `@param` syntax
- `IntervalExpr` - INTERVAL expressions
- `OnDuplicateKeyUpdate` - MySQL's ON DUPLICATE KEY UPDATE clause
- `ParenExpr` - Explicit parentheses (TiDB parser wraps expressions)

## Important Distinctions

### MySQL @variable vs sqlc @param
- sqlc named parameters (`@param`) never reach the AST: `internal/sql/preprocess`
  replaces them with the dialect's native placeholder before parsing, so engines
  only ever see a `ParamRef`
- `@name` is only sqlc syntax where the preprocessor's dialect says so. MySQL is
  excluded, so `@user_id` there stays a user variable (`VariableExpr`)

### Type Modifiers
- `TypeName.Typmods` holds type modifiers like `varchar(255)`
- For MySQL, only populate Typmods for types where length is user-specified:
  - VARCHAR, CHAR, VARBINARY, BINARY - need length
  - DATETIME, TIMESTAMP, DATE - internal flen should NOT be output
