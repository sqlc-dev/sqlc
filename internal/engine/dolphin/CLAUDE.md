# Dolphin Engine (MySQL) - Claude Code Guide

The dolphin engine handles MySQL parsing and AST conversion using the TiDB parser.

## Architecture

### Parser Flow
```
SQL String → TiDB Parser → TiDB AST → sqlc AST → Analysis/Codegen
```

### Key Files
- `convert.go` - Converts TiDB AST nodes to sqlc AST nodes
- `format.go` - MySQL-specific formatting (identifiers, types, parameters)
- `parse.go` - Entry point for parsing MySQL SQL

## TiDB Parser

The TiDB parser (`github.com/pingcap/tidb/pkg/parser`) is used for MySQL parsing:

```go
import (
    pcast "github.com/pingcap/tidb/pkg/parser/ast"
    "github.com/pingcap/tidb/pkg/parser/mysql"
    "github.com/pingcap/tidb/pkg/parser/types"
)
```

### Common TiDB Types
- `pcast.SelectStmt`, `pcast.InsertStmt`, etc. - Statement types
- `pcast.ColumnNameExpr` - Column reference
- `pcast.FuncCallExpr` - Function call
- `pcast.BinaryOperationExpr` - Binary expression
- `pcast.VariableExpr` - MySQL user variable (@var)
- `pcast.Join` - JOIN clause with Left, Right, On, Using

## Conversion Pattern

Each TiDB node type has a corresponding converter method:

```go
func (c *cc) convertSelectStmt(n *pcast.SelectStmt) *ast.SelectStmt {
    return &ast.SelectStmt{
        FromClause:  c.convertTableRefsClause(n.From),
        WhereClause: c.convert(n.Where),
        // ...
    }
}
```

The main `convert()` method dispatches to specific converters:
```go
func (c *cc) convert(node pcast.Node) ast.Node {
    switch n := node.(type) {
    case *pcast.SelectStmt:
        return c.convertSelectStmt(n)
    case *pcast.InsertStmt:
        return c.convertInsertStmt(n)
    // ...
    }
}
```

## Key Conversions

### Column References
```go
func (c *cc) convertColumnNameExpr(n *pcast.ColumnNameExpr) *ast.ColumnRef {
    var items []ast.Node
    if schema := n.Name.Schema.String(); schema != "" {
        items = append(items, NewIdentifier(schema))
    }
    if table := n.Name.Table.String(); table != "" {
        items = append(items, NewIdentifier(table))
    }
    items = append(items, NewIdentifier(n.Name.Name.String()))
    return &ast.ColumnRef{Fields: &ast.List{Items: items}}
}
```

### JOINs
```go
func (c *cc) convertJoin(n *pcast.Join) *ast.List {
    if n.Right != nil && n.Left != nil {
        return &ast.List{
            Items: []ast.Node{&ast.JoinExpr{
                Jointype:    ast.JoinType(n.Tp),
                Larg:        c.convert(n.Left),
                Rarg:        c.convert(n.Right),
                Quals:       c.convert(n.On),
                UsingClause: convertUsing(n.Using),
            }},
        }
    }
    // No join - just return tables
    // ...
}
```

### MySQL User Variables
MySQL user variables (`@var`) are different from sqlc's `@param` syntax:
```go
func (c *cc) convertVariableExpr(n *pcast.VariableExpr) ast.Node {
    // Use VariableExpr to preserve as-is (NOT A_Expr which would be treated as sqlc param)
    return &ast.VariableExpr{
        Name:     n.Name,
        Location: n.OriginTextPosition(),
    }
}
```

### Type Casts (CAST AS)
```go
func (c *cc) convertFuncCastExpr(n *pcast.FuncCastExpr) ast.Node {
    typeName := types.TypeStr(n.Tp.GetType())
    // Handle UNSIGNED/SIGNED specially
    if typeName == "bigint" {
        if mysql.HasUnsignedFlag(n.Tp.GetFlag()) {
            typeName = "bigint unsigned"
        } else {
            typeName = "bigint signed"
        }
    }
    return &ast.TypeCast{
        Arg:      c.convert(n.Expr),
        TypeName: &ast.TypeName{Name: typeName},
    }
}
```

### Column Definitions
```go
func convertColumnDef(def *pcast.ColumnDef) *ast.ColumnDef {
    typeName := &ast.TypeName{Name: types.TypeToStr(def.Tp.GetType(), def.Tp.GetCharset())}

    // Only add Typmods for types where length is meaningful
    tp := def.Tp.GetType()
    flen := def.Tp.GetFlen()
    switch tp {
    case mysql.TypeVarchar, mysql.TypeString, mysql.TypeVarString:
        if flen >= 0 {
            typeName.Typmods = &ast.List{
                Items: []ast.Node{&ast.Integer{Ival: int64(flen)}},
            }
        }
    // Don't add for DATETIME, TIMESTAMP - internal flen is not user-specified
    }
    // ...
}
```

### Multi-Table DELETE
MySQL supports `DELETE t1, t2 FROM t1 JOIN t2 ...`:
```go
func (c *cc) convertDeleteStmt(n *pcast.DeleteStmt) *ast.DeleteStmt {
    if n.IsMultiTable && n.Tables != nil {
        // Convert targets (t1.*, t2.*)
        targets := &ast.List{}
        for _, table := range n.Tables.Tables {
            // Build ColumnRef for each target
        }
        stmt.Targets = targets

        // Preserve JOINs in FromClause
        stmt.FromClause = c.convertTableRefsClause(n.TableRefs).Items[0]
    } else {
        // Single-table DELETE
        stmt.Relations = c.convertTableRefsClause(n.TableRefs)
    }
}
```

## MySQL-Specific Formatting

### format.go
```go
func (p *Parser) TypeName(ns, name string) string {
    switch name {
    case "bigint unsigned":
        return "UNSIGNED"
    case "bigint signed":
        return "SIGNED"
    }
    return name
}

func (p *Parser) Param(n int) string {
    return "?"  // MySQL uses ? for all parameters
}
```

## Common Issues and Solutions

### Issue: Panic in Walk/Apply
**Cause**: New AST node type not handled in `astutils/walk.go` or `astutils/rewrite.go`
**Solution**: Add case for the node type in both files

### Issue: sqlc.arg() not converted in ON DUPLICATE KEY UPDATE
**Cause**: `InsertStmt` case in `rewrite.go` didn't traverse `OnDuplicateKeyUpdate`
**Solution**: Add `a.apply(n, "OnDuplicateKeyUpdate", nil, n.OnDuplicateKeyUpdate)`

### Issue: MySQL @variable being treated as parameter
**Cause**: Converting `VariableExpr` to `A_Expr` with `@` operator
**Solution**: Use `ast.VariableExpr` instead, which is not detected by `named.IsParamSign()`

### Issue: Type length appearing incorrectly (e.g., datetime(39))
**Cause**: Using internal `flen` for all types
**Solution**: Only populate `Typmods` for types where length is user-specified (varchar, char, etc.)

## Testing

### TestFormat
Tests that SQL can be:
1. Parsed
2. Formatted back to SQL
3. Re-parsed
4. Re-formatted to match

### TestReplay
Tests the full sqlc pipeline:
1. Parse schema and queries
2. Analyze
3. Generate code
4. Compare with expected output
