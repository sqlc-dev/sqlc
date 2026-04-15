# AST Utilities Package - Claude Code Guide

This package provides utilities for traversing and transforming AST nodes.

## Key Functions

### Walk
`Walk(f Visitor, node ast.Node)` traverses the AST depth-first, calling `f.Visit()` on each node.

```go
type Visitor interface {
    Visit(node ast.Node) Visitor
}
```

**Important**: When adding new AST node types, you MUST add a case to the switch statement in `walk.go`, otherwise you'll get a panic:
```
panic: walk: unexpected node type *ast.YourNewType
```

### Apply (Rewrite)
`Apply(root ast.Node, pre, post ApplyFunc) ast.Node` traverses and optionally transforms the AST.

```go
type ApplyFunc func(*Cursor) bool
```

The `Cursor` provides:
- `Node()` - current node
- `Parent()` - parent node
- `Name()` - field name in parent
- `Index()` - index if in a list
- `Replace(node)` - replace current node

**Important**: When adding new AST node types, you MUST add a case to the switch statement in `rewrite.go`, otherwise you'll get a panic:
```
panic: Apply: unexpected node type *ast.YourNewType
```

### Search
`Search(root ast.Node, fn func(ast.Node) bool) *ast.List` finds all nodes matching a predicate.

### Join
`Join(list *ast.List, sep string) string` joins string nodes with a separator.

## Adding Support for New AST Nodes

When you create a new AST node type, you must update BOTH `walk.go` and `rewrite.go`:

### In walk.go
Add a case that walks all child nodes:
```go
case *ast.YourNewType:
    if n.ChildField != nil {
        Walk(f, n.ChildField)
    }
    if n.ChildList != nil {
        Walk(f, n.ChildList)
    }
```

For leaf nodes with no children:
```go
case *ast.YourNewType:
    // Leaf node - no children to traverse
```

### In rewrite.go
Add a case that applies to all child nodes:
```go
case *ast.YourNewType:
    a.apply(n, "ChildField", nil, n.ChildField)
    a.apply(n, "ChildList", nil, n.ChildList)
```

For leaf nodes:
```go
case *ast.YourNewType:
    // Leaf node - no children to traverse
```

## Common Patterns

### Finding All Tables in a Statement
```go
var tv tableVisitor
astutils.Walk(&tv, stmt.FromClause)
// tv.list now contains all RangeVar nodes
```

### Replacing Named Parameters
The `rewrite/parameters.go` uses Apply to replace `sqlc.arg()` calls with `ParamRef`:
```go
astutils.Apply(root, func(cr *astutils.Cursor) bool {
    if named.IsParamFunc(cr.Node()) {
        cr.Replace(&ast.ParamRef{Number: nextParam()})
    }
    return true
}, nil)
```

## Node Types That Must Be Handled

All node types in `internal/sql/ast/` must have cases in both walk.go and rewrite.go. Key MySQL-specific nodes:
- `IntervalExpr` - INTERVAL expressions
- `OnDuplicateKeyUpdate` - MySQL ON DUPLICATE KEY UPDATE
- `ParenExpr` - Parenthesized expressions
- `VariableExpr` - MySQL user variables (@var)

## Debugging Tips

If you see a panic like:
```
panic: walk: unexpected node type *ast.SomeType
```

Check that `SomeType` has a case in both `walk.go` and `rewrite.go`.
