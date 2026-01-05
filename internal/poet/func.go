package poet

import (
	"go/ast"
	"go/token"
)

// FuncBuilder helps build function declarations.
type FuncBuilder struct {
	name     string
	recv     *ast.FieldList
	params   *ast.FieldList
	results  *ast.FieldList
	body     []ast.Stmt
	comment  string
}

// Func creates a new function builder.
func Func(name string) *FuncBuilder {
	return &FuncBuilder{name: name}
}

// Receiver sets the receiver for a method.
func (b *FuncBuilder) Receiver(name string, typ ast.Expr) *FuncBuilder {
	b.recv = &ast.FieldList{
		List: []*ast.Field{{
			Names: []*ast.Ident{ast.NewIdent(name)},
			Type:  typ,
		}},
	}
	return b
}

// Params sets the function parameters.
func (b *FuncBuilder) Params(params ...*ast.Field) *FuncBuilder {
	b.params = &ast.FieldList{List: params}
	return b
}

// Results sets the function return types.
func (b *FuncBuilder) Results(results ...*ast.Field) *FuncBuilder {
	b.results = &ast.FieldList{List: results}
	return b
}

// ResultTypes sets the function return types from expressions.
func (b *FuncBuilder) ResultTypes(types ...ast.Expr) *FuncBuilder {
	var fields []*ast.Field
	for _, t := range types {
		fields = append(fields, &ast.Field{Type: t})
	}
	b.results = &ast.FieldList{List: fields}
	return b
}

// Body sets the function body.
func (b *FuncBuilder) Body(stmts ...ast.Stmt) *FuncBuilder {
	b.body = stmts
	return b
}

// Comment sets the doc comment for the function.
func (b *FuncBuilder) Comment(comment string) *FuncBuilder {
	b.comment = comment
	return b
}

// Build creates the function declaration.
func (b *FuncBuilder) Build() *ast.FuncDecl {
	decl := &ast.FuncDecl{
		Name: ast.NewIdent(b.name),
		Recv: b.recv,
		Type: &ast.FuncType{
			Params:  b.params,
			Results: b.results,
		},
		Body: &ast.BlockStmt{List: b.body},
	}
	if b.comment != "" {
		decl.Doc = &ast.CommentGroup{
			List: []*ast.Comment{{Text: b.comment}},
		}
	}
	return decl
}

// Param creates a function parameter field.
func Param(name string, typ ast.Expr) *ast.Field {
	if name == "" {
		return &ast.Field{Type: typ}
	}
	return &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(name)},
		Type:  typ,
	}
}

// Params creates a list of parameters with the same type.
func Params(typ ast.Expr, names ...string) *ast.Field {
	var idents []*ast.Ident
	for _, name := range names {
		idents = append(idents, ast.NewIdent(name))
	}
	return &ast.Field{
		Names: idents,
		Type:  typ,
	}
}

// Result creates a named return value field.
func Result(name string, typ ast.Expr) *ast.Field {
	if name == "" {
		return &ast.Field{Type: typ}
	}
	return &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(name)},
		Type:  typ,
	}
}

// FieldList creates an ast.FieldList from fields.
func FieldList(fields ...*ast.Field) *ast.FieldList {
	return &ast.FieldList{List: fields}
}

// Const creates a constant declaration.
func Const(name string, typ ast.Expr, value ast.Expr) *ast.GenDecl {
	spec := &ast.ValueSpec{
		Names:  []*ast.Ident{ast.NewIdent(name)},
		Values: []ast.Expr{value},
	}
	if typ != nil {
		spec.Type = typ
	}
	return &ast.GenDecl{
		Tok:   token.CONST,
		Specs: []ast.Spec{spec},
	}
}

// ConstGroup creates a grouped constant declaration.
func ConstGroup(specs ...*ast.ValueSpec) *ast.GenDecl {
	var astSpecs []ast.Spec
	for _, s := range specs {
		astSpecs = append(astSpecs, s)
	}
	return &ast.GenDecl{
		Tok:    token.CONST,
		Lparen: 1,
		Specs:  astSpecs,
	}
}

// ConstSpec creates a constant specification for use in ConstGroup.
func ConstSpec(name string, typ ast.Expr, value ast.Expr) *ast.ValueSpec {
	spec := &ast.ValueSpec{
		Names:  []*ast.Ident{ast.NewIdent(name)},
		Values: []ast.Expr{value},
	}
	if typ != nil {
		spec.Type = typ
	}
	return spec
}

// Var creates a variable declaration.
func Var(name string, typ ast.Expr, value ast.Expr) *ast.GenDecl {
	spec := &ast.ValueSpec{
		Names: []*ast.Ident{ast.NewIdent(name)},
	}
	if typ != nil {
		spec.Type = typ
	}
	if value != nil {
		spec.Values = []ast.Expr{value}
	}
	return &ast.GenDecl{
		Tok:   token.VAR,
		Specs: []ast.Spec{spec},
	}
}

// VarGroup creates a grouped variable declaration.
func VarGroup(specs ...*ast.ValueSpec) *ast.GenDecl {
	var astSpecs []ast.Spec
	for _, s := range specs {
		astSpecs = append(astSpecs, s)
	}
	return &ast.GenDecl{
		Tok:    token.VAR,
		Lparen: 1,
		Specs:  astSpecs,
	}
}

// VarSpec creates a variable specification for use in VarGroup.
func VarSpec(name string, typ ast.Expr, value ast.Expr) *ast.ValueSpec {
	spec := &ast.ValueSpec{
		Names: []*ast.Ident{ast.NewIdent(name)},
	}
	if typ != nil {
		spec.Type = typ
	}
	if value != nil {
		spec.Values = []ast.Expr{value}
	}
	return spec
}
