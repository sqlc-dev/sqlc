# Named Parameters Package - Claude Code Guide

This package provides utilities for identifying sqlc's named parameter syntax.

## Named Parameter Styles

sqlc supports two styles of named parameters:

### 1. Function-style: `sqlc.arg(name)`, `sqlc.narg(name)`, `sqlc.slice(name)`
Identified by `IsParamFunc()`:
```go
func IsParamFunc(node ast.Node) bool {
    call, ok := node.(*ast.FuncCall)
    if !ok {
        return false
    }
    return call.Func.Schema == "sqlc" &&
           (call.Func.Name == "arg" || call.Func.Name == "narg" || call.Func.Name == "slice")
}
```

### 2. At-sign style: `@param_name` (PostgreSQL only)
Identified by `IsParamSign()`:
```go
func IsParamSign(node ast.Node) bool {
    expr, ok := node.(*ast.A_Expr)
    return ok && astutils.Join(expr.Name, ".") == "@"
}
```

## Important Distinction: sqlc @param vs MySQL @variable

**sqlc named parameters** (`@param` in PostgreSQL queries):
- Represented as `A_Expr` with `Kind=A_Expr_Kind_OP` and `Name=["@"]`
- Detected by `IsParamSign()`
- Replaced with positional parameters (`$1`, `$2` for PostgreSQL, `?` for MySQL)

**MySQL user variables** (`@user_id` in MySQL queries):
- Represented as `VariableExpr`
- NOT detected by `IsParamSign()` (it checks for `A_Expr`, not `VariableExpr`)
- Preserved as-is in the output SQL

This distinction is critical:
```sql
-- PostgreSQL with sqlc @param syntax:
SELECT * FROM users WHERE id = @user_id
-- Becomes: SELECT * FROM users WHERE id = $1

-- MySQL with user variable:
SELECT * FROM users WHERE id != @user_id
-- Stays: SELECT * FROM users WHERE id != @user_id
```

## Usage in Parameter Rewriting

The `rewrite/parameters.go` package uses these functions to find and replace named parameters:

```go
// Find all named parameters
params := astutils.Search(root, func(node ast.Node) bool {
    return named.IsParamFunc(node) || named.IsParamSign(node)
})

// Replace with positional parameters
astutils.Apply(root, func(cr *astutils.Cursor) bool {
    if named.IsParamFunc(cr.Node()) || named.IsParamSign(cr.Node()) {
        cr.Replace(&ast.ParamRef{Number: nextParam()})
    }
    return true
}, nil)
```

## Converting MySQL @variable Correctly

When converting TiDB's `VariableExpr` in `dolphin/convert.go`:

```go
// CORRECT - preserves MySQL user variable as-is
func (c *cc) convertVariableExpr(n *pcast.VariableExpr) ast.Node {
    return &ast.VariableExpr{
        Name:     n.Name,
        Location: n.OriginTextPosition(),
    }
}

// WRONG - would be treated as sqlc named parameter
func (c *cc) convertVariableExpr(n *pcast.VariableExpr) ast.Node {
    return &ast.A_Expr{
        Kind: ast.A_Expr_Kind_OP,
        Name: &ast.List{Items: []ast.Node{&ast.String{Str: "@"}}},
        Rexpr: &ast.String{Str: n.Name},
    }
}
```
