# AST Package - Claude Code Guide

This package defines the Abstract Syntax Tree (AST) nodes used by sqlc to represent SQL statements across all supported databases (PostgreSQL, MySQL, SQLite).

## Key Concepts

### Node Interface
All AST nodes implement the `Node` interface with:
- `Pos() int` - returns the source position
- `Format(buf *TrackedBuffer)` - formats the node back to SQL

### TrackedBuffer
The `TrackedBuffer` type (`pg_query.go`) handles SQL formatting with dialect-specific behavior:
- `astFormat(node Node)` - formats any AST node
- `join(list *List, sep string)` - joins list items with separator
- `WriteString(s string)` - writes raw SQL
- `QuoteIdent(name string)` - quotes identifiers (dialect-specific)
- `TypeName(ns, name string)` - formats type names (dialect-specific)

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
- MySQL user variables (`@user_id`) use `VariableExpr` - preserved as-is in output
- sqlc named parameters (`@param`) use `A_Expr` with `@` operator - replaced with `?`
- The `named.IsParamSign()` function checks for `A_Expr` with `@` operator

### Type Modifiers
- `TypeName.Typmods` holds type modifiers like `varchar(255)`
- For MySQL, only populate Typmods for types where length is user-specified:
  - VARCHAR, CHAR, VARBINARY, BINARY - need length
  - DATETIME, TIMESTAMP, DATE - internal flen should NOT be output
