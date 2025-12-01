# SQL Rewrite Package - Claude Code Guide

This package handles AST transformations, primarily for parameter handling.

## Key Functions

### NamedParameters
`NamedParameters(engine config.Engine, raw *ast.RawStmt, ...) (*ast.RawStmt, map[int]Parameter, error)`

Finds and replaces named parameters (`sqlc.arg()`, `@param`) with positional parameters.

The function:
1. Searches for named parameters using `named.IsParamFunc()` and `named.IsParamSign()`
2. Extracts parameter names and types
3. Replaces them with `ast.ParamRef` nodes
4. Returns a map of parameter positions to their metadata

### Expand
`Expand(raw *ast.RawStmt, expected int) error`

Expands `sqlc.slice()` parameters into the correct number of positional parameters.

## How Parameter Rewriting Works

### Step 1: Find Named Parameters
```go
refs := astutils.Search(raw.Stmt, func(node ast.Node) bool {
    return named.IsParamFunc(node) || named.IsParamSign(node)
})
```

### Step 2: Replace with ParamRef
```go
astutils.Apply(raw.Stmt, func(cr *astutils.Cursor) bool {
    if named.IsParamFunc(cr.Node()) {
        // Extract name from sqlc.arg(name)
        call := cr.Node().(*ast.FuncCall)
        name := extractName(call.Args)

        cr.Replace(&ast.ParamRef{
            Number:   nextParam(),
            Location: call.Location,
        })
    }
    return true
}, nil)
```

## Important: AST Node Requirements

For parameter rewriting to work correctly, the AST must be walkable. This means:

1. All node types must have cases in `astutils/walk.go`
2. All node types must have cases in `astutils/rewrite.go`
3. New container types (like `OnDuplicateKeyUpdate`) must be traversed

### Example: OnDuplicateKeyUpdate

MySQL's `ON DUPLICATE KEY UPDATE` clause can contain `sqlc.arg()`:
```sql
INSERT INTO t (a) VALUES (sqlc.arg(val))
ON DUPLICATE KEY UPDATE a = sqlc.arg(new_val)
```

For the parameter in `ON DUPLICATE KEY UPDATE` to be found and replaced:

1. `InsertStmt` in `rewrite.go` must traverse `OnDuplicateKeyUpdate`:
```go
case *ast.InsertStmt:
    a.apply(n, "Relation", nil, n.Relation)
    a.apply(n, "Cols", nil, n.Cols)
    a.apply(n, "SelectStmt", nil, n.SelectStmt)
    a.apply(n, "OnConflictClause", nil, n.OnConflictClause)
    a.apply(n, "OnDuplicateKeyUpdate", nil, n.OnDuplicateKeyUpdate)  // Critical!
    a.apply(n, "ReturningList", nil, n.ReturningList)
    a.apply(n, "WithClause", nil, n.WithClause)
```

2. `OnDuplicateKeyUpdate` must have its own case:
```go
case *ast.OnDuplicateKeyUpdate:
    a.apply(n, "List", nil, n.List)
```

## Debugging Parameter Issues

If a `sqlc.arg()` isn't being converted to `?`:

1. Check that the containing node type has a case in `rewrite.go`
2. Check that the case traverses all child fields
3. Add debug logging to see if the node is being visited:
```go
case *ast.YourType:
    fmt.Printf("Visiting YourType with fields: %+v\n", n)
    a.apply(n, "ChildField", nil, n.ChildField)
```

## Parameter Output Format by Engine

- PostgreSQL: `$1`, `$2`, `$3`, ...
- MySQL: `?`, `?`, `?`, ...
- SQLite: `?`, `?`, `?`, ...

The format is determined by the `Dialect.Param()` method in each engine.
