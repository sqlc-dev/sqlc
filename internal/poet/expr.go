package poet

import (
	"go/ast"
	"go/token"
	"strconv"
)

// Ident creates an identifier expression.
func Ident(name string) *ast.Ident {
	return ast.NewIdent(name)
}

// Sel creates a selector expression (x.Sel).
func Sel(x ast.Expr, sel string) *ast.SelectorExpr {
	return &ast.SelectorExpr{X: x, Sel: ast.NewIdent(sel)}
}

// SelName creates a selector from two identifier names (pkg.Name).
func SelName(pkg, name string) *ast.SelectorExpr {
	return &ast.SelectorExpr{X: ast.NewIdent(pkg), Sel: ast.NewIdent(name)}
}

// Star creates a pointer type (*X).
func Star(x ast.Expr) *ast.StarExpr {
	return &ast.StarExpr{X: x}
}

// Addr creates an address-of expression (&X).
func Addr(x ast.Expr) *ast.UnaryExpr {
	return &ast.UnaryExpr{Op: token.AND, X: x}
}

// Deref creates a dereference expression (*X).
func Deref(x ast.Expr) *ast.StarExpr {
	return &ast.StarExpr{X: x}
}

// Index creates an index expression (X[Index]).
func Index(x, index ast.Expr) *ast.IndexExpr {
	return &ast.IndexExpr{X: x, Index: index}
}

// Slice creates a slice expression (X[Low:High]).
func Slice(x, low, high ast.Expr) *ast.SliceExpr {
	return &ast.SliceExpr{X: x, Low: low, High: high}
}

// SliceFull creates a full slice expression (X[Low:High:Max]).
func SliceFull(x, low, high, max ast.Expr) *ast.SliceExpr {
	return &ast.SliceExpr{X: x, Low: low, High: high, Max: max, Slice3: true}
}

// Call creates a function call expression.
func Call(fun ast.Expr, args ...ast.Expr) *ast.CallExpr {
	return &ast.CallExpr{Fun: fun, Args: args}
}

// CallEllipsis creates a function call with ellipsis (f(args...)).
func CallEllipsis(fun ast.Expr, args ...ast.Expr) *ast.CallExpr {
	return &ast.CallExpr{Fun: fun, Args: args, Ellipsis: 1}
}

// MethodCall creates a method call expression (recv.Method(args)).
func MethodCall(recv ast.Expr, method string, args ...ast.Expr) *ast.CallExpr {
	return &ast.CallExpr{
		Fun:  Sel(recv, method),
		Args: args,
	}
}

// Binary creates a binary expression.
func Binary(x ast.Expr, op token.Token, y ast.Expr) *ast.BinaryExpr {
	return &ast.BinaryExpr{X: x, Op: op, Y: y}
}

// Unary creates a unary expression.
func Unary(op token.Token, x ast.Expr) *ast.UnaryExpr {
	return &ast.UnaryExpr{Op: op, X: x}
}

// Paren creates a parenthesized expression ((X)).
func Paren(x ast.Expr) *ast.ParenExpr {
	return &ast.ParenExpr{X: x}
}

// TypeAssert creates a type assertion (X.(Type)).
func TypeAssert(x, typ ast.Expr) *ast.TypeAssertExpr {
	return &ast.TypeAssertExpr{X: x, Type: typ}
}

// Composite creates a composite literal ({elts}).
func Composite(typ ast.Expr, elts ...ast.Expr) *ast.CompositeLit {
	return &ast.CompositeLit{Type: typ, Elts: elts}
}

// KeyValue creates a key-value expression for composite literals.
func KeyValue(key, value ast.Expr) *ast.KeyValueExpr {
	return &ast.KeyValueExpr{Key: key, Value: value}
}

// FuncLit creates a function literal.
func FuncLit(params, results *ast.FieldList, body ...ast.Stmt) *ast.FuncLit {
	return &ast.FuncLit{
		Type: &ast.FuncType{Params: params, Results: results},
		Body: &ast.BlockStmt{List: body},
	}
}

// ArrayType creates an array type expression ([size]elt).
func ArrayType(size ast.Expr, elt ast.Expr) *ast.ArrayType {
	return &ast.ArrayType{Len: size, Elt: elt}
}

// SliceType creates a slice type expression ([]elt).
func SliceType(elt ast.Expr) *ast.ArrayType {
	return &ast.ArrayType{Elt: elt}
}

// MapType creates a map type expression (map[key]value).
func MapType(key, value ast.Expr) *ast.MapType {
	return &ast.MapType{Key: key, Value: value}
}

// ChanType creates a channel type expression.
func ChanType(dir ast.ChanDir, value ast.Expr) *ast.ChanType {
	return &ast.ChanType{Dir: dir, Value: value}
}

// FuncType creates a function type expression.
func FuncType(params, results *ast.FieldList) *ast.FuncType {
	return &ast.FuncType{Params: params, Results: results}
}

// InterfaceType creates an interface type expression.
func InterfaceType(methods ...*ast.Field) *ast.InterfaceType {
	return &ast.InterfaceType{Methods: &ast.FieldList{List: methods}}
}

// StructType creates a struct type expression.
func StructType(fields ...*ast.Field) *ast.StructType {
	return &ast.StructType{Fields: &ast.FieldList{List: fields}}
}

// Ellipsis creates an ellipsis type (...elt).
func Ellipsis(elt ast.Expr) *ast.Ellipsis {
	return &ast.Ellipsis{Elt: elt}
}

// Literals

// String creates a string literal.
func String(s string) *ast.BasicLit {
	return &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(s)}
}

// RawString creates a raw string literal.
func RawString(s string) *ast.BasicLit {
	return &ast.BasicLit{Kind: token.STRING, Value: "`" + s + "`"}
}

// Int creates an integer literal.
func Int(i int) *ast.BasicLit {
	return &ast.BasicLit{Kind: token.INT, Value: strconv.Itoa(i)}
}

// Int64 creates an int64 literal.
func Int64(i int64) *ast.BasicLit {
	return &ast.BasicLit{Kind: token.INT, Value: strconv.FormatInt(i, 10)}
}

// Float creates a float literal.
func Float(f float64) *ast.BasicLit {
	return &ast.BasicLit{Kind: token.FLOAT, Value: strconv.FormatFloat(f, 'f', -1, 64)}
}

// Nil returns the nil identifier.
func Nil() *ast.Ident {
	return ast.NewIdent("nil")
}

// True returns the true identifier.
func True() *ast.Ident {
	return ast.NewIdent("true")
}

// False returns the false identifier.
func False() *ast.Ident {
	return ast.NewIdent("false")
}

// Blank returns the blank identifier (_).
func Blank() *ast.Ident {
	return ast.NewIdent("_")
}
