package poet

import (
	"go/ast"
	"go/token"
)

// InterfaceBuilder helps build interface type declarations.
type InterfaceBuilder struct {
	name    string
	comment string
	methods []*ast.Field
}

// Interface creates a new interface builder.
func Interface(name string) *InterfaceBuilder {
	return &InterfaceBuilder{name: name}
}

// Comment sets the doc comment for the interface.
func (b *InterfaceBuilder) Comment(comment string) *InterfaceBuilder {
	b.comment = comment
	return b
}

// Method adds a method to the interface.
func (b *InterfaceBuilder) Method(name string, params, results *ast.FieldList) *InterfaceBuilder {
	method := &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(name)},
		Type: &ast.FuncType{
			Params:  params,
			Results: results,
		},
	}
	b.methods = append(b.methods, method)
	return b
}

// MethodWithComment adds a method with a doc comment to the interface.
func (b *InterfaceBuilder) MethodWithComment(name string, params, results *ast.FieldList, comment string) *InterfaceBuilder {
	method := &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(name)},
		Type: &ast.FuncType{
			Params:  params,
			Results: results,
		},
	}
	if comment != "" {
		method.Doc = &ast.CommentGroup{
			List: []*ast.Comment{{Text: comment}},
		}
	}
	b.methods = append(b.methods, method)
	return b
}

// Build creates the interface type declaration.
func (b *InterfaceBuilder) Build() *ast.GenDecl {
	spec := &ast.TypeSpec{
		Name: ast.NewIdent(b.name),
		Type: &ast.InterfaceType{
			Methods: &ast.FieldList{List: b.methods},
		},
	}
	decl := &ast.GenDecl{
		Tok:   token.TYPE,
		Specs: []ast.Spec{spec},
	}
	if b.comment != "" {
		decl.Doc = &ast.CommentGroup{
			List: []*ast.Comment{{Text: b.comment}},
		}
	}
	return decl
}

// StructBuilder helps build struct type declarations.
type StructBuilder struct {
	name    string
	comment string
	fields  []*ast.Field
}

// Struct creates a new struct builder.
func Struct(name string) *StructBuilder {
	return &StructBuilder{name: name}
}

// Comment sets the doc comment for the struct.
func (b *StructBuilder) Comment(comment string) *StructBuilder {
	b.comment = comment
	return b
}

// Field adds a field to the struct.
func (b *StructBuilder) Field(name string, typ ast.Expr) *StructBuilder {
	field := &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(name)},
		Type:  typ,
	}
	b.fields = append(b.fields, field)
	return b
}

// FieldWithTag adds a field with a struct tag to the struct.
func (b *StructBuilder) FieldWithTag(name string, typ ast.Expr, tag string) *StructBuilder {
	field := &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(name)},
		Type:  typ,
	}
	if tag != "" {
		field.Tag = &ast.BasicLit{Kind: token.STRING, Value: "`" + tag + "`"}
	}
	b.fields = append(b.fields, field)
	return b
}

// FieldWithComment adds a field with a doc comment to the struct.
func (b *StructBuilder) FieldWithComment(name string, typ ast.Expr, comment string) *StructBuilder {
	field := &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(name)},
		Type:  typ,
	}
	if comment != "" {
		field.Doc = &ast.CommentGroup{
			List: []*ast.Comment{{Text: comment}},
		}
	}
	b.fields = append(b.fields, field)
	return b
}

// FieldFull adds a field with all options.
func (b *StructBuilder) FieldFull(name string, typ ast.Expr, tag, comment string) *StructBuilder {
	field := &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(name)},
		Type:  typ,
	}
	if tag != "" {
		field.Tag = &ast.BasicLit{Kind: token.STRING, Value: "`" + tag + "`"}
	}
	if comment != "" {
		field.Doc = &ast.CommentGroup{
			List: []*ast.Comment{{Text: comment}},
		}
	}
	b.fields = append(b.fields, field)
	return b
}

// AddField adds a pre-built field to the struct.
func (b *StructBuilder) AddField(field *ast.Field) *StructBuilder {
	b.fields = append(b.fields, field)
	return b
}

// Build creates the struct type declaration.
func (b *StructBuilder) Build() *ast.GenDecl {
	spec := &ast.TypeSpec{
		Name: ast.NewIdent(b.name),
		Type: &ast.StructType{
			Fields: &ast.FieldList{List: b.fields},
		},
	}
	decl := &ast.GenDecl{
		Tok:   token.TYPE,
		Specs: []ast.Spec{spec},
	}
	if b.comment != "" {
		decl.Doc = &ast.CommentGroup{
			List: []*ast.Comment{{Text: b.comment}},
		}
	}
	return decl
}

// TypeAlias creates a type alias declaration (type Name = Alias).
func TypeAlias(name string, typ ast.Expr) *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name:   ast.NewIdent(name),
				Assign: 1, // Non-zero means alias
				Type:   typ,
			},
		},
	}
}

// TypeDef creates a type definition (type Name underlying).
func TypeDef(name string, typ ast.Expr) *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(name),
				Type: typ,
			},
		},
	}
}

// TypeDefWithComment creates a type definition with a comment.
func TypeDefWithComment(name string, typ ast.Expr, comment string) *ast.GenDecl {
	decl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(name),
				Type: typ,
			},
		},
	}
	if comment != "" {
		decl.Doc = &ast.CommentGroup{
			List: []*ast.Comment{{Text: comment}},
		}
	}
	return decl
}
